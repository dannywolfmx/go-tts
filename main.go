package main

import (
	"fmt"
	"time"

	"github.com/dannywolfmx/go-tts/tts"
)

func main() {
	myTTS := tts.NewTTS()
	player, err := myTTS.Play("es", "Hola mundo como stan? Hola mundo como stan? Hola mundo como stan?")

	if err != nil {
		fmt.Printf("Error %s \n", err)
	}
	player2, err := myTTS.Play("es", "Segunda prueba de control, segunda prueba")
	if err != nil {
		fmt.Printf("Error %s \n", err)
	}

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("Skip")
		player.Skip()
	}()

	go func() {
		time.Sleep(7 * time.Second)
		fmt.Println("Skip")
		player2.Skip()
		myTTS.Stop()
	}()

	myTTS.Run()
}
