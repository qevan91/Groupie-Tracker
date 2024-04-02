package app

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"gpo/data"
)

func PerformPostJsonRequest(w fyne.Window, artistName string) {
	artist, err := data.GetArtistByName(artistName)
	if err != nil {
		searchEntry := widget.NewEntry()

		searchButton := widget.NewButton("Recherche", func() {
			searchTerm := searchEntry.Text
			PerformPostJsonRequest(w, searchTerm)
		})

		if err == data.ErrArtistNotFound {
			w.SetContent(container.NewVBox(
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

		Favoris := widget.NewCheck("", nil)
		Favoris.SetText("Favoris")

		numMembersContainer := container.New(layout.NewHBoxLayout(),
			Favoris,
		)

		Favoris.OnChanged = func(checked bool) {
			if checked {
				data.AddFavoris(artist.Name)
				dialog.ShowInformation("Success", artist.Name+" has been added to your favorite list", w)
			} else {
				data.DeleteArtist(artist.Name)
				dialog.ShowInformation("Success", artist.Name+" has been removed from your favorite list", w)
			}
		}

		artistContainer := container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Name: %s", data.GetName())),
			widget.NewLabel(fmt.Sprintf("Members: %v", data.GetMembers())),
			widget.NewLabel(fmt.Sprintf("Creation Date: %v", data.GetCreationDate())),
			widget.NewLabel(fmt.Sprintf("First Album: %s", data.GetFirstAlbum())),
			widget.NewLabel(fmt.Sprintf("Locations: %s", data.GetLocations())),
			widget.NewLabel(fmt.Sprintf("Concert Dates: %s", data.GetConcertDates())),
			widget.NewLabel(fmt.Sprintf("Relations: %v", relation)),
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

		w.SetContent(container.NewVBox(
			searchEntry,
			searchButton,
			numMembersContainer,
			content,
			artistContainer,
		))

		return
	}
}
