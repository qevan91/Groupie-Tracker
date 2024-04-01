package app

import (
	"fmt"

	//"gpo/structs"

	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"gpo/data"
	//"image/color"
	//"os"
	//"fyne.io/fyne/v2/layout"
)

func PerformPostJsonRequest(w fyne.Window, artistName string) {
	artist, err := data.GetArtistByName(artistName)
	if err != nil {
		buttonHome := widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() {
			MainMenu()
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

	relation, err := data.GetRelationsByID(data.GetArtisteID())
	if err != nil {
		fmt.Println("Erreur lors de la récupération des relations:", err)
	}

	if artist != nil {

		buttonHome := widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() {
			MainMenu()
			//w.MainMenu()
			//fmt.Print("retour au menu 2") //fonction plus a cause de Ontapped
		})

		artistContainer := container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Name: %s", data.GetName())),
			widget.NewLabel(fmt.Sprintf("Members: %v", data.GetMembers())),
			widget.NewLabel(fmt.Sprintf("Creation Date: %v", data.GetCreationDate())),
			widget.NewLabel(fmt.Sprintf("First Album: %s", data.GetFirstAlbum())),
			widget.NewLabel(fmt.Sprintf("Locations: %s", data.GetLocations())),
			widget.NewLabel(fmt.Sprintf("Concert Dates: %s", data.GetConcertDates())),
			widget.NewLabel(fmt.Sprintf("Location: %v", relation)),
		)
		container.NewCenter(artistContainer)
		imgResource, err := fyne.LoadResourceFromURLString(data.GetImage())
		if err != nil {
			panic(err)
		}
		img := canvas.NewImageFromResource(imgResource)

		img.SetMinSize(fyne.NewSize(200, 200))
		img.FillMode = canvas.ImageFillContain

		content := container.NewVBox(
			img,
		)
		searchEntry := widget.NewEntry()

		searchButton := widget.NewButton("Recherche", func() {
			searchTerm := searchEntry.Text
			PerformPostJsonRequest(w, searchTerm)
		})

		handleHomeButton := func() {
			MainMenu()
		}

		buttonHome.OnTapped = handleHomeButton

		w.SetContent(container.NewVBox(
			buttonHome,
			searchEntry,
			searchButton,
			content,
			artistContainer,
		))

		return
	}
}
