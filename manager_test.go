package manager

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
)

type EntryTester struct {
	Name  string
	Email string
	Alias string
}

func (e *EntryTester) AddAfter() {
	_, _ = fmt.Fprintln(os.Stdout, "AddAfter:", e)
}

func (e *EntryTester) DeleteAfter() {
	_, _ = fmt.Fprintln(os.Stdout, "DeleteAfter:", e)
}

func (e *EntryTester) UpdateAfter() {
	_, _ = fmt.Fprintln(os.Stdout, "UpdateAfter:", e)
}

func (e *EntryTester) Key() interface{} {
	return e.Email
}

func (e *EntryTester) Copy(n Entry) {
	if ne, ok := n.(*EntryTester); ok {
		e.Name = ne.Name
		e.Email = ne.Email
		e.Alias = ne.Alias
	}
}

func Sort(entries []Entry) []Entry {
	sort.SliceStable(entries, func(i, j int) bool {
		if strings.Compare(entries[i].Key().(string), entries[j].Key().(string)) <= 0 {
			return true
		} else {
			return false
		}
	})
	return entries
}

var (
	tests = []*EntryTester{
		{Name: "Wang Qiang", Email: "wq@qq.com",},
		{Name: "Li Hong", Email: "lh@qq.com",},
		{Name: "Zhang Lei", Email: "zl@qq.com",},
	}
)

func TestMain(m *testing.M) {
	NotifyRegisterHandler(func(ch <-chan Notify) {
		for n := range ch {
			_, _ = fmt.Fprintln(os.Stdout, "receive:", n.Operate, n.Entry)
		}
	})
	SortRegisterHandler(Sort)
	m.Run()
}

func TestAdd(t *testing.T) {
	for _, e := range tests {
		Add(e)
		t.Log("Get:", Get(e.Key()))
	}
	for _, e := range GetAll() {
		t.Log(e)
	}
}

func TestUpdate(t *testing.T) {
	for _, e := range tests {
		Add(e)
		t.Log("Get:", Get(e.Key()))
	}
	for _, e := range tests {
		o, n := *e, *e
		n.Alias = "This is alias test"
		Update(&n)
		t.Log("old", o, "new:", Get(e.Key()))
	}
}

func TestDelete(t *testing.T) {
	for _, e := range tests {
		Add(e)
		t.Log("Get:", Get(e.Key()))
	}
	for _, e := range tests {
		o := Delete(e.Key())
		t.Log("Delete:", o)
	}
}
