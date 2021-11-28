# Clip Saver

## Motivation

When I'm watching movies on Windows and feeling like taking a screenshot, I prefer `Windows+PrtScn` over `PrtScn` because the former can automatically save the screenshot for me, while regarding the latter, I have to manually paste the screenshot somewhere else before taking next screenshot. However, if I press `Windows+PrtScn` when I'm watching Netflix, the pressing of `Windows` makes Netflix show the playback control options, which will appear in the screenshot, and I don't want that. As a result, I have to resort to `PrtScn`, but I also want to preserve the convenience provided by `Windows+PrtScn` at the same time, so this project is born.

## How It Works

1. Start this program and specify an directory to store clips.
1. Start watching the movie.
1. Every time an image is copied to the system clipboard, a corresponding image file is created in the specified directory.
1. Finish watching the movie.
1. Terminate the program.
1. Enjoy your screenshots!

## Prerequisites

[Install Go](https://golang.org/doc/install).

### Installation

```sh
go install -v github.com/davidhsingyuchen/clip-saver@latest
```

### Build From Source

```sh
git clone git@github.com:davidhsingyuchen/clip-saver.git
cd clip-saver
go build -v -o clip-saver
```
