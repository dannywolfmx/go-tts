package tts

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type TTS struct {
}

func NewTTS() *TTS {
	return &TTS{}
}

//Google wants a cache system to don't ban this client, so we need to add it
func (t *TTS) GetAudio(lang, text string) (*os.File, error) {
	hashText := hash(text)

	file, err := os.Open(hashText)
	if err == nil {
		fmt.Println("Usando cache")
		return file, nil
	}

	url := fmt.Sprintf("https://translate.google.com/translate_tts?ie=UTF-8&tl=%s&client=tw-ob&q=%s", lang, url.QueryEscape(text))
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buff, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	file, err = os.Create(hashText)
	defer file.Close()

	if err != nil {
		fmt.Println("Entre")
		return nil, err
	}

	n, err := file.Write(buff)

	if err != nil {
		return nil, err
	}

	if len(buff) != n {
		return nil, errors.New("error on len buff write")
	}

	return file, nil
}
