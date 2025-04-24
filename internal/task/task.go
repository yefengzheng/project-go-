package task

type ImageTask struct {
	ImageName string
}

var DownloadQueue chan ImageTask
var ScanQueue chan ImageTask

func InitQueues(queueSize int) {
	DownloadQueue = make(chan ImageTask, queueSize)
	ScanQueue = make(chan ImageTask, queueSize)

}
