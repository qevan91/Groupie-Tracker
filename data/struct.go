package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var ErrArtistNotFound = errors.New("artist not found")
var ErrArtistRelationsNotFound = errors.New("relation not found")
var Art *Artist
var ArtistList []Artist
var Rel *Relation
var RelationList []Relation

type Location struct {
	ID    int      `json:"id"`
	Loc   []string `json:"locations"`
	Dates string   `json:"dates"`
}

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

type Query struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	ID                string    `json:"id"`
	Type              string    `json:"type"`
	PlaceType         []string  `json:"place_type"`
	Relevance         float64   `json:"relevance"`
	Properties        Property  `json:"properties"`
	Text              string    `json:"text"`
	PlaceName         string    `json:"place_name"`
	MatchingText      string    `json:"matching_text,omitempty"`
	MatchingPlaceName string    `json:"matching_place_name,omitempty"`
	Bbox              []float64 `json:"bbox,omitempty"`
	Center            []float64 `json:"center,omitempty"`
	Geometry          Geometry  `json:"geometry,omitempty"`
	Context           []Context `json:"context,omitempty"`
}

type Property struct {
	MapboxID   string `json:"mapbox_id,omitempty"`
	Wikidata   string `json:"wikidata,omitempty"`
	ShortCode  string `json:"short_code,omitempty"`
	Foursquare string `json:"foursquare,omitempty"`
	Landmark   bool   `json:"landmark,omitempty"`
	Category   string `json:"category,omitempty"`
	Address    string `json:"address,omitempty"`
}

type Geometry struct {
	Type        string    `json:"type,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty"`
}

type Context struct {
	ID        string `json:"id,omitempty"`
	MapboxID  string `json:"mapbox_id,omitempty"`
	Wikidata  string `json:"wikidata,omitempty"`
	ShortCode string `json:"short_code,omitempty"`
	Text      string `json:"text,omitempty"`
}

func float64ArrayToString(arr [][]float64) string {
	var result []string

	for _, innerArr := range arr {
		for _, num := range innerArr {
			result = append(result, strconv.FormatFloat(num, 'f', -1, 64))
		}
	}

	return strings.Join(result, ", ")
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

func GetLocations() string {

	apiKey := "pk.eyJ1IjoiZ3JwdHJrIiwiYSI6ImNsdHIzdXo0YzA4djYya3VsaHYzbWFtYWUifQ.UGOVoLVD4F0i-R8LFBcfvw" // Acces Token pour API Mapbox
	url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%d", Art.ID)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	var loc Location
	if err := json.NewDecoder(resp.Body).Decode(&loc); err != nil {
		fmt.Println(err)
	}
	var returnTab [][]float64
	for _, l := range loc.Loc {
		// Construction de l'URL de requête
		url := fmt.Sprintf("https://api.mapbox.com/geocoding/v5/mapbox.places/%s.json?access_token=%s", l, apiKey)
		response, err := http.Get(url)
		if err != nil {
			fmt.Println("Erreur lors de la requête:", err)
		}
		defer response.Body.Close()

		var datas Query
		if err := json.NewDecoder(response.Body).Decode(&datas); err != nil {
			fmt.Println(err)
		}
		returnTab = append(returnTab, datas.Features[0].Center)
	}
	return float64ArrayToString(returnTab)
}

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
