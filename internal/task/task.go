package task

import (
	"log"
	"project-go-/internal/config"
	"project-go-/internal/database"

	"time"
)

type ImageTask struct {
	ImageName string
}

var DownloadQueue chan *ImageTask
var ScanQueue chan *ImageTask

func InitQueues(cfg *config.Config, dbCtx *database.Context) {
	DownloadQueue = make(chan *ImageTask, cfg.Worker.QueueSize)
	ScanQueue = make(chan *ImageTask, cfg.Worker.QueueSize)
	StartDownloadWorkers(cfg.Worker.DownloadWorkerCount, DownloadQueue, ScanQueue, dbCtx)
	StartScanWorkers(cfg.Worker.ScanWorkerCount, ScanQueue, dbCtx)
}
func StartDownloadWorkers(n int, in <-chan *ImageTask, out chan<- *ImageTask, dbCtx *database.Context) {
	for i := 0; i < n; i++ {
		go func(id int) {
			for t := range in {
				log.Printf("[Downloader-%d] downloading image: %s", id, t.ImageName)
				//TODO set redis record (second redis op)
				dbCtx.RedisContext.SetKeyValue(t.ImageName, "up", 10*time.Minute)
				time.Sleep(5 * time.Second)
				log.Printf("[Downloader-%d] download finished: %s", id, t.ImageName)
				out <- t
			}
		}(i)
	}
}

func StartScanWorkers(n int, in <-chan *ImageTask, dbCtx *database.Context) {
	for i := 0; i < n; i++ {
		go func(id int) {
			for t := range in {
				log.Printf("[Scanner-%d] scanning image: %s", id, t.ImageName)
				//TODO delete redis (3rd redis op)
				time.Sleep(2 * time.Second)
				dbCtx.RedisContext.DeleteKey(t.ImageName)
				log.Printf("[Scanner-%d] scan complete: %s", id, t.ImageName)
			}
		}(i)
	}
}

func CreateNewTask(imageName string) {
	task := ImageTask{imageName}
	DownloadQueue <- &task
}
