package jam

import (
	"math/big"
	"reflect"
	"testing"
)

func TestMat(t *testing.T) {
	r1 := Mat(B(100))
	c1 := [2]*big.Int{B(13), B(6456)}
	if !reflect.DeepEqual(r1[0], c1[0]) && !reflect.DeepEqual(r1[1], c1[1]) {
		t.Errorf("expected %s got %s", c1, r1)
	}

	r2 := Mat(B(100000000))
	c2 := [2]*big.Int{B(37), B(102400000736)}
	if !reflect.DeepEqual(r2[0], c2[0]) && !reflect.DeepEqual(r2[1], c2[1]) {
		t.Errorf("expected %s got %s", c2, r2)
	}

	b1 := B(0)
	b1.SetString("100000000000000000001", 10)
	b2 := B(0)
	b2.SetString("1638400000000000000017280", 10)
	r3 := Mat(b1)
	c3 := [2]*big.Int{B(18), b2}
	if !reflect.DeepEqual(r3[0], c3[0]) && !reflect.DeepEqual(r3[1], c3[1]) {
		t.Errorf("expected %s got %s", c3, r3)
	}

	r4 := Mat(B(0))
	c4 := [2]*big.Int{B(1), B(1)}
	if !reflect.DeepEqual(r4[0], c4[0]) && !reflect.DeepEqual(r4[1], c4[1]) {
		t.Errorf("expected %s got %s", c4, r4)
	}

	r5 := Mat(B(1))
	c5 := [2]*big.Int{B(3), B(6)}
	if !reflect.DeepEqual(r5[0], c5[0]) && !reflect.DeepEqual(r5[1], c5[1]) {
		t.Errorf("expected %s got %s", c5, r5)
	}
}
