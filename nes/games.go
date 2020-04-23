package nes

import (
	"fmt"
	"strconv"
)
// TODO: determine how to best design this for multiple games. One large file will get unwieldy fast
// TODO: generate minimum action spaces
func (nes *NES)convertScore(address uint16) (convertedVal int) {
	raw := nes.CPU.Memory.Fetch(address)
	convertedVal, _ = strconv.Atoi(fmt.Sprintf("%x", raw))
	return
}
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

	ones := float32(nes.convertScore(0x07FC)) * 10
	hundreds := float32(nes.convertScore(0x07FD)) * 100
	thousands := float32(nes.convertScore(0x07FE)) * 1000
	score = thousands + hundreds + ones
	//log.Printf("score at frame %v: %v", nes.frameCount, score)
	return
}

// TODO: add more games (SMB3, Mario bros and Donkey kong as good starting points)

//Worth noting that the formula for score recovery is close to identical in both Mario Brothers and Castlevania
func (nes *NES) GetMarioBrothersScore() (score float32) {
	/*
	At present, only gets score for player 1, will initially set up the environment to have the game
	only operate in single player mode
	https://datacrystal.romhacking.net/wiki/Mario_Bros.:RAM_map
	 */
	ones := float32(nes.convertScore(0x0095)) * 10
	hundreds := float32(nes.convertScore(0x0096)) * 100
	thousands := float32(nes.convertScore(0x0097)) * 1000
	score = thousands + hundreds + ones
	//log.Printf("score at frame %v: %v", nes.frameCount, score)
	return
}