package fonction

import (
	"encoding/json"
	"fmt"
	"net/http"

	"fyne.io/fyne"
	//"fyne.io/fyne/app"

	//"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// Artist représente la structure des données d'un artiste
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

func PerformPostJsonRequest(w fyne.Window, artistNAME string) {

	const myurl = "https://groupietrackers.herokuapp.com/api/artists"
	response, err := http.Get(myurl)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	var artists []Artist
	err = json.NewDecoder(response.Body).Decode(&artists)
	if err != nil {
		panic(err)
	}

	var Artist *Artist
	for _, artist := range artists {
		if artist.Name == artistNAME {
			Artist = &artist
			break
		}
	}

	if Artist == nil {
		fmt.Println("Artist not found")
		return
	}

	artistContainer := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Image: %v", Artist.Image)),
		widget.NewLabel(fmt.Sprintf("Name: %s", Artist.Name)),
		widget.NewLabel(fmt.Sprintf("Members: %v", Artist.Members)),
		widget.NewLabel(fmt.Sprintf("Creation Date: %d", Artist.CreationDate)),
		widget.NewLabel(fmt.Sprintf("First Album: %s", Artist.FirstAlbum)),
		widget.NewLabel(fmt.Sprintf("Locations: %s", Artist.Locations)),
		widget.NewLabel(fmt.Sprintf("Concert Dates: %s", Artist.ConcertDates)),
		widget.NewLabel(fmt.Sprintf("Relation: %s", Artist.Relation)),
	)

	buttonHome := widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() {
		w.MainMenu()
	})

	w.SetContent(container.NewVBox(
		buttonHome,
		widget.NewLabel("Artist Details"),
		artistContainer,
	))
}
