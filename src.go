package main

import (
	"github.com/ajstarks/svgo"
	"math/rand"
	"os"
	"time"
)

func ri(n int) int { return rand.Intn(n) }

func main() {
	canvas := svg.New(os.Stdout)
	width := 500
	height := 500
	nstars := 250
	style := "font-size:48pt;fill:white;text-anchor:middle"

	rand.Seed(time.Now().Unix())
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height)
	for i := 0; i < nstars; i++ {
		canvas.Circle(ri(width), ri(height), ri(3), "fill:white")
	}
	canvas.Circle(width/2, height, width/2, "fill:rgb(77, 117, 232)")
	canvas.Text(width/2, height*4/5, "hello, world", style)
	canvas.End()
}
