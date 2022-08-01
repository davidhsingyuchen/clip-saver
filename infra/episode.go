package infra

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

const (
	minEpisodeNo = 1
	minSeasonNo  = 1

	abbrSeason  = "S"
	abbrEpisode = "E"
)

type episode struct {
	episodeNo, seasonNo int
}

func newEpisode(s string) (*episode, error) {
	if !strings.HasPrefix(s, abbrSeason) {
		return nil, fmt.Errorf("%q does not begin with the season abbreviation %q", s, abbrSeason)
	}
	idx := strings.Index(s, abbrEpisode)
	if idx == -1 {
		return nil, fmt.Errorf("%q does not contain the episode abbreviation %q", s, abbrEpisode)
	}

	seasonStr, episodeStr := s[1:idx], s[idx+1:]
	seasonNo, err := strconv.Atoi(seasonStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the season: %w", err)
	}
	episodeNo, err := strconv.Atoi(episodeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the episode: %w", err)
	}

	return &episode{
		episodeNo: episodeNo,
		seasonNo:  seasonNo,
	}, nil
}

func minEpisode() *episode {
	return &episode{
		episodeNo: minEpisodeNo,
		seasonNo:  minSeasonNo,
	}
}

func lessThanEpisode(e1, e2 *episode) bool {
	if e1.seasonNo == e2.seasonNo {
		return e1.episodeNo < e2.episodeNo
	}
	return e1.seasonNo < e2.seasonNo
}

type episodeWithLock struct {
	// episodesPerSeason should be initialized by the constructor and never be changed thereafter.
	episodesPerSeason int

	e  *episode
	mu *sync.RWMutex
}

func newEpisodeWithLock(e *episode, episodesPerSeason int) *episodeWithLock {
	return &episodeWithLock{
		episodesPerSeason: episodesPerSeason,
		e:                 e,
		mu:                &sync.RWMutex{},
	}
}

func (ewl *episodeWithLock) increment() {
	ewl.mu.Lock()
	defer ewl.mu.Unlock()

	if ewl.e.episodeNo < ewl.episodesPerSeason {
		ewl.e.episodeNo++
	} else {
		ewl.e.seasonNo++
		ewl.e.episodeNo = minEpisodeNo
	}
}

func (ewl *episodeWithLock) decrement() {
	ewl.mu.Lock()
	defer ewl.mu.Unlock()

	if ewl.e.episodeNo > minEpisodeNo {
		ewl.e.episodeNo--
	} else {
		ewl.e.seasonNo--
		ewl.e.episodeNo = ewl.episodesPerSeason
	}
}

func (ewl *episodeWithLock) String() string {
	ewl.mu.RLock()
	defer ewl.mu.RUnlock()
	return fmt.Sprintf("%s%d%s%d", abbrSeason, ewl.e.seasonNo, abbrEpisode, ewl.e.episodeNo)
}
