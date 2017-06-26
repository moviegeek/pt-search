/*
Package movie provides basic data structure for a single movie, including
all useful informations
*/
package movie

// Item represents a single movie's information
type Item struct {
	Title string
}

// List is a list of movie objects
type List []Item
