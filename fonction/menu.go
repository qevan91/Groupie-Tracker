package fonction

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func createSearchEntry() *widget.Entry {
	// Crée un champ de texte pour la saisie de recherche
	searchEntry := widget.NewEntry()
	return searchEntry
}

func createSearchButton(searchEntry *widget.Entry) *widget.Button {
	// Crée un bouton de recherche avec une fonction de callback
	searchButton := widget.NewButton("Recherche", func() {
		// Récupère le terme de recherche saisi
		searchTerm := searchEntry.Text

		// Envoie une notification avec les résultats de la recherche
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Résultats de la recherche",
			Content: "Vous avez cherché : " + searchTerm,
		})
	})
	return searchButton
}

func CreateUIComponents() fyne.CanvasObject {
	// Crée les composants de l'interface utilisateur
	searchEntry := createSearchEntry()
	searchButton := createSearchButton(searchEntry)

	// Crée un conteneur pour organiser les éléments de l'interface utilisateur
	content := container.NewVBox(
		widget.NewLabel("Barre de Recherche"),
		searchEntry,
		searchButton,
	)
	return content
}

func SetupWindowAndRun(content fyne.CanvasObject) {
	// Crée une nouvelle instance de l'application Fyne
	myApp := app.New()

	// Crée une nouvelle fenêtre avec le titre "Hello"
	w := myApp.NewWindow("Hello")

	// Définit le contenu de la fenêtre
	w.SetContent(content)

	// Affiche la fenêtre et démarre l'application
	w.ShowAndRun()
}

func main() {
	// Crée les composants de l'interface utilisateur
	content := CreateUIComponents()

	// Configure la fenêtre et exécute l'application
	SetupWindowAndRun(content)
}
