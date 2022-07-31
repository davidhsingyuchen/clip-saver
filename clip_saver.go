package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"golang.design/x/clipboard"
)

const (
	fileExt = "png"
)

// FilenameGenerator generates the filenames to be used when writing screenshots to the disk.
type FilenameGenerator interface {
	// Gen generates a filename to be used.
	// Note that calling this method may change the internal state of the corresponding FilenameGenerator.
	Gen() string
}

func saveClips(ctx context.Context, dir string, filenameGenerator FilenameGenerator) error {
	ch := clipboard.Watch(ctx, clipboard.FmtImage)
	log.Println("Start to watch for clips...")
	for {
		select {
		case <-ctx.Done():
			return nil
		case img := <-ch:
			fileName := fmt.Sprintf("%s.%s", filenameGenerator.Gen(), fileExt)

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
	}
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