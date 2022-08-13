package tts

import (
	"context"
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

var hashes = make(map[string]bool)

type TTS struct {
	sync.Mutex
	player  player.Player
	playing chan player.Player
}

func NewTTS() *TTS {
	playing := make(chan player.Player)
	return &TTS{
		playing: playing,
	}
}

func (t *TTS) Playing() chan player.Player {
	return t.playing
}

func (t *TTS) Run() {
	for player := range t.Playing() {
		if err := player.Play(); err != nil {
			log.Printf("Error: %s", player.Play())
		}
	}
}

func (t *TTS) Stop() {
	var wg sync.WaitGroup
	//Delete cache files
	for key := range hashes {
		wg.Add(1)
		go func(key string, wg *sync.WaitGroup) {
			defer wg.Done()
			os.Remove(key)
		}(key, &wg)
	}

	wg.Wait()
	close(t.playing)
}

func (t *TTS) Skip() {
	t.Lock()
	defer t.Unlock()
	if t.player != nil {
		t.player.Skip()
	}
}

//Google wants a cache system to don't ban this client, so we need to add it
func (t *TTS) Play(lang, text string) error {
	ctx, cancel := context.WithCancel(context.Background())
	player := player.NewNativePlayer(ctx, cancel)

	hashText := hash(text)
	t.Lock()
	hashes[hashText] = true
	t.Unlock()

	var err error
	player.Buff, err = os.ReadFile(hashText)
	if err == nil {
		go func() {
			t.playing <- player
			//Important to set the player after the channel send the message
			t.Lock()
			t.player = player
			t.Unlock()
		}()

		return nil
	}

	url := fmt.Sprintf("https://translate.google.com/translate_tts?ie=UTF-8&tl=%s&client=tw-ob&q=%s", lang, url.QueryEscape(text))
	resp, err := http.Get(url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	player.Buff, err = io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	file, err := os.Create(hashText)
	defer file.Close()

	if err != nil {
		return err
	}

	n, err := file.Write(player.Buff)

	if err != nil {
		return err
	}

	if len(player.Buff) != n {
		return errors.New("error on len buff write")
	}
	go func() {
		t.playing <- player
		//Important to set the player after the channel send the message
		t.Lock()
		fmt.Println("Fijando")
		t.player = player
		t.Unlock()
	}()

	return nil
}
