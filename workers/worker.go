package workers

import (
	"log/slog"
	"sync"
)

// Job adalah unit pekerjaan yang akan dieksekusi di background
type Job struct {
	Name string
	Execute func() error
}

var (
	jobQueue chan Job
	wg sync.WaitGroup
)

// Start menjalankan worker pool dengan jumlah worker tertentu
func Start(workerCount int, queueSize int) {
	jobQueue = make(chan Job, queueSize)
	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go worker(i)
	}

	slog.Info("Worker pool started", "workers", workerCount, "queueSize", queueSize)
}

// Enqueue menambahkan job ke antrian
func Enqueue(job Job) {
	select {
		case jobQueue <- job:
			slog.Info("Job enqueued", "name", job.Name)
		default:
			slog.Warn("Job queue full, dropping job", "job", job.Name)
	}
}

// Stop menunggu semua job selesai lalu menutup worker pool
func Stop() {
	close(jobQueue)
	wg.Wait()
	slog.Info("Worker pool stopped")
}

func worker(id int) {
	defer wg.Done()

	for job := range jobQueue {
		slog.Info("Worker processing job", "workerId", id, "job", job.Name)
		if err := job.Execute(); err != nil {
			slog.Error("Job failed", "workerId", id, "job", job.Name, "error", err)
		} else {
			slog.Info("Job completed", "workerId", id, "job", job.Name)
		}
	}
}