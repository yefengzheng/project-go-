package main

import (
	"log"
	"os"
	"os/signal"
	"project-go-/internal/config"
	"project-go-/internal/database"
	"project-go-/internal/rest"
	"project-go-/internal/task"
	"project-go-/internal/worker"
	"time"
)

func main() {
	errCh := make(chan error, 1)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	dbCtx, err := database.CreateNewDbContext(cfg, 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	task.InitQueues(cfg, dbCtx)
	worker.StartDownloadWorkers(cfg.Worker.DownloadWorkerCount, task.DownloadQueue, task.ScanQueue, dbCtx)
	worker.StartScanWorkers(cfg.Worker.ScanWorkerCount, task.ScanQueue, dbCtx)
	restService := rest.NewRestService(cfg, dbCtx)
	go func() {
		errCh <- restService.Start()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errCh:
		log.Printf("Sservice Failed " + err.Error())
	case sig := <-sigs:
		log.Printf("Interrupt signal " + sig.String())
	}
	log.Println("Shutting down...")
}
