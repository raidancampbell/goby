package cpu

// tests the given bit of the given byte
// all bit opcodes alter flags Z01, where Z is the main meaning of the opcodes
func testBit(bit, byt byte) {
	isSet := byt & (1 << bit) > 0
	c.setFlag(flagZero,  isSet)
	c.setFlag(flagSubtract, false)
	c.setFlag(flagHalfCarry, true)
}

var cbTable = map[byte]opcode{
	0x7C: opcb7c,
	0x4F: opcb4f,
	0x11: opcb11,
}

//opcb7c flips bit 7 of the H register
// TODO: verify. why the heck does this have a length of 2?
var opcb7c = opcode{
	length:  2,
	cycles4: 8,
	label:   "BIT 7,H",
	value:   0x7c,
	impl: func() {
		//Z01
		testBit(0x7, c.hlREG[0])
		c.pc++
	},
}

var opcb4f = opcode{
	length:  2,
	cycles4: 8,
	label:   "BIT 1, A",
	value:   0x4f,
	impl: func() {
		testBit(0x1, c.accFlagReg[0])
		c.pc++
	},
}

//opcb11 rotates left on register C through the carry flagg
var opcb11 = opcode{
	length:  2,
	cycles4: 8,
	label:   "RL C",
	value:   0x11,
	impl: func() {
		val := c.bcREG[1]
		var carry byte
		var rot byte
		carry = val >> 7
		rot = (val<<1)&0xFF | carry
		c.bcREG[1] = rot
		c.setFlag(flagZero, rot == 0)
		c.setFlag(flagSubtract, false)
		c.setFlag(flagHalfCarry, false)
		c.setFlag(flagCarry, carry == 1)
		c.pc++
	},
}