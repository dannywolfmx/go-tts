package player

import (
	"bytes"
	"fmt"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

type Native struct {
	Buff   []byte
	text   string
	otoCtx *oto.Context
	player oto.Player
	pause  bool
}

func NewNativePlayer(text string, otoCtx *oto.Context) *Native {
	return &Native{
		text:   text,
		Buff:   []byte{},
		otoCtx: otoCtx,
	}
}

func (n *Native) Play() error {
	n.pause = false
	if n.player == nil {
		reader := bytes.NewReader(n.Buff)

		// Decode file
		decodedMp3, err := mp3.NewDecoder(reader)
		if err != nil {
			return fmt.Errorf("mp3.NewDecoder failed: %s", err)
		}

		n.player = n.otoCtx.NewPlayer(decodedMp3)
	}

	n.player.Play()

	for n.player.IsPlaying() || n.pause {
		time.Sleep(time.Millisecond)
	}

	return n.player.Close()
}

func (n *Native) Pause() {
	if n.player != nil {
		n.pause = true
		n.player.Pause()
	}
}

func (n *Native) Stop() {
	if n.player != nil {
		n.player.Close()
	}
}

func (n *Native) GetText() string {
	return n.text
}
