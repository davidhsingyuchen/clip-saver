package infra

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/davidhsingyuchen/clip-saver/service"
)

const (
	cmdPrevEpisode        = "<"
	cmdNextEpisode        = ">"
	episodeSeqNoSeparator = "-"
)

// SeriesFilenameGenerator generates filenames for screenshots of a TV series.
type SeriesFilenameGenerator struct {
	seqNo int
	ewl   *episodeWithLock
}

var _ service.FilenameGenerator = (*SeriesFilenameGenerator)(nil)

// NewSeriesFilenameGenerator expects dir to be the directory containing the screenshots.
func NewSeriesFilenameGenerator(dir string, episodesPerSeason int) (*SeriesFilenameGenerator, error) {
	if episodesPerSeason <= 0 {
		return nil, fmt.Errorf("episodesPerSeason must be nonnegative, but got %d", episodesPerSeason)
	}

	g := &SeriesFilenameGenerator{}
	seqNos, episodes, err := g.parseSeqNosAndEpisodes(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sequence numbers and episodes from dir: %w", err)
	}
	g.seqNo = nextSeqNo(seqNos)
	g.ewl = g.nextEpisode(episodes, episodesPerSeason)

	log.Printf("Next filename that will be used: %s", g.gen())
	go g.watchEpisodeChanges()
	return g, nil
}

func (g *SeriesFilenameGenerator) Gen() string {
	ret := g.gen()
	g.seqNo++
	return ret
}

func (g *SeriesFilenameGenerator) gen() string {
	return fmt.Sprintf("%s%s%d", g.ewl, episodeSeqNoSeparator, g.seqNo)
}

func (g *SeriesFilenameGenerator) parseSeqNosAndEpisodes(dir string) ([]int, []*episode, error) {
	fs, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var seqNos []int
	var episodes []*episode
	// We omit errors (e.g., -1) because there may be files other than the screenshots.
	for _, f := range fs {
		filename := trimExtension(f.Name())

		idx := strings.Index(filename, episodeSeqNoSeparator)
		if idx == -1 {
			continue
		}
		strEpisode, strSeqNo := filename[:idx], filename[idx+1:]

		episode, err := newEpisode(strEpisode)
		if err != nil {
			continue
		}
		seqNo, err := strconv.Atoi(strSeqNo)
		if err != nil {
			continue
		}

		seqNos = append(seqNos, seqNo)
		episodes = append(episodes, episode)
	}

	return seqNos, episodes, nil
}

func (g *SeriesFilenameGenerator) nextEpisode(episodes []*episode, episodesPerSeason int) *episodeWithLock {
	if len(episodes) == 0 {
		return newEpisodeWithLock(minEpisode(), episodesPerSeason)
	}

	e := newEpisodeWithLock(max(episodes, lessThanEpisode), episodesPerSeason)
	e.increment()
	return e
}

func (g *SeriesFilenameGenerator) watchEpisodeChanges() {
	const methodName = "watchEpisodeChanges"
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		switch scanner.Text() {
		case cmdPrevEpisode:
			g.ewl.decrement()
		case cmdNextEpisode:
			g.ewl.increment()
		default:
			continue
		}
		log.Printf("%s: current episode: %s", methodName, g.ewl)
	}

	if scanner.Err() != nil {
		log.Printf("%s: failed to scan from stdin: %v", methodName, scanner.Err())
	}
}
