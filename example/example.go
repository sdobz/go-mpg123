package main

import (
	"fmt"
	"github.com/sdobz/go-mpg123"
	"golang.org/x/mobile/exp/audio"
	"time"
)

func main() {
	mpg123.Initialize()
	mp3, err := mpg123.Open("example.mp3")
	if err != nil {
		panic(err)
	}

	rate, channels, encoding, format := mp3.Format()
	fmt.Printf("Rate: %i Channels: %i Encoding: %i Format: %s\n", rate, channels, encoding, format)

	p, err := audio.NewPlayer(fakers, audio.Format(format), rate)
	if err != nil {
		panic(err)
	}

	p.Play()
	for p.State() == audio.Playing {
		time.Sleep(time.Second)
	}
	mpg123.Exit()
}
