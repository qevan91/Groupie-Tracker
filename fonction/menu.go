package fonction

import (
	"fmt"
	//"fyne.io/fyne/v2"
	//"fyne.io/fyne/v2/container"
	//"fyne.io/fyne/v2/app"
	//"fyne.io/fyne/v2/widget"
)

func Saucisse() {
	fmt.Print("démarage")
}

/*
func CreateUI() fyne.Window {
	// Crée une nouvelle instance de l'application Fyne
	a := app.New()

	// Crée une nouvelle fenêtre avec le titre "Hello"
	w := a.NewWindow("Hello")

	// Crée un champ de texte pour la saisie de recherche
	searchEntry := widget.NewEntry()

	// Crée un bouton de recherche
	searchButton := widget.NewButton("Recherche", func() {
		// Récupère le terme de recherche saisi
		searchTerm := searchEntry.Text

		// Envoie une notification avec les résultats de la recherche
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Résultats de la recherche",
			Content: "Vous avez cherché : " + searchTerm,
		})
	})

	// Crée un conteneur pour organiser les éléments de l'interface utilisateur
	content := container.NewVBox(
		widget.NewLabel("Barre de Recherche"),
		searchEntry,
		searchButton,
	)

	// Définit le contenu de la fenêtre
	w.SetContent(content)

	return w
}

func Test() {
	a := app.New()
	w := a.NewWindow("Hello")

	w.ShowAndRun()
}
*/
