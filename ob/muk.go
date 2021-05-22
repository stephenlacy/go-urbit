package ob

import (
	"math/big"

	"github.com/twmb/murmur3"
)

var u_256 = big.NewInt(256)
var ux_FF = big.NewInt(0xFF)
var ux_FF00 = big.NewInt(0xFF00)

func Muk(seed uint32, arg *big.Int) *big.Int {
	b := arg.Bytes()
	a := []byte{0, 0}
	bl := len(b)

	if bl > 2 {
		panic("muk: Arg name big.Int should be less than 3 bytes")
	} else if bl == 1 {
		a = []byte{b[0], 0}
	} else if bl == 2 {
		// flip bits - murmur3 expects LittleEndian
		a = []byte{b[1], b[0]}
	}
	u := int64(murmur3.SeedSum32(seed, a))
	return big.NewInt(u)
}
