# Clip Saver

## Motivation

TL;DR When I press `PrtScn`, I want the screenshot to be automatically stored to the specified directory.

When I'm watching movies on Windows and feeling like taking a screenshot, I prefer `Windows+PrtScn` over `PrtScn` because the former can automatically save the screenshot for me, while regarding the latter, I have to manually paste the screenshot somewhere else before taking the next screenshot. However, if I press `Windows+PrtScn` when I'm watching Netflix, the pressing of `Windows` makes Netflix show the playback control options, which will appear in the screenshot, and I don't want that. As a result, I have to resort to `PrtScn`, but I also want to preserve the convenience provided by `Windows+PrtScn` at the same time, so this project is born.

## How It Works

1. Start this program.
1. Start watching the movie.
1. Every time an image is copied to the system clipboard, a corresponding image file is created in the specified directory.
1. Finish watching the movie.
1. Terminate the program.
1. Enjoy your screenshots!

For more details:

```sh
$ ./clip-saver --help
Usage: clip-saver [--dir <dir>] [--conf <conf-path>]

Automatically save images from the system clipboard to the specified directory.

If there is no screenshot (i.e., filenames conforming to the screenshot format) at the moment, the screenshots will be saved in <dir> as 1.png, 2.png, and so on; if there are some existing screenshots with the latest one being 5.png, the screenshots will be saved in <dir> as 6.png, 7.png, and so on.

  -conf string
        optional; the path to the configuration file. For more details regarding the format of it and a working example, please check https://github.com/davidhsingyuchen/clip-saver/blob/main/clip-saver.yml. If this flag is not present, <dir>/clip-saver.yaml is used. If the file does not exist, movie mode instead of TV series mode is assumed.
  -dir string
        optional; the directory to save the clipped images. If it does not exist yet, this program will attempt to create it. If it is not specified, the current working directory will be used.
  -version
        print version information
```

## Installation

1. [Install Go](https://golang.org/doc/install).
1. `go install -v github.com/davidhsingyuchen/clip-saver@latest`

## Release Process

1. Pull the open `release-please` PR ([example](https://github.com/davidhsingyuchen/clip-saver/pull/5)).
1. Update the `version` constant in `main.go` according to the title of the PR and push the changes.
1. Merge the updated `release-please` PR.

## Other Notes

I personally use [Sandboxie](https://sandboxie-plus.com/downloads/) to take screenshots on Netflix. I'm not affiliated with it in any way.
