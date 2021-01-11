package cpu

import (
	"encoding/binary"
	"fmt"
	"github.com/raidancampbell/goby/mem"
)

type REG [2]byte

//toUint16 renders the given register tuple as a big endian uint16.
//
// registers are big endian, but loads from ROM are little endian
// this is done to keep memory access simple
// e.g. H: 9F, L: FF yields uint16 9FFF, and can be directly accessed in memory that way
func (r *REG) toUint16() uint16 {
	return binary.BigEndian.Uint16([]byte{r[0], r[1]})
}

//fromUint16 sets the register tuple as the given big endian uint16
//see toUint16 for logic on why big endian was chosen
func (r *REG) fromUint16(u uint16) {
	var tmp = make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, u)
	r[0] = tmp[0]
	r[1] = tmp[1]
}

// CPU is the main brains of the operation
// it executes the opcodes and writes to memory through the memory controller
type CPU struct {
	pc                  uint16 // program counter
	sp                  uint16 // stack pointer
	accFlagReg          REG    // AF
	bcREG, deREG, hlREG REG    // BC, DE, HL
	// for BC/DE/HL, the first register is considered bits 8-15, and the second is bits 0-7
	// e.g. 9FFF is stored as (9F FF) in registers, whereas in ROM it is (FF 9F)
	ram                 mem.RAM
	interruptEnabled    bool
}

//GetRAM returns a pointer to the memory.  This is used for loading a cartridge
//TODO: invert responsibility here. pass the cartridge to the memory controller
func GetRAM() *mem.RAM {
	return &c.ram
}

//Run begins reading the memory and executing opcodes
// the bootrom isn't necessary assuming the program counter begins at 0x0100
func Run() {
	for i := 0; i < 90000; i++ {
		opByte := c.ram.ReadByte(c.pc)
		newOp, ok := table[opByte]
		fmt.Printf("executing opcode %x at location %x, execution number %v\t%s\n", opByte, c.pc, i, newOp.label)
		if !ok {
			panic(fmt.Sprintf("unable to find opcode %x", c.ram.ReadByte(c.pc)))
		}
		newOp.impl()
	}
}

//setFlag sets the given bit in the flag register to the given value (i.e. setFlag can clear a bit)
func (c *CPU) setFlag(flag uint8, isSet bool) {
	if isSet {
		c.accFlagReg[1] |= 1 << flag
	} else {
		var mask byte = ^(1 << flag)
		c.accFlagReg[1] &= mask
	}
}

//getFlag returns whether the given flag is set in the flag register
func (c *CPU) getFlag(flag uint8) bool {
	return (c.accFlagReg[1] & (1 << flag)) > 0
}

//popWord pops a little-endian uint16 off the stack
//TODO: shift this to big endian
func (c *CPU) popWord() uint16 {
	val := uint16(c.ram.ReadByte(c.sp)) + (uint16(c.ram.ReadByte(c.sp+1)) << 8)
	fmt.Printf("popped %x from stack addr %x\n", val, c.sp)
	c.sp += 2
	return val
}

//pushWord pushes a given word onto the stack.  The input word is assumed to be little endian
//TODO: should it be big endian? see popWord
func (c *CPU) pushWord(word uint16) {
	high := byte(word >> 8)
	low := byte(word & 0x00FF)
	c.pushBytes(low, high)
	fmt.Printf("pushed %x to stack addr %x\n", word, c.sp)
}

func (c *CPU) pushBytes(low, high byte) {
	c.sp--
	c.ram.WriteByte(c.sp, high)
	c.sp--
	c.ram.WriteByte(c.sp, low)
}

const (
	flagZero      = 0x7 //Z
	flagSubtract  = 0x6 //N
	flagHalfCarry = 0x5 //H
	flagCarry     = 0x5 //C
)

var c *CPU

func init() {
	c = &CPU{}
	c.accFlagReg = [2]byte{0x01, 0xb0}
	c.bcREG = [2]byte{0x00, 0x13}
	c.deREG = [2]byte{0x00, 0xd8}
	c.hlREG = [2]byte{0x01, 0x4d}
	c.pc = 0x0100
	c.sp = 0xFFFE
}

//InitPCForBootrom resets the program counter to 0, indicating that the bootrom should execute
// by default the program counter is initialized to 0x0100, the beginning of the cartridge ROM
func InitPCForBootrom() {
	c.pc = 0x0000
}
