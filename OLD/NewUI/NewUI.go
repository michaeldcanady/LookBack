//MAJOR WORK IN PROGRESS, ABLE TO USE GUI INSTEAD OF CMD FOR USER INTERACTION

package main

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
  "fyne.io/fyne"
)

type diagonal struct {
}

func (d *diagonal) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := 0, 0
	for _, o := range objects {
		childSize := o.MinSize()

		w += childSize.Width
		h += childSize.Height
	}
	return fyne.NewSize(w, h)
}

func (d *diagonal) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, containerSize.Height - d.MinSize(objects).Height)
	for _, o := range objects {
		size := o.MinSize()
		o.Resize(size)
		o.Move(pos)

		pos = pos.Add(fyne.NewPos(size.Width, size.Height))
	}
}

func main() {
	app := app.New()

	w := app.NewWindow("Look Back")
  text1 := widget.NewLabel("Username")
  text2 := widget.NewEntry()
	w.SetContent(fyne.NewContainerWithLayout(&diagonal{}, text1, text2))

	w.ShowAndRun()
}
