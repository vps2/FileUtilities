# FFMPEGConv

**ffmpegconv** - a wrapper for the ['ffmpeg'](https://www.ffmpeg.org/) file conversion utility that monitors the appearance of files in the specified folder and launches it to convert the appeared file with the specified parameters. If the conversion was successful the original file is deleted. The 'ffpmpeg' utility should be located either in the same folder with the program or in the 'path' environment variable.

```sh
ffmpegconv.exe -h
Usage of ffmpegconv.exe:
  -s, --src-dir string      the folder where new files are tracked
  -d, --dst-dir string      the folder where converted files from the source folder will be placed
  -t, --timeout duration    the timeout between polls of the source directory (default 1m0s)
  -i, --ifile-opts string   input file options for ffmpeg
  -o, --ofile-opts string   output file options for ffmpeg
  -e, --ofile-ext string    output file extension
```

### Usage example:

```sh
ffmpegconv.exe -s c:\temp\1 -d c:\temp\2 -t 15s -o "-n -vf scale=640:480 -c:v libx264 -crf 24"
```