package models

type Restaurant struct {
	ID   int    `json:"-"`
	Name string `json:"name"`
}

// Restaurants is a slice of Restaurant pointers.
type Restaurants []*Restaurant
