package player

import (
	"sync"
	"time"

	"github.com/hajimehoshi/oto/v2"
)

type Native struct {
	text   string
	player oto.Player
	sync.Mutex
}

func NewNativePlayer(text string, player oto.Player) *Native {
	return &Native{
		text:   text,
		player: player,
	}
}

func (n *Native) Play() error {
	n.player.Play()

	for n.player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	return n.player.Close()
}

func (n *Native) Stop() {
	n.player.Close()
}

func (n *Native) GetText() string {
	return n.text
}
