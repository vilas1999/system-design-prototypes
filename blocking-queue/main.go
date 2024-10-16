package main

import (
	"fmt"
	"sync"
)

type BlockingQueue struct {
	mu       sync.Mutex
	cond     sync.Cond
	data     []interface{}
	capacity int
}

func CreateNewBlockingQueue(capacity int) *BlockingQueue {
	q := &BlockingQueue{
		data:     make([]interface{}, 0, capacity),
		capacity: capacity,
	}
	q.cond = sync.Cond{L: &q.mu}
	return q
}

func (q *BlockingQueue) Put(item interface{}) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	for len(q.data) >= q.capacity {
		fmt.Println("Waiting for space to add data", q.data)
		q.cond.Wait()
	}

	q.data = append(q.data, item)
	q.cond.Signal()
}

func (q *BlockingQueue) Take() interface{} {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	for len(q.data) == 0 {
		fmt.Println("Waiting for data to be added", q.data)
		q.cond.Wait()
	}

	item := q.data[0]
	q.data = q.data[1:]
	q.cond.Signal()
	return item
}

func addItems(q *BlockingQueue, thread int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 5; i++ {
		formattedString := fmt.Sprintf("%d - %d", thread, i)
		fmt.Println("Added", formattedString)
		q.Put(formattedString)
	}

}

func removeItems(q *BlockingQueue, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		item := q.Take()
		fmt.Println("Removed item: ", item)
		if item == "end" {
			break
		}
	}
}

func main() {

	var removeWaitGroup, addWaitGroup sync.WaitGroup

	q := CreateNewBlockingQueue(2)

	addWaitGroup.Add(1)
	go addItems(q, 1, &addWaitGroup)
	addWaitGroup.Add(1)
	go addItems(q, 2, &addWaitGroup)
	addWaitGroup.Add(1)
	go addItems(q, 3, &addWaitGroup)

	go func() {
		addWaitGroup.Wait()
		q.Put("end")
	}()

	removeWaitGroup.Add(1)
	go removeItems(q, &removeWaitGroup)
	removeWaitGroup.Wait()
}
