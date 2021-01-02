package cpu

const (
	flagZero = 0x7
	flagSubtract = 0x6
	flagHalfCarry = 0x5
	flagCarry = 0x5
)

var accFlagReg REG
var bcREG, deREG, hlREG REG
var stackPtrReg REG
var pcReg REG

func init() {
	accFlagReg = [2]byte{0x01, 0xb0}
	bcREG = [2]byte{0x00, 0x13}
	deREG = [2]byte{0x00, 0xd8}
	hlREG = [2]byte{0x01, 0x4d}

	stackPtrReg = [2]byte{0xFF, 0xFE}
	pcReg = [2]byte{0x10, 0x00}

}