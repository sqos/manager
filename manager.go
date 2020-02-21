package manager

import (
	"sync"
)

type Manager struct {
	sync.Mutex
	sync.Map
	notifyChan chan Notify
	sort       func([]Entry) []Entry
}

func (c *Manager) notifyOperate(e Entry, op Operate) {
	if c.notifyChan == nil {
		return
	}
	c.notifyChan <- Notify{
		Entry:   e,
		Operate: op,
	}
}

func (c *Manager) NotifyChan() chan<- Notify {
	return c.notifyChan
}

func (c *Manager) notifyClose() {
	if c.notifyChan != nil {
		close(c.notifyChan)
		c.notifyChan = nil
	}
}

func (c *Manager) NotifyClose() {
	c.Lock()
	defer c.Unlock()

	c.notifyClose()
}

func (c *Manager) NotifyRegisterHandler(handler func(ch <-chan Notify)) {
	c.Lock()
	defer c.Unlock()

	c.notifyClose()
	c.notifyChan = make(chan Notify)
	if handler == nil {
		handler = func(ch <-chan Notify) {
			for range ch {
			}
		}
	}
	go handler(c.notifyChan)
}

func (c *Manager) SortRegisterHandler(sort func([]Entry) []Entry) {
	c.sort = sort
}

func (c *Manager) Get(key interface{}) Entry {
	v, ok := c.Map.Load(key)
	if !ok {
		return nil
	}
	if e, ok := v.(Entry); ok {
		return e
	} else {
		return nil
	}
}

func (c *Manager) GetAll() (entries []Entry) {
	c.Map.Range(func(k, v interface{}) bool {
		if s, ok := v.(Entry); ok {
			entries = append(entries, s)
		}
		return true
	})

	if sort := c.sort; sort != nil {
		return sort(entries)
	}
	return entries
}

func (c *Manager) update(e Entry) bool {
	if old := c.Get(e.Key()); old == nil {
		return false
	} else {
		old.Copy(e)
		e.UpdateAfter()
		c.notifyOperate(e, NotifyUpdate)
		return true
	}
}

func (c *Manager) Update(e Entry) bool {
	c.Lock()
	defer c.Unlock()

	return c.update(e)
}

func (c *Manager) Add(e Entry) bool {
	c.Lock()
	defer c.Unlock()

	if c.update(e) {
		return false
	}
	c.Map.Store(e.Key(), e)
	e.AddAfter()
	c.notifyOperate(e, NotifyAdd)
	return true
}

func (c *Manager) Delete(key interface{}) Entry {
	c.Lock()
	defer c.Unlock()

	e := c.Get(key) // just for other use after delete it
	c.Map.Delete(key)
	if e != nil {
		e.DeleteAfter()
		c.notifyOperate(e, NotifyDelete)
	}
	return e
}

var Default = &Manager{}

func NotifyChan() chan<- Notify {
	return Default.NotifyChan()
}

func SortRegisterHandler(sort func([]Entry) []Entry) {
	Default.SortRegisterHandler(sort)
}

func NotifyClose() {
	Default.NotifyClose()
}

func NotifyRegisterHandler(handler func(ch <-chan Notify)) {
	Default.NotifyRegisterHandler(handler)
}

func Get(key interface{}) Entry {
	return Default.Get(key)
}

func GetAll() (entries []Entry) {
	return Default.GetAll()
}

func Update(e Entry) bool {
	return Default.Update(e)
}

func Add(e Entry) bool {
	return Default.Add(e)
}

func Delete(key interface{}) Entry {
	return Default.Delete(key)
}
