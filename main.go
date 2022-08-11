package main

import (
	"fmt"
	"time"

	"github.com/dannywolfmx/go-tts/tts"
)

func main() {
	tts := tts.NewTTS()
	tts.Skip()

	err := tts.Play("es", "Hola")

	if err != nil {
		fmt.Printf("Error %s \n", err)
	}

	err = tts.Play("es", "Hola mundo como stan? Hola mundo como stan? Hola mundo como stan?")

	if err != nil {
		fmt.Printf("Error %s \n", err)
	}

	go func() {
		time.Sleep(10 * time.Second)
		fmt.Println("Skip")
		tts.Skip()
		tts.Stop()
	}()

	tts.Run()
}
