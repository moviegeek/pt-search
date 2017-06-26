package pt

import "bitbucket.org/laputa/movie-search/movie"

// Provider defines the interface for all PT torrent search providers
type Provider interface {
	FindAll(string) (movie.List, error)
}
