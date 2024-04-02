package app

import (
	"gpo/data"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func FavoriteGestion(w fyne.Window, favliste []string) {
	if len(favliste) == 0 {
		dialog.ShowInformation("Error", "Your favorite list is empty", w)
		return
	}

	title := canvas.NewText("Artist list", color.White)
	title.TextSize = 20
	titleContainer := container.NewCenter(title)

	var artistWidgets []fyne.CanvasObject
	for _, artist := range favliste {
		artistName := artist
		numMember := widget.NewCheck(artist, func(checked bool) {
			if checked {
				data.DeleteArtist(artistName)
				dialog.NewInformation("Success", artistName+" has been removed from your favorite list", w).Show()
			} else {
				data.AddFavoris(artistName)
				dialog.NewInformation("Success", artistName+" has been re-added to your favorite list", w).Show()
			}
		})
		artistWidgets = append(artistWidgets, numMember)
	}
	artistContainer := container.NewVBox(artistWidgets...)

	favcontent := container.NewVBox(
		titleContainer,
		artistContainer,
	)
	w.SetContent(favcontent)
}
