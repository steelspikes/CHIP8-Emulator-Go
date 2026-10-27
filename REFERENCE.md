# CHIP-8 Emulator Reference

This project is a CHIP-8 emulator written in Go. It follows the guide
[Guide to making a CHIP-8 emulator](https://tobiasvl.github.io/blog/write-a-chip-8-emulator/)
by Tobias V. Langhoff, which is used as the main reference for this document and for the implementation plan.

## Planned Future Work

- Write a custom CHIP-8 **assembler** that translates assembly mnemonics (for example `LD V0, 5`) into ROM binaries (`0x6005`).
- A disassembler (bytes to mnemonics) may come along with it, as a debugging aid.
- Support for configurable quirks, and possibly SUPER-CHIP / XO-CHIP extensions.

## Specifications

| Component | Description |
|---|---|
| Memory | 4 KB (4096 bytes), all writable. Programs are loaded at `0x200`. The area `0x000`-`0x1FF` is reserved for the interpreter and the font. |
| Display | 64x32 monochrome pixels (128x64 in SUPER-CHIP). |
| Registers | 16 eight-bit general-purpose registers, `V0`-`VF`. `VF` doubles as a flag register. |
| Index register `I` | 16-bit, used to point at memory addresses (only 12 bits are addressable). |
| Program counter `PC` | Points at the current instruction. |
| Stack | Stores 16-bit return addresses for subroutines. The original had about 16 entries. |
| Delay timer | 8-bit, decrements at 60 Hz while above 0. |
| Sound timer | 8-bit, decrements at 60 Hz while above 0. A beep plays while it is above 0. |

## Font

Each character is 4 pixels wide and 5 pixels tall. The font is stored in the reserved area of memory; `0x050`-`0x09F` is the common convention.

| Char | Bytes | Char | Bytes |
|---|---|---|---|
| 0 | `F0 90 90 90 F0` | 8 | `F0 90 F0 90 F0` |
| 1 | `20 60 20 20 70` | 9 | `F0 90 F0 10 F0` |
| 2 | `F0 10 F0 80 F0` | A | `F0 90 F0 90 90` |
| 3 | `F0 10 F0 10 F0` | B | `E0 90 E0 90 E0` |
| 4 | `90 90 F0 10 10` | C | `F0 80 80 80 F0` |
| 5 | `F0 80 F0 10 F0` | D | `E0 90 90 90 E0` |
| 6 | `F0 80 F0 90 F0` | E | `F0 80 F0 80 F0` |
| 7 | `F0 10 20 40 40` | F | `F0 80 F0 80 80` |

## Display

- Sprites are 1 to 15 bytes tall. Each bit is a pixel, with the most significant bit on the left.
- Pixels are drawn with **XOR**: zero bits are transparent, one bits flip the pixel on screen.
- The screen only needs to be redrawn after `00E0` or `DXYN`.

## Timing

- A speed of about **700 instructions per second** works for most programs. It should be configurable.
- The timers always tick at **60 Hz**, independently of the instruction speed.

## Keypad

Original COSMAC VIP layout (4x4):

```
1 2 3 C
4 5 6 D
7 8 9 E
A 0 B F
```

Suggested mapping on a QWERTY keyboard:

```
1 2 3 4
Q W E R
A S D F
Z X C V
```

Scancodes should be used so that non-QWERTY layouts keep working.

## Fetch / Decode / Execute

1. **Fetch:** read two bytes at `PC`, combine them into one 16-bit instruction, and increment `PC` by 2 immediately.
2. **Decode:** use the first nibble to choose the instruction group. Extract these fields up front:
   - `X`: second nibble (register index)
   - `Y`: third nibble (register index)
   - `N`: fourth nibble
   - `NN`: second byte (8-bit immediate)
   - `NNN`: second, third and fourth nibbles (12-bit address)
3. **Execute:** run the instruction. Always read `VX`/`VY` from the registers; never use the raw `X`/`Y` values as data.

## Instruction Set

### System and Control

| Opcode | Name | Behavior |
|---|---|---|
| `0NNN` | Machine routine | Original hardware only. Ignored. |
| `00E0` | Clear screen | Turn all pixels off. |
| `00EE` | Return | Pop an address from the stack into `PC`. |
| `1NNN` | Jump | Set `PC` to `NNN`. |
| `2NNN` | Call | Push the current `PC` to the stack, then set `PC` to `NNN`. |

### Conditional Skips

Each skip increments `PC` by 2 when its condition is true.

| Opcode | Name | Skips when |
|---|---|---|
| `3XNN` | Skip if equal | `VX == NN` |
| `4XNN` | Skip if not equal | `VX != NN` |
| `5XY0` | Skip if registers equal | `VX == VY` |
| `9XY0` | Skip if registers not equal | `VX != VY` |

### Set and Add

| Opcode | Name | Behavior |
|---|---|---|
| `6XNN` | Set | `VX = NN` |
| `7XNN` | Add | `VX += NN`. `VF` is not affected (no carry flag). |

### Logical and Arithmetic (`8XY_`)

| Opcode | Name | Behavior | `VF` |
|---|---|---|---|
| `8XY0` | Set | `VX = VY` | - |
| `8XY1` | Binary OR | `VX = VX OR VY` | - |
| `8XY2` | Binary AND | `VX = VX AND VY` | - |
| `8XY3` | Logical XOR | `VX = VX XOR VY` | - |
| `8XY4` | Add | `VX = VX + VY` | 1 if the result exceeds 255, else 0 |
| `8XY5` | Subtract | `VX = VX - VY` | 1 if `VX >= VY` (no borrow), else 0 |
| `8XY7` | Subtract reversed | `VX = VY - VX` | 1 if `VY >= VX` (no borrow), else 0 |
| `8XY6` | Shift right | Shift `VX` right by 1 bit | The bit shifted out |
| `8XYE` | Shift left | Shift `VX` left by 1 bit | The bit shifted out |

**Quirk (`8XY6` / `8XYE`):** the original COSMAC VIP copies `VY` into `VX` before shifting. CHIP-48 and SUPER-CHIP shift `VX` in place and ignore `Y`. This should be configurable.

### Index and Jump

| Opcode | Name | Behavior |
|---|---|---|
| `ANNN` | Set index | `I = NNN` |
| `BNNN` | Jump with offset | VIP: `PC = NNN + V0`. CHIP-48 / SUPER-CHIP: `PC = XNN + VX`. |

**Quirk (`BNNN`):** when picking a single behavior, use the original VIP one. It should be configurable.

### Random

| Opcode | Name | Behavior |
|---|---|---|
| `CXNN` | Random | `VX = random AND NN` |

### Display

| Opcode | Name | Behavior |
|---|---|---|
| `DXYN` | Draw sprite | Draw an `N`-row sprite read from memory at `I` at coordinates (`VX`, `VY`). `VF` is set to 1 if any lit pixel was turned off, else 0. |

Details:
- The starting coordinates wrap: `x = VX % 64`, `y = VY % 32`.
- The sprite itself is clipped at the screen edges and does not wrap.
- Row `n` is read from `Memory[I + n]`. `I`, `VX` and `VY` are not modified.
- Pixels are XORed onto the screen from the most significant bit to the least significant bit.

### Keyboard Input

| Opcode | Name | Behavior |
|---|---|---|
| `EX9E` | Skip if key pressed | Skip if the key with value `VX` is held down. |
| `EXA1` | Skip if key not pressed | Skip if the key with value `VX` is not held down. |
| `FX0A` | Wait for key | Block until a key is pressed and store its value in `VX`. Timers keep counting while waiting. On the original VIP the key registers when it is released. |

### Timers

| Opcode | Name | Behavior |
|---|---|---|
| `FX07` | Read delay timer | `VX = delay timer` |
| `FX15` | Set delay timer | `delay timer = VX` |
| `FX18` | Set sound timer | `sound timer = VX` |

### Index Arithmetic and Lookup

| Opcode | Name | Behavior |
|---|---|---|
| `FX1E` | Add to index | `I += VX`. On the VIP `VF` is not affected. Some interpreters set `VF` to 1 when `I` overflows past `0x0FFF` (the game *Spacefight 2091!* relies on this). |
| `FX29` | Font character | `I` = address of the font character for the value in `VX` (only the low nibble is used). |

### Binary-Coded Decimal

| Opcode | Name | Behavior |
|---|---|---|
| `FX33` | BCD conversion | Store the hundreds digit of `VX` at `I`, the tens digit at `I+1` and the ones digit at `I+2`. |

### Memory Store and Load

| Opcode | Name | Behavior |
|---|---|---|
| `FX55` | Store registers | Store `V0` through `VX` into memory starting at `I`. |
| `FX65` | Load registers | Load `V0` through `VX` from memory starting at `I`. |

**Quirk (`FX55` / `FX65`):** on the original VIP, `I` is incremented as registers are processed and ends at `I + X + 1`. From CHIP-48 / SUPER-CHIP onward, `I` is left unchanged. The test ROMs expect the modern behavior. This should be configurable.

## Implementation Plan

Following the order recommended by the guide:

1. Memory, registers, stack, timers and font.
2. Load a ROM into memory at `0x200`.
3. Fetch / decode / execute loop.
4. First instructions: `00E0`, `1NNN`, `6XNN`, `7XNN`, `ANNN`, `DXYN`. Verify with the IBM Logo ROM.
5. Subroutines and skips: `00EE`, `2NNN`, `3XNN`, `4XNN`, `5XY0`, `9XY0`.
6. Logical and arithmetic instructions (`8XY_`).
7. Remaining instructions: `BNNN`, `CXNN`, `EX9E`, `EXA1`, `FX__`.
8. Timers at 60 Hz, keypad input and sound.
9. Configurable quirks and instruction speed.
10. Validate with test ROMs (see below).
11. Future: custom assembler.

## Testing

- **IBM Logo:** minimal ROM that exercises the first six instructions.
- **Timendus `chip8-test-suite`:** opcode, flags, quirks, keypad and beep tests.
- **BC_test** and **corax89** test ROMs for additional opcode coverage.

Unknown opcodes should produce a clear error, and a step-by-step mode with register and memory inspection is useful for debugging.

## References

- Tobias V. Langhoff, [Guide to making a CHIP-8 emulator](https://tobiasvl.github.io/blog/write-a-chip-8-emulator/)
