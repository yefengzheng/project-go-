package main

import (
	"log"
	"os"
	"os/signal"
	"project-go-/internal/config"
	"project-go-/internal/rest"
	"project-go-/internal/task"
	"project-go-/internal/worker"
)

func main() {
	errCh := make(chan error, 1)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	task.InitQueues(cfg.Worker.QueueSize)
	worker.StartDownloadWorkers(cfg.Worker.DownloadWorkerCount)
	worker.StartScanWorkers(cfg.Worker.ScanWorkerCount)
	restService := rest.NewRestService(cfg)
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
