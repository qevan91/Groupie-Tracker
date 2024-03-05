package search

import (
	"fmt"
	"gpo/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func PerformPostJsonRequest(w fyne.Window, artistNAME string) {
	artist, err := data.GetArtistByName(artistNAME)
	if err != nil {
		buttonHome := widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() {
			w.MainMenu()
		})

		searchEntry := widget.NewEntry()

		searchButton := widget.NewButton("Recherche", func() {
			searchTerm := searchEntry.Text
			PerformPostJsonRequest(w, searchTerm)
		})
		if err == data.ErrArtistNotFound {
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
			widget.NewLabel(fmt.Sprintf("Image: %v", data.GetImage())),
			widget.NewLabel(fmt.Sprintf("Name: %s", data.GetName())),
			widget.NewLabel(fmt.Sprintf("Members: %v", data.GetMembers())),
			widget.NewLabel(fmt.Sprintf("Creation Date: %v", data.GetCreationDate())),
			widget.NewLabel(fmt.Sprintf("First Album: %s", data.GetFirstAlbum())),
			widget.NewLabel(fmt.Sprintf("Locations: %s", data.GetLocations())),
			widget.NewLabel(fmt.Sprintf("Concert Dates: %s", data.GetConcertDates())),
			widget.NewLabel(fmt.Sprintf("Relation: %s", data.GetRelation())),
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
			widget.NewLabel("Artist Details"),
			artistContainer,
		))
		return
	}
}
