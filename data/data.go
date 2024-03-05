package data

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

var ErrArtistNotFound = errors.New("artist not found")
var Art *Artist

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concert_dates"`
	Relation     string   `json:"relations"`
}

func GetID() int {
	return Art.ID
}

func GetImage() string {
	if Art == nil {
		return "Image not available"
	}
	return Art.Image
}

func GetName() string {
	return Art.Name
}

func GetMembers() []string {
	return Art.Members
}

func GetCreationDate() int {
	return Art.CreationDate
}

func GetFirstAlbum() string {
	return Art.FirstAlbum
}

func GetLocations() string {
	return Art.Locations
}

func GetConcertDates() string {
	return Art.ConcertDates
}

func GetRelation() string {
	return Art.Relation
}

func GetArtistByName(artistName string) (*Artist, error) {
	const myurl = "https://groupietrackers.herokuapp.com/api/artists"
	response, err := http.Get(myurl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var artists []Artist
	err = json.NewDecoder(response.Body).Decode(&artists)
	if err != nil {
		return nil, err
	}

	input := strings.ToLower(artistName)

	var artist *Artist
	for _, a := range artists {
		if strings.ToLower(a.Name) == input {
			artist = &a
			break
		}
	}

	if artist == nil {
		return nil, ErrArtistNotFound
	}

	Art = artist

	return artist, nil
}
