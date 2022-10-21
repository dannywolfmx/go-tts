package tts

import (
	"sync"
	"testing"
	"time"
)

const lang = "es"

func TestAdd(t *testing.T) {

	voice := NewTTS(lang)
	voice.Add("test")

	l := voice.QueueLen()
	if l != 1 {
		t.Fatalf("Queue is not 1, actual len is %d", l)
	}
}

// TestPause is a blocking thread test
func TestPause(t *testing.T) {
	voice := NewTTS(lang)
	var wg sync.WaitGroup

	wg.Add(1)
	go func(t *testing.T) {
		time.Sleep(2 * time.Second)
		voice.Pause()
		time.Sleep(2 * time.Second)
		voice.Play()
		wg.Done()
	}(t)

	voice.Add("test prueba, test prueba, test prueba, test prueba, test prueba, test prueba")
	voice.Play()

	wg.Wait()
}

// TestPlay is a blocking thread test
func TestNext(t *testing.T) {
	voice := NewTTS(lang)
	var wg sync.WaitGroup

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

	wg.Add(1)
	go func(t *testing.T) {
		time.Sleep(5 * time.Second)
		voice.Next()
		wg.Done()
	}(t)

	voice.Add("test prueba 1, test prueba 1, test prueba 1")
	voice.Add("test prueba 2, test prueba 2, test prueba 2")
	voice.Add("test prueba 3, test prueba 3, test prueba 3")
	voice.Play()

	wg.Wait()
}

// TestPlay is a blocking thread test
func TestPlay(t *testing.T) {
	voice := NewTTS(lang)

	voice.Add("test prueba")
	voice.Add("test prueba")
	voice.Add("test prueba")
	l := voice.QueueLen()
	if l != 3 {
		t.Fatalf("Queue is not 3, actual len is %d", l)
	}

	voice.Play()

	l = voice.QueueLen()
	if l != 0 {
		t.Fatalf("Queue is not empty, actual len is %d", l)
	}
}

func TestPlayActiveAutoplay(t *testing.T) {
	voice := NewTTS(lang)

	voice.Play()

	if !voice.autoplay {
		t.Fatal("Autoplay is not set to true")
	}

}

func TestQueueLen(t *testing.T) {
	voice := NewTTS(lang)

	l := voice.QueueLen()
	if l != 0 {
		t.Fatalf("Queue is not 0, actual len is %d", l)
	}

	voice.Add("test")

	l = voice.QueueLen()
	if l != 1 {
		t.Fatalf("Queue is not 1, actual len is %d", l)
	}
}
