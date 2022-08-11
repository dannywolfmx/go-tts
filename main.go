package main

import (
	"fmt"
	"time"

	"github.com/dannywolfmx/go-tts/tts"
)

func main() {
	myTTS := tts.NewTTS()
	myTTS.Skip()

	go func() {
		err := myTTS.Play("es", "Hola mundo como stan? Hola mundo como stan? Hola mundo como stan?")

		if err != nil {
			fmt.Printf("Error %s \n", err)
		}

	}()

	go func() {
		time.Sleep(1 * time.Second)
		err := myTTS.Play("es", "Segunda prueba de control, segunda prueba")
		if err != nil {
			fmt.Printf("Error %s \n", err)
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("Skip")
		myTTS.Skip()
	}()

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("Skip")
		myTTS.Skip()
	}()

	go func() {
		time.Sleep(4 * time.Second)
		fmt.Println("Skip")
		myTTS.Skip()
		myTTS.Stop()
	}()

	myTTS.Run()
}
