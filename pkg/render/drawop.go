package render

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leap-fish/clay/pkg/components/camera"
)

type DrawOp struct {
	mode  Mode
	index int

	scale    float64
	rotation float64

	posX, posY       float64
	originX, originY float64

	skewX, skewY float64

	image *ebiten.Image

	ops    *ebiten.DrawImageOptions
	filter ebiten.Filter
}

// Draw returns a new DrawOp which can be used to customize how the image is rendered.
// By using Draw instead of manual ebiten drawing, you get automatic handling of rotations and sprite origins.
// DrawOp makes sure the draw operations are performed in the correct order.
// Call QueueRender() to draw onto the render pass.
func Draw(image *ebiten.Image, mode Mode, index int) *DrawOp {
	return &DrawOp{
		scale: 1.0,

		index: index,
		image: image,
		mode:  mode,

		filter: ebiten.FilterNearest,
		ops:    &ebiten.DrawImageOptions{},
	}
}

func (d *DrawOp) Mode(mode Mode) *DrawOp {
	d.mode = mode
	return d
}

// CenterOrigin sets the origin of the sprite to its center.
func (d *DrawOp) CenterOrigin() *DrawOp {
	bounds := d.image.Bounds().Size()
	d.originX = float64(bounds.X / 2)
	d.originY = float64(bounds.Y / 2)
	return d
}

// Skew sets the skew of the sprite.
func (d *DrawOp) Skew(skewX, skewY float64) *DrawOp {
	d.skewX = skewX
	d.skewY = skewY
	return d
}

// Origin sets the origin of the sprite.
func (d *DrawOp) Origin(originX, originY float64) *DrawOp {
	d.originX = originX
	d.originY = originY
	return d
}

// OriginMul sets the origin of the sprite.
func (d *DrawOp) OriginMul(originX, originY float64) *DrawOp {
	bounds := d.image.Bounds()
	d.originX = originX * float64(bounds.Dx())
	d.originY = originY * float64(bounds.Dy())
	return d
}

// Scale sets the scale of the sprite.
func (d *DrawOp) Scale(scale float64) *DrawOp {
	d.scale = scale
	return d
}

// Rotation sets the rotation of the sprite.
func (d *DrawOp) Rotation(rotation float64) *DrawOp {
	d.rotation = rotation
	return d
}

// Position sets the position of the sprite.
func (d *DrawOp) Position(posX, posY float64) *DrawOp {
	d.posX = posX
	d.posY = posY
	return d
}

// Filter sets the draw filter mode.
func (d *DrawOp) Filter(filter ebiten.Filter) *DrawOp {
	d.filter = filter
	return d
}

// Draw is used internally to perform the actual rendering.
// Called by the render loop.
func (d *DrawOp) Draw(surface *ebiten.Image, camera *camera.Camera) {
	d.ops.GeoM.Translate(-d.originX, -d.originY)
	// Non-essential operations are checked first
	if d.scale != 1 {
		d.ops.GeoM.Scale(d.scale, d.scale)
	}
	d.ops.GeoM.Rotate(float64(d.rotation))
	d.ops.GeoM.Translate(d.originX, d.originY)
	spritePosX, spritePosY := d.posX-d.originX, d.posY-d.originY
	d.ops.GeoM.Translate(spritePosX, spritePosY)

	if d.skewX != 0 && d.skewY != 0 {
		d.ops.GeoM.Skew(d.skewX, d.skewY)
	}

	if d.filter != ebiten.FilterNearest {
		d.ops.Filter = d.filter
	}

	if d.mode == ModeCanvas {
		surface.DrawImage(d.image, d.ops)
		return
	}

	if d.mode == ModeWorld {
		// If we are in non-canvas mode we have to modify the image with data from our camera.
		camera.WorldMatrix(d.ops)
		surface.DrawImage(d.image, d.ops)
		return
	}
}
