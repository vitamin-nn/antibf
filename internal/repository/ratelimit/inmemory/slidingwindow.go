package inmemory

import (
	"sync"
)

type SlidingWindow struct {
	mutex *sync.Mutex
	data  map[string]map[int64]int
}

func NewSlidingWindow() *SlidingWindow {
	sw := new(SlidingWindow)
	sw.data = make(map[string]map[int64]int)
	return sw
}

func (sw *SlidingWindow) GetRate(key string, startTimestamp int64) int {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()

	var result int
	tsData, ok := sw.data[key]
	if !ok {
		return result
	}

	for ts, count := range tsData {
		if ts > startTimestamp {
			result += count
		} else {
			delete(sw.data[key], ts)
		}
	}
	return result
}

func (sw *SlidingWindow) Increment(key string, timestamp int64) {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()

	sw.data[key][timestamp]++
}

func (sw *SlidingWindow) Clear(key string) {
	sw.mutex.Lock()
	defer sw.mutex.Unlock()

	delete(sw.data, key)
}
