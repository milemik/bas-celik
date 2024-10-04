package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Group struct {
	widget.BaseWidget
	name    string
	objects []fyne.CanvasObject
}

type GroupRenderer struct {
	group          *Group
	nameText       *canvas.Text
	backgroundRect *canvas.Rectangle
	column         *fyne.Container
}

func NewGroup(name string, objects ...fyne.CanvasObject) *Group {
	group := &Group{
		name:    name,
		objects: objects,
	}
	group.ExtendBaseWidget(group)
	return group
}

func (g *Group) CreateRenderer() fyne.WidgetRenderer {
	nameText := canvas.NewText(g.name, color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0xFF})
	nameText.TextSize = 11

	nameText.Move(fyne.NewPos(theme.Padding(), -1.6*theme.Padding()))

	column := container.New(layout.NewVBoxLayout(), g.objects...)
	column.Move(fyne.NewPos(theme.Padding(), 0))

	backgroundRect := canvas.NewRectangle(theme.OverlayBackgroundColor())
	backgroundRect.StrokeWidth = 1
	backgroundRect.StrokeColor = theme.DisabledButtonColor()
	backgroundRect.Move(fyne.NewPos(0, 2*theme.Padding()))
	backgroundRect.Resize(fyne.NewSize(column.MinSize().Width+2*theme.Padding(), column.MinSize().Height+theme.Padding()))
	backgroundRect.CornerRadius = 3

	return &GroupRenderer{
		group:          g,
		nameText:       nameText,
		column:         column,
		backgroundRect: backgroundRect,
	}
}

func (r *GroupRenderer) Refresh() {}

func (r *GroupRenderer) Layout(s fyne.Size) {
	r.column.Move(fyne.Position{X: theme.Padding(), Y: 2 * theme.Padding()})
	bgSize := r.backgroundRect.Size()
	r.backgroundRect.Resize(fyne.NewSize(s.Width, bgSize.Height))
}

func (r *GroupRenderer) MinSize() fyne.Size {
	return fyne.NewSize(r.column.MinSize().Width+2*theme.Padding(), r.column.MinSize().Height+6*theme.Padding())
}

func (r *GroupRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.backgroundRect, r.column, r.nameText}
}

func (r *GroupRenderer) Destroy() {}
