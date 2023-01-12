package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewPrompt(title, message string) {
	prompt := fyne.CurrentApp().NewWindow(title)
	msg := widget.NewLabel(message)
	prompt.SetContent(msg)
	prompt.Show()
}
