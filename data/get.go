package data

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// GetArtistByName retrieves artist information by name.
func GetArtistByName(artistName string) (*Artist, error) {
	// Convert input to lowercase
	input := strings.ToLower(artistName)

	// Fetch list of artists
	artists, err := FetchArtists()
	if err != nil {
		return nil, err
	}

	var artist *Artist
	// Iterate through each artist
	for _, a := range artists {
		p := strings.Join(a.Members, " ")
		// Check if input matches artist name or any member name
		if strings.Contains(strings.ToLower(a.Name), input) || strings.Contains(strings.ToLower(p), input) {
			artist = &a
			break
		}
	}

	// Return error if artist is not found
	if artist == nil {
		return nil, ErrArtistNotFound
	}

	// Set global artist variable
	Art = artist

	return artist, nil
}

// GetRelations retrieves relations information by city or date.
func GetRelations(City string) (*Relation, error) {
	// Convert input to lowercase
	input := strings.ToLower(City)

	// Fetch list of relations
	relations, err := FetchRelations()
	if err != nil {
		return nil, err
	}

	// Iterate through each relation
	for _, r := range relations {
		// Check if city name or date matches input
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

	// Return error if relations are not found
	return nil, ErrArtistRelationsNotFound
}

// GetArtistByID retrieves artist information by ID.
func GetArtistByID(RelationID int) (*Artist, error) {
	// Fetch list of artists
	artists, err := FetchArtists()
	if err != nil {
		return nil, err
	}

	var artist *Artist
	// Iterate through each artist
	for _, a := range artists {
		// Check if artist ID matches provided ID
		if Art.ID == RelationID {
			artist = &a
			break
		}
	}

	// Return error if artist is not found
	if artist == nil {
		return nil, ErrArtistNotFound
	}

	// Set global artist variable
	Art = artist

	return artist, nil
}

// GetRelationsByID retrieves relations information by artist ID.
func GetRelationsByID(artistID int) (*Relation, error) {
	// Fetch list of relations
	relations, err := FetchRelations()
	if err != nil {
		return nil, err
	}

	// Iterate through each relation
	for _, r := range relations {
		// Check if relation ID matches provided ID
		if r.ID == artistID {
			Rel = &r
			return &r, nil
		}
	}

	// Return error if relations are not found
	return nil, ErrArtistRelationsNotFound
}

// GetArtistsByYear retrieves artists created in a specific year.
func GetArtistsByYear(year int) ([]Artist, error) {
	// Fetch list of artists
	artists, err := FetchArtists()
	if err != nil {
		return nil, err
	}

	var filteredArtists []Artist
	// Iterate through each artist
	for _, artist := range artists {
		// Check if artist creation year matches provided year
		if artist.CreationDate == year {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	// Return error if no artists are found for the year
	if len(filteredArtists) == 0 {
		return nil, ErrNoArtistsForYear
	}

	return filteredArtists, nil
}

// GetArtistsByFirstAlbumYear retrieves artists whose first album was released in a specific year.
func GetArtistsByFirstAlbumYear(year int) ([]Artist, error) {
	// Fetch list of artists
	artists, err := FetchArtists()
	if err != nil {
		return nil, err
	}

	var filteredArtists []Artist

	for _, artist := range artists {
		// Extract album year from artist's first album date
		albumYear, err := strconv.Atoi(strings.Split(artist.FirstAlbum, "-")[2])
		if err != nil {
			fmt.Println("Error converting album year:", err)
			continue
		}

		// Check if album year matches provided year
		if albumYear == year {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	// Return error if no artists are found for the album year
	if len(filteredArtists) == 0 {
		return nil, ErrNoArtistsAlbumYears
	}

	return filteredArtists, nil
}

// GetArtistByMember retrieves artists containing a specific member.
func GetArtistByMember(artistMember string) ([]string, error) {
	// Convert input to lowercase
	input := strings.ToLower(artistMember)

	// Fetch list of artists
	artists, err := FetchArtists()
	if err != nil {
		return nil, err
	}

	var matchingMembers []string
	// Iterate through each artist
	for _, a := range artists {
		// Check if any member name contains the input
		for _, member := range a.Members {
			if strings.Contains(strings.ToLower(member), input) {
				matchingMembers = append(matchingMembers, member)
			}
		}
	}

	// Return error if no artists are found with the member
	if len(matchingMembers) == 0 {
		return nil, ErrArtistNotFound
	}

	return matchingMembers, nil
}

// GetArtistsByName retrieves artists by name
func GetArtistsByName(artistName string) ([]string, error) {
	// Convert input to lowercase
	input := strings.ToLower(artistName)

	// Fetch list of artists
	artists, err := FetchArtists()
	if err != nil {
		return nil, err
	}

	var artistNames []string
	// Iterate through each artist
	for _, a := range artists {
		// Check if artist name contains the input
		if strings.Contains(strings.ToLower(a.Name), input) {
			artistNames = append(artistNames, a.Name)
		}
	}

	if len(artistNames) == 0 {
		return nil, ErrArtistNotFound
	}

	return artistNames, nil
}

// GetFavoris retrieves the list of favorite artists.
func GetFavoris() []string {
	return Favoris
}

// GetArtistsByCreationDate retrieves artists created in a specific year.
func GetArtistsByCreationDate(year int) ([]Artist, error) {
	// Fetch list of artists
	artists, err := FetchArtists()
	if err != nil {
		return nil, err
	}

	var filteredArtists []Artist

	for _, artist := range artists {
		// Extract creation year from Unix timestamp
		creationYear := time.Unix(int64(artist.CreationDate), 0).Year()

		// Check if creation year matches provided year
		if creationYear == year {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	if len(filteredArtists) == 0 {
		return nil, ErrNoArtistsCreationYears
	}

	return filteredArtists, nil
}
