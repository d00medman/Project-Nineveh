// +build !js
//Methodology cribbed from this https://medium.com/@ben.mcclelland/an-adventure-into-cgo-calling-go-code-with-c-b20aa6637e75

package main

//#include <stdio.h>
//#include <errno.h>
//#include <stdlib.h>
import "C"

import (
	"./nes"
	"bytes"
	"strconv"

	//"unsafe"
)

// Not clear if I NEED to use an array, can for sure see this approach being unstable, but need to just get something up and running
var handle []*nes.EmulatorInterface

//export NewNintendo
func NewNintendo(filename *C.char, displayMode *C.char, frameSkipRate C.int) (rc int) {
	options := &nes.Options{
		Region: "NTSC",
		Display: C.GoString(displayMode),
	}
	//nameString := C.GoString(filename)
	h := nes.NewEmulatorInterface(C.GoString(filename), int(frameSkipRate), options)
	handle = append(handle, h)
	return len(handle) - 1
}

//export RunNintendo
func RunNintendo() {
	emu := handle[0]
	emu.Start()
}

//export GetObservation
func GetObservation() *C.char {
	emu := handle[0]
	output := emu.Observe()
	var buffer bytes.Buffer
	for i := 0; i < len(output); i++ {
		buffer.WriteString(strconv.Itoa(int(output[i])))
		if i != len(output)-1 {
			buffer.WriteString(",")
		}
	}
	sendString := buffer.String()
	//fmt.Println(output)
	//outputBytes, _ := json.Marshal(output)
	//sendString := string(outputBytes[:])
	//fmt.Println(sendString)
	//fmt.Println("**end of go code**")
	cSend := C.CString(sendString)
	// Going to see what happens if I try this w/o defer first
	//defer C.free(unsafe.Pointer(cSend))
	return cSend
}

//export TakeAction
func TakeAction(btn int) C.float {
	emu := handle[0]
	reward := emu.Act(btn)
	return C.float(reward)
}

//export IsGameOver
func IsGameOver() bool {
	emu := handle[0]
	gameOver := emu.IsGameOver()
	return gameOver
}

//export CloseEmulator
func CloseEmulator()  {
	emu := handle[0]
	emu.Close()
}

//export Reset
func Reset() {
	emu := handle[0]
	emu.Reset()
}

//export OneFrameAdvance
func OneFrameAdvance() {
	emu := handle[0]
	emu.OneFrameAdvance()
}

//export OpenToStart
func OpenToStart() {
	emu := handle[0]
	emu.OpenToStart()
}

//export EndRecording
func EndRecording() {
	emu := handle[0]
	emu.EndRecording()
}

func main() {}