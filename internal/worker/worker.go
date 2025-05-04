package worker

import (
	"log"
	"project-go-/internal/Tcontrol"
	"project-go-/internal/database"
	"project-go-/internal/task"
	"time"
)

func ProcessRequest(req task.Request, dbCtx database.Context) {
	log.Printf("Start processing request with %d images\n", len(req.ImageNames))

	controller := &Tcontrol.TaskContro{}

	// Phase 1: Download all images
	for _, image := range req.ImageNames {
		lockKey := "image:" + image

		// Check Redis for ongoing scan
		status, err := dbCtx.RedisContext.GetValue(lockKey)
		if err == nil && status == "up" {
			log.Printf("Image %s is already being scanned. Skipping.\n", image)
			continue
		}

		// Mark as "up" before any processing
		_ = dbCtx.RedisContext.SetKeyValue(lockKey, "up", 10*time.Minute)

		log.Printf("Downloading: %s", image)
		time.Sleep(1 * time.Second) // Simulate download
		if err := dbCtx.RedisContext.SetKeyValue("download:"+image, "done", 5*time.Minute); err != nil {
			log.Printf("Failed to set download status for %s in Redis: %v", image, err)
		}
	}

	// Phase 2: Scan all images (after download is complete)
	for _, image := range req.ImageNames {
		controller.Begin()
		go func(img string) {
			defer controller.Done()
			log.Printf("Scanning: %s", img)
			time.Sleep(2 * time.Second) // Simulate scan
			log.Printf("Scan complete: %s", img)

			// Mark scan as finished
			_ = dbCtx.RedisContext.SetKeyValue("image:"+img, "down", 0)
		}(image)
	}

	controller.Wait()
	log.Println("Finished request\n")
}
