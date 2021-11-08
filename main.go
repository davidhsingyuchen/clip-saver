package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"golang.design/x/clipboard"
)

const (
	dirPerm = 0700
	imgPerm = 0600
	fileExt = ".png"
)

func main() {
	dir := flag.String("dir", "", "required; the directory to save the clipped images to"+
		"if it does not exist, this program will attempt to create it")
	startIdx := flag.Int("start-idx", 0, "optional; the starting index of the image file names")
	flag.Parse()

	if *dir == "" {
		log.Fatalln("--dir is required")
	}
	if err := os.MkdirAll(*dir, 0700); err != nil {
		log.Fatalf("failed to ensure that %q exists as a directory: %v", *dir, err)
	}

	if err := saveClips(context.Background(), *dir, *startIdx); err != nil {
		log.Fatalf("failed to save clips: %v", err)
	}
}

func saveClips(ctx context.Context, dir string, startIdx int) error {
	i := startIdx
	ch := clipboard.Watch(ctx, clipboard.FmtImage)
	log.Println("Start to watch for clips...")
	for img := range ch {
		fileName := filepath.Join(dir, strconv.Itoa(i)) + fileExt
		i++

		_, err := os.Stat(fileName)
		if !errors.Is(err, os.ErrNotExist) {
			if err != nil {
				return fmt.Errorf("failed to stat %q: %v", fileName, err)
			}
			return fmt.Errorf("file already exists: %q", fileName)
		}

		file, err := os.Create(fileName)
		if _, err := file.Write(img); err != nil {
			return fmt.Errorf("failed to write the clip to %q: %v", fileName, err)
		}
		log.Printf("Wrote to %q successfully!", fileName)
	}
	return nil
}
