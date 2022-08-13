package main

import (
	"fmt"
	"time"

	"github.com/dannywolfmx/go-tts/tts"
)

func main() {
	tts := tts.NewTTS("es")
	tts.Skip()

	err := tts.Play("Hola")

	if err != nil {
		fmt.Printf("Error %s \n", err)
	}

	err = tts.Play("Hola mundo como stan? Hola mundo como stan? Hola mundo como stan?")

	if err != nil {
		fmt.Printf("Error %s \n", err)
	}

	err = tts.Play("Hola mundo como stan? Hola mundo como stan? Hola mundo como stan?")

	if err != nil {
		fmt.Printf("Error %s \n", err)
	}
	err = tts.Play("Hola mundo como stan? Hola mundo como stan? Hola mundo como stan?")

	if err != nil {
		fmt.Printf("Error %s \n", err)
	}
	err = tts.Play("aaaaaaaaaa")

	if err != nil {
		fmt.Printf("Error %s \n", err)
	}

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("Skip")
		tts.Skip()
	}()

	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("Skip")
		tts.Skip()
		tts.Stop()
	}()

	tts.Run()
}
