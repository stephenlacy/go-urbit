package ames

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/stevelacy/go-ames/ob"
)

var ethAddr = "0x223c067f8cf28ae173ee5cafea60ca44c335fecb"
var apiAddr = "http://eth-mainnet.urbit.org:8545"
var ethMethod = "0x63fa9a87" // "points"

type LookupResponse struct {
	EncryptionKey    string
	EncryptionSecret string
	Sponsor          string
}

type ETHResponse struct {
	Result string `json:"result"`
}

func Lookup(name string) (LookupResponse, error) {
	hex, err := ob.Patp2hex(name)
	if err != nil {
		return LookupResponse{}, err
	}

	res, err := makeEthRequest(hex)
	if err != nil {
		return LookupResponse{}, err
	}
	// remove 0x prefix then split by 64 chars
	parts := ob.Chunks(res[2:], 64)
	resp := LookupResponse{
		EncryptionKey:    parts[0],
		EncryptionSecret: parts[1],
		Sponsor:          parts[5],
	}

	return resp, err
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
