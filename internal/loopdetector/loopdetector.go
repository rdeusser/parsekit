package loopdetector

import (
	"sync"
)

type Detector struct {
	curPos  int
	prevPos int
	count   int
	limit   int
}

func New() *Detector {
	return &Detector{
		limit: 10,
	}
}

func (d *Detector) Detect(mu *sync.Mutex, pos int) {
	mu.Lock()
	defer mu.Unlock()
	if d.curPos == pos {
		d.increment()
	} else {
		d.reset()
	}
	d.prevPos = d.curPos
	d.curPos = pos
}

func (d *Detector) IsLooping() bool {
	return d.count == d.limit
}
func (d *Detector) increment() { d.count++ }
func (d *Detector) reset()     { d.count = 0 }
