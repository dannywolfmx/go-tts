package tts

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/dannywolfmx/go-tts/player"
	"github.com/hajimehoshi/oto/v2"
)

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

func (t *TTS) Add(text string) {
	//Add the player to the Queue
	t.mu.Lock()
	t.queue = append(t.queue, player.NewNativePlayer(text, t.otoCtx))
	t.mu.Unlock()

	//Just autoplay when the queue is 1, to prevent
	// play multiples times, at the same time
	if t.autoplay && t.QueueLen() == 1 {
		t.Next()
	}
}

func (t *TTS) Next() {
	t.mu.Lock()
	fmt.Println("Next")
	fmt.Println("Lock")

	if t.QueueLen() == 0 {
		t.mu.Unlock()
		return
	}

	t.isNextActive = true
	if t.QueueLen() == 1 {
		//Stop de song
		t.queue[0].Stop()

		//clear the queue
		t.queue = nil
		t.mu.Unlock()
		fmt.Println("UNlock")
		return
	}

	//Stop de song
	t.queue[0].Stop()

	//Pass the next
	t.queue = t.queue[1:]
	player := t.queue[0]
	t.mu.Unlock()

	fmt.Println("UNlock")

	t.play(player)
}

func (t *TTS) OnPlayerEnds(action func(string)) {
	t.onPlayerEnds = action
}

func (t *TTS) OnPlayerStart(action func(string)) {
	t.onPlayerStart = action
}

func (t *TTS) Pause() {
	t.autoplay = false
	if t.QueueLen() > 0 {
		//Stop de song
		t.queue[0].Pause()
	}
}

func (t *TTS) Play() {
	t.autoplay = true

	if t.QueueLen() > 0 {
		t.Next()
	}
}

// [0, 1, 2]
func (t *TTS) play(player *player.Native) {
	//EmitEvent
	if t.onPlayerStart != nil {
		t.onPlayerStart(player.GetText())
	}

	var err error
	if player.Buff, err = getSpeech(player.GetText(), t.lang); err != nil {
		log.Printf("Error reading voice: %s", err)
	}

	//Play the song
	if err := player.Play(); err != nil {
		log.Printf("Error reading voice: %s", err)
	}

	if t.onPlayerEnds != nil {
		t.onPlayerEnds(player.GetText())
	}

	t.mu.Lock()
	if t.isNextActive {
		t.isNextActive = false
		t.mu.Unlock()
		return
	}
	t.mu.Unlock()
	t.Next()
}

func (t *TTS) QueueLen() int {
	return len(t.queue)
}

func (t *TTS) Stop() {
	t.mu.Lock()
	if t.QueueLen() == 0 {
		return
	}
	t.queue[0].Stop()
	t.queue = nil
	t.mu.Unlock()
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
