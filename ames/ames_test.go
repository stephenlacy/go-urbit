package ames

import (
	"testing"
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
