package ames

import (
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/stevelacy/go-ames/noun"
	. "github.com/stevelacy/go-ames/noun"
	"github.com/stevelacy/go-ames/urcrypt"
)

var ethAddr = "0x223c067f8cf28ae173ee5cafea60ca44c335fecb"
var apiAddr = "http://eth-mainnet.urbit.org:8545"
var ethMethod = "0x63fa9a87" // "points"

type LookupResponse struct {
	EncryptionKey     string
	AuthenticationKey string
	Sponsor           string
}

type ETHResponse struct {
	Result string `json:"result"`
}

func Lookup(name string) (LookupResponse, error) {
	hex, err := noun.Patp2hex(name)
	if err != nil {
		return LookupResponse{}, err
	}

	res, err := makeEthRequest(hex)
	if err != nil {
		return LookupResponse{}, err
	}
	// remove 0x prefix then split by 64 chars
	parts := noun.Chunks(res[2:], 64)
	resp := LookupResponse{
		EncryptionKey:     parts[0],
		AuthenticationKey: parts[1],
		Sponsor:           parts[5],
	}

	return resp, err
}

func ConstructPoke(path []string, mark string, data noun.Noun) noun.Noun {
	return noun.MakeNoun([]interface{}{"g", path, 0, "m", mark, data})
}

func SplitMessage(num int, blob noun.Noun) []noun.Noun {
	a := noun.Jam(blob)
	l := a.BitLen()
	if l == 0 {
		return []noun.Noun{noun.MakeNoun([]interface{}{num, 1, 0, blob})}
	}
	l = ((l - 1) >> 13) + 1
	acc := []noun.Noun{}
	for i := 0; i < l; i++ {
		n := noun.MakeNoun([]interface{}{num, l, i, noun.Cut(int64(i<<13), 1<<13, a)})
		acc = append(acc, n)
	}
	return acc
}

func FragmentToShutPacket(frag noun.Noun, bone int) noun.Noun {
	return noun.MakeNoun([]interface{}{bone, noun.Head(frag), 0, noun.Tail(frag)})
}

func EncodeShutPacket(pkt noun.Noun, symKey []byte, from *big.Int, to *big.Int, fromLife, toLife int64) (noun.Noun, error) {
	aVec := [][]byte{
		noun.BigToLittle(from),
		noun.BigToLittle(to),
		noun.BigToLittle(noun.B(fromLife)),
		noun.BigToLittle(noun.B(toLife)),
	}
	jPkt := noun.Jam(pkt)
	kHash := sha512.Sum512(symKey)
	err, ivs, cypherText := urcrypt.UrcryptAESSivcEn(jPkt, aVec, kHash)
	if err != nil {
		return noun.MakeNoun(0), err
	}

	rLen := noun.B(noun.ByteLen(cypherText))
	iv2 := make([]byte, 16)
	copy(iv2, ivs[:])

	siv := noun.LittleToBig(iv2)
	rLen2 := noun.B(0).Lsh(rLen, 128)
	cypherText2 := noun.B(0).Lsh(cypherText, 144)
	content := noun.B(0).Xor(siv, noun.B(0).Xor(rLen2, cypherText2))

	res := noun.MakeNoun([]interface{}{[]interface{}{from, to}, fromLife % 16, toLife % 16, 0, content})
	return res, nil
}

func EncodePacket(encoded Noun) []byte {
	sender, _ := AssertAtom(Head(Head(encoded)))
	receiver, _ := AssertAtom(Tail(Head(encoded)))

	senderTick, _ := AssertAtom(Head(Tail(encoded)))
	receiverTick, _ := AssertAtom(Head(Tail(Tail(encoded))))

	content, _ := AssertAtom(Tail(Tail(Tail(Tail(encoded)))))

	senderSize, senderRank := EncodeShipMetadata(sender)
	receiverSize, receiverRank := EncodeShipMetadata(receiver)

	body := B(0).Xor(
		senderTick.Value,
		B(0).Xor(
			B(0).Lsh(receiverTick.Value, 4),
			B(0).Xor(
				B(0).Lsh(sender.Value, 8),
				B(0).Xor(
					B(0).Lsh(receiver.Value, uint(8*(senderSize+1))),
					B(0).Lsh(content.Value, uint(8*(senderSize+receiverSize+1))),
				),
			),
		),
	)

	checksum := Mug(MakeNoun(body)) & 0xfffff

	header :=
		(0 << 0) ^
			(0 << 3) ^
			(0 << 4) ^
			(senderRank << 7) ^
			(receiverRank << 9) ^
			(checksum << 11) ^
			(1 << 31)

	b2 := B(0).Xor(B(int64(header)), B(0).Lsh(body, 32))
	return BigToLittle(b2)
}

func EncodeShipMetadata(name Noun) (uint32, uint32) {
	a, err := AssertAtom(name)
	if err != nil {
		return 0, 0
	}

	nLen := noun.ByteLen(a.Value)
	if nLen <= 2 {
		return 2, 0
	}
	if nLen <= 4 {
		return 4, 1
	}
	if nLen <= 8 {
		return 8, 2
	}
	return 16, 3
}

func SeedToEncKey(seed *big.Int) [32]byte {
	a1 := Cue(seed)
	a2 := Head(Tail(Tail(a1)))
	a3, _ := AssertAtom(a2)
	b1 := B(0).Rsh(a3.Value, 8+256)
	b2 := BigToLittle(b1)
	var b3 [32]byte
	copy(b3[:], b2)
	return b3
}

func MakeRequest(path []string, mark string, data Noun, num int, bone int, symKey []byte, from, to *big.Int, fromLife, toLife int64) ([]byte, error) {
	poke := ConstructPoke(path, mark, data)
	msg := SplitMessage(num, poke)
	pat := FragmentToShutPacket(msg[0], bone)
	pack, err := EncodeShutPacket(pat, symKey, from, to, fromLife, toLife)
	if err != nil {
		return []byte{}, err
	}
	packet := EncodePacket(pack)
	return packet, nil
}

func makeEthRequest(nameHex string) (string, error) {
	padName := padLeft(nameHex, 64, "0")
	str := `{"jsonrpc":"2.0","id":"0","method":"eth_call","params":[{"to": "` + ethAddr + `", "data": "` + ethMethod + padName + `"}, "latest"]}`

	body := bytes.NewReader([]byte(str))

	req, err := http.NewRequest("POST", apiAddr, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	res := ETHResponse{}
	json.NewDecoder(resp.Body).Decode(&res)

	return res.Result, nil
}

func padLeft(str string, length int, pad string) string {
	p := ""
	for len(p)+len(str) < length {
		p = p + pad
	}
	return p + str
}
