package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"github.com/vps2/futilities/internal/fs"

	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	srcDir, dstDir *string
	pollInterval   *time.Duration
)

func main() {
	log := createLogger().Sugar()
	defer log.Sync()
	log.Info("The application is starting...")
	defer log.Info("The application is stopped.")

	srcDir = flag.StringP("src-dir", "s", "", "the folder where new files are tracked")
	dstDir = flag.StringP("dst-dir", "d", "", "the folder where new files will be moved from the source folder")
	pollInterval = flag.DurationP("timeout", "t", 60*time.Second, "the timeout between polls of the source directory")
	help := flag.BoolP("help", "h", false, "show help")

	flag.CommandLine.MarkHidden("help")
	flag.CommandLine.SortFlags = false
	flag.Parse()

	if len(os.Args[1:]) == 0 || *help == true {
		flag.Usage()
		os.Exit(0)
	}

	if err := checkDirFlag("src-dir"); err != nil {
		log.Fatal(err)
	}
	if err := checkDirFlag("dst-dir"); err != nil {
		log.Fatal(err)
	}
	if *srcDir == *dstDir {
		log.Fatal("source and destination folders are the same")
	}

	dirReader := fs.NewDirReaderWithFilter(*srcDir, func(fileInfo os.FileInfo) bool { return fileInfo.Mode().IsRegular() })
	watcher := fs.NewDirWatcher(dirReader, *pollInterval)

	var wg sync.WaitGroup
	wg.Add(2)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer wg.Done()

		watcher.Watch(ctx)
	}()

	events := watcher.Events()
	errors := watcher.Errors()
	go func() {
		defer wg.Done()

	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case files := <-events:
				for _, file := range files {
					select {
					case <-ctx.Done():
						break loop
					default:
					}

					log.Infof("trying to move a file '%s' to folder '%s'", file.AbsolutePath(), *dstDir)
					if err := file.MoveTo(*dstDir); err != nil {
						log.Error(err)
					} else {
						log.Infof("the file '%s' was moved", file.AbsolutePath())
					}
				}
			case err := <-errors:
				if err != nil {
					log.Error(err)
				}
				break loop
			}
		}
	}()

	stopChan := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(stopChan, os.Interrupt)

	log.Info("The application is started.")

	// Ждём сигнала завершения от операционной системы или ошибки от watcher-ра
	select {
	case <-stopChan:
	case <-errors:
	}

	cancel()
	wg.Wait()
}

func checkDirFlag(name string) (err error) {
	isFlagFound := false
	flagValue := ""

	flag.Visit(func(f *flag.Flag) {
		if name == f.Name {
			isFlagFound = true
			flagValue = f.Value.String()
		}
	})

	if !isFlagFound || flagValue == "" {
		err = fmt.Errorf("'%s' flag is not set", name)

	} else {
		if err = checkDir(flagValue); err != nil {
			err = fmt.Errorf("the directory for the flag '%s' does not exist", name)
		}
	}

	return err
}

func checkDir(name string) error {
	stat, err := os.Stat(name)
	if err != nil {
		return fmt.Errorf("%s is not exists", name)
	}

	if !stat.IsDir() {
		return fmt.Errorf("%s is not a directory", name)
	}

	return nil
}

func createLogger() *zap.Logger {
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(filepath.Dir(os.Args[0]), "fmove.log"),
		MaxSize:    10, // megabytes
		MaxBackups: 3,
	})
	encoder := zap.NewProductionEncoderConfig()
	encoder.TimeKey = "time"
	encoder.EncodeTime = zapcore.RFC3339TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoder),
		writer,
		zap.InfoLevel,
	)

	return zap.New(core)
}
