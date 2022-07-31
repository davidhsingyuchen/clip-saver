package main

import "fmt"

const (
	abbrSeason  = "S"
	abbrEpisode = "E"
)

// SequentialFilenameGenerator generates filenames by following a numerical sequence.
type SequentialFilenameGenerator struct {
	// episodesPerSeason is set to be 0 if the screenshots come from a movie instead of a TV series.
	// Its value should be initialized by the constructor and never be changed thereafter.
	episodesPerSeason     int
	curSeason, curEpisode int
	curSeqIdx             int
}

var _ FilenameGenerator = &SequentialFilenameGenerator{}

// NewSequentialFilenameGenerator expects episodesPerSeason to be set to 0
// if the screenshots come from a movie instead of a TV series.
func NewSequentialFilenameGenerator(episodesPerSeason, curSeqIdx int) *SequentialFilenameGenerator {
	g := &SequentialFilenameGenerator{
		episodesPerSeason: episodesPerSeason,
		curSeqIdx:         curSeqIdx,
	}
	if episodesPerSeason > 0 {
		// TODO: Automatically detects the last episode.
		g.curSeason = 1
		g.curEpisode = 1
	}
	return g
}

func (g *SequentialFilenameGenerator) Gen() string {
	return fmt.Sprintf("%s%d%s%d-%d", abbrSeason, g.curSeason, abbrEpisode, g.curEpisode, g.curSeqIdx)
}

func (g *SequentialFilenameGenerator) Cleanup() {

}
