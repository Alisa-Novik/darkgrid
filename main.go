package main

import (
	"log"
	"math"
	"runtime"
	"strings"
	"time"

	"adagrad/internal/game"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	worldWidth  = 32
	worldHeight = 20
)

var (
	worldUp                = mgl32.Vec3{0, 1, 0}
	yaw            float32 = 45
	pitch          float32 = 55
	dist           float32 = 30
	center                 = mgl32.Vec3{float32(worldWidth) * 0.5, 0, float32(worldHeight) * 0.5}
	cameraPos              = mgl32.Vec3{0, 0, 0}
	hoverX, hoverZ         = -1, -1
)

func main() {
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	win, err := glfw.CreateWindow(1280, 800, "Sentient Dungeon RTS", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	win.MakeContextCurrent()
	world := game.NewGame(worldWidth, worldHeight)
	win.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	win.SetScrollCallback(onScroll)
	win.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft && action == glfw.Press {
			world.SelectTile(hoverX, hoverZ)
		}
	})
	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(0.08, 0.08, 0.1, 1)

	prog, err := newProgram(vertexSrc, fragmentSrc)
	if err != nil {
		log.Fatal(err)
	}
	lineProg, err := newProgram(lineVertSrc, lineFragSrc)
	if err != nil {
		log.Fatal(err)
	}

	planeVAO, planeEBO, planeCount := makeIndexedMesh(planeVertices(), planeIndices())
	cubeVAO, cubeEBO, cubeCount := makeIndexedMesh(cubeVertices(), cubeIndices())
	gridVAO, gridCount := makeGridLines(world.Width(), world.Height())

	mvpLoc := gl.GetUniformLocation(prog, gl.Str("mvp\x00"))
	colorLoc := gl.GetUniformLocation(prog, gl.Str("color\x00"))
	vpLoc := gl.GetUniformLocation(lineProg, gl.Str("vp\x00"))
	lineColorLoc := gl.GetUniformLocation(lineProg, gl.Str("lineColor\x00"))

	prev := time.Now()
	for !win.ShouldClose() {
		now := time.Now()
		dt := float32(now.Sub(prev).Seconds())
		prev = now

		tiles := world.Tiles()
		selX, selZ := world.SelectedTile()

		fbw, fbh := win.GetFramebufferSize()
		gl.Viewport(0, 0, int32(fbw), int32(fbh))
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		updateRTSCamera(win, dt)
		cameraPos = cameraFromOrbit(center, yaw, pitch, dist)
		view := mgl32.LookAtV(cameraPos, center, worldUp)
		proj := mgl32.Perspective(mgl32.DegToRad(60), float32(fbw)/float32(fbh), 0.1, 500)
		vp := proj.Mul4(view)

		mx, my := win.GetCursorPos()
		orig, dir := ScreenToWorldRay(mx, my, fbw, fbh, view, proj)
		hoverX, hoverZ = -1, -1
		if p, ok := RayHitY0(orig, dir); ok {
			hx := int(math.Floor(float64(p.X())))
			hz := int(math.Floor(float64(p.Z())))
			if world.InBounds(hx, hz) {
				hoverX, hoverZ = hx, hz
			}
		}

		gl.UseProgram(prog)
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		gl.Uniform3f(colorLoc, 0.20, 0.22, 0.26)
		gl.BindVertexArray(planeVAO)
		for z := 0; z < world.Height(); z++ {
			for x := 0; x < world.Width(); x++ {
				model := mgl32.Translate3D(float32(x), 0, float32(z))
				mvp := vp.Mul4(model)
				gl.UniformMatrix4fv(mvpLoc, 1, false, &mvp[0])
				gl.DrawElements(gl.TRIANGLES, planeCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
			}
		}

		if hoverX >= 0 && hoverZ >= 0 {
			gl.Uniform3f(colorLoc, 1.0, 1.0, 0.0)
			model := mgl32.Translate3D(float32(hoverX), 0.01, float32(hoverZ))
			mvp := vp.Mul4(model)
			gl.UniformMatrix4fv(mvpLoc, 1, false, &mvp[0])
			gl.BindVertexArray(planeVAO)
			gl.DrawElements(gl.TRIANGLES, planeCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
		}
		if selX >= 0 && selZ >= 0 {
			gl.Uniform3f(colorLoc, 1.0, 0.6, 0.1)
			model := mgl32.Translate3D(float32(selX), 0.02, float32(selZ))
			mvp := vp.Mul4(model)
			gl.UniformMatrix4fv(mvpLoc, 1, false, &mvp[0])
			gl.BindVertexArray(planeVAO)
			gl.DrawElements(gl.TRIANGLES, planeCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
		}

		gl.Uniform3f(colorLoc, 0.75, 0.75, 0.78)
		gl.BindVertexArray(cubeVAO)
		for z := 0; z < world.Height(); z++ {
			for x := 0; x < world.Width(); x++ {
				if tiles[z][x] == 1 {
					if x == hoverX && z == hoverZ {
						gl.Uniform3f(colorLoc, 1.0, 0.8, 0.2)
					} else if x == selX && z == selZ {
						gl.Uniform3f(colorLoc, 1.0, 0.6, 0.1)
					} else {
						gl.Uniform3f(colorLoc, 0.75, 0.75, 0.78)
					}
					model := mgl32.Translate3D(float32(x), 0, float32(z))
					mvp := vp.Mul4(model)
					gl.UniformMatrix4fv(mvpLoc, 1, false, &mvp[0])
					gl.DrawElements(gl.TRIANGLES, cubeCount, gl.UNSIGNED_INT, gl.PtrOffset(0))
				}
			}
		}

		gl.UseProgram(lineProg)
		gl.UniformMatrix4fv(vpLoc, 1, false, &vp[0])
		gl.Uniform3f(lineColorLoc, 0.1, 0.1, 0.1)
		gl.BindVertexArray(gridVAO)
		gl.DrawArrays(gl.LINES, 0, gridCount)

		win.SwapBuffers()
		glfw.PollEvents()
	}

	gl.DeleteVertexArrays(1, &planeVAO)
	gl.DeleteBuffers(1, &planeEBO)
	gl.DeleteVertexArrays(1, &cubeVAO)
	gl.DeleteBuffers(1, &cubeEBO)
}

func onScroll(w *glfw.Window, xoff, yoff float64) {
	dist -= float32(yoff) * 2
	if dist < 5 {
		dist = 5
	}
	if dist > 150 {
		dist = 150
	}
}

func updateRTSCamera(win *glfw.Window, dt float32) {
	ry := mgl32.DegToRad(yaw)
	forward := mgl32.Vec3{float32(math.Cos(float64(ry))), 0, float32(math.Sin(float64(ry)))}
	right := mgl32.Vec3{-forward.Z(), 0, forward.X()}
	speed := float32(8)
	if win.GetKey(glfw.KeyLeftShift) == glfw.Press {
		speed *= 2
	}
	move := speed * dt * (dist / 30)
	if win.GetKey(glfw.KeyW) == glfw.Press || win.GetKey(glfw.KeyUp) == glfw.Press {
		center = center.Sub(forward.Mul(move))
	}
	if win.GetKey(glfw.KeyS) == glfw.Press || win.GetKey(glfw.KeyDown) == glfw.Press {
		center = center.Add(forward.Mul(move))
	}
	if win.GetKey(glfw.KeyA) == glfw.Press || win.GetKey(glfw.KeyLeft) == glfw.Press {
		center = center.Sub(right.Mul(move))
	}
	if win.GetKey(glfw.KeyD) == glfw.Press || win.GetKey(glfw.KeyRight) == glfw.Press {
		center = center.Add(right.Mul(move))
	}
	rotSpeed := float32(60) * dt
	if win.GetKey(glfw.KeyQ) == glfw.Press {
		yaw -= rotSpeed
	}
	if win.GetKey(glfw.KeyE) == glfw.Press {
		yaw += rotSpeed
	}
	for yaw < 0 {
		yaw += 360
	}
	for yaw >= 360 {
		yaw -= 360
	}
}

func cameraFromOrbit(c mgl32.Vec3, yawDeg, pitchDeg, d float32) mgl32.Vec3 {
	ry := mgl32.DegToRad(yawDeg)
	rp := mgl32.DegToRad(pitchDeg)
	cp := float32(math.Cos(float64(rp)))
	sp := float32(math.Sin(float64(rp)))
	cy := float32(math.Cos(float64(ry)))
	sy := float32(math.Sin(float64(ry)))
	x := c.X() + d*cp*cy
	y := c.Y() + d*sp
	z := c.Z() + d*cp*sy
	return mgl32.Vec3{x, y, z}
}

func RayHitY0(orig, dir mgl32.Vec3) (mgl32.Vec3, bool) {
	if math.Abs(float64(dir.Y())) < 1e-6 {
		return mgl32.Vec3{}, false
	}
	t := -orig.Y() / dir.Y()
	if t <= 0 {
		return mgl32.Vec3{}, false
	}
	return orig.Add(dir.Mul(t)), true
}

func ScreenToWorldRay(mouseX, mouseY float64, winW, winH int, view, proj mgl32.Mat4) (mgl32.Vec3, mgl32.Vec3) {
	x := (2*float32(mouseX))/float32(winW) - 1
	y := 1 - (2*float32(mouseY))/float32(winH)
	rayClip := mgl32.Vec4{x, y, -1, 1}
	invProj := proj.Inv()
	rayEye := invProj.Mul4x1(rayClip)
	rayEye = mgl32.Vec4{rayEye.X(), rayEye.Y(), -1, 0}
	invView := view.Inv()
	rayWorld := invView.Mul4x1(rayEye).Vec3().Normalize()
	return cameraPos, rayWorld
}

func makeIndexedMesh(verts []float32, indices []uint32) (uint32, uint32, int32) {
	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*4, gl.Ptr(verts), gl.STATIC_DRAW)
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.BindVertexArray(0)
	return vao, ebo, int32(len(indices))
}

func makeGridLines(w, h int) (uint32, int32) {
	verts := make([]float32, 0, (w+h+2)*2*3)
	y := float32(0.01)
	for i := 0; i <= w; i++ {
		verts = append(verts, float32(i), y, 0, float32(i), y, float32(h))
	}
	for j := 0; j <= h; j++ {
		verts = append(verts, 0, y, float32(j), float32(w), y, float32(j))
	}
	var vao, vbo uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*4, gl.Ptr(verts), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.BindVertexArray(0)
	return vao, int32(len(verts) / 3)
}

func planeVertices() []float32 {
	return []float32{
		0, 0, 0,
		1, 0, 0,
		1, 0, 1,
		0, 0, 1,
	}
}

func planeIndices() []uint32 {
	return []uint32{0, 1, 2, 2, 3, 0}
}

func cubeVertices() []float32 {
	return []float32{
		0, 0, 0,
		1, 0, 0,
		1, 1, 0,
		0, 1, 0,
		0, 0, 1,
		1, 0, 1,
		1, 1, 1,
		0, 1, 1,
	}
}

func cubeIndices() []uint32 {
	return []uint32{
		4, 5, 6, 6, 7, 4,
		1, 0, 3, 3, 2, 1,
		0, 4, 7, 7, 3, 0,
		1, 2, 6, 6, 5, 1,
		3, 7, 6, 6, 2, 3,
		0, 5, 1, 0, 4, 5,
	}
}

func newProgram(vs, fs string) (uint32, error) {
	vert, err := compileShader(vs, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	frag, err := compileShader(fs, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vert)
	gl.AttachShader(prog, frag)
	gl.LinkProgram(prog)
	var status int32
	gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var l int32
		gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &l)
		logStr := strings.Repeat("\x00", int(l+1))
		gl.GetProgramInfoLog(prog, l, nil, gl.Str(logStr))
		return 0, err
	}
	gl.DeleteShader(vert)
	gl.DeleteShader(frag)
	return prog, nil
}

func compileShader(src string, t uint32) (uint32, error) {
	sh := gl.CreateShader(t)
	csources, free := gl.Strs(src + "\x00")
	gl.ShaderSource(sh, 1, csources, nil)
	free()
	gl.CompileShader(sh)
	var status int32
	gl.GetShaderiv(sh, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var l int32
		gl.GetShaderiv(sh, gl.INFO_LOG_LENGTH, &l)
		logStr := strings.Repeat("\x00", int(l+1))
		gl.GetShaderInfoLog(sh, l, nil, gl.Str(logStr))
		return 0, errShader(logStr)
	}
	return sh, nil
}

type errShader string

func (e errShader) Error() string { return string(e) }

var vertexSrc = `
#version 330 core
layout (location = 0) in vec3 vert;
uniform mat4 mvp;
void main() { gl_Position = mvp * vec4(vert, 1.0); }
`

var fragmentSrc = `
#version 330 core
out vec4 FragColor;
uniform vec3 color;
void main() { FragColor = vec4(color, 1.0); }
`

var lineVertSrc = `
#version 330 core
layout (location = 0) in vec3 vert;
uniform mat4 vp;
void main() { gl_Position = vp * vec4(vert, 1.0); }
`

var lineFragSrc = `
#version 330 core
out vec4 FragColor;
uniform vec3 lineColor;
void main() { FragColor = vec4(lineColor, 1.0); }
`
