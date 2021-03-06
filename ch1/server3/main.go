// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 21.

// Server3 is an "echo" server that displays request parameters.
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

var palette = []color.Color{color.White, color.Black}

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)

func main() {
	http.HandleFunc("/", lissajous)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(out http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		os.Exit(1)
	}
	fmt.Printf("Form size: %d\n", len(r.Form))
}

func lissajous(out http.ResponseWriter, r *http.Request) {
	const (
		res     = 0.001 // angular resolution
		size    = 500   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	if err := r.ParseForm(); err != nil {
		os.Exit(1)
	}
	for k, v := range r.Form {
		fmt.Printf("%s:%s\n", k, v)
	}
	cyclesInt, _ := strconv.Atoi(r.Form.Get("cycles")) // number of complete x oscillator revolutions
	cycles := float64(cyclesInt)
	fmt.Printf("Float is: %f\n", cycles)
	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5),
				blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}
