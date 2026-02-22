package render

import (
	"swarmlings/config"
	"swarmlings/sim"
	"bytes"
	"fmt"
	"image/color"
	"math"
	"strconv"
	"sync"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"
)

var (
	fontOnce   sync.Once
	regularSrc *text.GoTextFaceSource
	monoSrc    *text.GoTextFaceSource
)

func initFonts() {
	fontOnce.Do(func() {
		var err error
		regularSrc, err = text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
		if err != nil {
			panic(err)
		}
		monoSrc, err = text.NewGoTextFaceSource(bytes.NewReader(gomono.TTF))
		if err != nil {
			panic(err)
		}
	})
}

func BuildUI(world *sim.World, cfg *config.Config, scale float64) ebitenui.UI {
	initFonts()

	s := func(base int) int { return int(math.Round(float64(base) * scale)) }
	sf := func(base float64) float64 { return base * scale }

	// Each face must be its own variable (pointer semantics)
	var labelFace text.Face = &text.GoTextFace{Source: regularSrc, Size: sf(13)}
	var headerFace text.Face = &text.GoTextFace{Source: regularSrc, Size: sf(14)}
	var titleFace text.Face = &text.GoTextFace{Source: regularSrc, Size: sf(15)}
	var valueFace text.Face = &text.GoTextFace{Source: monoSrc, Size: sf(13)}

	panelBG := color.NRGBA{20, 20, 24, 220}
	panelBorder := color.NRGBA{60, 60, 70, 255}
	separatorColor := color.NRGBA{50, 50, 60, 255}
	textPrimary := color.NRGBA{220, 220, 230, 255}
	textDim := color.NRGBA{160, 160, 175, 255}
	headerColor := color.NRGBA{200, 200, 215, 255}
	accentTeal := color.NRGBA{0, 180, 200, 255}
	accentHover := color.NRGBA{0, 210, 230, 255}
	accentPressed := color.NRGBA{0, 150, 170, 255}
	trackColor := color.NRGBA{45, 45, 55, 255}
	inputBG := color.NRGBA{12, 12, 16, 255}
	inputBorder := color.NRGBA{50, 50, 60, 255}
	sliderTrack := &widget.SliderTrackImage{
		Idle:  image.NewNineSliceColor(trackColor),
		Hover: image.NewNineSliceColor(color.NRGBA{55, 55, 65, 255}),
	}
	sliderHandle := &widget.ButtonImage{
		Idle:    image.NewNineSliceColor(accentTeal),
		Hover:   image.NewNineSliceColor(accentHover),
		Pressed: image.NewNineSliceColor(accentPressed),
	}
	inputImage := &widget.TextInputImage{
		Idle: image.NewBorderedNineSliceColor(inputBG, inputBorder, 1),
	}
	inputColor := &widget.TextInputColor{
		Idle:  textPrimary,
		Caret: accentTeal,
	}

	makeSeparator := func() *widget.Container {
		return widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			)),
			widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(separatorColor)),
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.MinSize(0, 1),
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true}),
			),
		)
	}

	makeHeader := func(title string) *widget.Container {
		c := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
				widget.RowLayoutOpts.Padding(&widget.Insets{Top: s(2), Bottom: s(2)}),
			)),
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true}),
			),
		)
		c.AddChild(widget.NewText(widget.TextOpts.Text(title, &headerFace, headerColor)))
		return c
	}

	makeRow := func(label string, minVal, maxVal float64, cfgField *float64, formatStr string, apply func(float64)) *widget.Container {
		var input *widget.TextInput
		var slider *widget.Slider
		updating := false

		valToSlider := func(v float64) int {
			return int((v - minVal) / (maxVal - minVal) * 100)
		}
		sliderToVal := func(s int) float64 {
			return minVal + (maxVal-minVal)*float64(s)/100.0
		}

		input = widget.NewTextInput(
			widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.MinSize(s(65), s(24))),
			widget.TextInputOpts.Image(inputImage),
			widget.TextInputOpts.Face(&valueFace),
			widget.TextInputOpts.Color(inputColor),
			widget.TextInputOpts.Padding(&widget.Insets{Top: s(3), Left: s(4), Right: s(4), Bottom: s(3)}),
			widget.TextInputOpts.Placeholder("0"),
			widget.TextInputOpts.SubmitOnEnter(true),
			widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
				if updating {
					return
				}
				updating = true
				val, err := strconv.ParseFloat(input.GetText(), 64)
				if err == nil {
					if val < minVal {
						val = minVal
					} else if val > maxVal {
						val = maxVal
					}
					apply(val)
					*cfgField = val
					slider.Current = valToSlider(val)
					input.SetText(fmt.Sprintf(formatStr, val))
				}
				updating = false
			}),
		)
		input.SetText(fmt.Sprintf(formatStr, *cfgField))

		slider = widget.NewSlider(
			widget.SliderOpts.Orientation(widget.DirectionHorizontal),
			widget.SliderOpts.MinMax(0, 100),
			widget.SliderOpts.InitialCurrent(valToSlider(*cfgField)),
			widget.SliderOpts.WidgetOpts(widget.WidgetOpts.MinSize(s(160), s(20))),
			widget.SliderOpts.Images(sliderTrack, sliderHandle),
			widget.SliderOpts.FixedHandleSize(s(14)),
			widget.SliderOpts.TrackOffset(0),
			widget.SliderOpts.PageSizeFunc(func() int { return 1 }),
			widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
				if updating {
					return
				}
				updating = true
				val := sliderToVal(args.Current)
				apply(val)
				*cfgField = val
				input.SetText(fmt.Sprintf(formatStr, val))
				updating = false
			}),
		)

		labelContainer := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			)),
			widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(s(80), s(24))),
		)
		labelContainer.AddChild(widget.NewText(widget.TextOpts.Text(label, &labelFace, textDim)))

		row := widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
				widget.RowLayoutOpts.Spacing(s(8)),
			)),
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{Stretch: true}),
			),
		)
		row.AddChild(labelContainer)
		row.AddChild(slider)
		row.AddChild(input)
		return row
	}

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	pad := s(10)
	panel := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewBorderedNineSliceColor(panelBG, panelBorder, 1)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(&widget.Insets{Top: pad, Left: pad, Right: pad, Bottom: pad}),
			widget.RowLayoutOpts.Spacing(s(8)),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(s(280), 0),
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				StretchVertical:    true,
			}),
		),
	)

	panel.AddChild(widget.NewText(widget.TextOpts.Text("Ling Parameters", &titleFace, textPrimary)))
	panel.AddChild(makeSeparator())

	panel.AddChild(makeHeader("Forces"))
	panel.AddChild(makeRow("Avoid", 0, 2.0, &cfg.AvoidanceFactor, "%.3f", func(v float64) {
		world.AvoidanceFactor = v
	}))
	panel.AddChild(makeRow("Align", 0, 0.01, &cfg.AlignmentFactor, "%.4f", func(v float64) {
		world.AlignmentFactor = v
	}))
	panel.AddChild(makeRow("Gather", 0, 0.002, &cfg.GatheringFactor, "%.4f", func(v float64) {
		world.GatheringFactor = v
	}))

	panel.AddChild(makeSeparator())

	panel.AddChild(makeHeader("Radii"))
	panel.AddChild(makeRow("Avoid R", 0, 100, &cfg.AvoidanceRadius, "%.0f", func(v float64) {
		world.AvoidanceRadius = v
	}))
	panel.AddChild(makeRow("Detect R", 0, 300, &cfg.DetectionRadius, "%.0f", func(v float64) {
		world.DetectionRadius = v
	}))

	panel.AddChild(makeSeparator())

	panel.AddChild(makeHeader("Speed"))
	panel.AddChild(makeRow("Max", 0, 10, &cfg.MaxSpeed, "%.1f", func(v float64) {
		world.MaxSpeed = v
	}))

	panel.AddChild(makeSeparator())

	btnContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch:  true,
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)
	saveBtn := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(s(120), s(28))),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle:    image.NewNineSliceColor(accentTeal),
			Hover:   image.NewNineSliceColor(accentHover),
			Pressed: image.NewNineSliceColor(accentPressed),
		}),
		widget.ButtonOpts.Text("Save", &labelFace, &widget.ButtonTextColor{
			Idle: color.NRGBA{0, 0, 0, 255},
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			config.Save(*cfg)
		}),
	)
	btnContainer.AddChild(saveBtn)
	panel.AddChild(btnContainer)

	rootContainer.AddChild(panel)

	return ebitenui.UI{
		Container: rootContainer,
	}
}
