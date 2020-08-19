package fs

import (
	"context"
	"time"
)

//Watcher предназначен для отслеживания содержимого заданного каталога с заданным интервалом.
type Watcher struct {
	dirReader    DirReader
	pollInterval time.Duration
	events       chan []*File
	errors       chan error
}

//NewDirWatcher возвращает настроенный экземпляр Watcher.
func NewDirWatcher(dirReader DirReader, pollInterval time.Duration) *Watcher {
	return &Watcher{
		dirReader:    dirReader,
		pollInterval: pollInterval,
		events:       make(chan []*File),
		errors:       make(chan error, 1),
	}
}

//Watch начинает отслеживать (в бесконечном цикле) содержимое каталога.
//Агрумент ctx используется для остановки выполнения и выхода из метода.
//По окончании своей работы, метод закрывает каналы. Если завершение
//работы метода было вызвано ошибкой, то её в можно прочитать в канале ошибок.
func (w *Watcher) Watch(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	isTickerReset := false
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-ticker.C:
			entries, err := w.dirReader.Read()
			if err != nil {
				w.writeError(err)
				break loop
			}
			if len(entries) > 0 {
				w.writeEvent(entries)
			}

			if !isTickerReset {
				ticker.Reset(w.pollInterval)
				isTickerReset = true
			}
		}
	}

	ticker.Stop()
	close(w.events)
	close(w.errors)
}

func (w *Watcher) writeEvent(entries []*File) {
	select {
	case w.events <- entries:
	default:
	}
}

func (w *Watcher) writeError(error error) {
	w.errors <- error
}

//Events возвращает канал, в который пишется срез файлов, находящихся в отслеживаемой папке.
func (w *Watcher) Events() <-chan []*File {
	return w.events
}

//Errors возвращает канал, в который записывается ошибка не позволяющая экземпляру Watcher-ра
//выполнять свою работу (после появления ошибки в этом канале, работа завершается).
func (w *Watcher) Errors() <-chan error {
	return w.errors
}
