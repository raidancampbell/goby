package cpu

import "fmt"

// thank goodness macros exist.  don't build this table by hand!
var table = map[byte]opcode{
	0x31: op31,
	0xaf: opaf,
	0x21: op21,
	0x32: op32,
	0xcb: opcb,
	0x20: op20,
	0xfb: opfb,
	0x0e: op0e,
	0x3e: op3e,
	0xe2: ope2,
	0x0c: op0c,
	0x77: op77,
	0xe0: ope0,
	0x11: op11,
	0x1a: op1a,
	0xcd: opcd,
	0x13: op13,
	0x7b: op7b,
	0xfe: opfe,
	0x06: op06,
	0x23: op23,
	0x05: op05,
	0xea: opea,
	0x3d: op3d,
	0x28: op28,
	0x0d: op0d,
	0x2e: op2e,
	0x18: op18,
	0x67: op67,
	0x04: op04,
	0x1e: op1e,
	0xf0: opf0,
	0x1d: op1d,
	0x24: op24,
	0x7c: op7c,
	0x90: op90,
	0x42: op42,
	0x15: op15,
	0x16: op16,
	0x17: op17,
	0xc5: opc5,
	0xc1: opc1,
	0x22: op22,
	0xc9: opc9,
	0xce: opce,
	0x7d: op7d,
	0x78: op78,
	0x86: op86,
	0x50: op50,
}

// verify opcodes
func init() {
	// coarse timing verification
	for b := range table {
		if table[b].cycles4 < 4 || table[b].cycles4%4 != 0 {
			panic(fmt.Sprintf("malformed opcode: %+v", table[b]))
		}
	}
}

type opcode struct {
	length  uint8 // how many bytes long is the instruction
	cycles4 uint8 // 4MHz cycles. all opcodes should be divisible by 4,
	// as that's the clock rate used for executing the opcodes.
	// The 4MHz rate is used in the PPU
	label string // for human readability
	value byte   // what's the machine code value to invoke this instruction.  like 0x00 is a NOP
	impl  func() // the opcode implementation
	// ALL opcodes will change the CPU's program counter register
	// MOST opcodes will change other registers or memory
	// SOME opcodes will read ahead (e.g. opcodes that take more than one byte)
	// so opcodes need good accessibility to registers and memory
	// each instruction must leave the stack pointer in a position for the next instruction to be read
	// if a one byte arg is required, read at sp, then increment once complete
}

//op31 loads the given word into the stack pointer register
var op31 = opcode{
	length:  3,
	cycles4: 12,
	label:   "LD SP, d16",
	value:   0x31,
	impl: func() {
		// no flag changes
		// stack pointer is clobbered, so it doesn't matter
		c.sp = uint16(c.ram[c.sp+1]) | (uint16(c.ram[c.sp+2]) << 8)
	},
}

//opaf XORs the given A register (given by the opcode) with the implicit A register
// this simply clears the A register
var opaf = opcode{
	length:  1,
	cycles4: 4,
	label:   "XOR A",
	value:   0xAF,
	impl: func() {
		c.accFlagReg[0] = 0x00
		c.setFlag(flagZero, true)
		c.setFlag(flagSubtract, false)
		c.setFlag(flagCarry, false)
		c.setFlag(flagHalfCarry, false)
		c.sp++
	},
}

//op21 loads the given word into the HL register
var op21 = opcode{
	length:  3,
	cycles4: 12,
	label:   "LD HL,d16",
	value:   0x21,
	impl: func() {
		// no flag changes
		c.sp++
		c.hlREG[1] = c.ram[c.sp]
		c.sp++
		c.hlREG[0] = c.ram[c.sp]
		c.sp++
	},
}

//op32 loads A into address pointed to by HL.  HL is then decremented
// this is equivalent to loading A into the address pointed to by HL, then decrementing the value at HL
// LDD == LoaD, Decrement
var op32 = opcode{
	length:  1,
	cycles4: 8,
	label:   "LD (HL-),A",
	value:   0x32,
	impl: func() {
		// no flag changes
		c.ram.WriteByte(c.hlREG.toUint16(), c.accFlagReg[0])
		c.sp++
	},
}

//opcb is the prefix instruction to a secondary table of opcodes
var opcb = opcode{
	length:  1,
	cycles4: 1,
	label:   "PREFIX CB",
	value:   0xCB,
	impl: func() {
		c.sp++
		panic("TODO: CB opcode wrapper unimplemented")
	},
}

//op20 jumps to the given address if the Z flag is NOT set
//todo: r8 is relative 8-bit, and SIGNED, so this can likely jump forward or backward
var op20 = opcode{
	length:  2,
	cycles4: 8, //todo: 12 if jump is taken
	label:   "JR NZ,r8",
	value:   0x20,
	impl: func() {
		// no flag changes
		c.sp++
		if c.getFlag(flagZero) {
			relJump := c.ram[c.sp]
			if relJump < 0 {
				c.sp -= uint16(relJump)
			} else {
				c.sp += uint16(relJump)
			}
		}
	},
}

//opfb enables interrupts
var opfb = opcode{
	length:  1,
	cycles4: 4,
	label:   "ei",
	value:   0xFB,
	impl: func() {
		// no flag changes
		c.interruptEnabled = true
	},
}

//op0e loads the immediate 8-bit value d8 into the C register
var op0e = opcode{
	length:  2,
	cycles4: 8,
	label:   "LD C,d8",
	value:   0x0E,
	impl: func() {
		// no flag changes
		c.sp++
		c.bcREG[1] = c.ram[c.sp]
	},
}

//op3e loads the immediate 8-bit value d8 into the A register
var op3e = opcode{
	length:  2,
	cycles4: 8,
	label:   "LD A,d8",
	value:   0x3E,
	impl: func() {
		// no flag changes
		c.sp++
		c.accFlagReg[0] = c.ram[c.sp]
	},
}

//ope2 loads the value from A to the $(FF00 + $C) address
var ope2 = opcode{
	length:  1,
	cycles4: 8,
	label:   "LD (C),A",
	value:   0xE2,
	impl: func() {
		// no flag changes
		c.ram.WriteByte(0xFF00+uint16(c.bcREG[1]), c.accFlagReg[0])
	},
}

//op0c increments the C register
var op0c = opcode{
	length:  1,
	cycles4: 4,
	label:   "INC C",
	value:   0x0C,
	impl: func() {
		//Z0H flags
		c.bcREG[1]++
		c.setFlag(flagZero, c.bcREG[1] == 0x00)
		c.setFlag(flagSubtract, false)

		// take original lower nibble, add one, check if the result is greater than 0xF
		didHalfCarry := (c.bcREG[1]-1)&0xF+(1&0xF) > 0xF
		c.setFlag(flagHalfCarry, didHalfCarry)
	},
}

//op77 loads A into the address pointed to by HL
var op77 = opcode{
	length:  1,
	cycles4: 8,
	label:   "LD (HL), A",
	value:   0x77,
	impl: func() {
		// no flag changes
		c.ram.WriteByte(c.hlREG.toUint16(), c.accFlagReg[0])
	},
}

//ope0 loads A into the $(FF00+a8) address, where a8 is an immediate 8-bit value
var ope0 = opcode{
	length:  2,
	cycles4: 12,
	label:   "LDH (a8), A", // LD ($FF00+a8), A
	value:   0xE0,
	impl: func() {
		// no flag changes
		c.sp++
		c.ram.WriteByte(0xFF00+uint16(c.ram[c.sp]), c.accFlagReg[0])
	},
}

//op11 loads the immediate d16 word into the DE register-tuple
var op11 = opcode{
	length:  3,
	cycles4: 12,
	label:   "LD DE, d16",
	value:   0x11,
	impl: func() {
		// no flag changes
		c.deREG[1] = c.ram[c.sp]
		c.sp++
		c.hlREG[0] = c.ram[c.sp]
		c.sp++
	},
}

//op1a loads the value pointed to by the contents of the double-register DE
var op1a = opcode{
	length:  1,
	cycles4: 8,
	label:   "LD A, (DE)",
	value:   0x1A,
}

//opcd pushes the next address onto the stack, then jump to the given 16 bit address
//i.e. put PC+3 onto the stack, then jump to a16.
//sometime later a RET will pop off the stack and return to PC+3
var opcd = opcode{
	length:  3,
	cycles4: 24,
	label:   "CALL a16",
	value:   0xCD,
}

//op13 increments the word in the DE register-tuple
var op13 = opcode{
	length:  3,
	cycles4: 8,
	label:   "INC DE",
	value:   0x13,
}

//op7b loads the value from register E into register A
var op7b = opcode{
	length:  1,
	cycles4: 4,
	label:   "LD A, E",
	value:   0x7B,
}

//opfe compares A with the given immediate data byte
var opfe = opcode{
	length:  2,
	cycles4: 8,
	label:   "CP d8",
	value:   0xfe,
}

//op06 loads the given immediate data byte into register B
var op06 = opcode{
	length:  2,
	cycles4: 8,
	label:   "LD B, d8",
	value:   0x06,
}

//op23 increments the word in the HL register-tuple
var op23 = opcode{
	length:  1,
	cycles4: 8,
	label:   "INC HL",
	value:   0x23,
}

//op05 decrements the value in register B
var op05 = opcode{
	length:  1,
	cycles4: 4,
	label:   "DEC B",
	value:   0x05,
}

//opea loads the value from register A into the given 16-bit address
var opea = opcode{
	length:  3,
	cycles4: 16,
	label:   "LD (a16), A",
	value:   0xEA,
}

//op3d decrements the contents of register A
var op3d = opcode{
	length:  1,
	cycles4: 4,
	label:   "DEC A",
	value:   0x3d,
}

//op28 jumps to the given relative address if the Z flag is set
var op28 = opcode{
	length:  2,
	cycles4: 8, //12 if jump is taken
	label:   "JR Z, r8",
	value:   0x28,
}

//op0d decrements the contents of register C
var op0d = opcode{
	length:  1,
	cycles4: 4,
	label:   "DEC C",
	value:   0x0D,
}

//op2e loads the given immediate byte into register L
var op2e = opcode{
	length:  2,
	cycles4: 8,
	label:   "LD L, d8",
	value:   0x2e,
}

//op18 jumps to the given relative address
var op18 = opcode{
	length:  2,
	cycles4: 12,
	label:   "JR r8",
	value:   0x18,
}

//op67 loads the value from register A into register H
var op67 = opcode{
	length:  1,
	cycles4: 4,
	label:   "LD H, A",
	value:   0x67,
}

//op04 increments the contents of register B
var op04 = opcode{
	length:  1,
	cycles4: 4,
	label:   "INC B",
	value:   0x04,
}

//op1e loads the given immediate byte into register E
var op1e = opcode{
	length:  2,
	cycles4: 8,
	label:   "LD E, d8",
	value:   0x1E,
}

//opf0 loads the $(FF00+a8) address into A
//TODO: verify if this is the address itself, or the contents of the address
var opf0 = opcode{
	length:  2,
	cycles4: 12,
	label:   "LDH A, (a8)",
	value:   0xF0,
}

//op1d decrements the value in register E
var op1d = opcode{
	length:  1,
	cycles4: 4,
	label:   "DEC E",
	value:   0x05,
}

//op24 increments the value in register H
var op24 = opcode{
	length:  1,
	cycles4: 4,
	label:   "INC H",
	value:   0x24,
}

//op7c loads the contents of register H into register A
var op7c = opcode{
	length:  1,
	cycles4: 4,
	label:   "LD A, H",
	value:   0x7C,
}

//op90 subtracts the contents of register B from register A
var op90 = opcode{
	length:  1,
	cycles4: 4,
	label:   "SUB B",
	value:   0x90,
}

//op42 loads the contents of register D into register B
var op42 = opcode{
	length:  1,
	cycles4: 4,
	label:   "LD B, D",
	value:   0x42,
}

//op15 decrements the value in register D
var op15 = opcode{
	length:  1,
	cycles4: 4,
	label:   "DEC D",
	value:   0x15,
}

//op16 loads the given immediate byte into register D
var op16 = opcode{
	length:  2,
	cycles4: 8,
	label:   "LD D, d8",
	value:   0x16,
}

//op17 rotates register A left through the carry flag
var op17 = opcode{
	length:  1,
	cycles4: 4,
	label:   "RLA",
	value:   0x17,
}

//opc5 pushes the BC register-tuple onto the stack
var opc5 = opcode{
	length:  1,
	cycles4: 16,
	label:   "PUSH BC",
	value:   0xC5,
}

//opc1 pops a word off the stack pointer and places the value into the BC register-tuple
var opc1 = opcode{
	length:  1,
	cycles4: 12,
	label:   "POP BC",
	value:   0xC1,
}

//op22 loads A into address pointed to by HL.  HL is then incremented
var op22 = opcode{
	length:  1,
	cycles4: 8,
	label:   "LD (HL+), A",
	value:   0x22,
}

//opc9 returns to the caller by popping a word off the stack and returning to that address
var opc9 = opcode{
	length:  1,
	cycles4: 16,
	label:   "RET",
	value:   0xC9,
}

//opce adds the given immediate byte to register A, with carry
var opce = opcode{
	length:  2,
	cycles4: 8,
	label:   "ADC A, d8",
	value:   0xCE,
}

//op7d loads the contents of register L into register A
var op7d = opcode{
	length:  1,
	cycles4: 4,
	label:   "LD A, L",
	value:   0x7D,
}

//op7d loads the contents of register B into register A
var op78 = opcode{
	length:  1,
	cycles4: 4,
	label:   "LD A, B",
	value:   0x7D,
}

//op86 adds the value of the location pointed to by the HL register to the A register
var op86 = opcode{
	length:  1,
	cycles4: 8,
	label:   "ADD A, (HL)",
	value:   0x86,
}

//op50 loads the value from register B into register D
var op50 = opcode{
	length:  1,
	cycles4: 4,
	label:   "LD D, B",
	value:   0x50,
}
