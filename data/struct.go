package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var ErrArtistNotFound = errors.New("artist not found")
var ErrArtistRelationsNotFound = errors.New("relation not found")
var ErrNoArtistsForYear = errors.New("aucun artiste trouvé pour cette année")
var ErrNoArtistsAlbumYears = errors.New("aucun album trouvé pour cette année")
var ErrNoArtistsCreationYears = errors.New("aucun artist trouvé pour cette année")
var Art *Artist
var ArtistList []Artist
var Rel *Relation
var RelationList []Relation
var Favoris []string

// Struct of location
type Location struct {
	ID    int      `json:"id"`
	Loc   []string `json:"locations"`
	Dates string   `json:"dates"`
}

// Struct of the artist
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

// Struct of relation
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

// Struct for Spotify
type Track struct {
	Name    string
	Artists []struct {
		Name string
	}
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
		// Construction of the request URL
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
	url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%d", Art.ID)

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Erreur lors de la requête:", err)
	}
	defer response.Body.Close()

	var rel Relation
	if err := json.NewDecoder(response.Body).Decode(&rel); err != nil {
		fmt.Println(err)
	}

	returnString := ""
	i := 1
	for k, v := range rel.DatesLocations {
		if len(v) <= 1 {
			returnString += fmt.Sprintf("Date %d : "+"%s : %s\n", i, ConvertStringPropper(k), v[0])
			i++
		} else {
			stringup := ""
			for i := 0; i < len(v); i++ {
				stringup = fmt.Sprintf("%s, %s", stringup, v[i])
			}
			returnString += fmt.Sprintf("Date %d : "+"%s : %s\n", i, ConvertStringPropper(k), stringup)
			i++
		}
	}
	return "\n" + returnString + "\n" + "Don't miss it if you can still catch it !"
}

func GetConcertDates() string {
	return Art.ConcertDates
}

func GetArtists() ([]Artist, error) {
	artists, err := FetchArtists()
	if err != nil {
		return nil, err
	}
	ArtistList = artists
	return ArtistList, nil
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

func ConvertStringPropper(input string) string {
	// Convertir en majuscules
	upper := strings.ToUpper(input)

	// Supprimer les caractères spéciaux sauf les espaces
	regex := regexp.MustCompile("[^a-zA-Z0-9 ]+")
	result := regex.ReplaceAllString(upper, "")

	return result
}
