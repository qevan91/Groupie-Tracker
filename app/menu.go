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

	searchEntry := widget.NewEntry()

	searchButton := widget.NewButton("Recherche", func() {
		searchTerm := searchEntry.Text
		PerformPostJsonRequest(w, searchTerm)
	})

	filtregroupecont := container.NewHBox()

	updateArtistDisplay := func(allGroups []data.Artist, num int) {
		filtregroupecont.RemoveAll()
		for _, group := range allGroups {
			if len(group.Members) == num {
				groupInfo := fmt.Sprintf("Nom du groupe : %s,", group.Name)
				img := data.FetchImage(group.Image)
				img.FillMode = canvas.ImageFillContain
				img.SetMinSize(fyne.NewSize(100, 100))
				groupLabel := widget.NewLabel(groupInfo)
				filtregroupecont.Add(groupLabel)
				filtregroupecont.Add(img)
			}
		}
	}

	numMembers1 := widget.NewCheck("1", nil)
	numMembers2 := widget.NewCheck("2", nil)
	numMembers3 := widget.NewCheck("3", nil)
	numMembers4 := widget.NewCheck("4", nil)
	numMembers5 := widget.NewCheck("5", nil)
	numMembers6 := widget.NewCheck("6", nil)
	numMembers7 := widget.NewCheck("7", nil)
	numMembers8 := widget.NewCheck("8", nil)

	numMembersContainer := container.New(layout.NewHBoxLayout(),
		numMembers1, numMembers2, numMembers3, numMembers4, numMembers5, numMembers6, numMembers7, numMembers8,
	)

	numMembers1.OnChanged = func(checked bool) {
		if checked {
			updateArtistDisplay(artists, 1)
		} else {
			filtregroupecont.RemoveAll()
		}
	}

	numMembers2.OnChanged = func(checked bool) {
		if checked {
			updateArtistDisplay(artists, 2)
		} else {
			filtregroupecont.RemoveAll()
		}
	}

	numMembers3.OnChanged = func(checked bool) {
		if checked {
			updateArtistDisplay(artists, 3)
		} else {
			filtregroupecont.RemoveAll()
		}
	}

	numMembers4.OnChanged = func(checked bool) {
		if checked {
			updateArtistDisplay(artists, 4)
		} else {
			filtregroupecont.RemoveAll()
		}
	}

	numMembers5.OnChanged = func(checked bool) {
		if checked {
			updateArtistDisplay(artists, 5)
		} else {
			filtregroupecont.RemoveAll()
		}
	}

	numMembers6.OnChanged = func(checked bool) {
		if checked {
			updateArtistDisplay(artists, 6)
		} else {
			filtregroupecont.RemoveAll()
		}
	}

	numMembers7.OnChanged = func(checked bool) {
		if checked {
			updateArtistDisplay(artists, 7)
		} else {
			filtregroupecont.RemoveAll()
		}
	}

	numMembers8.OnChanged = func(checked bool) {
		if checked {
			updateArtistDisplay(artists, 8)
		} else {
			filtregroupecont.RemoveAll()
		}
	}

	filters := container.New(layout.NewGridLayout(2),
		numMembersContainer,
	)

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
		searchEntry,
		searchButton,
		filters,
		filtregroupecont,
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

	header := fyne.NewMainMenu(

		fyne.NewMenu("Settings",
			fyne.NewMenuItem("Dark theme", func() {
				a.Settings().SetTheme(theme.DarkTheme())
			}),

			fyne.NewMenuItem("Light theme", func() {
				a.Settings().SetTheme(theme.LightTheme())
			}),

			fyne.NewMenuItem("Full Screen", func() {
				if w.FullScreen() == true {
					w.SetFullScreen(false)
				} else {
					w.SetFullScreen(true)
				}
			}),
		),
		fyne.NewMenu("Home",
			fyne.NewMenuItem("Home", func() {
				w.SetContent(scroll)
			}),
		))

	w.SetMainMenu(header)

	w.ShowAndRun()
}
