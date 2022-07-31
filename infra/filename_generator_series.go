package infra

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

// SeriesFilenameGenerator generates filenames for screenshots of a TV series.
type SeriesFilenameGenerator struct {
	seqNo int

	// episodesPerSeason should be initialized by the constructor and never be changed thereafter.
	episodesPerSeason int

	// mu protects curSeason and curEpisode.
	mu                    sync.RWMutex
	curSeason, curEpisode int
}

// NewSeriesFilenameGenerator initializes and returns a new SeriesFilenameGenerator.
func NewSeriesFilenameGenerator(episodesPerSeason, seqNo int) (*SeriesFilenameGenerator, error) {
	if episodesPerSeason <= 0 {
		return nil, fmt.Errorf("episodesPerSeason must be nonnegative, but got %d", episodesPerSeason)
	}
	g := &SeriesFilenameGenerator{
		episodesPerSeason: episodesPerSeason,
		// TODO: Automatically detects the last idx and the last episode.
		seqNo:      seqNo,
		curSeason:  1,
		curEpisode: 1,
	}
	go g.adjustCurEpisode()
	return g, nil
}

func (g *SeriesFilenameGenerator) Gen() string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return fmt.Sprintf("%s%d%s%d-%d", abbrSeason, g.curSeason, abbrEpisode, g.curEpisode, g.seqNo)
}

func (g *SeriesFilenameGenerator) adjustCurEpisode() {
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
