package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/davidhsingyuchen/clip-saver/infra"
)

const (
	version = "v1.0.0"

	dirPerm         = 0700
	defaultConfName = "clip-saver.yaml"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: clip-saver --dir <dir>\n\n"+
			"If there is no screenshot (i.e., filenames conforming to the screenshot format) at the moment, "+
			"the screenshots will be saved in <dir> as 1.png, 2.png, and so on; "+
			"if there are some existing screenshots with the latest one being 5.png, "+
			"the screenshots will be saved in <dir> as 6.png, 7.png, and so on.\n\n"+
			`In "series" mode, the screenshots will be prefixed with the current episode, `+
			"and one can press < and > to move to the previous and the next episode respectively.\n\n",
		)
		flag.PrintDefaults()
	}
	dir := flag.String("dir", "", "required; the directory to save the clipped images; "+
		"if it does not exist yet, this program will attempt to create it")
	confPath := flag.String("conf", "", "optional; the path to the configuration file; "+
		fmt.Sprintf("if it is not provided, {$dir}/%s is used; ", defaultConfName)+
		"if the file does not exist, movie mode instead of TV series mode is assumed")
	printVer := flag.Bool("version", false, "print version information")
	flag.Parse()

	if *printVer {
		fmt.Println(version)
		os.Exit(0)
	}

	if *dir == "" {
		log.Fatalln("--dir is required")
	}
	if err := os.MkdirAll(*dir, 0700); err != nil {
		log.Fatalf("failed to ensure that %q exists as a directory: %v", *dir, err)
	}

	if *confPath == "" {
		*confPath = fmt.Sprintf("%s/%s", *dir, defaultConfName)
	}
	conf, err := NewConfig(*confPath)
	if err != nil {
		log.Fatalf("failed to create config: %v", err)
	}

	var filenameGenerator FilenameGenerator
	switch conf.Mode {
	case VideoModeMovie:
		filenameGenerator, err = infra.NewMovieFilenameGenerator(*dir)
		if err != nil {
			log.Fatalf("failed to create a filename generator for a movie: %v", err)
		}
	case VideoModeSeries:
		filenameGenerator, err = infra.NewSeriesFilenameGenerator(*dir, conf.EpisodesPerSeason)
		if err != nil {
			log.Fatalf("failed to create a filename generator for a series: %v", err)
		}
	default:
		log.Fatalf("unsupported mode: %v (supported modes: [%s, %s])", conf.Mode, VideoModeMovie, VideoModeSeries)
	}

	if err != nil {
		log.Fatalf("failed to create a new filename generator: %v", err)
	}
	if err := saveClips(context.Background(), *dir, filenameGenerator); err != nil {
		log.Fatalf("failed to save clips: %v", err)
	}
}
