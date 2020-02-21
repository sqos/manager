package manager

const (
	NotifyAdd Operate = iota + 1
	NotifyDelete
	NotifyUpdate
)

type Operate int

type Notify struct {
	Entry   Entry
	Operate Operate
}
