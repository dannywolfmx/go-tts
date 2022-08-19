package player

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

type Native struct {
	Buff   []byte
	ctx    context.Context
	cancel context.CancelFunc
	text   string
	otoCtx *oto.Context
}

func NewNativePlayer(text string, otoCtx *oto.Context) *Native {
	ctx, cancel := context.WithCancel(context.Background())
	return &Native{
		ctx:    ctx,
		cancel: cancel,
		text:   text,
		Buff:   []byte{},
		otoCtx: otoCtx,
	}
}

func (n *Native) Play() error {

	reader := bytes.NewReader(n.Buff)

	// Decode file
	decodedMp3, err := mp3.NewDecoder(reader)
	if err != nil {
		n.cancel()
		return fmt.Errorf("mp3.NewDecoder failed: %s", err)
	}

	player := n.otoCtx.NewPlayer(decodedMp3)

	player.Play()

	exit := make(chan error, 1)
	go func() {
		<-n.ctx.Done()
		exit <- player.Close()
	}()

	for {
		if player.IsPlaying() {
			time.Sleep(time.Millisecond)
		} else {
			n.cancel()
			break
		}
	}

	return <-exit
}

func (n *Native) Stop() {
	n.cancel()
}

func (n *Native) GetText() string {
	return n.text
}
