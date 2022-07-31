package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

const (
	abbrSeason     = "S"
	abbrEpisode    = "E"
	cmdPrevEpisode = "<"
	cmdNextEpisode = ">"
)

// SequentialFilenameGenerator generates filenames by following a numerical sequence.
// It reads from stdin if it's in TV series mode.
type SequentialFilenameGenerator struct {
	// episodesPerSeason is set to be 0 if the screenshots come from a movie instead of a TV series.
	// Its value should be initialized by the constructor and never be changed thereafter.
	episodesPerSeason int
	curSeqIdx         int

	curSeason, curEpisode int
	// mu protects curSeason and curEpisode.
	mu sync.RWMutex
}

var _ FilenameGenerator = &SequentialFilenameGenerator{}

// NewSequentialFilenameGenerator expects episodesPerSeason to be set to 0
// if the screenshots come from a movie instead of a TV series.
func NewSequentialFilenameGenerator(episodesPerSeason, curSeqIdx int) (*SequentialFilenameGenerator, error) {
	g := &SequentialFilenameGenerator{
		episodesPerSeason: episodesPerSeason,
		curSeqIdx:         curSeqIdx,
	}
	if episodesPerSeason > 0 {
		// TODO: Automatically detects the last episode.
		g.curSeason = 1
		g.curEpisode = 1

		go g.adjustCurEpisode()
	}
	return g, nil
}

func (g *SequentialFilenameGenerator) Gen() string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return fmt.Sprintf("%s%d%s%d-%d", abbrSeason, g.curSeason, abbrEpisode, g.curEpisode, g.curSeqIdx)
}

func (g *SequentialFilenameGenerator) adjustCurEpisode() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		g.mu.Lock()
		// Invariant: 1 <= curEpisode <= episodesPerSeason
		// The way curEpisode is initialized ensures the invariant.
		switch scanner.Text() {
		case cmdPrevEpisode:
			if g.curEpisode > 1 {
				g.curEpisode--
			} else {
				g.curSeason--
				g.curEpisode = g.episodesPerSeason
			}
		case cmdNextEpisode:
			if g.curEpisode < g.episodesPerSeason {
				g.curEpisode++
			} else {
				g.curSeason++
				g.curEpisode = 1
			}
		default:
			continue
		}
		log.Printf("adjustCurEpisode: current episode: %s%d%s%d", abbrSeason, g.curSeason, abbrEpisode, g.curEpisode)
		g.mu.Unlock()
	}

	if scanner.Err() != nil {
		log.Printf("adjustCurEpisode: failed to scan from stdin: %v", scanner.Err())
	}
}
