// todo: move the interface and all other nineveh files out of the nes package and into a new nineveh package
package nes

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
)

type EmulatorInterface struct {
	console *NES
	frameSkipRate int
	displayCount int
	//todo: use this to ensure more consistent output of screen state
	currentObservation []uint8
	//todo: delete me when done testing
	blankPixels []uint8
	blankTotal int
	actionsTaken int
}

// todo: frame skip rate as a settable value
func NewEmulatorInterface(filename string, frameSkipRate int, options *Options) (emulatorinterface *EmulatorInterface) {
	emulator, err := NewNES(filename, options)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Accept error: %s\n", err)
	}

	//todo: delete
	blank := make([]uint8, 61440)
	for i := 0; i < 61440; i++ {
		blank[i] = 0
	}

	emulatorinterface = &EmulatorInterface{
		console: emulator,
		frameSkipRate: frameSkipRate,
		// todo: delete when no longer needed
		blankPixels: blank,
		blankTotal: 0,
		actionsTaken: 0,
	}

	// why do I need to return here but did not in the other files?
	return
}

// Runs the emulator for a human user
func (emulatorinterface *EmulatorInterface) Start(){
	if err := emulatorinterface.console.Run(); err != nil {
		log.Printf("Error during run loop: %s \n", err)
	}
}

func (emulatorinterface *EmulatorInterface) Observe() (colors []uint8) {
	colors = emulatorinterface.currentObservation
	emulatorinterface.console.video.OutputScreenImage(colors)
	return
}

func getActionFields(btn int) (button Button, actionType string, actionName string) {
	switch btn {
	case 0:
		button = A
		actionType = "hold"
		actionName = "hold A"
	case 1:
		button = A
		actionType = "release"
		actionName = "release A"
	case 2:
		button = A
		actionType = "press"
		actionName = "press A"
	case 3:
		button = B
		actionType = "hold"
		actionName = "hold B"
	case 4:
		button = B
		actionType = "release"
		actionName = "release A"
	case 5:
		button = B
		actionType = "press"
		actionName = "press B"
	case 6:
		button = Select
		actionType = "hold"
		actionName = "hold Select"
	case 7:
		button = Select
		actionType = "release"
		actionName = "release Select"
	case 8:
		button = Select
		actionType = "press"
		actionName = "press Select"
	case 9:
		button = Start
		actionType = "release"
		actionName = "release Start"
	case 10:
		button = Start
		actionType = "release"
		actionName = "release Start"
	case 11:
		button = Start
		actionType = "press"
		actionName = "press Start"
	case 12:
		button = Up
		actionType = "hold"
		actionName = "hol Up"
	case 13:
		button = Up
		actionType = "release"
		actionName = "release Up"
	case 14:
		button = Up
		actionType = "press"
		actionName = "press Up"
	case 15:
		button = Down
		actionType = "hold"
		actionName = "hold Down"
	case 16:
		button = Down
		actionType = "release"
		actionName = "release Down"
	case 17:
		button = Down
		actionType = "press"
		actionName = "press Down"
	case 18:
		button = Left
		actionType = "hold"
		actionName = "hold Left"
	case 19:
		button = Left
		actionType = "release"
		actionName = "release Left"
	case 20:
		button = Left
		actionType = "press"
		actionName = "press Left"
	case 21:
		button = Right
		actionType = "hold"
		actionName = "hold Right"
	case 22:
		button = Right
		actionType = "release"
		actionName = "release Right"
	case 23:
		button = Right
		actionType = "press"
		actionName = "press Right"
	default:
		button = One
		actionName = "No action"
		actionType = "none"
	}
	return
}

// Method to input actions
func (emulatorinterface *EmulatorInterface) Act(btn int) (reward float32) {
	button, actionType, actionName := getActionFields(btn)
	log.Printf("%s selected for action %v\n", actionName, emulatorinterface.actionsTaken)
	// Don't want the agent accidentally pausing the thing. will eventually want this somewhere else
	if button == Start {
		//log.Printf("Use of start disallowed\n")
		return 0
	}

	frameForward := func(frames int) {
		for i := 0; i < frames; i++ {
			emulatorinterface.oneFrameAdvance(false)
		}
	}

	getButtonAction := func(button Button, down bool) *ControllerEvent {
		return &ControllerEvent{
			Controller: 0,
			Button: button,
			Down: down,
		}
	}

	//var nonAction bool = false

	switch actionType {
	case "press":
		var advancedFrames = 1
		if emulatorinterface.console.controllers.KeyIsDown(0, button) {
			//fmt.Printf("releasing %v from hold\n", button)
			emulatorinterface.console.events <- getButtonAction(button, false)
			emulatorinterface.oneFrameAdvance(false)
			advancedFrames++
		}
		emulatorinterface.console.events <- getButtonAction(button, true)
		// Skip by the frame rate minus the frames we advanced to ensure the button up events register
		frameForward(emulatorinterface.frameSkipRate - advancedFrames)
		emulatorinterface.console.events <- getButtonAction(button, false)
		emulatorinterface.oneFrameAdvance(false)
	case "hold":
		if emulatorinterface.console.controllers.KeyIsDown(0, button) {
			//fmt.Printf("%v is already being held: this is a non-action\n", button)
			//nonAction = true
		}
		emulatorinterface.console.events <- getButtonAction(button, true)
		frameForward(emulatorinterface.frameSkipRate)
	case "release":
		if !emulatorinterface.console.controllers.KeyIsDown(0, button) {
			//fmt.Printf("%v is not being held: this is a non-action\n", button)
			//nonAction = true
		}
		emulatorinterface.console.events <- getButtonAction(button, false)
		frameForward(emulatorinterface.frameSkipRate)
	default:
		frameForward(emulatorinterface.frameSkipRate)
	}
	emulatorinterface.actionsTaken++

	//if nonAction {
	//	// Penalize if the agent tries to take a non-action
	//	reward = -1.0
	//} else {
	//	reward = emulatorinterface.console.getReward()
	//}
	return emulatorinterface.console.getReward()
}

func (nes *NES) getReward() (reward float32) {
	switch nes.GameName {
	case "Castlevania":
		return nes.GetCastlevaniaScore()
	default:
		return rand.Float32()
	}
}

func (nes *NES) isGameOver() bool {
	switch nes.GameName {
	case "Castlevania":
		livesLeft := nes.CPU.Memory.Fetch(0x002A)
		//log.Printf("Lives left: %v\n", livesLeft)
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
	fmt.Println("Advancing first 15 frames to warm up system")
	for i := 0; i < 16; i++ {
		emulatorinterface.oneFrameAdvance(true)
	}
}

// Loads the starting state of the selected game
func (emulatorinterface *EmulatorInterface) OpenToStart() {
	emulatorinterface.actionsTaken = 0
	emulatorinterface.console.LoadState()
	for i := 0; i < 16; i++ {
		emulatorinterface.oneFrameAdvance(true)
	}
}

func (emulatorinterface *EmulatorInterface) oneFrameAdvance(warmup bool) {
	colors, _, err := emulatorinterface.console.ProcessToFrame()

	if err != nil {
		log.Printf("Error during process to frame: %s \n", err)
	}

	if !reflect.DeepEqual(colors, emulatorinterface.blankPixels) {
		emulatorinterface.currentObservation = colors
	} else if !warmup {
		// Can probably remove this block, as its logic no longer seems to really track frame misfires, which seem to have largely been ameliorated by the change to using a select statement
		emulatorinterface.blankTotal++
		//fmt.Printf("Blank output from frame processing at frame %v, which completed in %v cycles. %v blank outputs total\n", emulatorinterface.console.frameCount, cycleCount, emulatorinterface.blankTotal)
	}
}

func (emulatorinterface *EmulatorInterface) EndRecording() {
	emulatorinterface.console.video.OutputRunRecording()
}