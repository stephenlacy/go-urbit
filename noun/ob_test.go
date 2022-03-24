package noun

import (
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"testing"
)

var gName = "~zod"
var gName2 = "~fed"
var sName = "~fipfes"
var pName = "~litryl-tadmev"
var cName = "~libmer-bolnut-somteb-rapheb--fadneb-milsec-lissub-taddef"
var mName = "~dabhec-bitrux-lidluc-lidtyv"
var bName = "101110111100011011100110011"

var pHash = "e0500"
var cHash = "279f2d435959e414ce82d450094437b5"
var mHash = "b51bb67b72d20061"

func TestPatp2hex(t *testing.T) {
	r1, err := Patp2hex(pName)
	if err != nil {
		t.Errorf(err.Error())
	}
	if r1 != pHash {
		t.Errorf("'%s' does not match %s", r1, pHash)
	}
	r2, err := Patp2hex(cName)
	if r2 != cHash {
		t.Errorf("'%s' does not match %s", r2, cHash)
	}
	r3, err := Patp2hex(mName)
	if r3 != mHash {
		t.Errorf("'%s' does not match %s", r3, mHash)
	}
	r4, err := Patp2hex(gName)
	if r4 != "0" {
		t.Errorf("'%s' does not match %s", r4, "0")
	}
	r5, err := Patp2hex(gName2)
	if r5 != "ec" {
		t.Errorf("'%s' does not match %s", r5, "ec")
	}
	r6, err := Patp2hex(sName)
	if r6 != "ffff" {
		t.Errorf("'%s' does not match %s", r6, "ffff")
	}

}

func TestHex2patp(t *testing.T) {
	r1, err := Hex2patp(pHash)
	if err != nil {
		t.Errorf(err.Error())
	}
	if r1 != pName {
		t.Errorf("'%s' does not match %s", r1, pName)
	}
	r2, _ := Hex2patp(cHash)
	if r2 != cName {
		t.Errorf("'%s' does not match %s", r2, cName)
	}
}

func TestMakeAddr(t *testing.T) {
	res := makeAddr(pName)
	if res.Text(2) != bName {
		t.Errorf("'%s' does not match %s", res.Text(2), bName)
	}
}

func TestFynd(t *testing.T) {
	a1 := big.NewInt(918784)
	intr, _ := strconv.ParseInt(bName, 2, 64)
	res := Fynd(big.NewInt(intr), tail)
	if res.Cmp(a1) != 0 {
		t.Errorf("%s does not match %s", res, a1)
	}
}

func TestPatp2sys(t *testing.T) {
	a1 := []string{"lit", "ryl", "tad", "mev"}
	a2 := patp2syls(pName)
	if !reflect.DeepEqual(a1, a2) {
		t.Errorf("%s does not match %s", a2, a1)
	}
}

func TestIsValidPat(t *testing.T) {
	c1 := "~litryl-tadmev"
	c2 := "~zod"
	f1 := "litryl-tadmev"
	f2 := "lit"
	if !isValidPat(c1) {
		t.Errorf("%s should be valid", c1)
	}
	if !isValidPat(c2) {
		t.Errorf("%s should be valid", c2)
	}
	if isValidPat(f1) {
		t.Errorf("%s should not be valid", f1)
	}
	if isValidPat(f2) {
		t.Errorf("%s should not be valid", f2)
	}
}

func TestClan(t *testing.T) {
	g, _ := Clan("~zod")
	s, _ := Clan("~litryl")
	p, _ := Clan("~litryl-tadmev")
	m, _ := Clan("~mister-wicdev-wisryt")
	c, _ := Clan("~dabhec-bitrux-lidluc-lidtyv")

	out := [5]string{g, s, p, m, c}
	expected := [5]string{"galaxy", "star", "planet", "moon", "comet"}
	if out != expected {
		t.Errorf("expected: %v, got: %v", expected, out)
	}
}

func TestSein(t *testing.T) {
	g, _ := Sein("~zod")
	s, _ := Sein("~litryl")
	p, _ := Sein("~litryl-tadmev")
	m, _ := Sein("~mister-wicdev-wisryt")
	c, _ := Sein("~dabhec-bitrux-lidluc-lidtyv")

	out := [5]*big.Int{g, s, p, m, c}
	expected := [5]*big.Int{B(0), B(222), B(1280), B(65792), B(0)}
	if !reflect.DeepEqual(expected, out) {
		t.Errorf("expected: %v, got: %v", expected, out)
	}
}

func ExamplePatp2hex() {
	hexp, _ := Patp2hex("~ben")
	fmt.Println(hexp)
	// Output: 5c
}

func ExampleHex2patp() {
	hexp, _ := Hex2patp("ffff")
	fmt.Println(hexp)
	// Output: ~fipfes
}
