// Copyright 2018 The Ebiten Authors, Lei Yang
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"image"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var carrot *ebiten.Image
var bunny *ebiten.Image 

func init() {
	f, err := os.Open("carrot.png")
    if err != nil {
		panic(err)
    }
    defer f.Close()
    c, _, err := image.Decode(f)
    if err != nil {
		panic(err)
    }
	carrot = ebiten.NewImageFromImage(c)

	f2, err := os.Open("bunny.png")
    if err != nil {
		panic(err)
    }
    defer f2.Close()
    b, _, err := image.Decode(f2)
    if err != nil {
		panic(err)
    }
	bunny = ebiten.NewImageFromImage(b)
}

const (
//	screenWidth  = 640
//	screenWidth  = 220
//	screenHeight = 480
	screenWidth  = 330
	screenHeight = 360
)

type point struct {
	x int
	y int
}

type ringbuffer struct {
	buf []point
	head int
}

func (b *ringbuffer) write(x, y int) {
	p := point{x, y}
	b.head += 1
	if b.head >= len(b.buf) {
		b.head = 0
	}
	b.buf[b.head] = p
}

func (b *ringbuffer) read(offset int) (x, y int) {
	idx := b.head - offset
	if idx < 0 {
		idx += len(b.buf)
	}
	if idx < 0 {
		panic("offset too large")
	}
	return b.buf[idx].x, b.buf[idx].y
}

type Game struct {
	bunX int
	bunY int
	carrotX int
	carrotY int
	lastBun int
	currentTick int
	b *ringbuffer
}

func (g *Game) Update() error {
	g.currentTick += 1

	g.carrotX, g.carrotY = ebiten.CursorPosition()
	g.b.write(g.carrotX, g.carrotY)

	// 100 ticks per second
	if g.currentTick - g.lastBun >= 1 { // block time = 10ms
		g.lastBun = g.currentTick
		g.bunX, g.bunY = g.b.read(10) // network propagation time = 100ms
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.carrotX), float64(g.carrotY))
	screen.DrawImage(carrot, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.bunX+50), float64(g.bunY+10))
	screen.DrawImage(bunny, op)

	ebitenutil.DebugPrint(screen,
	fmt.Sprintf("Carrot (user input): (%d, %d)\nBunny (chain state): (%d, %d)", g.carrotX, g.carrotY, g.bunX, g.bunY))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetTPS(100)	// game engine ticks per second
	g := &Game{}
	g.b = &ringbuffer{make([]point, 3000), 0}	// 30s of memory

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Blocktime Demo")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}