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
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func MainMenu() {
	a := app.New()
	w := a.NewWindow("Artist Details")

	w.Resize(fyne.NewSize(1000, 500))

	w.SetFullScreen(true)

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
				groupInfo := fmt.Sprintf("Name : %s,", group.Name)
				img := data.FetchImage(group.Image)
				img.FillMode = canvas.ImageFillContain
				img.SetMinSize(fyne.NewSize(100, 100))
				groupLabel := widget.NewLabel(groupInfo)
				filtregroupecont.Add(groupLabel)
				filtregroupecont.Add(img)
			}
		}
	}

	numMembersContainer := container.New(layout.NewHBoxLayout())

	for i := 1; i <= 8; i++ {
		num := i
		check := widget.NewCheck(fmt.Sprintf("%d", num), func(checked bool) {
			if checked {
				updateArtistDisplay(artists, num)
			} else {
				filtregroupecont.RemoveAll()
			}
		})
		numMembersContainer.Add(check)
	}

	artistSugesstion := container.NewVBox()

	searchEntry.OnChanged = func(text string) {
		artistSugesstion.RemoveAll()
		mem, err := data.GetArtistByMember(text)
		art, errr := data.GetArtistsByName(text)
		if err != nil && errr != nil {
			artistSugesstion.Add(widget.NewLabel(err.Error()))
			return
		}
		if text == "" {
			artistSugesstion.RemoveAll()
			return
		}
		for _, artist := range art {
			artistSugesstion.Add(widget.NewLabel(fmt.Sprintf("Artist name: %v", artist)))
		}
		for _, member := range mem {
			artistSugesstion.Add(widget.NewLabel(fmt.Sprintf("Member name: %v", member)))
		}
	}

	grid := container.NewGridWithRows(0)
	count := 0
	listeArtist := container.NewVBox()

	for _, artist := range artists {
		art := artist
		img := data.FetchImage(artist.Image)
		img.FillMode = canvas.ImageFillContain
		img.SetMinSize(fyne.NewSize(100, 100))
		infoButton := widget.NewButton("", func() {
			PerformPostJsonRequest(w, art.Name)
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

	overlayfilter := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil))
	artistContainer := container.NewHBox()

	firstAlbumDateSliderLabel := widget.NewLabel("First album date: ")
	// Years Slider
	yearSlider := widget.NewSlider(0, 67)
	yearSlider.Step = 1
	yearSliderValue := widget.NewLabel("1957")
	yearSlider.OnChanged = func(value float64) {
		selectedYear := int(value) + 1957
		yearSliderValue.SetText(fmt.Sprintf("%d", selectedYear))
		artistContainer.RemoveAll()
		overlayfilter.RemoveAll()
		artists, err := data.GetArtistsByYear(selectedYear)
		if err != nil {
			return
		}

		for _, artist := range artists {
			img := data.FetchImage(artist.Image)
			img.FillMode = canvas.ImageFillContain
			img.SetMinSize(fyne.NewSize(100, 100))
			overlayfilter.Add(img)
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
		overlayfilter.RemoveAll()
		artists, err := data.GetArtistsByFirstAlbumYear(selectedFirstAlbumDate)
		if err != nil {
			return
		}

		for _, artist := range artists {
			img := data.FetchImage(artist.Image)
			img.FillMode = canvas.ImageFillContain
			img.SetMinSize(fyne.NewSize(100, 100))
			overlayfilter.Add(img)
			artistContainer.Add(widget.NewLabel(fmt.Sprintf("Name: %s", artist.Name)))
		}
	}

	content := container.NewVBox(
		titleContainer,
		searchEntry,
		searchButton,
		artistSugesstion,
		numMembersContainer,
		filtregroupecont,
		careerStartingYearLabel,
		yearSliderValue,
		yearSlider,
		firstAlbumDateSliderLabel,
		firstAlbumDateSliderValue,
		firstAlbumDateSlider,
		overlayfilter,
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

			fyne.NewMenuItem("Default theme", func() {
				a.Settings().SetTheme(theme.DefaultTheme())
			}),

			fyne.NewMenuItem("Light theme", func() {
				a.Settings().SetTheme(theme.LightTheme())
			}),

			fyne.NewMenuItem("Full Screen", func() {
				w.SetFullScreen(!w.FullScreen())
			}),

			fyne.NewMenuItem("favorites list", func() {
				FavoriteGestion(w, data.Favoris)
			}),

			fyne.NewMenuItem("ShortCut", func() {
				data.Shortcut(w)
			}),
		),
		fyne.NewMenu("Home",
			fyne.NewMenuItem("Home", func() {
				w.SetContent(scroll)
			}),
		))

	w.SetMainMenu(header)

	AltF4 := &desktop.CustomShortcut{KeyName: fyne.KeyF4, Modifier: fyne.KeyModifierAlt}
	w.Canvas().AddShortcut(AltF4, func(shortcut fyne.Shortcut) {
		os.Exit(0)
	})

	w.ShowAndRun()
}
