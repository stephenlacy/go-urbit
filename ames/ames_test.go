package ames

import (
	"reflect"
	"testing"

	"github.com/stevelacy/go-ames/noun"
	"github.com/stevelacy/go-ames/urcrypt"
)

var pName = "~litryl-tadmev"
var pEKey = "5a14c99c533ef2138de811430657957c1cdaabbac4d8c21e8785e62f994c99e7"
var pAKey = "8a6b789427cd0a03efc7f66ed4cc3841223f41c85f7488a75acd7517061b3ba0"

func TestLookup(t *testing.T) {
	res, _ := Lookup(pName)
	if res.EncryptionKey != pEKey {
		t.Errorf("expected %s got %s", pEKey, res.EncryptionKey)
	}
	if res.AuthenticationKey != pAKey {
		t.Errorf("expected %s got %s", pAKey, res.AuthenticationKey)
	}
}

func TestPadLeft(t *testing.T) {
	c1 := "00000000test"
	res := padLeft("test", 12, "0")
	if res != c1 {
		t.Errorf("expected %s got %s", c1, res)
	}
}

func TestConstructPoke(t *testing.T) {
	c1 := "[103 [25959 1685024616 0] 0 109 29669416873256296 1735289200]"
	r1 := ConstructPoke([]string{"ge", "hood"}, "helm-hi", noun.MakeNoun("ping"))
	if r1.String() != c1 {
		t.Errorf("expected %s got %s", c1, r1)
	}
}

func TestSplitMessage(t *testing.T) {
	n1 := ConstructPoke([]string{"ge", "hood"}, "helm-hi", noun.MakeNoun("ping"))
	c1 := "[5 1 0 5446293427400615627168770935011744630350584192948000500175511867329]"
	r1 := SplitMessage(5, n1)

	if c1 != r1[0].String() {
		t.Errorf("expected %s got %s", c1, r1)
	}
}

func TestEncodeShutPacket(t *testing.T) {
	n1 := ConstructPoke([]string{"ge", "hood"}, "helm-hi", noun.MakeNoun("ping"))
	c1 := "[[65792 139023796470016] 1 2 0 70270754126257173429024868609132679736029986866333602336505787903207239222796713666732498636112461130219095130694309422215675]"
	msg := SplitMessage(5, n1)
	pkt := FragmentToShutPacket(msg[0], 1)
	r1, _ := EncodeShutPacket(pkt, []byte{31}, noun.B(0x10100), noun.B(0x7e7100010100), 1, 2)
	if c1 != r1.String() {
		t.Errorf("expected %s got %s", c1, r1)
	}
}

func TestEncodePacket(t *testing.T) {
	c1 := []byte{128, 28, 112, 182, 33, 0, 1, 1, 0, 0, 1, 1, 0, 113, 126, 0, 0, 251, 177, 66, 74, 134, 147, 242, 188, 119, 57, 37, 27, 132, 153, 69, 253, 34, 0, 174, 98, 110, 181, 25, 144, 121, 192, 44, 232, 136, 22, 223, 146, 232, 23, 9, 200, 94, 235, 235, 169, 110, 64, 44, 233, 30, 17, 20, 94, 212, 254, 76, 106}
	n1 := ConstructPoke([]string{"ge", "hood"}, "helm-hi", noun.MakeNoun("ping"))
	msg := SplitMessage(5, n1)
	pkt := FragmentToShutPacket(msg[0], 1)
	r1, _ := EncodeShutPacket(pkt, []byte{31}, noun.B(0x10100), noun.B(0x7e7100010100), 1, 2)

	r2 := EncodePacket(r1)
	if !reflect.DeepEqual(c1, r2) {
		t.Errorf("expected %v got %v", c1, r2)
	}
}

func TestMakeRequest(t *testing.T) {
	seed := "10848742450084393055292986019175834315581274714688967213202092181691497678884554007131544538879740827205367656620731455195976623258197233159818107836502112858706829966037239382390357505"
	c1 := []byte{0, 43, 29, 179, 17, 0, 1, 1, 0, 113, 126, 0, 0, 0, 5, 14, 0, 216, 215, 247, 16, 11, 48, 190, 13, 179, 112, 116, 188, 23, 218, 212, 153, 38, 0, 161, 92, 192, 7, 0, 89, 184, 10, 53, 116, 84, 210, 26, 240, 194, 116, 202, 41, 177, 95, 229, 183, 0, 1, 216, 205, 230, 116, 108, 123, 179, 147, 107, 115, 185, 25, 90, 79}
	bSeed := noun.B(0)
	bSeed.SetString(seed, 10)

	planet := "~litryl-tadmev"
	ethRes, _ := Lookup(planet)

	pubKey := noun.B(0)
	pubKey.SetString(ethRes.EncryptionKey, 16)

	var pubKeyArr [32]byte
	copy(pubKeyArr[:], noun.BigToLittle(pubKey))

	privKey := SeedToEncKey(bSeed)
	symKey := urcrypt.UrcryptEdShar(pubKeyArr, privKey)
	from, _ := noun.Patp2bn("~mister-wicdev-wisryt")
	to, _ := noun.Patp2bn(planet)
	fromLife := int64(1)
	toLife := int64(1)

	res, _ := MakeRequest(
		[]string{"ge", "hood"},
		"helm-hi",
		noun.MakeNoun("urbit-go"),
		1,
		1,
		symKey,
		from,
		to,
		fromLife,
		toLife,
	)
	if !reflect.DeepEqual(c1, res) {
		t.Errorf("expected %v got %v", c1, res)
	}
}
