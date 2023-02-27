package tts

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

const lang = "es"

const sampleRate = 27000

// TestPause is a blocking thread test
func TestPause(t *testing.T) {
	voice := NewTTS(lang, sampleRate)
	voice.Play()
	var wg sync.WaitGroup

	wg.Add(1)
	go func(t *testing.T) {
		time.Sleep(2 * time.Second)
		voice.Pause()
		time.Sleep(2 * time.Second)
		voice.Continue()
		time.Sleep(4 * time.Second)
		wg.Done()
	}(t)

	voice.Add("Soy una frase de prueba, para demostrar el nivel de pitch en el audio y su velocidad")

	wg.Wait()
}

// TestPlay is a blocking thread test
func TestNext(t *testing.T) {
	voice := NewTTS(lang, sampleRate)
	var wg sync.WaitGroup
	voice.Play()

	wg.Add(1)
	go func(t *testing.T) {
		time.Sleep(1 * time.Second)
		voice.Next()
		wg.Done()
	}(t)

	wg.Add(1)
	go func(t *testing.T) {
		time.Sleep(2 * time.Second)
		voice.Next()
		wg.Done()
	}(t)

	wg.Add(1)
	go func(t *testing.T) {
		time.Sleep(3 * time.Second)
		voice.Next()
		wg.Done()
	}(t)

	wg.Add(1)
	go func(t *testing.T) {
		time.Sleep(4 * time.Second)
		voice.Next()
		wg.Done()
	}(t)

	voice.Add("test prueba 1, test prueba 1, test prueba 1")
	voice.Add("test prueba 2, test prueba 2, test prueba 2")
	voice.Add("test prueba 3, test prueba 3, test prueba 3")
	voice.Add("test prueba 4, test prueba 4, test prueba 4")
	voice.Add("test prueba 5, test prueba 5, test prueba 5")
	fmt.Println("Terminado")

	wg.Wait()
}

// TestPlay is a blocking thread test
func TestNextRaceCondition(t *testing.T) {
	voice := NewTTS(lang, sampleRate)
	var wg sync.WaitGroup
	voice.Play()

	wg.Add(1)
	go func(t *testing.T) {
		time.Sleep(2 * time.Second)
		voice.Next()
		wg.Done()
	}(t)

	wg.Add(1)
	go func(t *testing.T) {
		time.Sleep(2 * time.Second)
		voice.Next()
		wg.Done()
	}(t)

	wg.Add(1)
	go func(t *testing.T) {
		time.Sleep(2 * time.Second)
		voice.Next()
		wg.Done()
	}(t)

	wg.Add(1)
	go func(t *testing.T) {
		time.Sleep(2 * time.Second)
		voice.Next()
		wg.Done()
	}(t)

	voice.Add("test prueba 1, test prueba 1, test prueba 1")
	voice.Add("test prueba 2, test prueba 2, test prueba 2")
	voice.Add("test prueba 3, test prueba 3, test prueba 3")
	voice.Add("test prueba 4, test prueba 4, test prueba 4")
	voice.Add("test prueba 5, test prueba 5, test prueba 5")

	wg.Wait()
}

// TestPlay is a blocking thread test
func TestPlay(t *testing.T) {
	voice := NewTTS(lang, sampleRate)
	voice.Play()

	voice.Add("test prueba")
	voice.Add("test prueba")
	voice.Add("test prueba")
}

func TestPlayActiveAutoplay(t *testing.T) {
	voice := NewTTS(lang, sampleRate)
	voice.Play()

	if !voice.autoplay {
		t.Fatal("Autoplay is not set to true")
	}

}
