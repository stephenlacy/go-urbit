package ames

import (
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/stevelacy/go-ames/noun"
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
	return noun.MakeNoun([]interface{}{path, 0, "m", mark, data})
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

func EncodeShutPacket(pkt noun.Noun, symKey []byte, from *big.Int, to *big.Int, fromLife, toLife int64) (error, noun.Noun) {
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
		return err, noun.MakeNoun(0)
	}

	rLen := noun.B(int64((cypherText.BitLen()-1)/8 + 1))
	iv2 := make([]byte, 16)
	copy(iv2, ivs[:])

	siv := noun.LittleToBig(iv2)
	rLen2 := noun.B(0).Lsh(rLen, 128)
	cypherText2 := noun.B(0).Lsh(cypherText, 144)
	content := noun.B(0).Xor(siv, noun.B(0).Xor(rLen2, cypherText2))

	res := noun.MakeNoun([]interface{}{[]interface{}{from, to}, fromLife % 16, toLife % 16, 0, content})
	return nil, res
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
