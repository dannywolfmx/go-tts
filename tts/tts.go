package tts

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/dannywolfmx/go-tts/player"
	"github.com/hajimehoshi/oto/v2"
)

var Empty struct{}
var hashes = make(map[string]struct{})

type TTS struct {
	//The actual player in play mode
	lang          string
	queue         []*player.Native
	onPlayerStart func(text string)
	otoCtx        *oto.Context
}

func NewTTS(lang string) *TTS {
	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
	numOfChannels := 2

	// Bytes used by a channel to represent one sample. Either 1 or 2 (usually 2).
	audioBitDepth := 2

	// Remember that you should **not** create more than one context
	otoCtx, readyChan, err := oto.NewContext(24000, numOfChannels, audioBitDepth)
	if err != nil {
		log.Fatalf("Error on read oto: %s", err)
		return nil
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan

	return &TTS{
		lang:   lang,
		otoCtx: otoCtx,
	}
}

func (t *TTS) OnPlayerStart(action func(string)) {
	t.onPlayerStart = action
}

func (t *TTS) Add(text string) {
	//Add the player to the Queue
	t.queue = append(t.queue, player.NewNativePlayer(text, t.otoCtx))

	//The was empty and need to play
	if len(t.queue) == 1 {
		t.play()
	}
}

func (t *TTS) Next() {
	if len(t.queue) == 0 {
		return
	}

	if len(t.queue) == 1 {
		//Stop de song
		t.queue[0].Stop()

		//clear the queue
		t.queue = nil
		return
	}

	//Stop de song
	t.queue[0].Stop()

	//Pass the next
	t.queue = t.queue[1:]

	t.play()
}

func (t *TTS) Stop() {
	if len(t.queue) == 0 {
		return
	}
	t.queue[0].Stop()
	t.queue = nil
}

func (t *TTS) CleanCache() {
	for key := range hashes {
		if err := os.Remove(key); err != nil {
			log.Printf("On clean cache: %s \n", err)
		}
	}
}

func (t *TTS) EmitEvents(p *player.Native) {
	if t.onPlayerStart != nil {
		t.onPlayerStart(p.GetText())
	}
}

func (t *TTS) play() error {
	player := t.queue[0]

	//EmitEvent
	t.EmitEvents(player)

	//Play the song
	if err := play(t.lang, player); err != nil {
		return err
	}

	//Continue with the next
	t.Next()

	return nil
}

//Google wants a cache system to don't ban this client, so we need to add it
func play(lang string, player *player.Native) error {
	hash := GetHash(player.GetText())

	var err error
	if player.Buff, err = GetSpeech(player.GetText(), lang, hash); err != nil {
		return err
	}

	hashes[hash] = Empty

	return player.Play()
}

func GetHash(text string) string {
	hasher := sha1.New()

	hasher.Write([]byte(text))
	buff := hasher.Sum(nil)

	return hex.EncodeToString(buff)
}

func GetSpeech(text, lang, hash string) ([]byte, error) {
	var err error
	buff := []byte{}

	if buff, err = os.ReadFile(hash); err == nil {
		return buff, err
	}

	url := fmt.Sprintf("https://translate.google.com/translate_tts?ie=UTF-8&tl=%s&client=tw-ob&q=%s", lang, url.QueryEscape(text))
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	buff, err = io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	file, err := os.Create(hash)
	defer file.Close()

	if err != nil {
		return buff, err
	}

	n, err := file.Write(buff)

	if err != nil {
		return nil, err
	}

	if len(buff) != n {
		return buff, errors.New("error on len buff write")
	}

	return buff, err
}
