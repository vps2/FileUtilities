package ffmpeg

import (
	"bytes"
	"fmt"
	"futilities/internal/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var ffmpegPathName string

//FFMPEG оболочка для запуска внешнего конвертера ffmpeg
type FFMPEG struct {
	srcDir            string
	dstDir            string
	inputFileOptions  []string
	outputFileOptions []string
	outputFileExt     string
}

//Convert запускает ffmpeg для конвертации файла
func (f *FFMPEG) Convert(file *fs.File) error {
	dstFileExt := f.outputFileExt
	if dstFileExt == "" {
		dstFileExt = filepath.Ext(file.Name())
	}

	inputFileName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

	dstFileName := filepath.Join(f.dstDir, inputFileName)
	if dstFileExt != "" {
		dstFileName = dstFileName + dstFileExt
	}

	var args []string
	if len(f.inputFileOptions) != 0 {
		args = append(args, f.inputFileOptions...)
	}
	args = append(args, "-i", file.AbsolutePath())
	if len(f.outputFileOptions) != 0 {
		args = append(args, f.outputFileOptions...)
	}
	args = append(args, dstFileName)

	err := run(ffmpegPathName, args)
	if err != nil {
		if match, _ := regexp.MatchString(`File '.*?' already exists.`, err.Error()); !match {
			dstFile := fs.File{PathName: dstFileName}
			dstFile.Delete()
		}
	}

	return err
}

func run(command string, args []string) error {
	var buf bytes.Buffer

	cmd := exec.Command(command, args...)
	cmd.Stderr = &buf
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("%s: %w", buf.String(), err)
	}
	return err
}

//New создает новый экземпляр конвертера.
func New(srcDir, dstDir, inputFileOptions, outputFileOptions, outputFileExt string) (*FFMPEG, error) {
	var err error
	if ffmpegPathName, err = findFFMpeg(); err != nil {
		return nil, fmt.Errorf("ffmpeg converter was not found: %w", err)
	}

	ffmpeg := FFMPEG{
		srcDir:        srcDir,
		dstDir:        dstDir,
		outputFileExt: outputFileExt,
	}

	if inputFileOptions != "" {
		ffmpeg.inputFileOptions = strings.Split(inputFileOptions, " ")
	}
	if outputFileOptions != "" {
		ffmpeg.outputFileOptions = strings.Split(outputFileOptions, " ")
	}

	return &ffmpeg, nil

}

func findFFMpeg() (pathName string, retErr error) {
	var platform = runtime.GOOS

	const ffmpeg = "ffmpeg"

	switch platform {
	case "windows":
		pathName, retErr = testCommand("where", ffmpeg)
	default:
		pathName, retErr = testCommand("which", ffmpeg)
	}

	//попытка поиска в каталоге исполняемого файла
	if retErr != nil {
		if executablePath, err := os.Executable(); err == nil {
			dirReader := fs.NewDirReaderWithFilter(filepath.Dir(executablePath), func(fileInfo os.FileInfo) bool {
				return fileInfo.Mode().IsRegular() &&
					ffmpeg == strings.TrimSuffix(fileInfo.Name(), filepath.Ext(fileInfo.Name()))
			})
			if files, err := dirReader.Read(); err == nil && len(files) == 1 {
				pathName, retErr = files[0].AbsolutePath(), nil
			}
		}
	}

	return
}

func testCommand(command string, args ...string) (string, error) {
	var out bytes.Buffer

	cmd := exec.Command(command, args...)
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return out.String(), err
	}

	return strings.ReplaceAll(out.String(), lineSeparator(), ""), nil
}

func lineSeparator() string {
	switch runtime.GOOS {
	case "windows":
		return "\r\n"
	default:
		return "\n"
	}
}
