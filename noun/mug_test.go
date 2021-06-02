package noun

import (
	"math/big"
	"testing"
)

func TestMug(t *testing.T) {
	a1 := Atom{Value: B(12)}
	c1 := uint32(1850607025)
	res := Mug(a1)
	if res != c1 {
		t.Errorf("expected %d got %d", c1, res)
	}

	a := []interface{}{12, 16, 19, 23}
	b := []interface{}{12, 16, 19, 23}
	n2 := MakeNoun([]interface{}{a, b})
	c2 := uint32(1662846570)
	res2 := Mug(n2)
	if res2 != c2 {
		t.Errorf("expected %d got %d", c2, res2)
	}
}
func TestMuk(t *testing.T) {
	a1 := uint32(3744000282)
	res := Muk(0xb76d5eed, 2, big.NewInt(1501))
	if res != a1 {
		t.Errorf("%d does not match %d", res, a1)
	}
}
