package ames

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stevelacy/go-urbit/noun"
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

func TestDestructPath(t *testing.T) {
	a := []string{"one", "one", "two"}
	nPath := noun.MakeNoun(a)
	path, err := destructPath(nPath)

	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(path, a) {
		t.Errorf("expected %v got %v", a, path)
	}
}

/* func TestJoinMessage(t *testing.T) {
	num := 11
	poke := ConstructPoke([]string{"path"}, "mark", noun.MakeNoun(noun.B(0).Exp(noun.B(2), noun.B(7000), nil)))
	a := SplitMessage(num, poke)

	n, err := JoinMessage(a)
	fmt.Println(n, err)
} */

func TestShutPacketToFragment(t *testing.T) {
	num := 11
	bone := 9
	poke := ConstructPoke([]string{"path"}, "mark", noun.MakeNoun("data"))
	msg := SplitMessage(num, poke)
	pat := FragmentToShutPacket(msg[0], bone)
	b, n, isFrag, res, err := ShutPacketToMeat(pat)
	if err != nil {
		t.Error(err)
	}
	if b != bone {
		t.Errorf("expected %v got %v", bone, b)
	}
	if n != num {
		t.Errorf("expected %v got %v", n, num)
	}
	if !isFrag {
		t.Errorf("expected %v got %v", true, isFrag)
	}
	e1 := "[1 0 1139440589747613851334439309220300165109479259885505]"
	if res.String() != e1 {
		t.Errorf("expected %v got %v", e1, res.String())
	}
}

func TestDecodePacket(t *testing.T) {
	n1 := []byte{128, 28, 112, 182, 33, 0, 1, 1, 0, 0, 1, 1, 0, 113, 126, 0, 0, 251, 177, 66, 74, 134, 147, 242, 188, 119, 57, 37, 27, 132, 153, 69, 253, 34, 0, 174, 98, 110, 181, 25, 144, 121, 192, 44, 232, 136, 22, 223, 146, 232, 23, 9, 200, 94, 235, 235, 169, 110, 64, 44, 233, 30, 17, 20, 94, 212, 254, 76, 106}
	from, to, _, _, _, err := DecodePacket(n1)
	if err != nil {
		t.Error(err)
	}

	if noun.B(0x10100).Cmp(from) != 0 {
		t.Errorf("expected %v got %v", noun.B(0x10100), from)
	}
	if noun.B(0x7e7100010100).Cmp(to) != 0 {
		t.Errorf("expected %v got %v", noun.B(0x7e7100010100), to)
	}
}

func TestDecodeShutPacket(t *testing.T) {
	n1 := []byte{128, 28, 112, 182, 33, 0, 1, 1, 0, 0, 1, 1, 0, 113, 126, 0, 0, 251, 177, 66, 74, 134, 147, 242, 188, 119, 57, 37, 27, 132, 153, 69, 253, 34, 0, 174, 98, 110, 181, 25, 144, 121, 192, 44, 232, 136, 22, 223, 146, 232, 23, 9, 200, 94, 235, 235, 169, 110, 64, 44, 233, 30, 17, 20, 94, 212, 254, 76, 106}
	from, to, fromTick, toTick, content, err := DecodePacket(n1)

	fromLife := int64(1)
	toLife := int64(2)

	symKey := []byte{31}

	pkt, err := DecodeShutPacket(content, symKey, from, to, fromTick, toTick, fromLife, toLife)
	if err != nil {
		t.Error(err)
	}
	e1 := "[1 5 0 1 0 5446293427400615627168770935011744630350584192948000500175511867329]"
	if pkt.String() != e1 {
		t.Errorf("expected %v got %v", e1, pkt.String())
	}
}

func ExampleLookup() {
	res, err := Lookup("~zod")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.EncryptionKey)
}
