package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wabarc/imgbb"
)

var (
	key string
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n\n")
		fmt.Fprintf(os.Stderr, "  imgbb [options] [file1] ... [fileN]\n\n")

		flag.PrintDefaults()
	}
	var basePrint = func() {
		fmt.Print("A CLI tool help upload images to ImgBB.\n\n")
		flag.Usage()
		fmt.Fprint(os.Stderr, "\n")
	}

	flag.StringVar(&key, "k", "", "ImgBB api key, optional.")

	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		basePrint()
		os.Exit(0)
	}

}

func main() {
	files := flag.Args()

	i := imgbb.NewImgBB(nil, key)
	for _, path := range files {
		if _, err := os.Stat(path); err != nil {
			fmt.Fprintf(os.Stderr, "imgbb: %s: no such file or directory\n", path)
			continue
		}

		if url, err := i.Upload(path); err != nil {
			fmt.Fprintf(os.Stderr, "imgbb: %v\n", err)
		} else {
			fmt.Fprintf(os.Stdout, "%s  %s\n", url, path)
		}
	}
}
