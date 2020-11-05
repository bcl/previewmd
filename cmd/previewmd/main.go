package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/russross/blackfriday/v2"
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

func main() {
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	filename := flag.Arg(0)
	_, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Previewing %s on port %d\n", filename, cfg.Port)

	preview := func(w http.ResponseWriter, _ *http.Request) {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Printf("Error opening %s: %s\n", filename, err)
			return
		}

		// XXX Is this right inside a loop?
		defer f.Close()
		markdown, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Printf("Error reading %s: %s\n", filename, err)
			return
		}

		// Render html from markdown
		output := blackfriday.Run(markdown)
		w.Write(output)
	}

	http.HandleFunc("/", preview)

	listen := fmt.Sprintf(":%d", cfg.Port)
	log.Fatal(http.ListenAndServe(listen, nil))
}
