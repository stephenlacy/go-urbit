package argon2u

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	ship := 918784
	res, _ := Hash(HashOptions{
		Pass:        []byte("test"),
		Salt:        []byte(fmt.Sprintf("urbitkeygen%d", ship)),
		Type:        Argon2u,
		HashLen:     32,
		Parallelism: 4,
		Mem:         512000,
		Time:        1,
	})
	h1 := "149d312a05e3d89de2c4f24e179f92e4a76feb22f9e30b2109d1392285ee3c52"
	if res.HashHex != h1 {
		t.Errorf("expected: %s got: %s", h1, res.HashHex)
	}
}
