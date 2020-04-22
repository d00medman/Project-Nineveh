package nes

import (
	"fmt"
	"strconv"
)
// TODO: determine how to best design this for multiple games. One large file will get unwieldy fast

// functions for Castlevania
func (nes *NES) GetCastlevaniaScore() (score float32) {
	//fmt.Println("Getting current score of castlevania")
	/*
		Castlevania uses the hex representation of a number in the UI as a decimal, i.e to show '10', it holds '16' in
		memory. I will convert the score to the decimal represented to the player to better enable debugging

		Per https://datacrystal.romhacking.net/wiki/Castlevania:RAM_map, memory addresses of the score
		0x07FC	Ones/Tenths of points	In pseudo Decimal (Ex.: $08, $09, $10...)
		0x07FD	Hundreds/Thousands of points	Ditto.
		0x07FE	Ten Thousands/Hundred Thousands of points
	*/

	convertScore := func (address uint16) (convertedVal int) {
		raw := nes.CPU.Memory.Fetch(address)
		convertedVal, _ = strconv.Atoi(fmt.Sprintf("%x", raw))
		return
	}
	ones := float32(convertScore(0x07FC)) * 0.1
	hundreds := float32(convertScore(0x07FD))
	thousands := float32(convertScore(0x07FE)) * 100
	score = thousands + hundreds + ones
	//log.Printf("score at frame %v: %v", nes.frameCount, score)
	return
}

// TODO: add more games (SMB3, Mario bros and Donkey kong as good starting points)
