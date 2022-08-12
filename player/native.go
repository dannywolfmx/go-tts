package player

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

var wg sync.WaitGroup

type Native struct {
	Buff   []byte
	ctx    context.Context
	cancel context.CancelFunc
}

func NewNativePlayer(ctx context.Context, cancel context.CancelFunc) *Native {
	return &Native{
		ctx:    ctx,
		cancel: cancel,
		Buff:   []byte{},
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

	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
	numOfChannels := 2

	// Bytes used by a channel to represent one sample. Either 1 or 2 (usually 2).
	audioBitDepth := 2

	// Remember that you should **not** create more than one context
	otoCtx, readyChan, err := oto.NewContext(decodedMp3.SampleRate(), numOfChannels, audioBitDepth)
	if err != nil {
		n.cancel()
		return fmt.Errorf("oto.NewContext failed: %s", err)
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan

	player := otoCtx.NewPlayer(decodedMp3)

	player.Play()

	wg.Add(1)
	go func() {
		for {
			select {
			case <-n.ctx.Done():
				wg.Done()
				err = player.Close()
				if err != nil {
					panic("player.Close failed: " + err.Error())
				}
				return
			default:
				if player.IsPlaying() {
					time.Sleep(time.Millisecond)
				} else {
					n.cancel()
				}
			}
		}
	}()
	wg.Wait()
	return nil
}

func (n *Native) Skip() {
	n.cancel()
}
