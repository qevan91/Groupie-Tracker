package app

import (
	"gpo/data"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// Gestion of the favorite list
func FavoriteGestion(w fyne.Window, favliste []string) {
	// Check if the favorite list is empty
	if len(favliste) == 0 {
		dialog.ShowInformation("Error", "Your favorite list is empty", w)
		return
	}

	title := canvas.NewText("Artist list", color.White)
	title.TextSize = 20
	titleContainer := container.NewCenter(title)

	var artistWidgets []fyne.CanvasObject

	// Iterate through the favorite list
	for _, artist := range favliste {
		artistName := artist
		// Create a checkbox widget for each artist
		numMember := widget.NewCheck(artist, func(checked bool) {
			if checked {
				data.DeleteArtist(artistName)
				dialog.NewInformation("Success", artistName+" has been removed from your favorite list", w).Show()
			} else {
				data.AddFavoris(artistName)
				dialog.NewInformation("Success", artistName+" has been re-added to your favorite list", w).Show()
			}
		})
		// Add the checkbox widget to the list
		artistWidgets = append(artistWidgets, numMember)
	}
	// Create a vertical container for artist widgets
	artistContainer := container.NewVBox(artistWidgets...)

	favcontent := container.NewVBox(
		titleContainer,
		artistContainer,
	)

	w.SetContent(favcontent)

	// Shortcut
	ctrlF := &desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(ctrlF, func(shortcut fyne.Shortcut) {
		w.SetFullScreen(true)
	})

	ctrlE := &desktop.CustomShortcut{KeyName: fyne.KeyE, Modifier: fyne.KeyModifierControl}
	w.Canvas().AddShortcut(ctrlE, func(shortcut fyne.Shortcut) {
		w.SetFullScreen(false)
	})
}
