package store

type Store interface {
	Update(string, string)
	Lookup(string) (string, error)
}
