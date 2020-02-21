package manager

type Entry interface {
	AddAfter()
	DeleteAfter()
	UpdateAfter()

	Key() interface{}
	Copy(Entry)
}
