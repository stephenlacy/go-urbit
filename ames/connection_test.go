package ames

import (
	"fmt"
	"testing"

	"github.com/stevelacy/go-urbit/noun"
)

func TestConnectionMoon(t *testing.T) {

	seed := "0x5.f374.9e59.1ea2.fdfd.8165.5c0d.0e2c.0c5a.41b5.bb6b.8962.ab31.0a1c.c221.885d.b876.6f65.4a12.dbdc.645e.531e.5194.8d4d.4251.ad01.d96a.6f4a.c871.fb28.3a5a.d858.c9e1.080f.000e.3c3c.13ba.8007.0280.7a01"

	onPacket := func(c *Connection, pkt Packet) {
		fmt.Println("ames OnMessage", pkt.Data)
	}
	ames, err := NewAmes(seed, onPacket)
	if err != nil {
		t.Error(err)
	}

	to := "~litryl-tadmev"

	c1, err := ames.Connect(to)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("bone:", c1.bone)

	_, err = c1.Request([]string{"ge", "hood"}, "helm-hi", noun.MakeNoun("yay!"))

	if err != nil {
		t.Error(err)
	}

	c2, err := ames.Connect(to)
	if err != nil {
		t.Error(err)
	}

	_, err = c2.Request([]string{"ge", "hood"}, "helm-hi", noun.MakeNoun("second message"))
}

func ExampleNewAmes() {
	// Easiest way to connect with defaults
	seed := "0x5.f374.9e59.1ea2.fdfd.8165.5c0d.0e2c.0c5a.41b5.bb6b.8962.ab31.0a1c.c221.885d.b876.6f65.4a12.dbdc.645e.531e.5194.8d4d.4251.ad01.d96a.6f4a.c871.fb28.3a5a.d858.c9e1.080f.000e.3c3c.13ba.8007.0280.7a01"

	ames, err := NewAmes(seed, nil)
	if err != nil {
		panic(err)
	}
	to := "~litryl-tadmev"
	conn, err := ames.Connect(to)
	if err != nil {
		fmt.Println(err)
	}
	conn.Request([]string{"ge", "hood"}, "helm-hi", noun.MakeNoun("message here"))
}
