package main

import (
	chip8 "chip8_emulator/chip_8"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	chip8 chip8.CHIP8
}

func (g *Game) Update() error {
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
