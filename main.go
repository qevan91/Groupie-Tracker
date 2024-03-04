package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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

func main() {
	a := app.New()
	w := a.NewWindow("Artist Details")

	w.Resize(fyne.NewSize(1000, 500))

	buttonHome := widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() {
		w.MainMenu()
	})

	searchEntry := widget.NewEntry()

	searchButton := widget.NewButton("Recherche", func() {
		searchTerm := searchEntry.Text
		PerformPostJsonRequest(w, searchTerm)
	})

	content := container.NewVBox(
		buttonHome,
		searchEntry,
		searchButton,
	)

	w.SetContent(content)

	w.ShowAndRun()
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
	for _, artist := range artists {
		if artist.Name == artistNAME {
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

	// Définition du contenu de la fenêtre
	w.SetContent(container.NewVBox(
		buttonHome,
		widget.NewLabel("Détails de l'artiste"),
		content,
	))
}
