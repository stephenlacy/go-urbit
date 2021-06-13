package ames

import (
	"testing"

	"github.com/stevelacy/go-ames/noun"
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
	c1 := "[[25959 1685024616 0] 0 109 29669416873256296 1735289200]"
	r1 := ConstructPoke([]string{"ge", "hood"}, "helm-hi", noun.MakeNoun("ping"))
	if r1.String() != c1 {
		t.Errorf("expected %s got %s", c1, r1)
	}
}

func TestSplitMessage(t *testing.T) {
	n1 := ConstructPoke([]string{"ge", "hood"}, "helm-hi", noun.MakeNoun("ping"))
	c1 := "[5 1 0 83103842581186151537609419784725107274636599623840339663322629]"
	r1 := SplitMessage(5, n1)

	if c1 != r1[0].String() {
		t.Errorf("expected %s got %s", c1, r1)
	}
}

func TestEncodeShutPacket(t *testing.T) {
	n1 := ConstructPoke([]string{"ge", "hood"}, "helm-hi", noun.MakeNoun("ping"))
	c1 := "[[65792 139023796470016] 1 2 0 2506661347556068784616762579621156625882964019302253075174068806145117037309801655194662853082790689960251053528861610387]"
	msg := SplitMessage(5, n1)
	pkt := FragmentToShutPacket(msg[0], 1)
	_, r1 := EncodeShutPacket(pkt, []byte{31}, noun.B(0x10100), noun.B(0x7e7100010100), 1, 2)
	if c1 != r1.String() {
		t.Errorf("expected %s got %s", c1, r1)
	}
}
