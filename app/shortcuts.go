package app

import (
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func Shortcut(w fyne.Window) {
	title := canvas.NewText("Shortcut List", color.White)
	title.TextSize = 20
	titleContainer := container.NewCenter(title)

	// List of shortcut
	shortcontent := container.NewVBox(
		titleContainer,
		widget.NewLabel("ctrl + E Exit in fullscreen"),
		widget.NewLabel("ctrl + F Put in fullscreen"),
		widget.NewLabel("ctrl + L add to favotites list"),
		widget.NewLabel("alt + F4 Quit the application"),
	)

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

	w.SetContent(shortcontent)
}
