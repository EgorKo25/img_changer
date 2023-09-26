package main

import (
	"flag"
	"log"
	"os"

	"ImgChanger/pkg/resizer"
	"ImgChanger/pkg/rotater"
)

var (
	in string
)

func main() {

	flag.StringVar(&in, "in", "dataset/J_test/", "source dir")
	flag.Parse()

	out := in + "_changed"

	err := os.Mkdir(out, 0777)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	r, err := resizer.NewResizer()
	if err != nil {
		log.Fatalf("error constructor: %s", err)
	}

	rot, err := rotater.NewRotater()
	if err != nil {
		log.Fatalf("error constructor: %s", err)
	}
	err = rot.RotateAll(in, out)
	if err != nil {
		log.Fatalf("error rotate: %s", err)
	}

	err = r.ResizeAll(out, out)
	if err != nil {
		log.Fatalf("error resize: %s", err)
	}

}
