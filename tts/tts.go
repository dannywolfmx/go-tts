package tts

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/dannywolfmx/go-tts/player"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

var queue = make(chan *player.Native)
var actualPlaying *player.Native

type TTS struct {
	//The actual player in play mode
	autoplay      bool
	lang          string
	queue         []*player.Native
	onPlayerEnds  func(text string)
	onPlayerStart func(text string)
	otoCtx        *oto.Context
	isNextActive  bool
	mu            sync.Mutex
}

func NewTTS(lang string, sampleRate int) *TTS {
	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
	numOfChannels := 2

	// Bytes used by a channel to represent one sample. Either 1 or 2 (usually 2).
	audioBitDepth := 2

	// Remember that you should **not** create more than one context
	otoCtx, readyChan, err := oto.NewContext(sampleRate, numOfChannels, audioBitDepth)

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

func (t *TTS) Add(text string) error {

	otoPlayer, err := getVoiceText(text, t.lang, t.otoCtx)
	if err != nil {
		return err
	}

	queue <- player.NewNativePlayer(text, otoPlayer)

	return nil
}

func (t *TTS) Continue() {
	t.otoCtx.Resume()
}

func getVoiceText(text, lang string, ctx *oto.Context) (oto.Player, error) {
	buff, err := getSpeech(text, lang)
	if err != nil {
		return nil, fmt.Errorf("error reading voice: %s", err)
	}

	reader := bytes.NewReader(buff)

	// Decode file
	decoder, err := mp3.NewDecoder(reader)

	if err != nil {
		return nil, fmt.Errorf("mp3.NewDecoder failed: %s", err)
	}

	return ctx.NewPlayer(decoder), nil
}

func (t *TTS) Next() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if actualPlaying != nil {
		actualPlaying.Stop()
	}
}

func (t *TTS) OnPlayerEnds(action func(string)) {
	t.onPlayerEnds = action
}

func (t *TTS) OnPlayerStart(action func(string)) {
	t.onPlayerStart = action
}

func (t *TTS) Pause() {
	t.otoCtx.Suspend()
}

func (t *TTS) Play() {
	go play(t)
}

func play(t *TTS) {
	for player := range queue {
		t.mu.Lock()
		actualPlaying = player
		t.mu.Unlock()
		//EmitEvent
		if t.onPlayerStart != nil {
			t.onPlayerStart(player.GetText())
		}

		//Play the song
		if err := player.Play(); err != nil {
			log.Printf("Error reading voice: %s", err)
		}

		if t.onPlayerEnds != nil {
			t.onPlayerEnds(player.GetText())
		}
	}
}

func (t *TTS) Stop() {

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
