package worker

import (
	"log"
	"project-go-/internal/config"
	"project-go-/internal/database"
	"project-go-/internal/task"
	"time"
)

func StartDownloadWorkers(n int, in <-chan *task.ImageTask, out chan<- *task.ImageTask, dbCtx *database.Context) {
	for i := 0; i < n; i++ {
		go func(id int) {
			for t := range in {
				log.Printf("[Downloader-%d] downloading image: %s", id, t.ImageName)
				//TODO set redis record (second redis op)
				dbCtx.RedisContext.SetKeyValue(t.ImageName, "up", 10*time.Minute)
				t.StartTime = time.Now()
				time.Sleep(5 * time.Second)
				log.Printf("[Downloader-%d] download finished: %s", id, t.ImageName)
				out <- t
			}
		}(i)
	}
}

func StartScanWorkers(n int, in <-chan *task.ImageTask, dbCtx *database.Context) {
	for i := 0; i < n; i++ {
		go func(id int) {
			for t := range in {
				log.Printf("[Scanner-%d] scanning image: %s", id, t.ImageName)
				//TODO delete redis (3rd redis op)
				time.Sleep(2 * time.Second)
				dbCtx.RedisContext.DeleteKey(t.ImageName)

				//correctSHA, err := dbCtx.PgsqlContext.GetSHA256ByImageName(t.ImageName)
				//if err != nil {
				//	log.Printf("[Scanner-%d] error querying SHA256: %v", id, err)
				//	continue
				//}

				// Compare with the uploaded SHA256 (assume it's passed in task.ImageTask)
				//uploadedSHA := t.SHA256 // you’ll need to add SHA256 field to ImageTask

				result := "safe"
				malicious := "false"
				//if uploadedSHA != correctSHA {
				//	result = "different"
				//	malicious = true
				//}

				end := time.Now()
				scan := config.ScanRequest{
					ImageName:      t.ImageName,
					ScanStartTime:  t.StartTime.Format("2006-01-02 15:04:05"),
					ScanFinishTime: end.Format("2006-01-02 15:04:05"),
					ScanResult:     result,
					MaliciousFiles: []string{malicious}, // optional
				}
				if err := dbCtx.PgsqlContext.InsertScanResult(scan); err != nil {
					log.Printf("[Scanner-%d] failed to insert result: %v", id, err)
				} else {
					log.Printf("[Scanner-%d] result stored: %s", id, result)
				}
				log.Printf("[Scanner-%d] scan complete: %s", id, t.ImageName)
			}
		}(i)
	}
}
