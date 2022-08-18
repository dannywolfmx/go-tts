package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/dannywolfmx/go-tts/tts"
)

func main() {

	var wg sync.WaitGroup
	tts := tts.NewTTS("es")

	tts.Add("Hola")

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 4)
		fmt.Println("Next")
		tts.Next()
	}()

	tts.Add("Hola mundo como estan? Hola mundo como estan? Hola mundo como estan? Hola mundo como estan? Hola mundo como estan?")

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 5)
		fmt.Println("Next")
		tts.Next()
	}()

	tts.Add("AAAA AAA AAA AAAA AAA")
	tts.Add("Hola mundo como estan? Hola mundo como estan? Hola mundo como estan? Hola mundo como estan? Hola mundo como estan?")

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 6)
		fmt.Println("Next")
		tts.Next()
	}()

	wg.Wait()

	tts.CleanCache()
}
