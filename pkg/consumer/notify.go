package consumer

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-logr/logr"
)

func StartWatchingDir(ctx context.Context, dir string, callback func(filePath string)) error {
	log := logr.FromContextOrDiscard(ctx)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		waitFor := 1 * time.Second

		// Keep track of the timers, as [path]timer.
		timers := sync.Map{}

		for {
			select {
			case <-ctx.Done():
				log.V(1).Info("Stopping watcher")
				break
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error(err, "")
			case e, ok := <-watcher.Events:
				if !ok {
					return
				}

				// We just want to watch for file creation, so ignore everything outside Create and Write.
				if !e.Has(fsnotify.Create) && !e.Has(fsnotify.Write) {
					continue
				}
				log.V(2).Info("New Event", "name", e.Name, "op", e.Op)

				// Get timer.
				t, exists := timers.Load(e.Name)

				if !exists {
					t = time.AfterFunc(math.MaxInt64, func() {
						log.V(2).Info("Deleted timer", "name", e.Name)
						timers.Delete(e.Name)
						callback(e.Name)
					})
					t.(*time.Timer).Stop()
					timers.Store(e.Name, t)
				}
				// Reset the timer for this path, so it will start from 100ms again.
				t.(*time.Timer).Reset(waitFor)
			}
		}
	}()
	err = watcher.Add(dir)
	if err != nil {
		return fmt.Errorf("cannot start watcher: %w", err)
	}
	log.V(1).Info("Started watcher", "dir", dir)
	return nil
}
