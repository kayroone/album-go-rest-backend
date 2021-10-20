package album

import (
	"fmt"
)

type Album struct {
	ID     int    `json:"id" db:"id"`
	Title  string `json:"title" db:"title"`
	Artist string `json:"artist" db:"artist"`
	Price  int    `json:"price" db:"price"`
}

func (a Album) printAlbum() {

	fmt.Printf("%d, %s, %s, %d", a.ID, a.Title, a.Artist, a.Price)
}
