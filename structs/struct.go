package structs

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

var ErrArtistNotFound = errors.New("artist not found")
var ErrArtistRelationsNotFound = errors.New("relation not found")
var Art *Artist
var ArtistList []Artist
var Rel *Relation
var RelationList []Relation

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

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func GetArtisteID() int {
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

/*func GetLocations() string {

	apiKey := "pk.eyJ1IjoiZ3JwdHJrIiwiYSI6ImNsdHIzdXo0YzA4djYya3VsaHYzbWFtYWUifQ.UGOVoLVD4F0i-R8LFBcfvw" // Acces Token pour API Mapbox
	locData := Art.Locations

	client := &http.Client{} // Créer un client HTTP

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
	}

}*/

func GetRelation() string {
	return Art.Relation
}

func GetConcertDates() string {
	return Art.ConcertDates
}

func GetArtists() []Artist {
	return ArtistList
}

func GetRelationID() int {
	return Rel.ID
}

func GetRelationList() ([]Relation, error) {
	relations, err := FetchRelations()
	if err != nil {
		return nil, err
	}

	RelationList = relations

	return RelationList, nil
}

func GetDateLocations() map[string][]string {
	if GetArtisteID() == GetRelationID() {
		return Rel.DatesLocations
	}
	return Rel.DatesLocations
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

func FetchRelations() ([]Relation, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var relations []Relation

	var response struct {
		Index []Relation `json:"index"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	relations = response.Index

	return relations, nil
}

func GetArtistByName(artistName string) (*Artist, error) {
	input := strings.ToLower(artistName)

	artists, err := FetchArtists()
	if err != nil {
		return nil, err
	}

	var artist *Artist
	for _, a := range artists {
		p := strings.Join(a.Members, " ")
		if strings.Contains(strings.ToLower(a.Name), input) || strings.Contains(strings.ToLower(p), input) {
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

func GetRelations(City string) (*Relation, error) {
	input := strings.ToLower(City)

	relations, err := FetchRelations()
	if err != nil {
		return nil, err
	}

	for _, r := range relations {
		for city, dates := range r.DatesLocations {
			if strings.Contains(strings.ToLower(city), input) {
				Rel = &r
				return &r, nil
			}
			for _, date := range dates {
				if strings.Contains(strings.ToLower(date), input) {
					Rel = &r
					return &r, nil
				}
			}
		}
	}

	return nil, ErrArtistRelationsNotFound
}

func GetArtistByID(RelationID int) (*Artist, error) {
	artists, err := FetchArtists()
	if err != nil {
		return nil, err
	}

	var artist *Artist
	for _, a := range artists {
		if Art.ID == RelationID {
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

func GetRelationsByID(artistID int) (*Relation, error) {
	relations, err := FetchRelations()
	if err != nil {
		return nil, err
	}

	for _, r := range relations {
		if r.ID == artistID {
			Rel = &r
			return &r, nil
		}
	}

	return nil, ErrArtistRelationsNotFound
}
