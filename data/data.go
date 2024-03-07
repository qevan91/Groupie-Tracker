package main

import (
	"encoding/json"
	"fmt"
	"image"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Artist représente la structure des données récupérées de l'API
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

// fetchArtists récupère les données de l'API
func fetchArtists() ([]Artist, error) {
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

// fetchImage récupère une image depuis une URL et retourne un canvas.Image
func fetchImage(url string) *canvas.Image {
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

// showArtistDetails ouvre une nouvelle fenêtre avec les détails de l'artiste
func showArtistDetails(a fyne.App, artist Artist) {
	w := a.NewWindow(artist.Name)
	info := fmt.Sprintf("Name: %s\nMembers: %v\nFirst Album: %s\nCreation Date: %d", artist.Name, artist.Members, artist.FirstAlbum, artist.CreationDate)
	w.SetContent(widget.NewLabel(info))
	w.Show()
}

func main() {
	a := app.New()
	w := a.NewWindow("Groupie Tracker")

	artists, err := fetchArtists()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur lors de la récupération des artistes: %v\n", err)
		return
	}

	content := container.NewVBox()

	for _, artist := range artists {
		artistInfo := fmt.Sprintf("Name: %s\nFirst Album: %s", artist.Name, artist.FirstAlbum)
		artistLabel := widget.NewLabel(artistInfo)

		img := fetchImage(artist.Image)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(100, 100))

		infoButton := widget.NewButton("", func() {
			showArtistDetails(a, artist)
		})
		infoButton.Importance = widget.LowImportance // Rend le bouton transparent
		infoButton.SetIcon(nil)                      // Aucune icône
		overlay := container.NewMax(img, infoButton) // Superpose l'image et le bouton transparent

		hbox := container.NewHBox(overlay, artistLabel)
		content.Add(hbox)
	}

	scroll := container.NewScroll(content)
	w.SetContent(scroll)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}
