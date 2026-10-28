package main

import (
	chip8 "chip8_emulator/chip_8"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	chip8 chip8.CHIP8
}

func (g *Game) ComputeKeysState() {
	KEYS_LIB := []ebiten.Key{
		ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4,
		ebiten.KeyQ, ebiten.KeyW, ebiten.KeyE, ebiten.KeyR,
		ebiten.KeyA, ebiten.KeyS, ebiten.KeyD, ebiten.KeyF,
		ebiten.KeyZ, ebiten.KeyX, ebiten.KeyC, ebiten.KeyV,
	}

	for i, k := range KEYS_LIB {
		state := ebiten.IsKeyPressed(k)

		if g.chip8.WaitingKey && g.chip8.Keys[i] != state {
			fmt.Println("Pressed")
			g.chip8.LastKeyPressed = i
		}

		g.chip8.Keys[i] = state
	}
}

func (g *Game) Update() error {
	g.ComputeKeysState()
	// fmt.Println(g.chip8.Keys)

	g.chip8.DecrementDelayTimer()
	g.chip8.DecrementSoundTimer()
	g.chip8.DoCPUStep()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if g.chip8.Display[y][x] {
				screen.Set(x, y, color.White)
			} else {
				screen.Set(x, y, color.Black)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 64, 32 // Leer del CHIP
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Title")

	game := &Game{}
	game.chip8.Init()

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}

	// chip := chip8.CHIP8{}
	// chip.Init()

	// chip.LoadROM("2-ibm-logo.ch8")

	// for {
	// 	O, X, Y, N, NN, NNN := chip.DecodeCurrentInstruction()

	// 	fmt.Println(O, X, Y, N, NN, NNN)
	// 	time.Sleep(500 * time.Millisecond)

	// 	chip.PCStep()
	// }

}
