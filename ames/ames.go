package ames

import "github.com/stevelacy/go-ames/ob"

var ethAddr = "0x223c067f8cf28ae173ee5cafea60ca44c335fecb"

func Lookup(name string) (string, error) {
	hex, err := ob.Patp2hex(name)

	return hex, err
}
