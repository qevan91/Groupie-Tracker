package search

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
