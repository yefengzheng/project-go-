package task

type ImageTask struct {
	ImageName string
}

type Request struct {
	ImageNames []string `json:"image_names"`
}

var (
	RequestQueue  chan Request // ✅ Add this line
	DownloadQueue chan ImageTask
	ScanQueue     chan ImageTask // Optional if unused
)

func InitQueues(size int) {
	RequestQueue = make(chan Request, size) // ✅ Initialize it
	DownloadQueue = make(chan ImageTask, size)
	ScanQueue = make(chan ImageTask, size) // Optional if not used
}
