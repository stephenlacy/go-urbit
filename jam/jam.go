package jam

import (
	"fmt"
	"math/big"
)

type MatTupl [2]*big.Int

func B(i int64) *big.Int {
	return big.NewInt(i)
}

// Mat is Jam on atoms
func Mat(arg *big.Int) MatTupl {
	if arg.Cmp(B(0)) == 0 {
		return MatTupl{B(1), B(1)}
	}
	b := int64(arg.BitLen())
	c := int64(len(fmt.Sprintf("%b", b)))
	tup1 := B(b + c + c)

	d1 := 1 << c // 2 ** c
	var d2 int64 = b % (1 << (c - 1))
	d3 := B(0).Lsh(arg, uint(c-1))
	d4 := B(0).Xor(d3, B(d2))
	d5 := B(0).Lsh(d4, uint(len(fmt.Sprintf("%b", d1))))
	tup2 := B(0).Add(d5, B(int64(d1)))

	return MatTupl{tup1, tup2}
}
