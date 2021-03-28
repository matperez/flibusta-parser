package pool

import (
	flibusta2 "github.com/matperez/flibusta-parser/internal/flibusta"
	"gorm.io/gorm"
)

var WorkerChannel = make(chan chan Work)

type Collector struct {
	Work chan Work
	End  chan bool
}

func StartDispatcher(workerCount int, db *gorm.DB, flb flibusta2.Client) Collector {
	var i int
	var workers []Worker
	input := make(chan Work) // channel to receive work
	end := make(chan bool)   // channel to spin down workers
	collector := Collector{Work: input, End: end}

	for i < workerCount {
		i++
		worker := Worker{
			ID:            i,
			Channel:       make(chan Work),
			WorkerChannel: WorkerChannel,
			End:           make(chan bool)}
		worker.Start(db, flb)
		workers = append(workers, worker) // store worker
	}

	// start collector
	go func() {
		for {
			select {
			case <-end:
				for _, w := range workers {
					w.Stop() // stop worker
				}
				return
			case work := <-input:
				worker := <-WorkerChannel // wait for available channel
				worker <- work            // dispatch work to worker
			}
		}
	}()

	return collector
}
