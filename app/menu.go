package app

import (
	"fmt"
	"net/url"
	"strconv"

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
	// Create the app
	a := app.New()
	w := a.NewWindow("Artist Details")

	w.Resize(fyne.NewSize(1000, 500))

	// Put in full screen
	w.SetFullScreen(true)

	// Fetch list of artists
	artists, err := data.FetchArtists()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erreur lors de la récupération des artistes: %v\n", err)
		return
	}

	// Display a title
	title := canvas.NewText("Groupie Tracker", color.White)
	title.TextSize = 20
	titleContainer := container.NewCenter(title)

	searchEntry := widget.NewEntry()

	// Validate input
	searchButton := widget.NewButton("Search", func() {
		searchTerm := searchEntry.Text
		// Redirect to the search page
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
	// Create check box
	for i := 1; i <= 8; i++ {
		num := i
		check := widget.NewCheck(fmt.Sprintf("%d", num), func(checked bool) {
			if checked {
				// Display the artist with the number of member
				updateArtistDisplay(artists, num)
			} else {
				filtregroupecont.RemoveAll()
			}
		})
		numMembersContainer.Add(check)
	}

	artistSuggestion := container.NewVBox()

	// Suggestion
	searchEntry.OnChanged = func(text string) {
		artistSuggestion.RemoveAll()

		if text == "" {
			return
		}

		// Try to convert text into a number (year or creation year)
		year, err := strconv.Atoi(text)
		if err == nil {
			// Search for artists with the first album in the given year
			FirstAlbArtists, err := data.GetArtistsByFirstAlbumYear(year)
			if err != nil {
				artistSuggestion.Add(widget.NewLabel(fmt.Sprintf("Error fetching artists for year %d: %v", year, err)))
				return
			}
			for _, artist := range FirstAlbArtists {
				// Display the year and the name of the artist
				artistSuggestion.Add(widget.NewLabel(fmt.Sprintf("Artist with first album in %d: %s", year, artist.Name)))
			}

			// Search for artists created in the given year
			suggestedArtists, err := data.GetArtistsByCreationDate(year)
			if err != nil {
				artistSuggestion.Add(widget.NewLabel(fmt.Sprintf("Error fetching artists for creation year %d: %v", year, err)))
				return
			}
			for _, artist := range suggestedArtists {
				artistSuggestion.Add(widget.NewLabel(fmt.Sprintf("Artist created in %d: %s", year, artist.Name)))
			}
		} else {
			// If the text is neither a number nor a creation year,
			// show suggestions based on artist names and members
			mem, err := data.GetArtistByMember(text)
			if err != nil {
				artistSuggestion.Add(widget.NewLabel(fmt.Sprintf("Error: %v", err)))
			}

			art, err := data.GetArtistsByName(text)
			if err != nil {
				artistSuggestion.Add(widget.NewLabel(fmt.Sprintf("Error: %v", err)))
			}

			for _, artist := range art {
				artistSuggestion.Add(widget.NewLabel(fmt.Sprintf("Artist name: %v", artist)))
			}

			for _, member := range mem {
				artistSuggestion.Add(widget.NewLabel(fmt.Sprintf("Member name: %v", member)))
			}
		}

	}

	grid := container.NewGridWithRows(0)
	count := 0
	listeArtist := container.NewVBox()

	// Display all the artist
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
		artistSuggestion,
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

	// Incorporates a scroll
	scroll := container.NewScroll(content)
	w.SetContent(scroll)

	// Create a header
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
				// Put on or take off the full screen
				w.SetFullScreen(!w.FullScreen())
			}),

			fyne.NewMenuItem("favorites list", func() {
				// Redirect to favotites list page
				FavoriteGestion(w, data.Favoris)
			}),

			fyne.NewMenuItem("ShortCut", func() {
				// Redirect to Shortcut page
				Shortcut(w)
			}),
		),
		fyne.NewMenu("Login",
			fyne.NewMenuItem("Login", func() {
				data.Login(w)
			}),
			fyne.NewMenuItem("Sign", func() {
				data.Signup(w, widget.NewEntry(), widget.NewEntry(), widget.NewLabel(""))
				err := data.CreateUser(widget.NewEntry().Text, widget.NewEntry().Text, data.Favoris)
				if err != nil {
					fmt.Println("Erreur lors de la création de l'utilisateur:", err)
				} else {
					fmt.Println("Utilisateur créé avec succès!")
				}

			}),
		),
		fyne.NewMenu("Home",
			fyne.NewMenuItem("Home", func() {
				// Display Mainmenu
				w.SetContent(scroll)
			}),
		),
		fyne.NewMenu("Spotify",
			fyne.NewMenuItem("Spotify Embeds", func() {
				// Open link of spotify embeds documentation
				link, _ := url.Parse("https://developer.spotify.com/documentation/embeds")
				_ = fyne.CurrentApp().OpenURL(link)
			}),

			fyne.NewMenuItem("Spotify", func() {
				// Open Spotify
				Spotifylink, _ := url.Parse("https://open.spotify.com/intl-fr")
				_ = fyne.CurrentApp().OpenURL(Spotifylink)
			}),
		),
	)

	w.SetMainMenu(header)

	// Shortcut

	AltF4 := &desktop.CustomShortcut{KeyName: fyne.KeyF4, Modifier: fyne.KeyModifierAlt}
	w.Canvas().AddShortcut(AltF4, func(shortcut fyne.Shortcut) {
		os.Exit(0)
	})

	ctrlF := &desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(ctrlF, func(shortcut fyne.Shortcut) {
		w.SetFullScreen(true)
	})

	ctrlE := &desktop.CustomShortcut{KeyName: fyne.KeyE, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(ctrlE, func(shortcut fyne.Shortcut) {
		w.SetFullScreen(false)
	})

	w.ShowAndRun()
}

/*

 */
