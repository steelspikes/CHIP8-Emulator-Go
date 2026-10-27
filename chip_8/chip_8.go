package chip8

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type CHIP8 struct {
	Memory     [4096]byte
	Display    [32][64]bool
	PC         uint16 // Solo se pueden direccionar 12 bits pero no hay dato de 12 bits y se toma de 16 bits
	I          uint16
	Stack      []uint16
	DelayTimer uint8
	SoundTimer uint8
	R          [16]byte
}

func (c *CHIP8) AddToStack(v uint16) {
	c.Stack = append(c.Stack, v)
}

func (c *CHIP8) PopStack() (uint16, error) {
	if len(c.Stack) == 0 {
		return 0, errors.New("Empty Stack")
	}

	removed := c.Stack[len(c.Stack)-1]

	c.Stack = c.Stack[:len(c.Stack)-1]

	return removed, nil
}

func (c *CHIP8) PCStep() {
	c.PC += 2
}

func ReadFile(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join("bin", name))
}

func (c *CHIP8) LoadROM(name string) {
	data, err := ReadFile(name)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, b := range data {
		c.Memory[0x200+i] = b
	}
}

func (c *CHIP8) LoadFont() {
	data, err := ReadFile("font.bin")
	if err != nil {
		fmt.Println(err)
		return
	}

	copy(c.Memory[0x050:], data)
}

func (c *CHIP8) DecodeCurrentInstruction() (uint8, uint8, uint8, uint8, uint8, uint16) {
	i1 := c.Memory[c.PC]
	i2 := c.Memory[c.PC+1]

	O := i1 >> 4 // Recorrer 4 bits para mantener los 4 bits altos
	X := i1 & 0b00001111
	Y := i2 >> 4
	N := i2 & 0b00001111
	NN := i2
	NNN := uint16(X) << 8
	NNN = NNN | uint16(i2)

	return O, X, Y, N, NN, NNN
}

func (c *CHIP8) ClearScreen() {
	c.Display = [32][64]bool{}
}

func (c *CHIP8) ExecuteInstruction() {
	O, X, Y, N, NN, NNN := c.DecodeCurrentInstruction()
	switch O {
	case 0x0:
		{
			if NN == 0xE0 {
				c.ClearScreen()
			}
			if NN == 0xEE {
				saved_counter, err := c.PopStack()

				if err != nil {
					panic(err)
				}

				c.PC = saved_counter
			}
		}
	case 0x1:
		{
			c.PC = NNN
		}
	case 0x2:
		{
			c.AddToStack(c.PC + 2)
			c.PC = NNN
		}
	case 0x3:
		{
			if c.R[X] == NN {
				c.PCStep()
			}
		}
	case 0x4:
		{
			if c.R[X] != NN {
				c.PCStep()
			}
		}
	case 0x5:
		{
			if c.R[X] == c.R[Y] && N == 0 {
				c.PCStep()
			}
		}
	case 0x9:
		{
			if c.R[X] != c.R[Y] && N == 0 {
				c.PCStep()
			}
		}
	case 0x6:
		{
			c.R[X] = NN
		}
	case 0x7:
		{
			c.R[X] += NN
		}
	case 0x8:
		{
			switch N {
			case 0:
				{
					c.R[X] = c.R[Y]
				}
			case 1:
				{
					c.R[X] = c.R[X] | c.R[Y]
				}
			case 2:
				{
					c.R[X] = c.R[X] & c.R[Y]
				}
			case 3:
				{
					c.R[X] = c.R[X] ^ c.R[Y]
				}
			case 4:
				{
					c.R[X] = c.R[X] + c.R[Y]
				}
			case 5:
				{
					c.R[X] = c.R[X] - c.R[Y]
				}
			case 7:
				{
					c.R[X] = c.R[Y] ^ c.R[X]
				}
			}
		}
	case 0xA:
		{
			c.I = NNN
		}
	case 0xD:
		{
			vX := c.R[X] & 63
			vY := c.R[Y] & 31

			fmt.Println(vX, vY)

			c.R[0xF] = 0

			n := uint8(0)
			for n < N {
				sprite_index := c.I + uint16(n)
				sprite := c.Memory[sprite_index]

				fmt.Println(sprite)

				i := uint8(0)
				for i < 8 {
					posX := vX + i
					posY := vY + n

					if posX > 63 || posY > 31 {
						break
					}

					sprite_val := (sprite>>(7-i))&1 == 1

					screen_val := c.Display[vY+n][vX+i]

					if sprite_val && screen_val {
						c.Display[vY+n][vX+i] = false
						c.R[0xF] = 1
					}
					if sprite_val && !(screen_val) {
						c.Display[vY+n][vX+i] = true
					}

					i++
				}

				n++
			}

		}
	}

	alreadyIncrementedPC := O == 1 || O == 2 || (O == 0 && X == 0 && NN == 0xEE)

	if !alreadyIncrementedPC {
		c.PCStep()
	}

}

func (c *CHIP8) DoCPUStep() {
	c.ExecuteInstruction()
}

func (c *CHIP8) Init() {
	c.DelayTimer = 0xFF
	c.SoundTimer = 0xFF
	c.PC = 0x200
	// c.Display[10][10] = true

	c.LoadFont()
	// c.LoadROM("1-chip8-logo.ch8")
	c.LoadROM("2-ibm-logo.ch8")
	// c.LoadROM("test_opcode.ch8")
}
