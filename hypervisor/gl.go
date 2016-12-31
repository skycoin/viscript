package hypervisor

import (
	"fmt"
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/gfx"
	"github.com/go-gl/gl/v2.1/gl"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"os"
)

var (
	Texture   uint32
	rotationX float32
	rotationY float32
)

func setFrustum(r *app.Rectangle) {
	gl.Frustum(
		float64(r.Left),
		float64(r.Right),
		float64(r.Bottom),
		float64(r.Top), 1.0, 10.0)
}

func DrawScene() {
	//rotationX += 0.5
	//rotationY += 0.5
	gl.Viewport(0, 0, gfx.CurrAppWidth, gfx.CurrAppHeight) // OPTIMIZEME?  could set flag upon frame buffer size change event
	if *gfx.PrevFrustum != *gfx.CurrFrustum {
		*gfx.PrevFrustum = *gfx.CurrFrustum
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		setFrustum(gfx.CurrFrustum)
		fmt.Println("CHANGE OF FRUSTUM")
		fmt.Printf(".Panels[0].Rect.Right: %.2f\n", gfx.Rend.Panels[0].Rect.Right)
		fmt.Printf(".Panels[0].Rect.Top: %.2f\n", gfx.Rend.Panels[0].Rect.Top)
	}
	gl.MatrixMode(gl.MODELVIEW) //.PROJECTION)                   //.MODELVIEW)
	gl.LoadIdentity()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Translatef(0, 0, -gfx.Rend.DistanceFromOrigin)
	//gl.Rotatef(rotationX, 1, 0, 0)
	//gl.Rotatef(rotationY, 0, 1, 0)

	gl.BindTexture(gl.TEXTURE_2D, Texture)

	gl.Begin(gl.QUADS)
	gfx.Rend.DrawAll()
	gl.End()
}

func NewTexture(file string) uint32 {
	imgFile, err := os.Open(file)
	if err != nil {
		log.Fatalf("texture %q not found on disk: %v\n", file, err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic("unsupported stride")
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture
}

func destroyScene() {
}
