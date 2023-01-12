package widgets

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewTable(rows, cols int) *widget.Table {
	table := widget.NewTable(
		//define the dimensions of the table
		func() (int, int) { return rows, cols },
		//specifies the widget
		func() fyne.CanvasObject { return widget.NewLabel("...") },
		//update the table
		func(tci widget.TableCellID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(fmt.Sprintf("%d %d", tci.Col, tci.Row))
		},
	)

	return table
}
