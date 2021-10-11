package main

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

const (
	wMaxInvFPS = 1 / 30.0
)

type Game struct {
	ids     int
	idLock  sync.Mutex
	images  *glutil.Images
	fps     *debug.FPS
	glc     gl.Context
	lastTS  time.Time
	frameDt float64
	elapsed float64
	uEffect gl.Uniform
	uTime   gl.Uniform
	uModel  gl.Uniform
	uView   gl.Uniform
	uProj   gl.Uniform
	touchX  float32
	touchY  float32
	lastX   int
	lastY   int
	size    size.Event
	objects map[int]Object
	program gl.Program
	projf   []float32
	viewf   []float32
}

func (g *Game) Init(glctx gl.Context) {
	g.glc = glctx

	var err error
	g.program, err = glutil.CreateProgram(g.glc, vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}

	g.UpdateView()
	g.glc.Viewport(0, 0, g.size.WidthPx, g.size.HeightPx)

	g.uModel = g.glc.GetUniformLocation(g.program, "model")
	g.uView = g.glc.GetUniformLocation(g.program, "view")
	g.uProj = g.glc.GetUniformLocation(g.program, "projection")
	g.uEffect = g.glc.GetUniformLocation(g.program, "effect")
	g.uTime = g.glc.GetUniformLocation(g.program, "uTime")

	rand.Seed(time.Now().Unix())

	g.objects = make(map[int]Object)

	g.glc.CullFace(gl.BACK)
	g.glc.Enable(gl.DEPTH_TEST)
	g.glc.DepthFunc(gl.LESS)

	s2 := &Sprite{}
	s2.Init(0, 320, 0, 1.0, "bg3.png", g)
	g.AddObjects(s2)

	s1 := &Sprite{}
	s1.Init(210, 310, 1, 0.8, "logo.png", g)
	g.AddObjects(s1)

	g.images = glutil.NewImages(g.glc)
	g.fps = debug.NewFPS(g.images)
	g.lastTS = time.Now()
}

func (g *Game) Stop() {
	g.glc.DeleteProgram(g.program)
	g.fps.Release()
	g.images.Release()
}

func (g *Game) Draw() {
	dt := time.Since(g.lastTS).Seconds()
	g.frameDt += dt
	g.lastTS = time.Now()
	g.glc.Enable(gl.BLEND)
	g.glc.Enable(gl.DEPTH_TEST)
	g.glc.DepthFunc(gl.LESS)
	g.glc.BlendFunc(gl.BLEND_SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	g.glc.ClearColor(0.0, 0.0, 0.0, 1.0)
	g.glc.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for {
		if g.frameDt >= wMaxInvFPS {
			g.elapsed += wMaxInvFPS
			for k := range g.objects {
				g.objects[k].Update(float64(wMaxInvFPS))
			}
		} else {
			break
		}

		g.frameDt -= wMaxInvFPS
	}

	for k := range g.objects {
		g.objects[k].Draw(float64(wMaxInvFPS))
	}

	g.glc.Disable(gl.DEPTH_TEST)
	g.fps.Draw(g.size)
}

func (g *Game) Click(x, y float32) {
	return

	//  wh := float32(320)
	//  ww := float32(320)
	//  aspect := (float32(g.size.HeightPx) / wh) / 2
	//  touchX := x / (float32(g.size.WidthPx) / float32(ww))

	//  k := float32(g.size.HeightPx) - ((wh * aspect) / 2) // Start position
	//  p1 := float32(g.size.HeightPx) - k                  // Position low
	//  touchY := (p1-y)/aspect + wh                        // Actual touch

	// Make sure we don't generate too many clicks.
	if g.lastX == int(x) && g.lastY == int(y) {
		return
	}
	g.lastX = int(x)
	g.lastY = int(y)

	// TBD
}

func (g *Game) Resize(e size.Event) {
	g.size = e
	g.UpdateProjection()
	g.UpdateView()
}

func (g *Game) UpdateProjection() {
	ww := float32(320)
	wh := float32(320)

	projection := mgl32.Ortho2D(0, ww, -(ww/g.size.PixelsPerPt)/2, wh+(ww/g.size.PixelsPerPt)/2)

	projf := []float32{}
	c1, c2, c3, c4 := projection.Cols()
	projf = append(projf, c1.X())
	projf = append(projf, c1.Y())
	projf = append(projf, c1.Z())
	projf = append(projf, c1.W())
	projf = append(projf, c2.X())
	projf = append(projf, c2.Y())
	projf = append(projf, c2.Z())
	projf = append(projf, c2.W())
	projf = append(projf, c3.X())
	projf = append(projf, c3.Y())
	projf = append(projf, c3.Z())
	projf = append(projf, c3.W())
	projf = append(projf, c4.X())
	projf = append(projf, c4.Y())
	projf = append(projf, c4.Z())
	projf = append(projf, c4.W())

	g.projf = projf
}

func (g *Game) UpdateView() {
	view := mgl32.Translate3D(float32(0), float32(0), float32(0.0))
	p := []float32{}
	c1, c2, c3, c4 := view.Cols()
	p = append(p, c1.X())
	p = append(p, c1.Y())
	p = append(p, c1.Z())
	p = append(p, c1.W())
	p = append(p, c2.X())
	p = append(p, c2.Y())
	p = append(p, c2.Z())
	p = append(p, c2.W())
	p = append(p, c3.X())
	p = append(p, c3.Y())
	p = append(p, c3.Z())
	p = append(p, c3.W())
	p = append(p, c4.X())
	p = append(p, c4.Y())
	p = append(p, c4.Z())
	p = append(p, c4.W())
	g.viewf = p
}

func (g *Game) AddObjects(obj ...Object) {
	for i := range obj {
		if _, ok := g.objects[obj[i].GetID()]; !ok {
			g.objects[obj[i].GetID()] = obj[i]
		}
	}
}

func (g *Game) DeleteObject(obj Object) {
	for i := range g.objects {
		if obj.GetID() == g.objects[i].GetID() {
			g.objects[i].Delete()
			delete(g.objects, i)
			return
		}
	}
}

func (g *Game) NewID() int {
	g.idLock.Lock()
	defer g.idLock.Unlock()

	g.ids++
	return g.ids
}
