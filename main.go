package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont"
	"github.com/golang/freetype/truetype"
)

type config struct {
	EditWidget    *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile   fyne.URI
	SaveMenuItem  *fyne.MenuItem
}

var cfg config

func init() {
	//设置中文字体
	os.Setenv("FYNE_FONT", "Alibaba-PuHuiTi-Medium.ttf")
}
func init() {
	fontPath, err := findfont.Find("SIMYOU.TTF")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found 'arial.ttf' in '%s'\n", fontPath)

	// load the font with the freetype library
	fontData, err := os.ReadFile(fontPath)
	if err != nil {
		panic(err)
	}
	_, err = truetype.Parse(fontData)
	if err != nil {
		panic(err)
	}
	os.Setenv("FYNE_FONT", fontPath)
}
func main() {
	// create a fyne app
	a := app.New()
	// create a window for the app
	win := a.NewWindow("Markdown")
	// get the user interface
	edit, preview := cfg.makeUI()
	cfg.createMenuItems(win)
	// set the content of the window
	win.SetContent(container.NewHSplit(edit, preview))
	// show window and run app
	win.Resize(fyne.NewSize(800, 500))
	win.CenterOnScreen()
	win.ShowAndRun()
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
func (app *config) createMenuItems(win fyne.Window) {
	openMenuItem := fyne.NewMenuItem("Open...", app.openFunc(win))
	saveMenuItem := fyne.NewMenuItem("Save", app.saveFunc(win))
	app.SaveMenuItem = saveMenuItem
	// 禁用save
	app.SaveMenuItem.Disabled = true
	saveAsMenuItem := fyne.NewMenuItem("Save As...", app.saveAsFunc(win))
	// 按顺序展示
	fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)
	// 主菜单
	menu := fyne.NewMainMenu(fileMenu)
	win.SetMainMenu(menu)
}
// 过滤器
var filter = storage.NewExtensionFileFilter([]string{".md", ".MD"})
func (app *config) saveFunc(win fyne.Window) func(){
	return func(){
		if app.CurrentFile != nil{
			write, err := storage.Writer(app.CurrentFile)
			if err != nil{
				dialog.ShowError(err,win)
			}
			write.Write([]byte(app.EditWidget.Text))
			defer write.Close()
		}
	}
}
func (app *config) openFunc(win fyne.Window) func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			// 出错
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			// 文件为空
			if read == nil {
				return
			}
			defer read.Close()
			// 读取数据
			data, err := io.ReadAll(read)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			app.EditWidget.SetText(string(data))

			app.CurrentFile = read.URI()
			win.SetTitle(win.Title() + "-" + read.URI().Name())
			// 确保保存组件可用
			app.SaveMenuItem.Disabled = false
		}, win)
		// 只能打开md文件
		openDialog.SetFilter(filter)
		openDialog.Show()
		//
	}
}
func (app *config) saveAsFunc(win fyne.Window) func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if write == nil {
				// user cancelled
				return
			}
			// 把文件名限制为小写.md
			if !strings.HasSuffix(strings.ToLower(write.URI().String()), ".md"){
				dialog.ShowInformation("Error","Please name your file with a .md extension!",win)
				return
			}

			// save file
			write.Write([]byte(app.EditWidget.Text))
			app.CurrentFile = write.URI()

			defer write.Close()

			win.SetTitle(win.Title() + "-" + write.URI().Name())
			app.SaveMenuItem.Disabled = false
		}, win)
		// 文件名
		saveDialog.SetFileName("untitled.md")
		saveDialog.SetFilter(filter)
		saveDialog.Show()
	}
}
