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

var a1 = "~zod"
var a2 = "~litryl"
var a3 = "~litryl-tadmev"
var a4 = "~mister-wicdev-wisryt"
var a5 = "~dabhec-bitrux-lidluc-lidtyv"

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
	g, _ := Clan(gName)
	s, _ := Clan(sName)
	p, _ := Clan(pName)
	m, _ := Clan(mName)
	c, _ := Clan(cName)

	out := [5]string{g, s, p, m, c}
	expected := [5]string{"galaxy", "star", "planet", "moon", "comet"}
	if out != expected {
		t.Errorf("expected: %v, got: %v", expected, out)
	}
}

func TestSein(t *testing.T) {

	g, _ := Sein(gName)
	s, _ := Sein(sName)
	p, _ := Sein(pName)
	m, _ := Sein(mName)
	c, _ := Sein(cName)

	out := [5]*big.Int{g, s, p, m, c}
	expected := [5]*big.Int{B(0), B(255), B(1280), B(1926365281), B(0)}
	if !reflect.DeepEqual(expected, out) {
		t.Errorf("expected: %v, got: %v", expected, out)
	}
}

func TestBN2patp(t *testing.T) {
	// turn the patp to big.Int
	g, _ := Patp2bn(gName)
	s, _ := Patp2bn(sName)
	p, _ := Patp2bn(pName)
	m, _ := Patp2bn(mName)
	c, _ := Patp2bn(cName)

	// take the big.Int and convert to string
	g1, _ := BN2patp(g)
	s1, _ := BN2patp(s)
	p1, _ := BN2patp(p)
	m1, _ := BN2patp(m)
	c1, _ := BN2patp(c)

	out2 := [5]string{g1, s1, p1, m1, c1}
	expected2 := [5]string{gName, sName, pName, mName, cName}

	if out2 != expected2 {
		t.Errorf("expected: %v, got: %v", expected2, out2)
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

func ExamplePatp2bn() {
	bn, _ := Patp2bn(pName)

	patp, _ := BN2patp(bn)
	fmt.Println(patp)
	// Output: ~litryl-tadmev
}

func ExampleClan() {
	// Get the class of the patp
	clan, _ := Clan(sName)
	fmt.Println(clan)
	// Output: star
}

func ExampleSein() {
	// Get the parent of the patp
	clanB, _ := Sein(sName)
	clan, _ := BN2patp(clanB)
	fmt.Println(clan)
	// Output: ~fes
}
