package main

import (
	"UPCSToolkit/services"
	"UPCSToolkit/widgets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func main() {
	application := app.New()                        //creates the aplication
	window := application.NewWindow("UPCS Toolkit") //creates the window and give a title
	window.Resize(fyne.NewSize(800, 500))           // sets the dimensions
	window.SetFixedSize(true)                       //not resizeable

	var file_button fyne.CanvasObject = widgets.NewFileCSV(window, func(data [][]string) {
		if len(data) == 0 {
			return
		}
		if len(data[0]) == 0 {
			return
		}

		services.CreateReportCalidadAcademica(data)
	})

	// file_button := widget.NewButton("Select a File", func() {

	// 	file_dialog := dialog.NewFileOpen(func(r fyne.URIReadCloser, _ error) {
	// 		//read files
	// 		data, _ := ioutil.ReadAll(r)
	// 		var dataStr string = string(data)
	// 		fmt.Print(dataStr)
	// 	}, window)

	// 	file_dialog.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
	// 	file_dialog.Show()
	// })

	content := container.NewVBox(
		file_button,
	)

	window.SetContent(content)
	window.ShowAndRun()
}
