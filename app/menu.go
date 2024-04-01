package app

import (
	"fmt"

	"gpo/data"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func MainMenu() {
	a := app.New()
	w := a.NewWindow("Artist Details")

	w.Resize(fyne.NewSize(1000, 500))

	artists, err := data.FetchArtists()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur lors de la récupération des artistes: %v\n", err)
		return
	}

	title := canvas.NewText("Groupie Tracker", color.White)
	title.TextSize = 20
	titleContainer := container.NewCenter(title)

	buttonHome := widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() {
		w.MainMenu()
		fmt.Print("retour au menu")
	})

	searchEntry := widget.NewEntry()

	searchButton := widget.NewButton("Recherche", func() {
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
			fmt.Println("Problème")
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

	firstAlbumDateSliderLabel := widget.NewLabel("First album date: ")
	// Years Slider
	yearSlider := widget.NewSlider(0, 67)
	yearSlider.Step = 1
	yearSliderValue := widget.NewLabel("1957")
	yearSlider.OnChanged = func(value float64) {
		selectedYear := int(value) + 1957
		yearSliderValue.SetText(fmt.Sprintf("%d", selectedYear))
		artistContainer.RemoveAll()
		overlay.RemoveAll()
		artists, err := data.GetArtistsByYear(selectedYear)
		if err != nil {
			fmt.Println("Problème:", err)
			return
		}

		for _, artist := range artists {
			img := data.FetchImage(artist.Image)
			img.FillMode = canvas.ImageFillContain
			img.SetMinSize(fyne.NewSize(100, 100))
			overlay.Add(img)
			artistContainer.Add(widget.NewLabel(fmt.Sprintf("Name: %s", artist.Name)))
		}
	}

	careerStartingYearLabel := widget.NewLabel("Career Starting Year: ")

	// First album date Years
	firstAlbumDateSlider := widget.NewSlider(0, 62)
	firstAlbumDateSlider.Step = 1
	firstAlbumDateSliderValue := widget.NewLabel("1962")
	firstAlbumDateSlider.OnChanged = func(value float64) {
		selectedFirstAlbumDate := int(value) + 1962
		firstAlbumDateSliderValue.SetText(fmt.Sprintf("%d", selectedFirstAlbumDate))
		artistContainer.RemoveAll()
		overlay.RemoveAll()
		artists, err := data.GetArtistsByFirstAlbumYear(selectedFirstAlbumDate)
		if err != nil {
			fmt.Println("Problème:", err)
			return
		}

		for _, artist := range artists {
			img := data.FetchImage(artist.Image)
			img.FillMode = canvas.ImageFillContain
			img.SetMinSize(fyne.NewSize(100, 100))
			overlay.Add(img)
			artistContainer.Add(widget.NewLabel(fmt.Sprintf("Name: %s", artist.Name)))
		}
	}

	content := container.NewVBox(
		titleContainer,
		buttonHome,
		searchEntry,
		searchButton,
		careerStartingYearLabel,
		yearSliderValue,
		yearSlider,
		firstAlbumDateSliderLabel,
		firstAlbumDateSliderValue,
		firstAlbumDateSlider,
		overlay,
		artistContainer,
		listeArtist,
	)

	scroll := container.NewScroll(content)
	w.SetContent(scroll)
	w.Resize(fyne.NewSize(800, 600))
	w.ShowAndRun()
}

/*
1973
1967
1963
*/
