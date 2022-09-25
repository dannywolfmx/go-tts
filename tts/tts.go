package tts

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/dannywolfmx/go-tts/player"
	"github.com/hajimehoshi/oto/v2"
)

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

func (t *TTS) EmitEvents(p *player.Native) {
	if t.onPlayerStart != nil {
		t.onPlayerStart(p.GetText())
	}
}

func (t *TTS) play() {
	player := t.queue[0]

	//EmitEvent
	t.EmitEvents(player)

	var err error
	if player.Buff, err = getSpeech(player.GetText(), t.lang); err != nil {
		log.Printf("Error reading voice: %s", err)
	}

	//Play the song
	if err := player.Play(); err != nil {
		log.Printf("Error reading voice: %s", err)
	}

	//Continue with the next
	t.Next()
}

func getSpeech(text, lang string) ([]byte, error) {

	url := fmt.Sprintf("https://translate.google.com/translate_tts?ie=UTF-8&tl=%s&client=tw-ob&q=%s", lang, url.QueryEscape(text))
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
