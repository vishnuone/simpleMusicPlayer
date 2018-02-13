package main

import (
	"os"
	"github.com/faiface/beep"
	// "github.com/faiface/beep/wav"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"time"
	"log"
	"io/ioutil"
	"strings"
	"fmt"
)

func main() {

	if len(os.Args) == 1 {
		log.Fatalln("Missing argument 1. Expecting directory or file path. use -help for doc")
		return
	}

	path := os.Args[1]
	if path == ""{
		log.Fatalln("Missing argument 1. Expecting directory or file path. use -help for doc")
		return
	}

	if path == "-help"{
		showDoc()
		return
	}

	fi, err := os.Stat(path)

	if os.IsNotExist(err) {
		log.Fatalln("path : " + path + " file or directory not found. use -help for doc")
		return
	} else if err != nil {
		log.Fatalln(err)
		return
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		playFromDir(path)
	case mode.IsRegular():
		done := make(chan struct{})
		playFile(path, done)
		<-done
	}
}

func playFromDir(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, f := range files {
		// Change this to check against a list of supported file extensions.
		if strings.Contains(f.Name(), ".mp3") || strings.Contains(f.Name(), ".MP3") {
			done := make(chan struct{})

			if path[len(path)-1:] != "/" {
				path = path + "/"
			}

			playFile(path + f.Name(), done)
			<-done
		}
	}
}

func playFile(filePath string, done chan struct{}) {
	var s beep.StreamCloser
	var format beep.Format


	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
		return
	}

	// make code support .wav too
	//if(filepath == "wav") { // add check here for .wav
	//	s, format, err = wav.Decode(f)
	//} else {
		// s, format, err = mp3.Decode(f)
	//}

	s, format, err = mp3.Decode(f)
	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Println("Playing: "+ filePath)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		close(done)
	})))
}

func showDoc() {
	fmt.Println("Player only supports .mp3 format.")
	fmt.Println("")
	fmt.Println("Usage:-")
	fmt.Println("To play single file: simpleMusicPlayer example.mp3")
	fmt.Println("To play all .mp3 files in folder: simpleMusicPlayer pathToMusicDir")
	fmt.Println("")
}