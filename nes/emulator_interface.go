// todo: move the interface and all other nineveh files out of the nes package and into a new nineveh package
package nes

import (
	"fmt"
	"log"
	"math/rand"
	"os"
)

type EmulatorInterface struct {
	console *NES
	frameSkipRate int
	displayCount int
	//todo: use this to ensure more consistent output of screen state
	currentPixels []uint8
}

// todo: frame skip rate as a settable value
func NewEmulatorInterface(filename string, frameSkipRate int, options *Options) *EmulatorInterface {
	emulator, err := NewNES(filename, options)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Accept error: %s\n", err)
	}

	emulatorinterface := &EmulatorInterface{
		console: emulator,
		frameSkipRate: frameSkipRate,
	}

	// why do I need to return here but did not in the other files?
	return emulatorinterface
}

// Runs the emulator for a human user
func (emulatorinterface *EmulatorInterface) Start(){
	if err := emulatorinterface.console.Run(); err != nil {
		log.Printf("Error during run loop: %s \n", err)
	}
}

func (emulatorinterface *EmulatorInterface) Observe() []uint8 {
	// Getting observations still requires processing a single frame forward. This is likely a problem
	colors, _ := emulatorinterface.console.ProcessToFrame(true)
	return colors
}

// Method to input actions
func (emulatorinterface *EmulatorInterface) Act(btn int) (reward float32) {
	var hold bool
	button := One
	var actionName string

	// todo: refactor to have a release option as well
	switch btn {
	case 0:
		button = A
		hold = true
		actionName = "hold A"
	case 1:
		button = A
		hold = false
		actionName = "press A"
	case 2:
		button = B
		hold = true
		actionName = "hold B"
	case 3:
		button = B
		hold = false
		actionName = "press B"
	case 4:
		button = Select
		hold = true
		actionName = "hold Select"
	case 5:
		button = Select
		hold = false
		actionName = "press Select"
	case 6:
		button = Start
		hold = true
		actionName = "hold Start"
	case 7:
		button = Start
		hold = false
		actionName = "press Start"
	case 8:
		button = Up
		hold = true
		actionName = "hol Up"
	case 9:
		button = Up
		hold = false
		actionName = "press Up"
	case 10:
		button = Down
		hold = true
		actionName = "hold Down"
	case 11:
		button = Down
		hold = false
		actionName = "press Down"
	case 12:
		button = Left
		hold = true
		actionName = "hold Left"
	case 13:
		button = Left
		hold = false
		actionName = "press Left"
	case 14:
		button = Right
		hold = true
		actionName = "hold Right"
	case 15:
		button = Right
		hold = false
		actionName = "press Right"
	default:
		actionName = "No action"
	}

	log.Printf("%s selected\n", actionName)

	frameForward := func() {
		for i := 0; i < emulatorinterface.frameSkipRate; i++ {
			if _, err := emulatorinterface.console.ProcessToFrame(false); err != nil {
				log.Printf("Error during process to frame: %s \n", err)
			}
		}
	}

	getButtonAction := func(button Button, down bool) *ControllerEvent {
		return &ControllerEvent{
			Controller: 0,
			Button: button,
			Down: down,
		}
	}

	// If this button is already being held, we want to first release the button before taking action on it
	if emulatorinterface.console.controllers.KeyIsDown(0, button) {
		fmt.Println("Releasing held button")
		emulatorinterface.console.events <- getButtonAction(button, false)
	}

	if hold {
		emulatorinterface.console.events <- getButtonAction(button, true)
		frameForward()
	} else {
		emulatorinterface.console.events <- getButtonAction(button, true)
		frameForward()
		emulatorinterface.console.events <- getButtonAction(button, false)
		// Process for one more frame to ensure that the button is actually up. Irregular behavior between the two cases strikes me as off
		if _, err := emulatorinterface.console.ProcessToFrame(false); err != nil {
			log.Printf("Error during process to frame: %s \n", err)
		}
	}

	return emulatorinterface.console.getReward()
}

func (nes *NES) getReward() (reward float32) {
	switch nes.GameName {
	case "Castlevania":
		//fmt.Println("Getting current score of castlevania")
		/*
		Per https://datacrystal.romhacking.net/wiki/Castlevania:RAM_map, memory addresses of the score
		0x07FC	Ones/Tenths of points	In pseudo Decimal (Ex.: $08, $09, $10...)
		0x07FD	Hundreds/Thousands of points	Ditto.
		0x07FE	Ten Thousands/Hundred Thousands of points
		 */
		ones := float32(nes.CPU.Memory.Fetch(0x07FC))
		hundreds := float32(nes.CPU.Memory.Fetch(0x07FD)) * 100
		thousands := float32(nes.CPU.Memory.Fetch(0x07FE)) * 1000
		//log.Printf("result of fetching for memory addresses thousands (0x07FE): %v, hundreds (0x07FD): %v, ones (0x07FC) %v \n", thousands, hundreds, ones)
		/*
		todo: This isn't actually the reward, its really just the score. Need to come up with some methodology for counting the score
		 */
		return thousands + hundreds + ones
	default:
		return rand.Float32()
	}
}

func (nes *NES) isGameOver() bool {
	switch nes.GameName {
	case "Castlevania":
		livesLeft := nes.CPU.Memory.Fetch(0x002A)
		log.Printf("Lives left: %v\n", livesLeft)
		return livesLeft <= 0
	default:
		return true
	}
}

func (emulatorinterface *EmulatorInterface) IsGameOver() bool {
	return emulatorinterface.console.isGameOver()
}

func (emulatorinterface *EmulatorInterface) Close() {
	quit := &QuitEvent{}
	emulatorinterface.console.events <- quit
}

func (emulatorinterface *EmulatorInterface) Reset() {
	emulatorinterface.console.Reset()

	// Currently, the first 15 frames of running seem to be blank. going to forward through these for the time being
	fmt.Println("Advancing first 15 frames")
	for i := 0; i < 16; i++ {
		if _, err := emulatorinterface.console.ProcessToFrame(false); err != nil {
			log.Printf("Error during process to frame: %s \n", err)
		}
	}
}

// Loads the starting state of the selected game
func (emulatorinterface *EmulatorInterface) OpenToStart() {
	emulatorinterface.console.LoadState()
}

// Primarily for debugging at present. Should not be surfaced; remove when work is through
func (emulatorinterface *EmulatorInterface) OneFrameAdvance() {
	//colors, err := emulatorinterface.console.ProcessToFrame()
	//if err != nil {
	//	log.Printf("Error on output file creation: %s \n", err)
	//}
	//fmt.Println(colors)

	if _, err := emulatorinterface.console.ProcessToFrame(true); err != nil {
		log.Printf("Error during process to frame: %s \n", err)
	}

	fmt.Println("OneFrameAdvance finished")
	fmt.Println("******")
}

func (emulatorinterface *EmulatorInterface) EndRecording() {
	emulatorinterface.console.video.OutputRunRecording()
}