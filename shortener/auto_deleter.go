package shortener

import (
	"context"
	"gorm.io/gorm"
	"log"
	"time"
)

/*
type priorityQueue []*time.Time

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].Before(*pq[j])
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue) Push(x interface{}) {
	item := x.(*time.Time)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *priorityQueue) update(item *ShortUrl, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}
*/

func DeleteExpired(ctx context.Context, db *gorm.DB) {
	dbUpdateDuration := 1 * time.Minute
	go func() {
		for {
			dbUpdateTimeout := time.After(dbUpdateDuration)
			select {
			case <-dbUpdateTimeout:
				log.Printf("Deleting expired URLs")
				err := db.Transaction(func(tx *gorm.DB) error {
					tx.Where("expires_at < localtime").Delete(ShortUrl{})
					return tx.Error
				})
				if err != nil {
					log.Printf("Could not delete expired URLs: %v", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
