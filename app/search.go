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

	"image/color"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/layout"
)

func PerformPostJsonRequest(w fyne.Window, artistName string) {
	artist, err := data.GetArtistByName(artistName)
	if err != nil {
		buttonHome := widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() {
			//w.MainMenu()
			//fmt.Print("retour au menu") sa ne revoir rien
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
			//w.MainMenu()
			//fmt.Print("retour au menu 2") //fonction plus a cause de Ontapped
		})

		artistContainer := container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Name: %s", data.GetName())),
			widget.NewLabel(fmt.Sprintf("Members: %v", data.GetMembers())),
			widget.NewLabel(fmt.Sprintf("Creation Date: %v", data.GetCreationDate())),
			widget.NewLabel(fmt.Sprintf("First Album: %s", data.GetFirstAlbum())),
			//widget.NewLabel(fmt.Sprintf("Locations: %s", data.GetLocations())),
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
		/*
			handleHomeButton := func() {
				fmt.Print("tu as cliquer /")
				backToHomeContent := BackToHome()
				w.SetContent(backToHomeContent.content)
			}

			buttonHome.OnTapped = handleHomeButton
		*/
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

func BackToHome() fyne.Window {
	a := app.New()
	w := a.NewWindow("Artist Details")

	w.Resize(fyne.NewSize(1000, 500))

	artists, err := data.FetchArtists()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while fetching artists: %v\n", err)
		return w
	}

	title := canvas.NewText("Groupie Tracker", color.White)
	title.TextSize = 20
	titleContainer := container.NewCenter(title)

	buttonHome := widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() {
		w.MainMenu()
		fmt.Print("Return to menu")
	})

	searchEntry := widget.NewEntry()

	searchButton := widget.NewButton("Search", func() {
		searchTerm := searchEntry.Text
		PerformPostJsonRequest(w, searchTerm)
	})

	artistContainer := container.NewHBox()

	overlay := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil))

	searchEntry.OnChanged = func(text string) {
		artistContainer.RemoveAll()
		overlay.RemoveAll()
		art, err := data.GetArtistByName(text)
		if err != nil {
			fmt.Println("Problem")
			return
		}
		img := data.FetchImage(art.Image)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(100, 100))
		if text == "" {
			artistContainer.RemoveAll()
			overlay.RemoveAll()
			return
		}
		overlay.Add(img)
		artistContainer.Add(widget.NewLabel(fmt.Sprintf("Name: %s", art.Name)))
	}

	grid := container.NewGridWithRows(0)
	count := 0
	listeArtist := container.NewVBox()

	for _, artist := range artists {
		img := data.FetchImage(artist.Image)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(100, 100))
		infoButton := widget.NewButton("", func() {
			data.ShowArtistDetails(a, artist)
		})
		infoButton.Importance = widget.LowImportance
		infoButton.SetIcon(nil)
		artistInfo := fmt.Sprintf(artist.Name)
		artistLabel := widget.NewLabel(artistInfo)
		overlay := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil),
			img,
			artistLabel,
			infoButton,
		)
		grid.Add(overlay)
		count++
		if count == 4 {
			count = 0
			listeArtist.Add(grid)
			grid = container.NewGridWithRows(0)
		}
	}
	if count > 0 {
		listeArtist.Add(grid)
	}

	content := container.NewVBox(
		titleContainer,
		buttonHome,
		searchEntry,
		searchButton,
		overlay,
		artistContainer,
		listeArtist,
	)

	scroll := container.NewScroll(content)
	w.SetContent(scroll)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
	return w
}
