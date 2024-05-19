package cmd

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/leap-fish/clay"
	"github.com/leap-fish/clay/bundle"
	"github.com/leap-fish/clay/components/spatial"
	txt "github.com/leap-fish/clay/components/text"
	"github.com/leap-fish/clay/resource"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/debug"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/filter"
	"image/color"
	"time"
)

var DebugMarker = donburi.NewTag()

type ExampleRenderer struct {
	font     *text.GoTextFaceSource
	fontFace *text.GoTextFace
}

func (e *ExampleRenderer) Update(w donburi.World, dt time.Duration) {
	q := donburi.NewQuery(filter.Contains(txt.Component, spatial.TransformComponent))
	entry, exists := q.First(w)
	if !exists || entry == nil {
		return
	}

	t := txt.Component.Get(entry)
	t.Content.Reset()
	for _, c := range debug.GetEntityCounts(w) {
		t.Content.WriteString(fmt.Sprintf("> %s\n", c.String()))
	}
}

func (e *ExampleRenderer) Init(w donburi.World) {
	e.font = resource.Get[*text.GoTextFaceSource]("font:BaiJamjuree-Regular")
	e.fontFace = &text.GoTextFace{
		Source: e.font,
		Size:   14,
	}
	bigFontFace := &text.GoTextFace{
		Source: e.font,
		Size:   22,
	}

	DebugMarker.SetName("DebugMarker")

	textBundle := bundle.New().
		With(spatial.TransformComponent, spatial.DefaultTransform).
		With(txt.Component, txt.Text{
			FontFace:   e.fontFace,
			Content:    bytes.Buffer{},
			Size:       16,
			LineHeight: 1.0,
			Color:      color.RGBA{255, 255, 255, 255},
		}).
		With(DebugMarker, struct{}{})

	textBundle.Spawn(w)

	secondTextBuf := bytes.Buffer{}
	secondTextBuf.WriteString("Clay Engine in golang")
	secondText := bundle.New().
		With(spatial.TransformComponent, spatial.Transform{Position: math.NewVec2(200, 200), Scale: 1.0}).
		With(txt.Component, txt.Text{
			FontFace:   bigFontFace,
			Content:    secondTextBuf,
			Size:       16,
			LineHeight: 1.0,
			Color:      color.RGBA{255, 100, 100, 255},
		}).
		With(DebugMarker, struct{}{})

	secondText.Spawn(w)
}

func (e *ExampleRenderer) Render(rg *clay.RenderGraph, w donburi.World) {}
