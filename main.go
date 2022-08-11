package main

import (
	"fmt"

	"github.com/dannywolfmx/go-tts/tts"
)

func main() {
	myTTS := tts.NewTTS()
	_, err := myTTS.GetAudio("es", "Hola mundo como stan?")

	if err != nil {
		fmt.Printf("Error %s \n", err)
	}
}
