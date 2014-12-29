package inmem

import ( 
	"sync"
)

type Listener chan *Page;

type PageListenerList struct {
	mutex          sync.Mutex // public functions are threadsafe
	listeners []Listener;
}

func (pll *PageListenerList) Add(l Listener) {
	pll.mutex.Lock()
	pll.listeners = append(pll.listeners, l);
	pll.mutex.Unlock()
}

// remove any occurances of l from this set of listeners
func (pll *PageListenerList) Remove(l Listener) {
	pll.mutex.Lock()
	new := make([]Listener, len(pll.listeners)-1)
	for _,value := range pll.listeners {
		if value != l {
			new = append(new, value);
		}
	}
	pll.listeners = new
	pll.mutex.Unlock()
}

// this can block, if one of the listener queues is full, so you may
// want to use "go pll.Notify(page)"
func (pll *PageListenerList) Notify(page *Page) {
	pll.mutex.Lock()
	snapshot := make([]Listener, len(pll.listeners))
	copy(snapshot, pll.listeners)
	pll.mutex.Unlock()
	for _,l := range snapshot {
		l <- page
	}
}
