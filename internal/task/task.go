package task

import (
	"project-go-/internal/config"
	"project-go-/internal/database"
	"time"
)

type ImageTask struct {
	ImageName string
	SHA256    string
	StartTime time.Time
}

var DownloadQueue chan *ImageTask
var ScanQueue chan *ImageTask

func InitQueues(cfg *config.Config, dbCtx *database.Context) {
	DownloadQueue = make(chan *ImageTask, cfg.Worker.QueueSize)
	ScanQueue = make(chan *ImageTask, cfg.Worker.QueueSize)
}

func CreateNewTask(imageName string, SHA256 string) {
	task := ImageTask{ImageName: imageName, SHA256: SHA256}
	DownloadQueue <- &task
}
