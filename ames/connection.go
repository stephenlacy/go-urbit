package ames

import (
	"errors"
	"fmt"
	"math/big"
	"net"
	"sync"

	"github.com/stevelacy/go-urbit/noun"
	"github.com/stevelacy/go-urbit/urcrypt"
)

type Connection struct {
	pubKey    *big.Int
	privKey   [32]byte
	mut       sync.Mutex
	from, to  *big.Int
	bone, num int
	conn      *net.UDPConn
	rAddr     *net.UDPAddr
	symKey    []byte
	breach    bool
}

func NewConnection(from, to, seed string) (*Connection, error) {
	conn := Connection{
		breach: true,
	}
	_, err := conn.Connect(from, to, seed)
	return &conn, err
}

// Connect initiates a connection with ames
func (c *Connection) Connect(from string, to string, seed string) ([]byte, error) {

	raddr, err := net.ResolveUDPAddr("udp", zodAddr)
	if err != nil {
		return nil, err
	}
	c.rAddr = raddr

	bSeed, ok := hexSeedToBig(seed)
	if !ok {
		return nil, errors.New("Invalid seed value or encoding provided")
	}
	c.privKey = SeedToEncKey(bSeed)
	c.from, err = noun.Patp2bn(from)
	if err != nil {
		return nil, err
	}

	zodPatp, err := noun.Patp2bn("~zod")
	if err != nil {
		return nil, err
	}

	// query zod on eth
	ethZodRes, err := Lookup("~zod")
	if err != nil {
		return nil, err
	}

	// query to addr on eth
	ethRes, err := Lookup(to)
	if err != nil {
		return nil, err
	}

	// create local listener with random port
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		return nil, err
	}
	c.conn = conn

	// handle all incoming packets
	go c.handleConn()

	// ping zod

	zodPubkey := noun.B(0)
	zodPubkey.SetString(ethZodRes.EncryptionKey, 16)
	var pubZodKeyArr [32]byte
	copy(pubZodKeyArr[:], noun.BigToLittle(zodPubkey))
	zodSymKey := urcrypt.UrcryptEdShar(pubZodKeyArr, c.privKey)

	pkt, err := CreatePacket(
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

	_, err = c.SendPacket(pkt)

	if err != nil {
		return nil, err
	}

	// ping target

	c.pubKey = noun.B(0)
	c.pubKey.SetString(ethRes.EncryptionKey, 16)

	var pubKeyArr [32]byte
	copy(pubKeyArr[:], noun.BigToLittle(c.pubKey))

	c.symKey = urcrypt.UrcryptEdShar(pubKeyArr, c.privKey)

	c.to, err = noun.Patp2bn(to)
	if err != nil {
		return nil, err
	}

	fromLife := int64(1)
	toLife := int64(1)

	// set bone to 1
	c.bone = 1
	c.num = 1

	pkt, err = CreatePacket(
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
	if err != nil {
		return nil, err
	}
	_, err = c.SendPacket(pkt)

	return pkt, err
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

	pkt, err := CreatePacket(path, mark, data, c.num, c.bone, c.symKey, c.from, c.to, fromLife, toLife)
	if err != nil {
		return nil, err
	}
	_, err = c.SendPacket(pkt)
	return pkt, err
}

func (c *Connection) handleConn() {
	buffer := make([]byte, 1024)
	n, addr, err := c.conn.ReadFromUDP(buffer)
	fmt.Println(n, addr, err, buffer)

	/* buf := make([]byte, 0, 4096)
	tmp := make([]byte, 512)
	for {
		a, err := bufio.NewReader(c.conn).Read(tmp)
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

func (c *Connection) SendPacket(pkt []byte) (int, error) {
	return c.conn.WriteToUDP(pkt, c.rAddr)
}
