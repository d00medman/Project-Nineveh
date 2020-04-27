// +build !js
//Methodology cribbed from this https://medium.com/@ben.mcclelland/an-adventure-into-cgo-calling-go-code-with-c-b20aa6637e75

//todo: determine how the ALE included its games; might need to zip rom files and game states to be viable

package main

//#include <stdio.h>
//#include <errno.h>
//#include <stdlib.h>
import "C"

import (
	"./nes"
	"bytes"
	"strconv"
)

// todo: see if there is a "more elegant" way to retain this reference
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
	return  C.CString(buffer.String())
}

//export TakeAction
func TakeAction(btn int) C.float {
	emu := handle[0]
	//reward := emu.Act(btn)
	return C.float(emu.Act(btn))
}

//export IsGameOver
func IsGameOver() bool {
	emu := handle[0]
	//gameOver := emu.IsGameOver()
	return emu.IsGameOver()
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

//export OpenToStart
func OpenToStart() {
	emu := handle[0]
	emu.OpenToStart()
}

//export EndRecording
func EndRecording(filename *C.char) {
	emu := handle[0]
	emu.EndRecording(C.GoString(filename))
}

func main() {}