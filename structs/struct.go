package structs

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var ErrArtistNotFound = errors.New("artist not found")
var Art *Artist
var ArtistList []Artist

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

	apiKey := "pk.eyJ1IjoiZ3JwdHJrIiwiYSI6ImNsdHIzdXo0YzA4djYya3VsaHYzbWFtYWUifQ.UGOVoLVD4F0i-R8LFBcfvw" // Acces Token pour API Mapbox
	locData := Art.Locations

	// Structure pour stocker les données JSON
	var data map[string]interface{}

	// Décodage des données JSON dans la structure
	if err := json.Unmarshal([]byte(locData), &data); err != nil {
		fmt.Println("Erreur lors du décodage JSON:", err)
	}

	for key := range data {
		fmt.Println("Clé:", key)

		// Construction de l'URL de requête
		url := fmt.Sprintf("https://api.mapbox.com/geocoding/v5/mapbox.places/%s.json?access_token=%s", key, apiKey)

		response, err := http.Get(url)
    	if err != nil {
        	fmt.Println("Erreur lors de la requête:", err)
    	}
		defer response.Body.Close()

		var dataGeoc map[string]interface{}
        if err := json.NewDecoder(response.Body).Decode(&dataGeoc); err != nil {
            fmt.Println("Erreur lors de l'analyse de la réponse:", err)
            continue
        }


	}

}

func GetConcertDates() string {
	return Art.ConcertDates
}

func GetRelation() string {
	return Art.Relation
}

func GetArtists() []Artist {
	return ArtistList
}

func FetchArtists() ([]Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []Artist
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, err
	}

	return artists, nil
}

func GetArtistByName(artistName string) (*Artist, error) {
	input := strings.ToLower(artistName)

	artists, err := FetchArtists()
	if err != nil {
		return nil, err
	}

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
