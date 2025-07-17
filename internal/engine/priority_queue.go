package engine

import (
	"container/heap"
	"sync"

	"github.com/itxiaohu001/taskmaster/pkg/task"
)

// PriorityQueue 优先级队列
type PriorityQueue struct {
	mu    sync.RWMutex
	items []*queueItem
	index map[string]int // 任务ID到索引的映射
}

// queueItem 队列项
type queueItem struct {
	task     task.Task
	priority int
	index    int
}

// NewPriorityQueue 创建新的优先级队列
func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{
		items: make([]*queueItem, 0),
		index: make(map[string]int),
	}
	heap.Init(pq)
	return pq
}

func (pq *PriorityQueue) Len() int {
	return len(pq.items)
}

func (pq *PriorityQueue) Less(i, j int) bool {
	return pq.items[i].priority < pq.items[j].priority
}

func (pq *PriorityQueue) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
	pq.index[pq.items[i].task.ID()] = i
	pq.index[pq.items[j].task.ID()] = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*queueItem)
	item.index = len(pq.items)
	pq.items = append(pq.items, item)
	pq.index[item.task.ID()] = item.index
}

func (pq *PriorityQueue) Pop() interface{} {
	old := pq.items
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	pq.items = old[0 : n-1]
	delete(pq.index, item.task.ID())
	return item
}

// 业务方法
func (pq *PriorityQueue) PushTask(t task.Task) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	item := &queueItem{
		task:     t,
		priority: t.GetPriority(),
	}
	heap.Push(pq, item)
}

func (pq *PriorityQueue) PopTask() task.Task {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if pq.Len() == 0 {
		return nil
	}
	item := heap.Pop(pq).(*queueItem)
	return item.task
}

// Remove 从队列中移除指定任务
func (pq *PriorityQueue) Remove(taskID string) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if index, exists := pq.index[taskID]; exists {
		heap.Remove(pq, index)
		delete(pq.index, taskID)
	}
}
