package data

import (
	"fmt"
	"gpo/structs"
	"image"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

func FetchImage(url string) *canvas.Image {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed to load image:", err)
		return nil
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		fmt.Println("Failed to decode image:", err)
		return nil
	}

	return canvas.NewImageFromImage(img)
}

func ShowArtistDetails(a fyne.App, artist structs.Artist) {
	w := a.NewWindow(artist.Name)
	info := fmt.Sprintf("Name: %s\nMembers: %v\nFirst Album: %s\nCreation Date: %d", artist.Name, artist.Members, artist.FirstAlbum, artist.CreationDate)
	w.SetContent(widget.NewLabel(info))
	w.Show()
}
