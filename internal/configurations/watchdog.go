package configurations

import (
	"log"
	"os"
	"time"
)

// FileWatch sets a watchdog on a file and calls the callback function on file changes
func FileWatch(path string, callback func(string) bool) {
	oldStat, _ := os.Stat(path)
	for {
		stat, _ := os.Stat(path)
		if oldStat.ModTime() != stat.ModTime() {
			oldStat = stat
			if changes := callback(path); changes {
				log.Println("Configuration file reloaded")
			}
		}
		time.Sleep(time.Second)
	}
}
