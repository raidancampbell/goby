package cpu

//opcb7c flips bit 7 of the H register
// TODO: verify. why the heck does this have a length of 2?
var opcb7c = opcode{
	length:  2,
	cycles4: 8,
	label:   "BIT 7,H",
	value:   0x7c,
}

var opcb4f = opcode{
	length:  2,
	cycles4: 8,
	label:   "BIT 1, A",
	value:   0x4f,
}

//opcb11 rotates left on register C through the carry flagg
var opcb11 = opcode{
	length:  2,
	cycles4: 8,
	label:   "RL C",
	value:   0x11,
}