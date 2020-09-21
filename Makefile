.PHONY: build-all build-fmove build-ffmpegconv clean help

GOOS = $(shell go env GOOS)

## build-all: создать исполняемые файлы всех утилит
build-all : build-fmove build-ffmpegconv

## build-fmove: создать исполняемый файл утилиты fmove
build-fmove:
ifeq ($(GOOS),windows)
	go build -o bin/fmove.exe -ldflags "-s -w" cmd/fmove/main.go
else
	go build -o bin/fmove -ldflags "-s -w" cmd/fmove/main.go
endif

## build-ffmpegconv: создать исполняемый файл утилиты ffmpegconv
build-ffmpegconv:
ifeq ($(GOOS),windows)
	go build -o bin/ffmpegconv.exe -ldflags "-s -w" cmd/ffmpegconv/main.go
else
	go build -o bin/ffmpegconv -ldflags "-s -w" cmd/ffmpegconv/main.go
endif

## clean: удалить содержимое папки bin
clean:
	rm -f bin/*

help: Makefile
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
