package urcrypt

import (
	"math/big"
	"reflect"
	"testing"
)

func TestBuild(t *testing.T) {
	a := [32]byte{4}
	b := [32]byte{6}
	c1 := []byte{186, 34, 105, 120, 127, 104, 221, 180, 156, 197, 25, 25, 94, 74, 231, 55, 137, 221, 77, 66, 235, 236, 209, 171, 187, 184, 169, 153, 84, 230, 9, 7}

	r1 := UrcryptEdShar(a, b)
	if !reflect.DeepEqual(r1, c1) {
		t.Errorf("expected %s got %s", c1, r1)
	}
}

func TestUrcryptAESSivcEn(t *testing.T) {
	data := [][]byte{{2}}
	c1 := big.NewInt(159)
	c2 := [16]byte{137, 220, 113, 64, 125, 70, 107, 227, 161, 165, 106, 235, 180, 199, 98, 161}

	_, iv, r1 := UrcryptAESSivcEn(big.NewInt(2), data, [64]byte{4})
	if !reflect.DeepEqual(r1, c1) {
		t.Errorf("expected %d got %d", c1, r1)
	}
	if !reflect.DeepEqual(c2, iv) {
		t.Errorf("expected %d got %d", c2, iv)
	}

	data2 := [][]byte{{2}, {1}}
	c3 := big.NewInt(42471)
	c4 := [16]byte{43, 66, 159, 14, 114, 106, 238, 44, 82, 27, 175, 196, 130, 12, 13, 67}

	_, iv2, r2 := UrcryptAESSivcEn(big.NewInt(256), data2, [64]byte{48})
	if !reflect.DeepEqual(r2, c3) {
		t.Errorf("expected %d got %d", c3, r2)
	}
	if !reflect.DeepEqual(c4, iv2) {
		t.Errorf("expected %d got %d", c4, iv2)
	}
}

func TestUrcryptAESSivcDe(t *testing.T) {
	data := [][]byte{{2}}
	c1 := big.NewInt(2)
	c2 := [16]byte{137, 220, 113, 64, 125, 70, 107, 227, 161, 165, 106, 235, 180, 199, 98, 161}
	ba1 := big.NewInt(159)

	_, r1 := UrcryptAESSivcDe(ba1, data, [64]byte{4}, c2)

	if !reflect.DeepEqual(c1, r1) {
		t.Errorf("expected %d got %d", c1, r1)
	}

}
