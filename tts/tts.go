package tts

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/dannywolfmx/go-tts/player"
)

var Empty struct{}
var hashes = make(map[string]struct{})

type TTS struct {
	sync.Mutex
	player  player.Player
	playing chan player.Player
	lang    string
	Done    chan bool
}

func NewTTS(lang string) *TTS {
	playing := make(chan player.Player)
	done := make(chan bool)
	return &TTS{
		playing: playing,
		lang:    lang,
		Done:    done,
	}
}

func (t *TTS) GetSpeech(p *player.Native, text, hash string) ([]byte, error) {
	var err error
	buff := []byte{}

	if buff, err = os.ReadFile(hash); err == nil {
		return buff, err
	}

	url := fmt.Sprintf("https://translate.google.com/translate_tts?ie=UTF-8&tl=%s&client=tw-ob&q=%s", t.lang, url.QueryEscape(text))
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

//Google wants a cache system to don't ban this client, so we need to add it
func (t *TTS) Play(text string) error {
	ctx, cancel := context.WithCancel(context.Background())
	player := player.NewNativePlayer(ctx, cancel, text)

	hash := GetHash(text)

	var err error
	if player.Buff, err = t.GetSpeech(player, text, hash); err != nil {
		return err
	}
	//Save hash after GetSpeech

	t.Lock()
	hashes[hash] = Empty
	t.Unlock()

	t.play(player)

	return nil
}

func (t *TTS) play(p player.Player) {
	go func(p player.Player) {
		t.playing <- p
		//Important to set the player after the channel send the message
		t.Lock()
		t.player = p
		t.Unlock()
	}(p)
}

func (t *TTS) Playing() chan player.Player {
	return t.playing
}

func (t *TTS) Run() {
	for p := range t.Playing() {
		select {
		case <-t.Done:
			return
		default:
			if err := p.Play(); err != nil {
				log.Printf("Error: %s", p.Play())
			}
		}
	}
}

func (t *TTS) Skip() {
	t.Lock()
	defer t.Unlock()
	if t.player != nil {
		t.player.Skip()
	}
}

func (t *TTS) Stop() {
	var wg sync.WaitGroup
	//Delete cache files
	t.Lock()
	for key := range hashes {
		wg.Add(1)
		go func(key string, wg *sync.WaitGroup) {
			defer wg.Done()
			if err := os.Remove(key); err != nil {
				log.Printf("Error on delete %s", err)
			}
		}(key, &wg)
	}
	t.Unlock()
	wg.Wait()
	t.Done <- true
}

func GetHash(text string) string {
	hasher := sha1.New()

	hasher.Write([]byte(text))
	buff := hasher.Sum(nil)

	return hex.EncodeToString(buff)
}
