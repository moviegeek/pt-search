package pt

type PTMovie struct {
	From  string
	Title string
	Age   string
	Size  string
	Seeds string
}

// Provider defines the interface for all PT torrent search providers
type Provider interface {
	FindAll(string) ([]PTMovie, error)
}
