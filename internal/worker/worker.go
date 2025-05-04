package worker

import (
	"log"
	"project-go-/internal/Tcontrol"
	"project-go-/internal/database"
	"project-go-/internal/task"
	"time"
)

// ProcessRequest handles one request (download all images, then scan them)
func ProcessRequest(req task.Request, dbCtx database.Context) {
	log.Printf("Start processing request with %d images\n", len(req.ImageNames))

	controller := &Tcontrol.TaskContro{}

	for _, image := range req.ImageNames {
		log.Printf("Downloading: %s", image)
		// simulate download
		time.Sleep(1 * time.Second)
		dbCtx.RedisContext.SetKeyValue("download:"+image, "done", 5*time.Minute)

		// Begin a scan task
		controller.Begin()
		go func(img string) {
			defer controller.Done()
			log.Printf("Scanning: %s", img)
			time.Sleep(2 * time.Second)
			log.Printf("Scan complete: %s", img)
		}(image)
	}

	controller.Wait()
	log.Println("Finished request\n")
}
