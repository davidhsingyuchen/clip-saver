package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/davidhsingyuchen/clip-saver/infra"
	"github.com/davidhsingyuchen/clip-saver/service"
)

const (
	version = "v1.1.1"

	dirPerm         = 0700
	defaultConfName = "clip-saver.yaml"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: clip-saver [--dir <dir>] [--conf <conf-path>] \n\n"+
			"Automatically save images from the system clipboard to the specified directory.\n\n"+
			"If there is no screenshot (i.e., filenames conforming to the screenshot format) at the moment, "+
			"the screenshots will be saved in <dir> as 1.png, 2.png, and so on; "+
			"if there are some existing screenshots with the latest one being 5.png, "+
			"the screenshots will be saved in <dir> as 6.png, 7.png, and so on.\n\n",
		)
		flag.PrintDefaults()
	}
	dir := flag.String("dir", "", "optional; the directory to save the clipped images. "+
		"If it does not exist yet, this program will attempt to create it. "+
		"If it is not specified, the current working directory will be used.")
	confPath := flag.String("conf", "", "optional; the path to the configuration file. "+
		"For more details regarding the format of it and a working example, please check "+
		"https://github.com/davidhsingyuchen/clip-saver/blob/main/clip-saver.yml. "+
		fmt.Sprintf("If this flag is not present, <dir>/%s is used. ", defaultConfName)+
		"If the file does not exist, movie mode instead of TV series mode is assumed.")
	printVer := flag.Bool("version", false, "print version information")
	flag.Parse()

	if *printVer {
		fmt.Println(version)
		os.Exit(0)
	}

	if *dir == "" {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("failed to get the current working directory: %v", err)
		}
		*dir = wd
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

	var filenameGenerator service.FilenameGenerator
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
	if err := service.SaveClips(context.Background(), *dir, filenameGenerator); err != nil {
		log.Fatalf("failed to save clips: %v", err)
	}
}
