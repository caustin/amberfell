package main

import (
    "github.com/banthar/Go-SDL/sdl"
    "github.com/banthar/gl"
    "github.com/banthar/glu"
    "math/rand"  
    "flag"
    "fmt"
    "github.com/iand/amberfell/af"
    "time"
    "math"
    
)    

const piover180 = 0.0174532925



var printInfo = flag.Bool("info", false, "print GL implementation information")

var T0 uint32 = 0
var Frames uint32 = 0


var view_rotx float64 = 50.0
var view_roty float64 = 50.0
var view_rotz float64 = 0.0
var gear1, gear2, gear3 uint
var angle float64 = 0.0



var (
    player *af.Player
    wolf *af.Wolf
    world af.World
    screenWidth, screenHeight int
    tileWidth = 48
    screenScale int = int(5 * float64(tileWidth) / 2)
    ShowOverlay bool

    timeOfDay float32 = 19

    lightpos af.Vector


)

func main() {
    flag.Parse()
    rand.Seed(71)   
    var done bool
    var keys []uint8
    player = new(af.Player)
    player.Init(0, 10, 10, af.GroundLevel+1)




    world.Init(20,20,30)
    
    sdl.Init(sdl.INIT_VIDEO)

    lightpos[af.XAXIS] = -0.5
    lightpos[af.YAXIS] = -0.5
    lightpos[af.ZAXIS] = -0.5


    var screen = sdl.SetVideoMode(800, 600, 32, sdl.OPENGL|sdl.RESIZABLE)

    if screen == nil {
        sdl.Quit()
        panic("Couldn't set GL video mode: " + sdl.GetError() + "\n")
    } 

    if gl.Init() != 0 {
        panic("gl error")   
    }

    sdl.WM_SetCaption("Amberfell", "amberfell")

    init2()
    reshape(int(screen.W), int(screen.H))

    var currentTime, accumulator int64 = 0, 0
    var t, dt int64 = 0, 1e9/40
    var drawFrame, computeFrame int64 = 0, 0
    fps := new(af.Timer)
    fps.Start()

    update := new(af.Timer)
    update.Start()

    done = false
    for !done {

        for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
            switch e.(type) {
            case *sdl.ResizeEvent:
                re := e.(*sdl.ResizeEvent)
                screen = sdl.SetVideoMode(int(re.W), int(re.H), 16,
                    sdl.OPENGL|sdl.RESIZABLE)
                if screen != nil {
                    reshape(int(screen.W), int(screen.H))
                } else {
                    panic("we couldn't set the new video mode??")
                }
                break

            case *sdl.MouseButtonEvent:
                re := e.(*sdl.MouseButtonEvent)
                if re.Button == 1 && re.State == 1 { // LEFT, DOWN
                    // println("Click:", re.X, re.Y, re.State, re.Button, re.Which)

                    // MOUSEBUTTONDOWNMASK
                    xv, yv := int(re.X), screenHeight - int(re.Y)
                    data := [4]uint8{0, 0, 0, 0}

                    Draw(true)
                    gl.ReadPixels(xv, yv, 1, 1, gl.RGBA, &data[0])
                    Draw(false)

                    fmt.Printf("pixel data: %d, %d, %d, %d\n", data[0], data[1], data[2], data[3])

                    id := uint16(data[0]) + uint16(data[1]) * 256
                    if id != 0 {
                        face := data[2]
                        dx, dy, dz := af.BlockIdToRelativeCoordinate(id)
                        fmt.Printf("id: %d, dx: %d, dy: %d, dz: %d, face: %d\n", id, dx, dy, dz, face)
                        if ! (dx == 0 && dy == 0 && dz == 0) {
                            pos := af.IntPosition(player.Position())
                            pos[af.XAXIS] += dx
                            pos[af.YAXIS] += dy
                            pos[af.ZAXIS] += dz
                            if face == af.TOP_FACE { // top
                                pos[af.YAXIS]++
                            } else if face == af.BOTTOM_FACE { // bottom
                                pos[af.YAXIS]--
                            } else if face == af.FRONT_FACE { // front
                                pos[af.ZAXIS]++
                            } else if face == af.BACK_FACE { // back
                                pos[af.ZAXIS]--
                            } else if face == af.LEFT_FACE { // left
                                pos[af.XAXIS]++
                            } else if face == af.RIGHT_FACE { // right
                                pos[af.XAXIS]--
                            }
                            world.Set(pos[af.XAXIS], pos[af.YAXIS], pos[af.ZAXIS], 2)
                        }
                    }
                }




            case *sdl.QuitEvent:
                done = true
                break
            }
        }
        keys = sdl.GetKeyState()

        if keys[sdl.K_ESCAPE] != 0 {
            ShowOverlay = !ShowOverlay

            // Overlay

        
        }
        if keys[sdl.K_UP] != 0 {
            
            if view_rotx < 75 {
                view_rotx += 5.0
            }
        }
        if keys[sdl.K_DOWN] != 0 {
            if view_rotx > 15.0 {
                view_rotx -= 5.0
            }
        }
        if keys[sdl.K_LEFT] != 0 {
            view_roty += 9
            //println("view_roty:", view_roty)
        }
        if keys[sdl.K_RIGHT] != 0 {
            view_roty -= 9
        }

        player.HandleKeys(keys)

        if keys[sdl.K_F3] != 0 {
            if af.DebugMode == true {
                af.DebugMode = false
            } else {
                af.DebugMode = true
            }
        }               

        if keys[sdl.K_o] != 0 {
            timeOfDay += 0.1
            if timeOfDay > 24 {
                timeOfDay -= 24
            }
        }
        if keys[sdl.K_l] != 0 {
            timeOfDay -= 0.1
            if timeOfDay < 0 {
                timeOfDay += 24
            }
        }

        if keys[sdl.K_i] != 0 {
            if af.ViewRadius < 90 {
                af.ViewRadius++
                println("ViewRadius: ", af.ViewRadius)
            }
        }
        if keys[sdl.K_k] != 0 {
            if af.ViewRadius > 10 {
                af.ViewRadius--
                println("ViewRadius: ", af.ViewRadius)
            }
        }

        if af.DebugMode {
            fmt.Printf("x:%f, z:%f\n", player.X(), player.Z())
        }

        newTime := time.Now().UnixNano()
        deltaTime := newTime - currentTime
        currentTime = newTime
        if deltaTime > 1e9/4 {
            deltaTime = 1e9/4
        }

        accumulator += deltaTime

        for accumulator > dt {
            accumulator -= dt

            world.ApplyForces(player, float64(dt) / 1e9)
            
            player.Update(float64(dt) / 1e9)
            world.Simulate(float64(dt) / 1e9)

            computeFrame++
            t += dt
        }

        //interpolate(previous, current, accumulator/dt)


        Draw(false)
        drawFrame++

        if update.GetTicks() > 1e9/2 {
            fmt.Printf("draw fps: %f\n", float64(drawFrame) / (float64(update.GetTicks()) / float64(1e9)) )
            fmt.Printf("compute fps: %f\n", float64(computeFrame) / (float64(update.GetTicks()) / float64(1e9)) )

            timeOfDay += 0.02
            if timeOfDay > 24 {
                timeOfDay -= 24
            }

            drawFrame, computeFrame = 0, 0
            update.Start()
        }

    }
    sdl.Quit()
    return

}

func screenToView(xs uint16, ys uint16) (xv float64, yv float64) {
    // xs = 0 => -float64(screenWidth) / screenScale
    // xs = screenWidth => float64(screenWidth) / screenScale

    viewWidth := 2 * float64(screenWidth) / float64(screenScale)
    xv = (-viewWidth / 2 + viewWidth * float64(xs) / float64(screenWidth))

    viewHeight := 2 * float64(screenHeight) / float64(screenScale)
    yv = (-viewHeight / 2 + viewHeight * float64(ys) / float64(screenHeight))

    return
}

/* new window size or exposure */
func reshape(width int, height int) {
    screenWidth = width
    screenHeight = height

    gl.Viewport(0, 0, width, height)
    gl.MatrixMode(gl.PROJECTION)
    gl.LoadIdentity()
    
    xmin, ymin := screenToView(0, 0)
    xmax, ymax := screenToView(uint16(width), uint16(height))
    
    gl.Ortho(float64(xmin), float64(xmax), float64(ymin), float64(ymax), -40, 40)
    gl.MatrixMode(gl.MODELVIEW)
    gl.LoadIdentity()
    glu.LookAt(0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0)
}


func init2() {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    // gl.ShadeModel(gl.FLAT)    
    gl.ShadeModel(gl.SMOOTH)    
    gl.Enable(gl.LIGHTING)
    gl.Enable(gl.LIGHT0)
    gl.Enable(gl.LIGHT1)


    // gl.ColorMaterial ( gl.FRONT_AND_BACK, gl.EMISSION )
    // gl.ColorMaterial ( gl.FRONT_AND_BACK, gl.AMBIENT_AND_DIFFUSE )
    // gl.Enable ( gl.COLOR_MATERIAL )



    gl.MatrixMode(gl.PROJECTION)
    gl.LoadIdentity()
    gl.Ortho(-12.0, 12.0, -12.0, 12.0, -10, 10.0)
    gl.MatrixMode(gl.MODELVIEW)
    gl.LoadIdentity()
    glu.LookAt(0.0, 0.0, 0.0, 0.0, 0.0, -1.0, 0.0, 1.0, 0.0)


    gl.ClearDepth(1.0)                         // Depth Buffer Setup
    gl.Enable(gl.DEPTH_TEST)                        // Enables Depth Testing
    gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.FASTEST)

    gl.Enable(gl.TEXTURE_2D)
    af.LoadMapTextures()
    af.LoadTerrainCubes()


}











func Draw(selectMode bool) {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)


    if selectMode {
        gl.Disable(gl.TEXTURE_2D)
        gl.Disable(gl.FOG)
        gl.Disable(gl.LIGHTING)
        gl.Disable(gl.LIGHT0)
        gl.Disable(gl.LIGHT1)
        gl.Disable ( gl.COLOR_MATERIAL )
    } else {
        gl.Color4ub(255, 255, 255, 255)
        gl.Enable(gl.TEXTURE_2D)
        //gl.Enable(gl.FOG)
        gl.Enable(gl.LIGHTING)
        gl.Enable(gl.LIGHT0)
        gl.Enable ( gl.COLOR_MATERIAL )

        if timeOfDay < 5.3 || timeOfDay > 20.7 {
            gl.Enable(gl.LIGHT1)
        } else {
            gl.Disable(gl.LIGHT1)
        }


    }


    af.CheckGLError()
    gl.LoadIdentity()
    gl.Rotated(view_rotx, 1.0, 0.0, 0.0)
    gl.Rotated(-player.Heading() + view_roty, 0.0, 1.0, 0.0)
    gl.Rotated(0, 0.0, 0.0, 1.0)


    center := player.Position()




    // Sun
    var daylightIntensity float32 = 0.4
    if timeOfDay < 5 || timeOfDay > 21 {
        daylightIntensity = 0.00
        gl.LightModelfv(gl.LIGHT_MODEL_AMBIENT, []float32{0.1, 0.1, 0.1,1})
    } else if timeOfDay < 6 {
        daylightIntensity = 0.4 * (timeOfDay - 5)
    } else if timeOfDay > 20 {
        daylightIntensity = 0.4 * (21 - timeOfDay)
    }



    gl.Lightfv(gl.LIGHT0, gl.POSITION, []float32{0.5, 1, 1, 0})
    gl.Lightfv(gl.LIGHT0, gl.AMBIENT, []float32{daylightIntensity, daylightIntensity, daylightIntensity,1})
    // gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, []float32{daylightIntensity, daylightIntensity, daylightIntensity,1})
    // gl.Lightfv(gl.LIGHT0, gl.SPECULAR, []float32{daylightIntensity, daylightIntensity, daylightIntensity,1})
    gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, []float32{daylightIntensity, daylightIntensity, daylightIntensity,1})
    gl.Lightfv(gl.LIGHT0, gl.SPECULAR, []float32{daylightIntensity, daylightIntensity, daylightIntensity,1})

    // Torch
    ambient := float32(0.6)
    specular := float32(0.6)
    diffuse := float32(1)


    gl.Lightfv(gl.LIGHT1, gl.POSITION, []float32{0,1,0, 1})
    gl.Lightfv(gl.LIGHT1, gl.AMBIENT, []float32{ambient, ambient, ambient,1})
    gl.Lightfv(gl.LIGHT1, gl.SPECULAR, []float32{specular,specular,specular,1})
    gl.Lightfv(gl.LIGHT1, gl.DIFFUSE, []float32{diffuse,diffuse,diffuse,1})
    gl.Lightf(gl.LIGHT1, gl.CONSTANT_ATTENUATION, 1.5)
    gl.Lightf(gl.LIGHT1, gl.LINEAR_ATTENUATION, 0.5)
    gl.Lightf(gl.LIGHT1, gl.QUADRATIC_ATTENUATION, 0.01)
    gl.Lightf(gl.LIGHT1, gl.SPOT_CUTOFF, 35)
    gl.Lightf(gl.LIGHT1, gl.SPOT_EXPONENT, 2.0)
    gl.Lightfv(gl.LIGHT1, gl.SPOT_DIRECTION, []float32{float32(math.Cos(player.Heading() * math.Pi/180)),float32(-0.7), -float32(math.Sin(player.Heading() * math.Pi/180))})

    af.CheckGLError()

    player.Draw(center, selectMode)
    af.CheckGLError()

    world.Draw(center, selectMode)
    af.CheckGLError()

    // // var mousex, mousey int
    // // mouseState := sdl.GetMouseState(&mousex, &mousey)
    // gl.PushMatrix()
    // gl.Translatef(float32(center[af.XAXIS]),float32(center[af.YAXIS])-1,float32(center[af.ZAXIS]))
    // //print ("i:", i, "j:", j, "b:", world.At(i, j, groundLevel))
    // af.HighlightCuboidFace(1, 1, 1, af.TOP_FACE)
    // gl.PopMatrix()

    //gl.Translatef( 0.0, -20.0, -5.0 )

    if ShowOverlay {
        gl.PushMatrix();
        gl.LoadIdentity();
        gl.Color4f(0, 0, 0, 0.25);
        gl.Begin(gl.QUADS);
        gl.Vertex2f(0, 0);
        gl.Vertex2f(float32(screenWidth), 0);
        gl.Vertex2f(float32(screenWidth), float32(screenHeight));
        gl.Vertex2f(0, float32(screenHeight));
        gl.End();
        gl.PopMatrix();
    }

    if !selectMode {
        // gl.Finish()
        // gl.Flush()
        sdl.GL_SwapBuffers()
    }    
}