package ames

import (
	"fmt"
	"os"
	"testing"

	"github.com/stevelacy/go-urbit/noun"
)

func TestConnectionMoon(t *testing.T) {
	seed := os.Getenv("MOON_SEED")
	if seed == "" {
		t.Errorf("Please define env var MOON_SEED")
	}

	onPacket := func(c *Connection, pkt Packet) {
		fmt.Println("ames OnPacket", pkt.Data)
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
	seed := os.Getenv("MOON_SEED")

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
