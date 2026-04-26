.PHONY: build
build:
	go build -v -o clip-saver

.PHONY: gen
gen:
	go run github.com/golang/mock/mockgen -source=service/clip_saver.go -destination=service/mock_filename_generator.go -package=service FilenameGenerator
