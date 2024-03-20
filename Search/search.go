package search

/*
import (
	"fmt"
	"gpo/structs"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func PerformPostJsonRequest(w fyne.Window, artistNAME string) {
	artist, err := structs.GetArtistByName(artistNAME)
	if err != nil {
		buttonHome := widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() {
			w.MainMenu()
		})

		searchEntry := widget.NewEntry()

		searchButton := widget.NewButton("Recherche", func() {
			searchTerm := searchEntry.Text
			PerformPostJsonRequest(w, searchTerm)
		})
		if err == structs.ErrArtistNotFound {
			w.SetContent(container.NewVBox(
				buttonHome,
				searchEntry,
				searchButton,
				widget.NewLabel("Artist not found"),
			))
		} else {
			fmt.Println("Error:", err)
			return
		}
		return
	}

	if artist != nil {

		buttonHome := widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() {
			w.MainMenu()
		})

		artistContainer := container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Image: %v", structs.GetImage())),
			widget.NewLabel(fmt.Sprintf("Name: %s", structs.GetName())),
			widget.NewLabel(fmt.Sprintf("Members: %v", structs.GetMembers())),
			widget.NewLabel(fmt.Sprintf("Creation Date: %v", structs.GetCreationDate())),
			widget.NewLabel(fmt.Sprintf("First Album: %s", structs.GetFirstAlbum())),
			widget.NewLabel(fmt.Sprintf("Locations: %s", structs.GetLocations())),
			widget.NewLabel(fmt.Sprintf("Concert Dates: %s", structs.GetConcertDates())),
			widget.NewLabel(fmt.Sprintf("Relation: %s", structs.GetRelation())),
		)

		searchEntry := widget.NewEntry()

		searchButton := widget.NewButton("Recherche", func() {
			searchTerm := searchEntry.Text
			PerformPostJsonRequest(w, searchTerm)
		})

		w.SetContent(container.NewVBox(
			buttonHome,
			searchEntry,
			searchButton,
			artistContainer,
		))
		return
	}
}
*/

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

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
	// URL de l'API pour récupérer les détails des artistes
	const myurl = "https://groupietrackers.herokuapp.com/api/artists"

	// Récupération des données depuis l'API
	response, err := http.Get(myurl)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// Décodage des données JSON dans une structure d'artiste
	var artists []Artist
	err = json.NewDecoder(response.Body).Decode(&artists)
	if err != nil {
		panic(err)
	}

	// Recherche de l'artiste par son nom
	var Artist *Artist
	input := strings.ToLower(artistNAME)
	for _, artist := range artists {
		if strings.ToLower(artist.Name) == input {
			Artist = &artist
			break
		}
	}

	// Vérification si l'artiste est trouvé
	if Artist == nil {
		fmt.Println("Artiste non trouvé")
		return
	}

	// Création du conteneur pour les détails de l'artiste
	artistDetails := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Nom: %s", Artist.Name)),
		widget.NewLabel(fmt.Sprintf("Membres: %v", Artist.Members)),
		widget.NewLabel(fmt.Sprintf("Date de création: %d", Artist.CreationDate)),
		widget.NewLabel(fmt.Sprintf("Premier album: %s", Artist.FirstAlbum)),
		widget.NewLabel(fmt.Sprintf("Lieux: %s", Artist.Locations)),
		widget.NewLabel(fmt.Sprintf("Dates de concert: %s", Artist.ConcertDates)),
		widget.NewLabel(fmt.Sprintf("Relation: %s", Artist.Relation)),
	)

	// Chargement de l'image depuis l'URL
	imgResource, err := fyne.LoadResourceFromURLString(Artist.Image) // En supposant que Artist.Image contient l'URL de l'image
	if err != nil {
		panic(err)
	}
	img := canvas.NewImageFromResource(imgResource)

	// Redimensionnement de l'image
	img.SetMinSize(fyne.NewSize(150, 150))
	img.FillMode = canvas.ImageFillContain

	// Création du conteneur pour l'image et les détails de l'artiste
	content := container.NewVBox(
		img,           // L'image
		artistDetails, // Les détails de l'artiste
	)

	// Bouton pour retourner à la page d'accueil
	buttonHome := widget.NewButtonWithIcon("Accueil", theme.HomeIcon(), func() {
		w.MainMenu()
	})

	searchEntry := widget.NewEntry()

	searchButton := widget.NewButton("Recherche", func() {
		searchTerm := searchEntry.Text
		PerformPostJsonRequest(w, searchTerm)
	})

	// Définition du contenu de la fenêtre
	w.SetContent(container.NewVBox(
		buttonHome,
		searchEntry,
		searchButton,
		widget.NewLabel("Détails de l'artiste"),
		content,
	))
}
