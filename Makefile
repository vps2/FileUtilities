.PHONY: all build-fmove build-ffmpegconv clean help

## all: создать исполняемые файлы всех утилит
all : build-fmove build-ffmpegconv

## build-fmove: создать исполняемый файл утилиты fmove
build-fmove:
	go build -o bin/fmove.exe -ldflags "-s -w" cmd/fmove/main.go

## build-ffmpegconv: создать исполняемый файл утилиты ffmpegconv
build-ffmpegconv:
	go build -o bin/ffmpegconv.exe -ldflags "-s -w" cmd/ffmpegconv/main.go

## clean: удалить содержимое папки bin
clean:
	rm -f bin/*.*

help: Makefile
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'