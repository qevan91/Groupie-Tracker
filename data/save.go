package data

import (
	"bufio"
	"fmt"

	//"gpo/data"
	"image/color"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Login(w fyne.Window) {
	title := canvas.NewText("Login", color.White)
	title.TextSize = 20
	titleContainer := container.NewCenter(title)

	label := widget.NewLabel("")
	entryUsername := widget.NewEntry()
	entryPassword := widget.NewPasswordEntry()
	form := widget.NewForm(
		widget.NewFormItem("Username", entryUsername),
		widget.NewFormItem("Password", entryPassword),
	)

	form.OnCancel = func() {
		label.Text = "Canceled"
		label.Refresh()
	}

	form.OnSubmit = func() {
		username := entryUsername.Text
		password := entryPassword.Text
		ok, err := VerifyUser(username, password)
		if err != nil {
			label.Text = "Error: " + err.Error()
			label.Refresh()
			return
		}
		if ok {
			label.Text = "Connected"
			label.Refresh()

		} else {
			label.Text = "Invalid username or password"
			label.Refresh()
		}
	}

	logintcontent := container.NewVBox(
		titleContainer,
		form,
		label,
	)
	w.SetContent(logintcontent)
}

func CreateUser(Log string, Pass string) error {
	file, err := os.OpenFile("users.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	userData := fmt.Sprintf("Username: %s, Password: %s\n", Log, Pass)

	_, err = file.WriteString(userData)
	if err != nil {
		return err
	}

	fmt.Println("Utilisateur créé avec succès!")
	return nil
}

func Signup(w fyne.Window, entryUsername *widget.Entry, entryPassword *widget.Entry, label *widget.Label) {
	form := widget.NewForm(
		widget.NewFormItem("Username", entryUsername),
		widget.NewFormItem("Password", entryPassword),
	)

	form.OnSubmit = func() {
		username := entryUsername.Text
		password := entryPassword.Text

		exists, err := VerifyUser(username, password)
		if err != nil {
			label.Text = "Error: " + err.Error()
			label.Refresh()
			return
		}
		if exists {
			label.Text = "Utilisateur déjà existant"
			label.Refresh()
			return
		}

		err = CreateUser(username, password)
		if err != nil {
			label.Text = "Erreur lors de la création de l'utilisateur: " + err.Error()
			label.Refresh()
			return
		}

		label.Text = "Utilisateur créé avec succès!"
		label.Refresh()
	}
	w.SetContent(container.NewVBox(form, label))
}

func VerifyUser(username string, password string) (bool, error) {
	file, err := os.Open("users.txt")
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Username: "+username) && strings.Contains(line, "Password: "+password) {
			return true, nil
		}
	}
	return false, nil
}
