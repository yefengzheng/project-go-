package task

type Request struct {
	ImageNames []string
}

var RequestQueue chan Request

func InitQueues(queueSize int) {
	RequestQueue = make(chan Request, queueSize)
}
