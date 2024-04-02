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

func ShowArtistDetails(a fyne.App, artist Artist) {
	w := a.NewWindow(artist.Name)
	info := fmt.Sprintf("Name: %s\nMembers: %v\nFirst Album: %s\nCreation Date: %d", artist.Name, artist.Members, artist.FirstAlbum, artist.CreationDate)
	w.SetContent(widget.NewLabel(info))
	w.Show()
}

func OpenLinkByID(id int) {
	file, err := os.Open("SpotifySongs.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	counter := 0
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "https://open.spotify.com/intl-fr/track/") {
			counter++
			if counter == Art.ID {
				err := exec.Command("cmd", "/c", "start", scanner.Text()).Run()
				if err != nil {
					fmt.Println("Error during command execution:", err)
					return
				}
				return
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Println("ID not found")
}
