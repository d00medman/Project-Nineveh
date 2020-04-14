// adapted from github.com/scottferg/Go-SDL/gfx/framerate.go

package nes

import (
	//"fmt"
	"time"
)

const (
	DefaultFPSNTSC float64 = 60.0988
	DefaultFPSPAL  float64 = 50.0070
)

type FPS struct {
	enabled bool
	frames  float64
	rate    float64
	ticks   uint64
}

func NewFPS(rate float64) *FPS {
	fps := &FPS{}

	fps.SetRate(rate)

	return fps
}

func (fps *FPS) Enable() {
	fps.enabled = true
}

func (fps *FPS) Disable() {
	fps.enabled = false
}

func (fps *FPS) Resumed() {
	fps.frames = 0
	fps.ticks = uint64(time.Now().UnixNano()) / 1e6
}

func (fps *FPS) SetRate(rate float64) {
	fps.Enable()
	fps.Resumed()
	fps.rate = 1000.0 / rate
}

func (fps *FPS) Delay() {
	// next frame
	fps.frames++

	//fmt.Printf("hits fps delay, frame: %v \n", fps.frames)

	// get/calc ticks
	current := uint64(time.Now().UnixNano()) / 1e6
	target := fps.ticks + uint64(fps.frames*fps.rate)

	if fps.enabled && current <= target {
		//fmt.Println("hits fps enabled sleep if")
		time.Sleep(time.Duration((target - current) * 1e6))
	} else {
		//fmt.Println("hits fps enabled else")
		fps.frames = 0.0
		fps.ticks = uint64(time.Now().UnixNano()) / 1e6
	}
}
