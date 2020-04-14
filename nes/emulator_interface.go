package nes

import (
	"fmt"
	"log"
	"os"
)

type EmulatorInterface struct {
	console *NES
	frameSkipRate int
	displayCount int
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
func (emulatorinterface *EmulatorInterface) Act(btn int) {
	var hold bool
	button := One
	var actionName string

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