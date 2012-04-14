package af

import (    
    "math"
    "github.com/banthar/gl"
   // "fmt"
)



type Player struct {
    Bounce float64
    heading float64
    position Vector
    velocity Vector
    falling bool
}

func (p *Player) Init(heading float64, x float32, z float32, y float32) {
    p.heading = heading
    p.position[XAXIS] = float64(x)
    p.position[YAXIS] = float64(y)
    p.position[ZAXIS] = float64(z)
}

func (p *Player) W() float64 { return 0.8 }
func (p *Player) H() float64 { return 1.9 }
func (p *Player) D() float64 { return 0.6 }

func (p *Player) Heading() float64 { return p.heading }
func (p *Player) X() float32 { return float32(p.position[XAXIS]) }
func (p *Player) Y() float32 { return float32(p.position[YAXIS]) }
func (p *Player) Z() float32 { return float32(p.position[ZAXIS]) }
func (p *Player) Velocity() Vector { return p.velocity }
func (p *Player) Position() Vector { return p.position }
func (p *Player) IntPosition() IntVector { 
    return IntVector{ int16(math.Floor(p.position[XAXIS] + 0.5)),
                      int16(math.Floor(p.position[YAXIS] + 0.5)),
                      int16(math.Floor(p.position[ZAXIS] + 0.5))}
}

func (p *Player) FrontBlock() IntVector { 
    ip := p.IntPosition()
    if p.heading > 337.5 || p.heading <= 22.5 {
        ip[XAXIS]++
    } else if p.heading > 22.5 && p.heading <= 67.5 {
        ip[XAXIS]++
        ip[ZAXIS]--
    } else if p.heading > 67.5 && p.heading <= 112.5 {
        ip[ZAXIS]--
    } else if p.heading > 112.5 && p.heading <= 157.5 {
        ip[XAXIS]--
        ip[ZAXIS]--
    } else if p.heading > 157.5 && p.heading <= 202.5 {
        ip[XAXIS]--
    } else if p.heading > 202.5 && p.heading <= 247.5 {
        ip[XAXIS]--
        ip[ZAXIS]++
    } else if p.heading > 247.5 && p.heading <= 292.5 {
        ip[ZAXIS]++
    } else if p.heading > 292.5 && p.heading <= 337.5 {
        ip[XAXIS]++
        ip[ZAXIS]++
    }

    return ip
}


func (p *Player) SetFalling(b bool) { p.falling = b }

func (p *Player) Rotate(angle float64) {
    p.heading += angle
    if p.heading < 0 {
        p.heading += 360
    }
    if p.heading > 360 {
        p.heading -= 360
    }
}

func (p *Player) Update(dt float64) {
    p.position[XAXIS] += p.velocity[XAXIS] * dt
    p.position[YAXIS] += p.velocity[YAXIS] * dt
    p.position[ZAXIS] += p.velocity[ZAXIS] * dt
    // fmt.Printf("position: %s\n", p.position)
}

func (p *Player) Accelerate(v Vector) {
    p.velocity[XAXIS] += v[XAXIS]
    p.velocity[YAXIS] += v[YAXIS]
    p.velocity[ZAXIS] += v[ZAXIS]
}

func (p *Player) IsFalling() bool {
    return p.falling
}

func (p *Player) Snapx(x float64, vx float64) {
    p.position[XAXIS] = x
    p.velocity[XAXIS] = vx
}

func (p *Player) Snapz(z float64, vz float64) {
    p.position[ZAXIS] = z
    p.velocity[ZAXIS] = vz
}

func (p *Player) Snapy(y float64, vy float64) {
    p.position[YAXIS] = y
    p.velocity[YAXIS] = vy
}

func (p *Player) Setvx(vx float64) {
    p.velocity[XAXIS] = vx
}

func (p *Player) Setvz(vz float64) {
    p.velocity[ZAXIS] = vz
}

func (p *Player) Setvy(vy float64) {
    p.velocity[YAXIS] = vy
}

func (p *Player) Act(dt float64) {
    // noop
}


func (player *Player) Draw(pos Vector, selectMode bool) {

    gl.PushMatrix()
    //stepHeight := float32(math.Sin(player.Bounce * piover180)/10.0)
    gl.Rotated(player.Heading(), 0.0, 1.0, 0.0)

    h := float32(player.H()) / 2
    gl.Translatef(0.0, h / 2 ,0.0)
    w := float32(player.W()) / 2
    d := float32(player.D()) / 2
    //gl.Translatef(0.0,-h,0.0)
    MapTextures[33].Bind(gl.TEXTURE_2D)
    //topTexture.Bind(gl.TEXTURE_2D)
    gl.Begin(gl.QUADS)                  // Start Drawing Quads
        //gl.Color3f(0.3,0.3,0.6)
        // Front face
        gl.Normal3f( 1.0, 0.0, 0.0)
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f( d, -h, -w)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f( d,  h, -w)  // Top Right Of The Texture and Quad
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f( d,  h,  w)  // Top Left Of The Texture and Quad
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f( d, -h,  w)  // Bottom Left Of The Texture and Quad

    gl.End()

    MapTextures[32].Bind(gl.TEXTURE_2D)

    // dirtTexture.Bind(gl.TEXTURE_2D)
    gl.Begin(gl.QUADS)                  // Start Drawing Quads
        // Left Face
        gl.Normal3f( 0.0, 0.0, -1.0)
        gl.TexCoord2f(1.0, 0.0)        
        gl.Vertex3f(-d, -h, -w)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f(-d,  h, -w)  // Top Right Of The Texture and Quad
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f( d,  h, -w)  // Top Left Of The Texture and Quad
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f( d, -h, -w)  // Bottom Left Of The Texture and Quad


        // Right Face
        //gl.Color3f(0.5,0.5,1.0)              // Set The Color To Blue One Time Only
        gl.Normal3f( 0.0, 0.0, 1.0)
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f( -d, -h,  w)  // Bottom Left Of The Texture and Quad
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f(  d, -h,  w)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f(  d,  h,  w)  // Top Right Of The Texture and Quad
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f( -d,  h,  w)  // Top Left Of The Texture and Quad


        // Back Face
        gl.Normal3f( -1.0, 0.0, 0.0)
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f(-d, -h, -w)  // Bottom Left Of The Texture and Quad
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f(-d, -h,  w)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f(-d,  h,  w)  // Top Right Of The Texture and Quad
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f(-d,  h, -w)  // Top Left Of The Texture and Quad

     //gl.Color3f(0.3,1.0,0.3)
        // Top Face
        gl.Normal3f( 0.0, 1.0, 0.0)
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f(-d,  h, -w)  // Top Left Of The Texture and Quad
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f(-d,  h,  w)  // Bottom Left Of The Texture and Quad
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f( d,  h,  w)  // Bottom Right Of The Texture and Quad
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f( d,  h, -w)  // Top Right Of The Texture and Quad

        // Bottom Face
        gl.Normal3f( 0.0, -1.0, 0.0)
        gl.TexCoord2f(1.0, 1.0)
        gl.Vertex3f(-d, -h, -w)  // Top Right Of The Texture and Quad
        gl.TexCoord2f(0.0, 1.0)
        gl.Vertex3f( d, -h, -w)  // Top Left Of The Texture and Quad
        gl.TexCoord2f(0.0, 0.0)
        gl.Vertex3f( d, -h,  w)  // Bottom Left Of The Texture and Quad
        gl.TexCoord2f(1.0, 0.0)
        gl.Vertex3f(-d, -h,  w)  // Bottom Right Of The Texture and Quad

    gl.End();   
    gl.PopMatrix()
}
