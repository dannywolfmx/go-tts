package tts

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/dannywolfmx/go-tts/player"
)

var hashes = make(map[string]bool)

var c = make(chan player.Player)

type TTS struct{}

func NewTTS() *TTS {
	return &TTS{}
}

func (t *TTS) Run() {
	for player := range c {
		player.Play()
	}

}

func (t *TTS) Stop() {
	//Delete cache files
	for key := range hashes {
		os.Remove(key)
	}

	close(c)
}

//Google wants a cache system to don't ban this client, so we need to add it
func (t *TTS) Play(lang, text string) (player.Player, error) {
	ctx, cancel := context.WithCancel(context.Background())
	player := player.NewNativePlayer(ctx, cancel)

	hashText := hash(text)
	hashes[hashText] = true

	var err error
	player.Buff, err = os.ReadFile(hashText)
	if err == nil {
		go func() {
			c <- player
		}()

		return player, nil
	}

	url := fmt.Sprintf("https://translate.google.com/translate_tts?ie=UTF-8&tl=%s&client=tw-ob&q=%s", lang, url.QueryEscape(text))
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	player.Buff, err = io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	file, err := os.Create(hashText)
	defer file.Close()

	if err != nil {
		fmt.Println("Entre")
		return nil, err
	}

	n, err := file.Write(player.Buff)

	if err != nil {
		return nil, err
	}

	if len(player.Buff) != n {
		return nil, errors.New("error on len buff write")
	}
	go func() {
		c <- player
	}()

	return player, nil
}
