package pt

type PTMovie struct {
	ID     string `json:"id"`
	From   string `json:"from"`
	Title  string `json:"title"`
	Age    string `json:"age"`
	Size   string `json:"size"`
	Seeder string `json:"seeder"`
	URL    string `json:"url"`
}

// Provider defines the interface for all PT torrent search providers
type Provider interface {
	FindAll(string) ([]PTMovie, error)
}
