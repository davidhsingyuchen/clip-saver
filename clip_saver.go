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

			if _, err := os.Stat(fileName); err == nil {
				return fmt.Errorf("file already exists: %q", fileName)
			} else if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("failed to stat %q: %w", fileName, err)
			}

			if err := os.WriteFile(fileName, img, 0644); err != nil {
				return fmt.Errorf("failed to write the clip to %q: %w", fileName, err)
			}
			log.Printf("Wrote to %q successfully!", fileName)
		}
	}
}
