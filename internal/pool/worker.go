package pool

import (
	flibusta2 "github.com/matperez/flibusta-parser/internal/flibusta"
	"github.com/matperez/flibusta-parser/internal/work"
	"gorm.io/gorm"
	"log"
)

type Work struct {
	ID     int
	BookID int
}

type Worker struct {
	ID            int
	WorkerChannel chan chan Work
	Channel       chan Work
	End           chan bool
}

// start worker
func (w *Worker) Start(db *gorm.DB, flb flibusta2.Client) {
	log.Printf("worker [%d] is starting", w.ID)
	go func() {
		for {
			w.WorkerChannel <- w.Channel
			select {
			case job := <-w.Channel:
				// do work
				work.DoWork(db, flb, job.BookID, w.ID)
			case <-w.End:
				return
			}
		}
	}()
}

// end worker
func (w *Worker) Stop() {
	log.Printf("worker [%d] is stopping", w.ID)
	w.End <- true
}
