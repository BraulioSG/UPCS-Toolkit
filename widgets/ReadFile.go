package widgets

import (
	"encoding/csv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func NewFileCSV(parent fyne.Window, callback func([][]string)) (inputFileButton fyne.CanvasObject) {
	inputFileButton = widget.NewButton("Select CSV file", func() {
		var file_dialog *dialog.FileDialog = dialog.NewFileOpen(func(r fyne.URIReadCloser, _ error) {
			//reads the csv
			data, error := csv.NewReader(r).ReadAll()

			//if an error happens send a prototype to the callback
			if error != nil {
				var prototype [][]string
				callback(prototype)
				return
			}

			callback(data)

		}, parent)

		//just allow .csv files
		file_dialog.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
		file_dialog.Show()
	})
	return
}
