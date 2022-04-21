package ames

import (
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"strings"

	"github.com/stevelacy/go-urbit/noun"
	. "github.com/stevelacy/go-urbit/noun"
	"github.com/stevelacy/go-urbit/urcrypt"
)

var zodAddr = "zod.urbit.org:13337"
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

func DestructPoke(n noun.Noun) ([]string, string, noun.Noun, error) {
	path, err := destructPath(Head(Tail(n)))
	if err != nil {
		return []string{""}, "", MakeNoun(0), err
	}
	mark, err := AssertAtom(Snag(n, 4))
	if err != nil {
		return []string{""}, "", MakeNoun(0), err
	}
	data := Slag(n, 5)
	return path, string(BigToLittle(mark.Value)), data, nil
}

func destructPath(n noun.Noun) ([]string, error) {
	var strs []string

	cur := n
	cont := true

	for cont {
		switch cur.(type) {
		case Cell:
			hd := Head(cur)
			st, err := AssertAtom(hd)
			if err != nil {
				return []string{}, err
			}
			strs = append(strs, string(BigToLittle(st.Value)))
			cur = Tail(cur)
		default:
			cont = false
			break
		}
	}
	return strs, nil
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

// JoinMessage is the reverse of SplitMessage
func JoinMessage(n []noun.Noun) (int, noun.Noun, error) {
	a, err := AssertAtom(Head(n[0]))
	if err != nil {
		return 0, MakeNoun(0), err
	}
	num := int(a.Value.Int64())
	msg := B(0)

	for k, v := range n {
		// check if the message num is in correct order
		n1, err := AssertAtom(Head(Tail(Tail(v))))
		if err != nil {
			return 0, MakeNoun(0), err
		}
		nIndex := int(n1.Value.Int64())

		if k != nIndex {
			return 0, MakeNoun(0), errors.New("out of order error")
		}

		frag, err := AssertAtom(Tail(Tail(Tail(v))))

		if err != nil {
			return 0, MakeNoun(0), err
		}

		msg = CatLen(msg, frag.Value, uint(k<<13))
	}

	full := Cue(msg)

	return num, full, nil
}

func FragmentToShutPacket(frag noun.Noun, bone int) noun.Noun {
	return noun.MakeNoun([]interface{}{bone, noun.Head(frag), 0, noun.Tail(frag)})
}

func ShutPacketToFragment(n Noun) (Noun, int, int, error) {
	bn1, err := AssertAtom(Head(n))
	if err != nil {
		return noun.MakeNoun(0), 0, 0, nil
	}
	bn2 := BigToLittle(bn1.Value)
	bone := int(B(0).SetBytes(bn2).Int64())

	nm1, err := AssertAtom(Head(Tail(n)))
	if err != nil {
		return noun.MakeNoun(0), 0, 0, nil
	}
	nm2 := BigToLittle(nm1.Value)
	num := int(B(0).SetBytes(nm2).Int64())
	frag := MakeNoun([]interface{}{Head(Tail(n)), Tail(Tail(Tail(n)))})

	return frag, int(bone), int(num), nil
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

	// xor leftshift magic
	header :=
		(0 << 0) ^ // padding
			(0 << 3) ^ // is-ames
			(0 << 4) ^ // version of zero
			(senderRank << 7) ^
			(receiverRank << 9) ^
			(checksum << 11) ^
			(1 << 31)

	b2 := B(0).Xor(B(int64(header)), B(0).Lsh(body, 32))
	return BigToLittle(b2)
}

// EncodeShipMetadata returns size, rank of given name
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

// DecodeShipMetadata returns the bit length of rank
func DecodeShipMetadata(rank byte) int64 {
	if rank == 0 {
		return 16
	}
	if rank == 1 {
		return 32
	}
	if rank == 2 {
		return 64
	}
	return 128
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

func hexSeedToBig(seed string) (*big.Int, bool) {
	b := B(0)
	if strings.Contains(seed, ".") {
		return b.SetString(strings.ReplaceAll(seed, ".", ""), 0)
	}
	return b.SetString(seed, 10)
}

func CreatePacket(path []string, mark string, data Noun, num int, bone int, symKey []byte, from, to *big.Int, fromLife, toLife int64) ([]byte, error) {
	poke := ConstructPoke(path, mark, data)
	msg := SplitMessage(num, poke)
	// TODO: create each packet vs msg[0]
	pat := FragmentToShutPacket(msg[0], bone)
	pack, err := EncodeShutPacket(pat, symKey, from, to, fromLife, toLife)
	if err != nil {
		return []byte{}, err
	}
	packet := EncodePacket(pack)

	return packet, nil
}

// ParsePacket is the reverse of CreatePacket
func ParsePacket(pkt []byte, symKey []byte, fromLife, toLife int64) ([]string, string, Noun, int, int, *big.Int, *big.Int, int64, int64, error) {
	from, to, fromTick, toTick, content, err := DecodePacket(pkt)
	if err != nil {
		return []string{""}, "", noun.MakeNoun(0), 0, 0, B(0), B(0), 0, 0, err
	}

	pat, err := DecodeShutPacket(content, symKey, from, to, fromTick, toTick, fromLife, toLife)
	msg1, bone, num, err := ShutPacketToFragment(pat)
	if err != nil {
		return []string{""}, "", noun.MakeNoun(0), 0, 0, B(0), B(0), 0, 0, err
	}
	_, msg, err := JoinMessage([]noun.Noun{msg1})
	path, mark, data, err := DestructPoke(msg)

	return path, mark, data, num, bone, to, from, fromLife, toLife, nil
}

func DecodeShutPacket(content *big.Int, symKey []byte, from, to, fromTick, toTick *big.Int, fromLife, toLife int64) (noun.Noun, error) {
	siv := noun.Cut(0, 128, content)
	l := noun.Cut(128, 16, content)
	len1 := l.Int64()
	cypherText := noun.Cut(144, len1*8, content)

	aVec := [][]byte{
		noun.BigToLittle(from),
		noun.BigToLittle(to),
		noun.BigToLittle(noun.B(fromLife)),
		noun.BigToLittle(noun.B(toLife)),
	}

	iv1 := noun.BigToLittle(siv)

	var ivs [16]byte
	copy(ivs[:], iv1)

	kHash := sha512.Sum512(symKey)

	decoded, err := urcrypt.UrcryptAESSivcDe(cypherText, aVec, kHash, ivs)

	return noun.Cue(decoded), err
}

func DecodePacket(pkt []byte) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, error) {
	header := pkt[:4]
	body := pkt[4:]

	var checksum uint32

	isAmes := (header[0] >> 3 & 0b1)
	version := header[0] >> 4 & 0b111
	if isAmes != 0 || version != 0 {
		return B(0), B(0), B(0), B(0), B(0), errors.New("error: version invalid")
	}
	senderRank := (header[0] >> 7 & 0b1) ^ ((header[1] >> 0 & 0b1) << 1)
	receiverRank := header[1] >> 1 & 0b11
	checksum = (uint32(header[1]) >> 3 & 0b1111) ^ (uint32(header[2]) << 5) ^ ((uint32(header[3]) >> 0 & 0b1111111) << 13)
	isRelay := (header[3] >> 7 & 0b1) == 0

	lBody := LittleToBig(body)
	if isRelay {
		lBody.Rsh(lBody, uint(6))
	}
	nBody := noun.MakeNoun(lBody)

	if Mug(nBody)&0xfffff != checksum {
		return B(0), B(0), B(0), B(0), B(0), errors.New("error: checksum does not match")
	}

	senderSize := DecodeShipMetadata(senderRank)
	receiverSize := DecodeShipMetadata(receiverRank)

	senderTick := noun.Cut(0, 4, lBody)
	receiverTick := noun.Cut(4, 4, lBody)
	senderValue := noun.Cut(8, senderSize, lBody)
	receiverValue := noun.Cut(8+senderSize, receiverSize, lBody)

	content := B(0)
	content.Rsh(lBody, uint(8+senderSize+receiverSize))

	return senderValue, receiverValue, senderTick, receiverTick, content, nil
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
