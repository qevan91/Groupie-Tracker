package main

import (
	"gpo/fonction"
)

func main() {
	// Crée les composants de l'interface utilisateur
	content := fonction.CreateUIComponents()

	// Configure la fenêtre et exécute l'application
	fonction.SetupWindowAndRun(content)
}
