package infra

import "fmt"

// MovieFilenameGenerator generates filenames for screenshots of a movie.
type MovieFilenameGenerator struct {
	seqNo int
}

// NewMovieFilenameGenerator initializes and returns a new MovieFilenameGenerator.
func NewMovieFilenameGenerator(seqNo int) *MovieFilenameGenerator {
	// TODO: Automatically detects the last index used.
	return &MovieFilenameGenerator{
		seqNo: seqNo,
	}
}

func (g *MovieFilenameGenerator) Gen() string {
	return fmt.Sprint(g.seqNo)
}
