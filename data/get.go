package data

import (
	"fmt"
	"strconv"
	"strings"
)

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

// filtre artiste par carriere
func GetArtistsByYear(year int) ([]Artist, error) {
	artists, err := FetchArtists()
	if err != nil {
		return nil, err
	}

	var filteredArtists []Artist
	for _, artist := range artists {
		if artist.CreationDate == year {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	if len(filteredArtists) == 0 {
		return nil, ErrNoArtistsForYear
	}

	return filteredArtists, nil
}

func GetArtistsByFirstAlbumYear(year int) ([]Artist, error) {
	artists, err := FetchArtists()
	if err != nil {
		return nil, err
	}

	var filteredArtists []Artist

	for _, artist := range artists {
		albumYear, err := strconv.Atoi(strings.Split(artist.FirstAlbum, "-")[2])
		if err != nil {
			fmt.Println("Erreur de conversion de l'année de l'album:", err)
			continue
		}

		if albumYear == year {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	if len(filteredArtists) == 0 {
		return nil, ErrNoArtistsAlbumYears
	}

	return filteredArtists, nil
}

func GetArtistByMember(artistMember string) ([]string, error) {
	input := strings.ToLower(artistMember)

	artists, err := FetchArtists()
	if err != nil {
		return nil, err
	}

	var matchingMembers []string
	for _, a := range artists {
		for _, member := range a.Members {
			if strings.Contains(strings.ToLower(member), input) {
				matchingMembers = append(matchingMembers, member)
			}
		}
	}

	if len(matchingMembers) == 0 {
		return nil, ErrArtistNotFound
	}

	return matchingMembers, nil
}

func GetArtistsByName(artistName string) ([]string, error) {
	input := strings.ToLower(artistName)

	artists, err := FetchArtists()
	if err != nil {
		return nil, err
	}

	var artistNames []string
	for _, a := range artists {
		if strings.Contains(strings.ToLower(a.Name), input) {
			artistNames = append(artistNames, a.Name)
		}
	}

	if len(artistNames) == 0 {
		return nil, ErrArtistNotFound
	}

	return artistNames, nil
}

func GetFavoris() []string {
	return Favoris
}

func AddFavoris(artistName string) {
	if artistName != "" {
		p := strings.Join(Favoris, "")
		if strings.Contains(p, artistName) {
			fmt.Println("Artiste déjà dans la liste")
		} else {
			Favoris = append(Favoris, artistName)
		}
	}
}

func DeleteArtist(artistName string) {
	if artistName != "" {
		for i, artist := range Favoris {
			if artist == artistName {
				Favoris = append(Favoris[:i], Favoris[i+1:]...)
				return
			}
		}
	}
}
