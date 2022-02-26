package ames

import (
	"testing"
)

func TestConnection(t *testing.T) {

	seed := "10848742450084393055292986019175834315581274714688967213202092181691497678884554007131544538879740827205367656620731455195976623258197233159818107836502112858706829966037239382390357505"
	from := "~mister-wicdev-wisryt"
	to := "~litryl-tadmev"
	// to := "~zod"
	connection := Connection{}
	_, err := connection.Connect(from, to, seed)
	if err != nil {
		t.Error(err)
	}

	/* data := noun.MakeNoun([]interface{}{noun.MakeNoun(0x10100), 0})
	res, err := connection.Request([]string{"ge", "hood"}, "helm-send-hi", data)
	fmt.Println("pkt:", res)
	if err != nil {
		t.Error(err)
	} */
	select {}
}
