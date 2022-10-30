package main

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type config struct {
	EditWidget    *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile   fyne.URI
	SaveMenuItem  *fyne.MenuItem
}

func main() {
	// create a fyne app
	a := app.New()
	// create a window for the app
	win := a.NewWindow("Markdown")
	// get the user interface
	edit, preview := cfg.makeUI()

	// set the content of the window
	win.SetContent(container.NewHSplit(edit, preview))
	// show window and run app
}
func (app *config) makeUI() (*widget.Entry, *widget.RichText) {
	// 允许多行
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("")
	app.EditWidget = edit
	app.PreviewWidget = preview
	edit.OnChanged = preview.ParseMarkdown
	return edit, preview
}
