package app

import (
	"context"
	"fmt"
	"image/color"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"gpo/data"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

// PerformPostJsonRequest retrieves artist data based on the provided artist name
func PerformPostJsonRequest(w fyne.Window, artistName string) {
	// Retrieve artist data by searching the name
	artist, err := data.GetArtistByName(artistName)

	// Handle artist not found or errors
	if err != nil {
		searchEntry := widget.NewEntry()

		searchButton := widget.NewButton("Search", func() {
			searchTerm := searchEntry.Text
			PerformPostJsonRequest(w, searchTerm)
		})

		artistSugesstion := container.NewVBox()

		searchEntry.OnChanged = func(text string) {
			artistSugesstion.RemoveAll()

			if text == "" {
				artistSugesstion.RemoveAll()
				return
			}

			// Check if entered text is a number
			year, err := strconv.Atoi(text)
			if err == nil {
				// If text is a number, show suggestions based on years
				suggestedArtists, err := data.GetArtistsByFirstAlbumYear(year)
				if err != nil {
					artistSugesstion.Add(widget.NewLabel(fmt.Sprintf("Error fetching artists for year %d: %v", year, err)))
					return
				}
				for _, artist := range suggestedArtists {
					artistSugesstion.Add(widget.NewLabel(fmt.Sprintf("Artist with first album in %d: %s", year, artist.Name)))
				}
				return
			}

			// Check if entered text is a creation year
			creationYear, err := strconv.Atoi(text)
			if err == nil {
				// If text is a creation year, show suggestions based on that year
				suggestedArtists, err := data.GetArtistsByCreationDate(creationYear)
				if err != nil {
					artistSugesstion.Add(widget.NewLabel(fmt.Sprintf("Error fetching artists for creation year %d: %v", creationYear, err)))
					return
				}
				for _, artist := range suggestedArtists {
					artistSugesstion.Add(widget.NewLabel(fmt.Sprintf("Artist created in %d: %s", creationYear, artist.Name)))
				}
				return
			}

			mem, err := data.GetArtistByMember(text)
			art, errr := data.GetArtistsByName(text)
			if err != nil && errr != nil {
				artistSugesstion.Add(widget.NewLabel(fmt.Sprintf("Error: %v", err)))
				artistSugesstion.Add(widget.NewLabel(fmt.Sprintf("Error: %v", errr)))
				return
			}

			for _, artist := range art {
				artistSugesstion.Add(widget.NewLabel(fmt.Sprintf("Artist name: %v", artist)))
			}

			for _, member := range mem {
				artistSugesstion.Add(widget.NewLabel(fmt.Sprintf("Member name: %v", member)))
			}
		}

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

	if artist != nil {

		Favoris := widget.NewCheck("", nil)
		Favoris.SetText("Favorites")

		Music := widget.NewCheck("", nil)
		Music.SetText("Listen to their most popular music")

		// Add the artist to the favorite list
		Favoris.OnChanged = func(checked bool) {
			if checked {
				data.AddFavoris(artist.Name)
				dialog.ShowInformation("Success", artist.Name+" has been added to your favorite list", w)
			} else {
				data.DeleteArtist(artist.Name)
				dialog.ShowInformation("Success", artist.Name+" has been removed from your favorite list", w)
			}
		}

		// Listen the most popular music ofth artist
		Music.OnChanged = func(checked bool) {
			if checked {
				data.OpenLinkByID(artist.ID)
			}
		}

		numMembersContainer := container.New(layout.NewHBoxLayout(),
			Favoris, Music,
		)

		titleArt := canvas.NewText("Artist Details", color.White)
		titleArt.TextSize = 20
		titleArtContent := container.NewCenter(titleArt)

		// Writes artist api data
		artistContainer := container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Name: %s", data.GetName())),
			widget.NewLabel(fmt.Sprintf("Members: %v", data.GetMembers())),
			widget.NewLabel(fmt.Sprintf("Creation Date: %v", data.GetCreationDate())),
			widget.NewLabel(fmt.Sprintf("First Album: %s", data.GetFirstAlbum())),
			widget.NewLabel(fmt.Sprintf("Locations of the Concerts: %s", data.GetLocations())),
			widget.NewLabel(fmt.Sprintf("Concert Dates: %s", data.GetConcertDates())),
			widget.NewLabel(fmt.Sprintf("Relations: %v", data.GetRelation())),
		)
		artistContainer.Move(fyne.NewPos(50, 50))

		// Define the configuration for client credentials
		config := &clientcredentials.Config{
			ClientID:     "84bf3627f36f483a8d56ea10498bffdc",
			ClientSecret: "a0e9b4ad403f469eb3427d839ca71623",
			TokenURL:     spotify.TokenURL,
		}
		client := config.Client(context.Background())

		// Use the client to interact with the Spotify API
		spotifyClient := spotify.NewClient(client)

		result, err := spotifyClient.Search(artistName, spotify.SearchTypeArtist)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to search artist: %v\n", err)
			os.Exit(1)
		}

		SpotifyContent := container.NewVBox()

		// Check if there are artists found in the search results
		if len(result.Artists.Artists) > 0 {
			artistID := result.Artists.Artists[0].ID
			fmt.Println("Artist ID for", artistName, ":", artistID)

			// Search the artist top tracks
			tracks, err := spotifyClient.GetArtistsTopTracks(artistID, "US")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to retrieve artist's top tracks: %v\n", err)
				os.Exit(1)
			}

			// Processing results
			fmt.Println("Top tracks for artist:", artistID)
			for _, track := range tracks {
				fmt.Println(track.Name)
			}

			// Create a widget to display tracks
			tracksList := widget.NewLabel(fmt.Sprintf(" %s\n\n", artistName))
			for _, track := range tracks {
				tracksList.SetText(tracksList.Text + track.Name + "\n\n")
			}

			SpotifyContent.Add(tracksList)
		}

		title := canvas.NewText("Spotify", color.White)
		title.TextSize = 20
		titleContainer := container.NewCenter(title)

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

		searchButton := widget.NewButton("Search", func() {
			searchTerm := searchEntry.Text
			PerformPostJsonRequest(w, searchTerm)
		})

		artistSugesstion := container.NewVBox()

		searchEntry.OnChanged = func(text string) {
			artistSugesstion.RemoveAll()

			if text == "" {
				artistSugesstion.RemoveAll()
				return
			}

			// Check if entered text is a number
			year, err := strconv.Atoi(text)
			if err == nil {
				// If text is a number, show suggestions based on years
				suggestedArtists, err := data.GetArtistsByFirstAlbumYear(year)
				if err != nil {
					artistSugesstion.Add(widget.NewLabel(fmt.Sprintf("Error fetching artists for year %d: %v", year, err)))
					return
				}
				for _, artist := range suggestedArtists {
					artistSugesstion.Add(widget.NewLabel(fmt.Sprintf("Artist with first album in %d: %s", year, artist.Name)))
				}
				return
			}

			// Check if entered text is a creation year
			creationYear, err := strconv.Atoi(text)
			if err == nil {
				// If text is a creation year, show suggestions based on that year
				suggestedArtists, err := data.GetArtistsByCreationDate(creationYear)
				if err != nil {
					artistSugesstion.Add(widget.NewLabel(fmt.Sprintf("Error fetching artists for creation year %d: %v", creationYear, err)))
					return
				}
				for _, artist := range suggestedArtists {
					artistSugesstion.Add(widget.NewLabel(fmt.Sprintf("Artist created in %d: %s", creationYear, artist.Name)))
				}
				return
			}

			// Search a member
			mem, err := data.GetArtistByMember(text)
			// Search an Artsit
			art, errr := data.GetArtistsByName(text)
			if err != nil && errr != nil {
				artistSugesstion.Add(widget.NewLabel(fmt.Sprintf("Error: %v", err)))
				artistSugesstion.Add(widget.NewLabel(fmt.Sprintf("Error: %v", errr)))
				return
			}

			// Display Artist found
			for _, artist := range art {
				artistSugesstion.Add(widget.NewLabel(fmt.Sprintf("Artist name: %v", artist)))
			}

			// Display Member found
			for _, member := range mem {
				artistSugesstion.Add(widget.NewLabel(fmt.Sprintf("Member name: %v", member)))
			}
		}

		// Shortcuts
		ctrlL := &desktop.CustomShortcut{KeyName: fyne.KeyL, Modifier: fyne.KeyModifierControl}
		w.Canvas().AddShortcut(ctrlL, func(shortcut fyne.Shortcut) {
			Favoris.SetChecked(!Favoris.Checked)
		})

		ctrlF := &desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: fyne.KeyModifierControl}
		w.Canvas().AddShortcut(ctrlF, func(shortcut fyne.Shortcut) {
			w.SetFullScreen(true)
		})

		ctrlE := &desktop.CustomShortcut{KeyName: fyne.KeyE, Modifier: fyne.KeyModifierControl}
		w.Canvas().AddShortcut(ctrlE, func(shortcut fyne.Shortcut) {
			w.SetFullScreen(false)
		})

		AltF4 := &desktop.CustomShortcut{KeyName: fyne.KeyF4, Modifier: fyne.KeyModifierAlt}
		w.Canvas().AddShortcut(AltF4, func(shortcut fyne.Shortcut) {
			os.Exit(0)
		})

		scrollContainer := container.NewScroll(
			container.NewVBox(
				searchEntry,
				searchButton,
				artistSugesstion,
				numMembersContainer,
				titleArtContent,
				content,
				artistContainer,
				titleContainer,
				SpotifyContent,
			),
		)

		w.SetContent(scrollContainer)
		return
	}
}
