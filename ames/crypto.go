package ames

import (
	"math/big"
)

// DeriveSymetricKey returns a symetric key
func DeriveSymetricKey(pub, priv []byte) []byte {

	return []byte{}
}

// EncryptPacket symetric key, vectors, poayload
func EncryptPacket(symKey []byte, iv []uint32, payload []byte) (*big.Int, []byte) {

	return big.NewInt(0), []byte{}
}

func DencryptPacket(symKey []byte, iv []uint32, siv *big.Int, cypherText []byte) []byte {
	return []byte{}

}
