package fuzzsearch

import (
	"fmt"
	"log"

	"github.com/ktr0731/go-fuzzyfinder"
)

type Track struct {
	Name      string
	AlbumName string
	Artist    string
}

var tracks = []Track{
	{"foo", "album1", "artist1"},
	{"bar", "album1", "artist1"},
	{"foo", "album2", "artist1"},
	{"baz", "album2", "artist2"},
	{"baz", "album3", "artist2"},
}

func Search() {
	idx, err := fuzzyfinder.Find(
		tracks,
		func(i int) string {
			return tracks[i].Name
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("selected: %v\n", idx)
}
