package conf

import (
	"github.com/fsnotify/fsnotify"
	"time"
)

// FileWatcher watches a file on a set interval, and preforms de-duplication of write
// events such that only 1 write event is reported even if multiple writes
// happened during the specified duration.
type FileWatcher struct {
	fsNotify *fsnotify.Watcher
	interval time.Duration
	done     chan struct{}
	callback func()
}

// Begin watching a file with a specific interval and action
func WatchFile(path string, interval time.Duration, action func()) (*FileWatcher, error) {
	// initialize a new fsnotify watcher
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// Add the file to be watched
	fsWatcher.Add(path)

	// initialize FileWatcher
	watcher := &FileWatcher{
		fsNotify: fsWatcher,
		interval: interval,
		done:     make(chan struct{}, 1),
		callback: action,
	}

	// Launch a new goroutine to watch the file
	watcher.run()

	// return the initialized instances
	return watcher, err
}

func (this *FileWatcher) run() {
	// Check for write events at this interval
	tick := time.Tick(this.interval)

	// Last write event from fsNotify
	var lastWriteEvent *fsnotify.Event
	for {
		select {
		case event := <-this.fsNotify.Events:
			// Kubernetes Usage:
			// When a ConfigMap update occurs kubernetes AtomicWriter() creates a new directory;
			// writing the updated ConfigMap contents to the new directory. Once the write is
			// complete it removes the original file symlink and replaces it with a new symlink
			// pointing to the contents of the newly created directory. It does this to achieve
			// atomic ConfigMap updates. But it also means the file we were monitoring for write
			// events never got them and was instead deleted.

			// The correct way to handle this would be to monitor the symlink instead of the
			// actual file for events. However fsnotify.v1 does not allow us to pass in the
			// IN_DON'T_FOLLOW flag to notify which would allow us to monitor the
			// symlink for changes instead of the de-referenced file. This is not likely to
			// change as fsnotify is designed as cross platform and not all platforms support
			// symlinks.
			if event.Op == fsnotify.Remove {
				// Since the symlink was removed, we must
				// re-register the file to be watched
				this.fsNotify.Remove(event.Name)
				this.fsNotify.Add(event.Name)
				lastWriteEvent = &event
			}
			// If it was a write event
			if event.Op == fsnotify.Write {
				lastWriteEvent = &event
			}
		case <-tick:
			// No events during this interval
			if lastWriteEvent == nil {
				continue
			}
			// Execute the callback
			this.callback()
			// Reset the last event
			lastWriteEvent = nil
		case <-this.done:
			goto Close
		}
	}
	// Close label just close the channel
Close:
	close(this.done)
}

func (this *FileWatcher) Close() {
	close(this.done)
	this.fsNotify.Close()
}
