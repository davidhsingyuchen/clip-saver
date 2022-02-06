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
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: clip-saver --dir <dir> [--start-idx <idx>]\n\n"+
			"Assume that <start-idx> is set to be 6, "+
			"then the screenshots will be saved in <dir> as 6.png, 7.png, and so on.\n\n",
		)
		flag.PrintDefaults()
	}
	dir := flag.String("dir", "", "required; the directory to save the clipped images; "+
		"if it does not exist yet, this program will attempt to create it")
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
				return fmt.Errorf("failed to stat %q: %w", fileName, err)
			}
			return fmt.Errorf("file already exists: %q", fileName)
		}

		if err := writeImgToFile(img, fileName); err != nil {
			return fmt.Errorf("failed to write the clip to %q: %w", fileName, err)
		}
		log.Printf("Wrote to %q successfully!", fileName)
	}
	return nil
}

func writeImgToFile(img []byte, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create the file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(img); err != nil {
		return fmt.Errorf("failed to write to the file: %w", err)
	}
	return nil
}
