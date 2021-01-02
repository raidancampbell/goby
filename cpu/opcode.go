package cpu

type Opcode interface {
	Exec() uint16 // execute the opcode and return the new program counter value
	alterFlags()
}

var table = map[byte]opcode {

}

type opcode struct {
	length uint8 // how many bytes long is the instruction
	cycles4 uint8 // 4MHz cycles. all opcodes should be divisible by 4,
	// as that's the clock rate used for executing the opcodes.
	// The 4MHz rate is used in the PPU
	label string // for human readability
	value byte // what's the machine code value to invoke this instruction.  like 0x00 is a NOP
}

//op31 loads the given word into the stack pointer register
var op31 = opcode {
	length: 3,
	cycles4: 12,
	label: "LD SP, d16",
	value: 0x31,
}

//opaf XORs the given A register (given by the opcode) with the implicit A register
// this simply clears the A register
var opaf = opcode {
	length: 1,
	cycles4: 4,
	label: "XOR A",
	value: 0xAF,
}

//op21 loads the given word into the HL register
var op21 = opcode {
	length: 3,
	cycles4: 12,
	label: "LD HL,d16",
	value: 0x21,
}

//op32 loads A into address pointed to by HL.  HL is then decremented
// this is equivalent to loading A into the address pointed to by HL, then decrementing the value at HL
// LDD == LoaD, Decrement
var op32 = opcode {
	length: 1,
	cycles4: 8,
	label: "LD (HL-),A",
	value: 0x32,
}

//opcb is the prefix instruction to a secondary table of opcodes
var opcb = opcode {
	length: 1,
	cycles4: 1,
	label: "PREFIX CB",
	value: 0xCB,
}

//op20 jumps to the given address if the Z flag is set
//todo: r8 is relative 8-bit, and SIGNED, so this can likely jump forward or backward
var op20 = opcode {
	length: 2,
	cycles4: 8,//todo: 12 if jump is taken
	label: "JR NZ,r8",
	value: 0x20,
}

//opfb enables interrupts
var opfb = opcode {
	length: 1,
	cycles4: 4,
	label: "ei",
	value: 0xFB,
}

//op0e loads the immediate 8-bit value d8 into the C register
var op0e = opcode {
	length: 2,
	cycles4: 8,
	label: "LD C,d8",
	value: 0x0E,
}

//op3e loads the immediate 8-bit value d8 into the A register
var op3e = opcode {
	length: 2,
	cycles4: 8,
	label: "LD A,d8",
	value: 0x3E,
}

//ope2 loads the value from A to the $(FF00 + $C) address
//TODO: why is this length 2?
//	some references have it as length1. keeping as length 1
var ope2 = opcode {
	length: 1,
	cycles4: 8,
	label: "LD (C),A",
	value: 0xE2,
}

//op0c increments the C register
var op0c = opcode {
	length: 1,
	cycles4: 4,
	label: "INC C",
	value: 0x0C,
}

//op77 loads A into the address pointed to by HL
var op77 = opcode {
	length: 1,
	cycles4: 8,
	label: "LD (HL), A",
	value: 0x77,
}

//ope0 loads A into the $(FF00+a8) address, where a8 is an immediate 8-bit value
var ope0 = opcode {
	length: 1,
	cycles4: 12,
	label: "LDH (a8), A", // LD ($FF00+a8), A
	value: 0xE0,
}

//op11 loads the immediate d16 word into the DE register-tuple
var op11 = opcode {
	length: 3,
	cycles4: 12,
	label: "LD DE, d16",
	value: 0x11,
}

//op1a loads the value pointed to by the contents of the double-register DE
var op1a = opcode {
	length: 1,
	cycles4: 8,
	label: "LD A, (DE)",
	value: 0x1A,
}

//opcd pushes the next address onto the stack, then jump to the given 16 bit address
//i.e. put PC+3 onto the stack, then jump to a16.
//sometime later a RET will pop off the stack and return to PC+3
var opcd = opcode {
	length: 3,
	cycles4: 24,
	label: "CALL a16",
	value: 0xCD,
}

//op13 increments the word in the DE register-tuple
var op13 = opcode {
	length: 3,
	cycles4: 8,
	label: "INC DE",
	value: 0x13,
}

//op7b loads the value from register E into register A
var op7b = opcode{
	length:  1,
	cycles4: 4,
	label:   "LD A, E",
	value:   0x7B,
}