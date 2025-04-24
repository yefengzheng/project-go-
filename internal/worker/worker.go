package worker

import (
	"log"
	"project-go-/internal/task"
)

func StartDownloadWorkers(n int, in <-chan *task.ImageTask, out chan<- *task.ImageTask) {
	for i := 0; i < n; i++ {
		go func(id int) {
			for t := range task.DownloadQueue {
				log.Printf("[Downloader-%d] downloading image: %s", id, t.ImageName)
				
				log.Printf("[Downloader-%d] download finished: %s", id, t.ImageName)
				task.ScanQueue <- t
			}
		}(i)
	}
}

func StartScanWorkers(n int) {
	for i := 0; i < n; i++ {
		go func(id int) {
			for t := range task.ScanQueue {
				log.Printf("[Scanner-%d] scanning image: %s", id, t.ImageName)

				log.Printf("[Scanner-%d] scan complete: %s", id, t.ImageName)
			}
		}(i)
	}
}
