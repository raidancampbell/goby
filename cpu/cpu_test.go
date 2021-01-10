package cpu

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestREG_toUint16(t *testing.T) {
	tests := []struct {
		name string
		r    REG
		want uint16
	}{
		{"happy",
			REG{0x9F, 0xFF},
		binary.LittleEndian.Uint16([]byte{0xFF,0x9F}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.toUint16(); got != tt.want {
				t.Errorf("toUint16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestREG_roundTrip(t *testing.T) {
	r := REG{0x9F, 0xFF}

	newReg := REG{}
	newReg.fromUint16(r.toUint16())

	assert.Equal(t, r[0], newReg[0])
	assert.Equal(t, r[1], newReg[1])
	assert.Equal(t, uint16(0x9FFF), newReg.toUint16())
}