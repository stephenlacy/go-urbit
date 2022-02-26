package ames

import (
	"fmt"
	"math/big"
	"net"
	"sync"

	"github.com/stevelacy/go-ames/noun"
	"github.com/stevelacy/go-ames/urcrypt"
)

type Connection struct {
	pubKey       *big.Int
	privKey      [32]byte
	mut          sync.Mutex
	from, to     *big.Int
	bone, num    int
	rConn, lConn *net.UDPConn
	rAddr        *net.UDPAddr
	symKey       []byte
}

// Connect initiates a connection with ames
func (c *Connection) Connect(from string, to string, seed string) ([]byte, error) {

	raddr, err := net.ResolveUDPAddr("udp", zodAddr)
	if err != nil {
		return nil, err
	}
	c.rAddr = raddr

	bSeed := noun.B(0)
	bSeed.SetString(seed, 10)
	c.privKey = SeedToEncKey(bSeed)

	zodPatp, err := noun.Patp2bn("~zod")
	if err != nil {
		return nil, err
	}

	// query to addr on eth
	ethRes, err := Lookup(to)
	if err != nil {
		return nil, err
	}

	// create local listener with random port
	lConn, err := net.ListenUDP("udp", nil)
	if err != nil {
		return nil, err
	}
	c.lConn = lConn

	// read all incoming packets
	go c.handleConn()

	// ping zod

	// query zod on eth
	ethZodRes, err := Lookup("~zod")
	if err != nil {
		return nil, err
	}

	zodPubkey := noun.B(0)
	zodPubkey.SetString(ethZodRes.EncryptionKey, 16)
	var pubZodKeyArr [32]byte
	copy(pubZodKeyArr[:], noun.BigToLittle(zodPubkey))
	zodSymKey := urcrypt.UrcryptEdShar(pubZodKeyArr, c.privKey)

	_, err = MakeConnRequest(
		c,
		[]string{"ge", "hood"},
		"helm-hi",
		noun.MakeNoun("urbit-go"),
		1, // zod num (1)
		6, // zod bone (6)
		zodSymKey,
		c.from,
		zodPatp,
		int64(1), // zod from life (1)
		int64(1), // zod to life (1)
	)

	if err != nil {
		return nil, err
	}

	// ping target

	c.pubKey = noun.B(0)
	c.pubKey.SetString(ethRes.EncryptionKey, 16)

	var pubKeyArr [32]byte
	copy(pubKeyArr[:], noun.BigToLittle(c.pubKey))

	c.symKey = urcrypt.UrcryptEdShar(pubKeyArr, c.privKey)

	c.from, err = noun.Patp2bn(from)
	if err != nil {
		return nil, err
	}

	c.to, err = noun.Patp2bn(to)
	if err != nil {
		return nil, err
	}

	fromLife := int64(1)
	toLife := int64(1)

	// set bone to 1
	c.bone = 1
	c.num = 1

	res, err := MakeConnRequest(
		c,
		[]string{"ge", "hood"},
		"helm-hi",
		noun.MakeNoun("urbit-go"),
		c.num,
		c.bone,
		c.symKey,
		c.from,
		c.to,
		fromLife,
		toLife,
	)

	return res, err
}

// Request sends a mark and noun to a connected ship
func (c *Connection) Request(path []string, mark string, data noun.Noun) ([]byte, error) {
	c.mut.Lock()
	defer c.mut.Unlock()
	var pubKeyArr [32]byte
	copy(pubKeyArr[:], noun.BigToLittle(c.pubKey))

	fromLife := int64(1)
	toLife := int64(1)

	c.num++

	return MakeConnRequest(c, path, mark, data, c.bone, c.num, c.symKey, c.from, c.to, fromLife, toLife)

}

func (c *Connection) handleConn() {
	buffer := make([]byte, 1024)
	n, addr, err := c.lConn.ReadFromUDP(buffer)
	fmt.Println(n, addr, err, buffer)

	/* buf := make([]byte, 0, 4096)
	tmp := make([]byte, 512)
	for {
		a, err := bufio.NewReader(c.lConn).Read(tmp)
		fmt.Println("for read", a, err)
		if err != nil {
			fmt.Println("close:", err)
			if err != io.EOF {
				fmt.Println("conn error:", err)
			}
			break
		}
		buf = append(buf, tmp[:a]...)
	}
	fmt.Println(buf) */
}
