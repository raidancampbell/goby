package cpu

type CPU struct {
	pc uint16 // program counter
	sp uint16 // stack pointer
	reg []REG
}