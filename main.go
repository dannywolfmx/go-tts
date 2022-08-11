package main

import (
	"fmt"
	"time"

	"github.com/dannywolfmx/go-tts/tts"
)

func main() {
	tts := tts.NewTTS()
	tts.Skip()

	go func() {
		err := tts.Play("es", "Hola mundo como stan? Hola mundo como stan? Hola mundo como stan?")

		if err != nil {
			fmt.Printf("Error %s \n", err)
		}

	}()

	go func() {
		time.Sleep(1 * time.Second)
		err := tts.Play("es", "Segunda prueba de control, segunda prueba")
		if err != nil {
			fmt.Printf("Error %s \n", err)
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("Skip")
		tts.Skip()
	}()

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("Skip")
		tts.Skip()
	}()

	go func() {
		time.Sleep(4 * time.Second)
		fmt.Println("Skip")
		tts.Skip()
		tts.Stop()
	}()

	tts.Run()
}
