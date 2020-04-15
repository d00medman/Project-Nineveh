package nes

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"log"
	"os"
	"sync"
)

var HeadlessPalette []color.Color = []color.Color{
	color.RGBA{0x66, 0x66, 0x66, 0xff},
	color.RGBA{0x00, 0x2A, 0x88, 0xff},
	color.RGBA{0x14, 0x12, 0xA7, 0xff},
	color.RGBA{0x3B, 0x00, 0xA4, 0xff},
	color.RGBA{0x5C, 0x00, 0x7E, 0xff},
	color.RGBA{0x6E, 0x00, 0x40, 0xff},
	color.RGBA{0x6C, 0x06, 0x00, 0xff},
	color.RGBA{0x56, 0x1D, 0x00, 0xff},
	color.RGBA{0x33, 0x35, 0x00, 0xff},
	color.RGBA{0x0B, 0x48, 0x00, 0xff},
	color.RGBA{0x00, 0x52, 0x00, 0xff},
	color.RGBA{0x00, 0x4F, 0x08, 0xff},
	color.RGBA{0x00, 0x40, 0x4D, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0xAD, 0xAD, 0xAD, 0xff},
	color.RGBA{0x15, 0x5F, 0xD9, 0xff},
	color.RGBA{0x42, 0x40, 0xFF, 0xff},
	color.RGBA{0x75, 0x27, 0xFE, 0xff},
	color.RGBA{0xA0, 0x1A, 0xCC, 0xff},
	color.RGBA{0xB7, 0x1E, 0x7B, 0xff},
	color.RGBA{0xB5, 0x31, 0x20, 0xff},
	color.RGBA{0x99, 0x4E, 0x00, 0xff},
	color.RGBA{0x6B, 0x6D, 0x00, 0xff},
	color.RGBA{0x38, 0x87, 0x00, 0xff},
	color.RGBA{0x0C, 0x93, 0x00, 0xff},
	color.RGBA{0x00, 0x8F, 0x32, 0xff},
	color.RGBA{0x00, 0x7C, 0x8D, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0xFF, 0xFE, 0xFF, 0xff},
	color.RGBA{0x64, 0xB0, 0xFF, 0xff},
	color.RGBA{0x92, 0x90, 0xFF, 0xff},
	color.RGBA{0xC6, 0x76, 0xFF, 0xff},
	color.RGBA{0xF3, 0x6A, 0xFF, 0xff},
	color.RGBA{0xFE, 0x6E, 0xCC, 0xff},
	color.RGBA{0xFE, 0x81, 0x70, 0xff},
	color.RGBA{0xEA, 0x9E, 0x22, 0xff},
	color.RGBA{0xBC, 0xBE, 0x00, 0xff},
	color.RGBA{0x88, 0xD8, 0x00, 0xff},
	color.RGBA{0x5C, 0xE4, 0x30, 0xff},
	color.RGBA{0x45, 0xE0, 0x82, 0xff},
	color.RGBA{0x48, 0xCD, 0xDE, 0xff},
	color.RGBA{0x4F, 0x4F, 0x4F, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0xFF, 0xFE, 0xFF, 0xff},
	color.RGBA{0xC0, 0xDF, 0xFF, 0xff},
	color.RGBA{0xD3, 0xD2, 0xFF, 0xff},
	color.RGBA{0xE8, 0xC8, 0xFF, 0xff},
	color.RGBA{0xFB, 0xC2, 0xFF, 0xff},
	color.RGBA{0xFE, 0xC4, 0xEA, 0xff},
	color.RGBA{0xFE, 0xCC, 0xC5, 0xff},
	color.RGBA{0xF7, 0xD8, 0xA5, 0xff},
	color.RGBA{0xE4, 0xE5, 0x94, 0xff},
	color.RGBA{0xCF, 0xEF, 0x96, 0xff},
	color.RGBA{0xBD, 0xF4, 0xAB, 0xff},
	color.RGBA{0xB3, 0xF3, 0xCC, 0xff},
	color.RGBA{0xB5, 0xEB, 0xF2, 0xff},
	color.RGBA{0xB8, 0xB8, 0xB8, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
}

type HeadlessVideo struct {
	input         chan []uint8
	width, height int
	palette       []color.Color
	events        chan Event
	framePool     *sync.Pool
	overscan      bool
	caption       string
	fps           float64
	displayCount int
	gif     *gif.GIF
}

func NewHeadlessVideo(caption string,
	events chan Event,
	framePool *sync.Pool,
	fps float64) (video *HeadlessVideo, err error) {
	fmt.Println("init headless video")

	runGif := &gif.GIF{
		Image:     []*image.Paletted{},
		Delay:     []int{},
		LoopCount: 0xfffffff,
	}

	video = &HeadlessVideo{
		input:     make(chan []uint8, 128),
		events:    events,
		framePool: framePool,
		palette:   HeadlessPalette,
		overscan:  true,
		caption:   caption,
		fps:       fps,
		displayCount: 0,
		gif: runGif,
	}

	return
}

func (video *HeadlessVideo) Run() {
	fmt.Println("running in headless mode")
	//iterations := 0
	for {
		done := make(chan bool)

		paletteSrc := image.NewRGBA(image.Rect(0, 0, 64, 2))
		for x, c := range HeadlessPalette {
			paletteSrc.Set(x, 0, c)
		}

		img := image.NewRGBA(image.Rect(0, 0, 256, 256))
		updateTex := func(colors []uint8) {
			for i, c := range colors {
				img.Pix[i<<2] = c
			}
			video.framePool.Put(colors)
		}

		go func() {
			for {
				select {
				case <-done:
					return
				case colors := <-video.input:
					updateTex(colors)
				}
			}
		}()
	}
}


func (video *HeadlessVideo) Events() chan Event {
	return video.events
}

func (video *HeadlessVideo) Input() chan []uint8 {
	return video.input
}

func (video *HeadlessVideo) SetCaption(caption string) {
	video.caption = caption
}

func (video *HeadlessVideo) GetScreenImage(colors []uint8) (frame *image.Paletted) {
	frame = image.NewPaletted(image.Rect(0, 0, 256, 240), HeadlessPalette)
	x, y := 0, 0
	for _, c := range colors {
		frame.Set(x, y, video.palette[c])

		switch x {
		case 255:
			x = 0
			y++
		default:
			x++
		}
	}
	return
}

func (video *HeadlessVideo) AddImageToRecording(colors []uint8) {
	frame := video.GetScreenImage(colors)
	video.gif.Image = append(video.gif.Image, frame)
	video.gif.Delay = append(video.gif.Delay, 5)
}

func (video *HeadlessVideo) OutputRunRecording() {
	if video.gif != nil {
		fileName := fmt.Sprintf("output/frame.gif")
		fmt.Printf("Output recording to %s\n", fileName)

		//todo: will want a degree of variance in name of gif, especially given that multiple environments will be at work at once; already did this for jpeg
		fo, _ := os.Create(fileName)
		w := bufio.NewWriter(fo)
		if err := gif.EncodeAll(w, video.gif); err != nil {
			log.Printf("Error in outputting gif: %s \n", err)
		}
	}
}

func (video *HeadlessVideo) OutputScreenImage(colors []uint8) {
	frame := video.GetScreenImage(colors)

	fileName := fmt.Sprintf("output/img_%v.jpg", video.displayCount)
	video.displayCount++
	fo, _ := os.Create(fmt.Sprintf(fileName))
	w := bufio.NewWriter(fo)
	if err := jpeg.Encode(w, frame, &jpeg.Options{Quality: 100}); err != nil {
		log.Printf("Error saving screen as jpeg: %s \n", err)
	}
}
