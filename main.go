package main

import (
	"os"
	"github.com/faiface/beep"
	"github.com/faiface/beep/wav"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"time"
	"log"
)

func main() {

	if len(os.Args) == 1 {
		log.Fatalln("Missing argument 1. Expecting directory or file path")
		return
	}


	path := os.Args[1]
	if path == ""{
		log.Fatalln("Missing argument 1. Expecting directory or file path.")
		return
	}

	fi, err := os.Stat(path)

	if os.IsNotExist(err) {
		log.Fatalln("path : " + path + " file or directory not found")
		return
	} else if err != nil {
		log.Fatalln(err)
		return
	}


	switch mode := fi.Mode(); {
	case mode.IsDir():

		log.Fatalln("directory")
	case mode.IsRegular():
		done := make(chan struct{})
		playFile(path, done)
		<-done
	}
}

func playFromDir(path string) {

}

func playFile(filepath string, done chan struct{}) {
	var s beep.StreamCloser
	var format beep.Format


	f, err := os.Open(filepath)
	if err != nil {
		log.Fatalln(err)
		return
	}

	if(filepath == "wav") {
		s, format, err = wav.Decode(f)
	} else {
		s, format, err = mp3.Decode(f)
	}

	if err != nil {
		log.Fatalln(err)
		return
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	speaker.Play(beep.Seq(s, beep.Callback(func() {
		close(done)
	})))
}