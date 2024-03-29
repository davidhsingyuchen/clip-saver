package infra

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/davidhsingyuchen/clip-saver/service"
)

// MovieFilenameGenerator generates filenames for screenshots of a movie.
type MovieFilenameGenerator struct {
	seqNo int
}

var _ service.FilenameGenerator = (*MovieFilenameGenerator)(nil)

// NewMovieFilenameGenerator expects dir to be the directory containing the screenshots.
func NewMovieFilenameGenerator(dir string) (*MovieFilenameGenerator, error) {
	g := &MovieFilenameGenerator{}
	seqNos, err := g.parseSeqNos(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to init the sequence number: %w", err)
	}
	g.seqNo = nextSeqNo(seqNos)
	log.Printf("Next filename that will be used: %s", g.gen())
	return g, nil
}

func (g *MovieFilenameGenerator) Gen() string {
	ret := g.gen()
	g.seqNo++
	return ret
}

func (g *MovieFilenameGenerator) gen() string {
	return fmt.Sprint(g.seqNo)
}

func (g *MovieFilenameGenerator) parseSeqNos(dir string) ([]int, error) {
	fs, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var seqNos []int
	for _, f := range fs {
		filename := trimExtension(f.Name())
		// We omit errors because there may be files other than the screenshots.
		if n, err := strconv.Atoi(filename); err == nil {
			seqNos = append(seqNos, n)
		}
	}

	return seqNos, nil
}
