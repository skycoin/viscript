package hypervisor

import (
	"fmt"
	"github.com/corpusc/viscript/app"
	"github.com/corpusc/viscript/gfx"
	//"log"
	"github.com/go-gl/gl/v2.1/gl"
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

func destroyScene() {
}
