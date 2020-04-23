package nes

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
)

type Packet struct {
	Tick uint64
	Ev   Event
}

type Event interface {
	Flag() uint
	Process(nes *NES)
	String() string
}

func init() {
	gob.Register(&FrameEvent{})
	gob.Register(&SampleEvent{})
	gob.Register(&ControllerEvent{})
	gob.Register(&PauseEvent{})
	gob.Register(&ResetEvent{})
	gob.Register(&RecordEvent{})
	gob.Register(&StopEvent{})
	gob.Register(&AudioRecordEvent{})
	gob.Register(&AudioStopEvent{})
	gob.Register(&QuitEvent{})
	gob.Register(&ShowBackgroundEvent{})
	gob.Register(&ShowSpritesEvent{})
	gob.Register(&CPUDecodeEvent{})
	gob.Register(&PPUDecodeEvent{})
	gob.Register(&SaveStateEvent{})
	gob.Register(&LoadStateEvent{})
	gob.Register(&FPSEvent{})
	gob.Register(&SavePatternTablesEvent{})
	gob.Register(&MuteEvent{})
	gob.Register(&MuteNoiseEvent{})
	gob.Register(&MuteTriangleEvent{})
	gob.Register(&MutePulse1Event{})
	gob.Register(&MutePulse2Event{})
	gob.Register(&HeartbeatEvent{})
}

const (
	EvGlobal uint = 1 << iota
	EvMaster
	EvSlave
)

type FrameEvent struct {
	Colors []uint8
}

func (e *FrameEvent) String() string {
	return "FrameEvent"
}

/*
At this point, the "frame event" is pretty wildly mangled. Going to need to do some finagling if I want this system
to be human controllable
 */
func (e *FrameEvent) Process(nes *NES) {
	if nes.state != Running {
		fmt.Println("state not running, no frame event")
		return
	}

	if nes.recorder != nil {
		fmt.Println("recording during frame event")
		nes.recorder.Input() <- e.Colors
	}
	nes.frameCount++
	// Can probably remove this block as well
	if nes.headless {
		//fmt.Sprintf("Frame event at frame %v, colors: %v", nes.frameCount, e.Colors[8000:8256])
		nes.framePool.Put(e.Colors)
		nes.video.Input() <- e.Colors
	} else {
		//fmt.Println("Headless false during frame event")
		nes.video.Input() <- e.Colors
	}
}

func (e *FrameEvent) Flag() uint {
	return EvMaster | EvSlave
}

type SampleEvent struct {
	Sample int16
}

func (e *SampleEvent) String() string {
	return "SampleEvent"
}

func (e *SampleEvent) Process(nes *NES) {
	if nes.state != Running {
		return
	}

	if nes.audioRecorder != nil {
		nes.audioRecorder.Input() <- e.Sample
	}

	//nes.audio.Input() <- e.Sample
}

func (e *SampleEvent) Flag() uint {
	return EvMaster | EvSlave
}

type ControllerEvent struct {
	Controller int
	Down       bool
	Button     Button
}

func (e *ControllerEvent) String() string {
	return "ControllerEvent"
}

func (e *ControllerEvent) Process(nes *NES) {
	if nes.state != Running {
		fmt.Println("not running on controller event")
		return
	}

	//nes.getReward() // Primarily here for debugging purposes

	//fmt.Printf("Execute controller event %v\n", e.Button)
	if e.Down {
		//fmt.Printf("key down for %v\n", e.Button)
		nes.controllers.KeyDown(e.Controller, e.Button)
	} else {
		//fmt.Printf("key up for %v\n", e.Button)
		nes.controllers.KeyUp(e.Controller, e.Button)
	}
}

func (e *ControllerEvent) Flag() uint {
	return EvGlobal | EvMaster | EvSlave
}

type PauseEvent struct {
	Request PauseRequest
	Changed chan bool
}

func (e *PauseEvent) String() string {
	return "PauseEvent"
}

func (e *PauseEvent) Process(nes *NES) {
	nes.audio.TogglePaused()
	nes.Paused = !nes.Paused
}

func (e *PauseEvent) Flag() uint {
	return EvGlobal | EvMaster | EvSlave
}

type ResetEvent struct{}

func (e *ResetEvent) String() string {
	return "ResetEvent"
}

func (e *ResetEvent) Process(nes *NES) {
	if nes.state != Running {
		return
	}

	nes.Reset()
}

func (e *ResetEvent) Flag() uint {
	return EvGlobal | EvMaster
}

type RecordEvent struct{}

func (e *RecordEvent) String() string {
	return "RecordEvent"
}

func (e *RecordEvent) Process(nes *NES) {
	fmt.Println("processing record event")

	if nes.recorder != nil {
		nes.recorder.Record()
	}
}

func (e *RecordEvent) Flag() uint {
	return EvMaster | EvSlave
}

type StopEvent struct{}

func (e *StopEvent) String() string {
	return "StopEvent"
}

func (e *StopEvent) Process(nes *NES) {
	if nes.recorder != nil {
		nes.recorder.Stop()
	}
}

func (e *StopEvent) Flag() uint {
	return EvMaster | EvSlave
}

type AudioRecordEvent struct{}

func (e *AudioRecordEvent) String() string {
	return "AudioRecordEvent"
}

func (e *AudioRecordEvent) Process(nes *NES) {
	if nes.audioRecorder != nil {
		nes.audioRecorder.Record()
	}
}

func (e *AudioRecordEvent) Flag() uint {
	return EvMaster | EvSlave
}

type AudioStopEvent struct{}

func (e *AudioStopEvent) String() string {
	return "AudioStopEvent"
}

func (e *AudioStopEvent) Process(nes *NES) {
	if nes.audioRecorder != nil {
		nes.audioRecorder.Stop()
	}
}

func (e *AudioStopEvent) Flag() uint {
	return EvMaster | EvSlave
}

type QuitEvent struct{}

func (e *QuitEvent) String() string {
	return "QuitEvent"
}

func (e *QuitEvent) Process(nes *NES) {
	nes.state = Quitting
}

func (e *QuitEvent) Flag() uint {
	return EvMaster | EvSlave
}

type ShowBackgroundEvent struct{}

func (e *ShowBackgroundEvent) String() string {
	return "ShowBackgroundEvent"
}

func (e *ShowBackgroundEvent) Process(nes *NES) {
	nes.PPU.ShowBackground = !nes.PPU.ShowBackground
	fmt.Println("*** Toggling show background =", nes.PPU.ShowBackground)
}

func (e *ShowBackgroundEvent) Flag() uint {
	return EvMaster | EvSlave
}

type ShowSpritesEvent struct{}

func (e *ShowSpritesEvent) String() string {
	return "ShowSpritesEvent"
}

func (e *ShowSpritesEvent) Process(nes *NES) {
	nes.PPU.ShowSprites = !nes.PPU.ShowSprites
	fmt.Println("*** Toggling show sprites =", nes.PPU.ShowSprites)
}

func (e *ShowSpritesEvent) Flag() uint {
	return EvMaster | EvSlave
}

type CPUDecodeEvent struct{}

func (e *CPUDecodeEvent) String() string {
	return "CPUDecodeEvent"
}

func (e *CPUDecodeEvent) Process(nes *NES) {
	fmt.Println("*** Toggling CPU decode =", nes.CPU.ToggleDecode())
}

func (e *CPUDecodeEvent) Flag() uint {
	return EvMaster | EvSlave
}

type PPUDecodeEvent struct{}

func (e *PPUDecodeEvent) String() string {
	return "PPUDecodeEvent"
}

func (e *PPUDecodeEvent) Process(nes *NES) {
	fmt.Println("*** Toggling PPU decode =", nes.PPU.ToggleDecode())
}

func (e *PPUDecodeEvent) Flag() uint {
	return EvMaster | EvSlave
}

type SaveStateEvent struct{}

func (e *SaveStateEvent) String() string {
	return "SaveStateEvent"
}

func (e *SaveStateEvent) Process(nes *NES) {
	nes.SaveState()
}

func (e *SaveStateEvent) Flag() uint {
	return EvGlobal | EvMaster
}

type LoadStateEvent struct {
	Data []byte
}

func (e *LoadStateEvent) String() string {
	return "LoadStateEvent"
}

func (e *LoadStateEvent) Process(nes *NES) {
	if e.Data == nil {
		if !nes.master {
			// Should not go here, NES.processEvents already filter the events.
			return
		}
		name := nes.GameName + ".nst"
		data, err := ioutil.ReadFile(name)
		if err != nil {
			return
		}
		e.Data = data
	}
	reader := bytes.NewReader(e.Data)
	nes.LoadStateFromReader(reader, int64(len(e.Data)))
}

func (e *LoadStateEvent) Flag() uint {
	return EvGlobal | EvMaster
}

type FPSEvent struct {
	Rate float64
}

func (e *FPSEvent) String() string {
	return "FPSEvent"
}

func (e *FPSEvent) Process(nes *NES) {
	nes.fps.SetRate(nes.DefaultFPS * e.Rate)

	nes.audio.SetSpeed(float32(e.Rate))
	fmt.Printf("*** Setting fps to %0.2f\n", e.Rate)
}

func (e *FPSEvent) Flag() uint {
	return EvGlobal | EvMaster
}

type SavePatternTablesEvent struct{}

func (e *SavePatternTablesEvent) String() string {
	return "SavePatternTablesEvent"
}

func (e *SavePatternTablesEvent) Process(nes *NES) {
	fmt.Println("*** Saving PPU pattern tables")
	nes.PPU.SavePatternTables()
}

func (e *SavePatternTablesEvent) Flag() uint {
	return EvMaster | EvSlave
}

type MuteEvent struct{}

func (e *MuteEvent) String() string {
	return "MuteEvent"
}

func (e *MuteEvent) Process(nes *NES) {
	nes.CPU.APU.Muted = !nes.CPU.APU.Muted
	fmt.Println("*** Toggling mute =", nes.CPU.APU.Muted)
}

func (e *MuteEvent) Flag() uint {
	return EvMaster | EvSlave
}

type MuteDMCEvent struct{}

func (e *MuteDMCEvent) String() string {
	return "MuteDMCEvent"
}

func (e *MuteDMCEvent) Process(nes *NES) {
	nes.CPU.APU.DMC.Muted = !nes.CPU.APU.DMC.Muted
	fmt.Println("*** Toggling mute DMC =", nes.CPU.APU.DMC.Muted)
}

func (e *MuteDMCEvent) Flag() uint {
	return EvMaster | EvSlave
}

type MuteNoiseEvent struct{}

func (e *MuteNoiseEvent) String() string {
	return "MuteNoiseEvent"
}

func (e *MuteNoiseEvent) Process(nes *NES) {
	nes.CPU.APU.Noise.Muted = !nes.CPU.APU.Noise.Muted
	fmt.Println("*** Toggling mute noise =", nes.CPU.APU.Noise.Muted)
}

func (e *MuteNoiseEvent) Flag() uint {
	return EvMaster | EvSlave
}

type MuteTriangleEvent struct{}

func (e *MuteTriangleEvent) String() string {
	return "MuteTriangleEvent"
}

func (e *MuteTriangleEvent) Process(nes *NES) {
	nes.CPU.APU.Triangle.Muted = !nes.CPU.APU.Triangle.Muted
	fmt.Println("*** Toggling mute triangle =", nes.CPU.APU.Triangle.Muted)
}

func (e *MuteTriangleEvent) Flag() uint {
	return EvMaster | EvSlave
}

type MutePulse1Event struct{}

func (e *MutePulse1Event) String() string {
	return "MutePulse1Event"
}

func (e *MutePulse1Event) Process(nes *NES) {
	nes.CPU.APU.Pulse1.Muted = !nes.CPU.APU.Pulse1.Muted
	fmt.Println("*** Toggling mute pulse1 =", nes.CPU.APU.Pulse1.Muted)
}

func (e *MutePulse1Event) Flag() uint {
	return EvMaster | EvSlave
}

type MutePulse2Event struct{}

func (e *MutePulse2Event) String() string {
	return "MutePulse2Event"
}

func (e *MutePulse2Event) Process(nes *NES) {
	nes.CPU.APU.Pulse2.Muted = !nes.CPU.APU.Pulse2.Muted
	fmt.Println("*** Toggling mute pulse2 =", nes.CPU.APU.Pulse2.Muted)
}

func (e *MutePulse2Event) Flag() uint {
	return EvMaster | EvSlave
}

type HeartbeatEvent struct{}

func (e *HeartbeatEvent) String() string {
	return "HeartbeatEvent"
}

func (e *HeartbeatEvent) Process(nes *NES) {
	// do nothing
}

func (e *HeartbeatEvent) Flag() uint {
	return EvGlobal | EvMaster
}
