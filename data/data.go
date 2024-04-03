package data

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// Transform a link in image
func FetchImage(url string) *canvas.Image {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed to load image:", err)
		return nil
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		fmt.Println("Failed to decode image:", err)
		return nil
	}

	return canvas.NewImageFromImage(img)
}

// Recovery the api Artist
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

// Recovery the api Relation
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

func ShowArtistDetails(a fyne.App, artist Artist) {
	w := a.NewWindow(artist.Name)
	info := fmt.Sprintf("Name: %s\nMembers: %v\nFirst Album: %s\nCreation Date: %d", artist.Name, artist.Members, artist.FirstAlbum, artist.CreationDate)
	w.SetContent(widget.NewLabel(info))
	w.Show()
}

// Open a link
func OpenLinkByID(id int) {
	// Open the file containing Spotify songs
	file, err := os.Open("SpotifySongs.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	counter := 0
	// Cheak each line in the file
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "https://open.spotify.com/intl-fr/track/") {
			counter++
			if counter == Art.ID {
				// Execute a command to open the Spotify link
				err := exec.Command("cmd", "/c", "start", scanner.Text()).Run()
				if err != nil {
					fmt.Println("Error during command execution:", err)
					return
				}
				return
			}
		}
	}

	// Check for any scanning errors
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Println("ID not found")
}

// Add an artist to the favorites list
func AddFavoris(artistName string) {
	if artistName != "" {
		p := strings.Join(Favoris, "")
		// Check if artist is already in the list of favorites
		if strings.Contains(p, artistName) {
			fmt.Println("Artist already in the list")
		} else {
			// Add artist to the list of favorites
			Favoris = append(Favoris, artistName)
		}
	}
}

// DeleteArtist removes an artist from the list of favorites.
func DeleteArtist(artistName string) {
	if artistName != "" {
		// Iterate through the list of favorites
		for i, artist := range Favoris {
			if artist == artistName {
				// Remove artist
				Favoris = append(Favoris[:i], Favoris[i+1:]...)
				return
			}
		}
	}
}
