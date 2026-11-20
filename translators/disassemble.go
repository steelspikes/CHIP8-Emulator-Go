package translators

import (
	"chip8_emulator/utils"
	"errors"
	"fmt"
)

func CreateLabelInTable(location uint16, lablesTable *map[uint16]string) string {
	lineNumber := (location - 0x200) / 2

	name := fmt.Sprintf("%v_%v", "label", lineNumber)
	(*lablesTable)[lineNumber] = name
	return name
}

func GetLabelInTable(location uint16, labelsTable *map[uint16]string) string {
	lineNumber := (location - 0x200) / 2
	return (*labelsTable)[lineNumber]
}

func SelectInstruction(data []byte, offset int, labelsTable *map[uint16]string) (string, error) {
	O, X, Y, N, NN, NNN := utils.DecodeInstruction(data[offset], data[offset+1])

	switch O {
	case 0x0:
		{
			if NN == 0xE0 {
				return "CLS", nil
			}
			if NN == 0xEE {
				return "RET", nil
			}
		}

	case 0x1:
		{
			label := CreateLabelInTable(NNN, labelsTable)
			return fmt.Sprintf("JMP %v", label), nil
		}

	case 0x2:
		{
			label := CreateLabelInTable(NNN, labelsTable)
			return fmt.Sprintf("CALL %v", label), nil
		}
	case 0x3:
		{
			return fmt.Sprintf("SKE r%d, 0x%X", X, NN), nil // Skip if equals  register to value
		}
	case 0x4:
		{
			return fmt.Sprintf("SKNE r%d, 0x%X", X, NN), nil // Skip if not equals value
		}
	case 0x5:
		{
			if N == 0 {
				return fmt.Sprintf("SKE r%d, r%d", X, Y), nil // Skip if VX and VY equals (registers)
			}
		}
	case 0x9:
		{
			if N == 0 {
				return fmt.Sprintf("SKNE r%d, r%d", X, Y), nil // Skip if VX and VY not equals (registers)
			}
		}

	case 0x6:
		{
			return fmt.Sprintf("SET r%v, 0x%X", X, NN), nil // Set register X the value NN
		}
	case 0x7:
		{
			return fmt.Sprintf("ADD r%v, 0x%X", X, NN), nil // Add value to a register
		}

	case 0x8:
		{
			switch N {
			case 0:
				{
					return fmt.Sprintf("SET r%v, r%v", X, Y), nil
				}
			case 1:
				{
					return fmt.Sprintf("OR r%v, r%v", X, Y), nil
				}
			case 2:
				{
					return fmt.Sprintf("AND r%v, r%v", X, Y), nil
				}
			case 3:
				{
					return fmt.Sprintf("XOR r%v, r%v", X, Y), nil
				}
			case 4:
				{
					return fmt.Sprintf("ADD r%v, r%v", X, Y), nil
				}
			case 5:
				{
					return fmt.Sprintf("SUB r%v, r%v", X, Y), nil
				}
			case 7:
				{
					return fmt.Sprintf("SUBR r%v, r%v", X, Y), nil
				}
			case 6:
				{
					return fmt.Sprintf("SHR r%v, r%v", X, Y), nil
				}
			case 0xE:
				{
					return fmt.Sprintf("SHL r%v, r%v", X, Y), nil
				}
			}
		}

	case 0xA:
		{
			return fmt.Sprintf("SET I, 0x%X", NNN), nil
		}

	case 0xB:
		{
			return fmt.Sprintf("JMPO 0x%X", NNN), nil
		}
	case 0xC:
		{
			return fmt.Sprintf("RAND r%d, 0x%X", X, NN), nil
		}
	case 0xD:
		{
			return fmt.Sprintf("DRAW r%d, r%d, %d", X, Y, N), nil
		}
	case 0xE:
		{
			switch NN {
			case 0x9E:
				{
					return fmt.Sprintf("SKPP r%d", X), nil
				}
			case 0xA1:
				{
					return fmt.Sprintf("SKPR r%d", X), nil
				}
			}
		}
	case 0xF:
		{
			switch NN {
			case 0x7:
				{
					return fmt.Sprintf("SET r%d, DT", X), nil
				}
			case 0x15:
				{
					return fmt.Sprintf("SET DT, r%d", X), nil
				}
			case 0x18:
				{
					return fmt.Sprintf("SET ST, r%d", X), nil
				}
			case 0x1E:
				{
					return fmt.Sprintf("ADD I, r%d", X), nil
				}
			case 0x29:
				{
					return fmt.Sprintf("SETF r%d", X), nil
				}
			case 0x33:
				{
					return fmt.Sprintf("BCDC r%d", X), nil
				}
			case 0x55:
				{
					return fmt.Sprintf("STR r%d", X), nil
				}
			case 0x65:
				{
					return fmt.Sprintf("LDR r%d", X), nil
				}
			case 0x0A:
				{
					return fmt.Sprintf("RKEY r%d", X), nil
				}
			}
		}
		// default:
		// 	{
		// 		*instructions = append(*instructions, fmt.Sprintf("**** %v %v %v %v %v %v", O, X, Y, N, NN, NNN))
		// 	}
	}

	return "", errors.New("Invalid opcode")
}

func Disassemble(file string) {
	data, err := utils.ReadFile(file)

	if err != nil {
		fmt.Println(err)
		return
	}

	instructions := []string{}
	labelsTable := map[uint16]string{}

	for i := 0; i < len(data); i += 2 {
		instruction, err := SelectInstruction(data, i, &labelsTable)
		if err != nil {
			instructions = append(instructions, fmt.Sprintf("EXEC 0x%02X%02X", data[i], data[i+1]))
		} else {
			instructions = append(instructions, instruction)
		}
	}

	for lineIdx := range labelsTable {
		name := labelsTable[lineIdx]
		instructions[lineIdx] = fmt.Sprintf("%v: %v", name, instructions[lineIdx])
	}

	for _, v := range instructions {
		fmt.Println(v)
	}

	// fmt.Println(labelsTable)
}
