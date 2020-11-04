package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type cmdlineArgs struct {
	Debug bool // Extra debug output
	Port  int
}

var cfg = cmdlineArgs{
	Debug: false,
	Port:  8080,
}

func init() {
	flag.BoolVar(&cfg.Debug, "debug", cfg.Debug, "Extra debug output")
	flag.IntVar(&cfg.Port, "port", cfg.Port, "Port to serve preview from")

	flag.Parse()
}

func PreviewMarkdown(update <-chan string, port int) {

	// Read file the first time
	// Render html from markdown

	// Setup http server
	for {
		filename := <-update
		fmt.Println(filename)

		// re-read the file

		// Render html from markdown
	}
}

func main() {
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	filename := flag.Arg(0)
	lastInfo, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Previewing %s on port %d\n", filename, cfg.Port)

	updatePreview := make(chan string, 1)
	go PreviewMarkdown(updatePreview, cfg.Port)

	// Loop, checking the file mtime every second
	// If it has been modified, tell preview to update
	for {
		info, err := os.Stat(filename)
		if err != nil {
			fmt.Println(err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if info.Size() != lastInfo.Size() || info.ModTime() != lastInfo.ModTime() {
			lastInfo = info

			// Send message to preview
			updatePreview <- filename
		}

		time.Sleep(1 * time.Second)
	}
}
