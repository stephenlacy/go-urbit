package ames

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"sync"
	"time"

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
	connected bool
	rAddr     *net.UDPAddr
	symKey    []byte
	breach    bool
}

func NewConnection(from, to, seed string) (*Connection, error) {
	conn := Connection{
		connected: false,
		breach:    true,
	}
	err := conn.Connect(from, to, seed)
	return &conn, err
}

// Connect initiates a connection with ames
func (c *Connection) Connect(from string, to string, seed string) error {

	raddr, err := net.ResolveUDPAddr("udp", zodAddr)
	if err != nil {
		return err
	}
	c.rAddr = raddr

	bSeed, ok := hexSeedToBig(seed)
	if !ok {
		return errors.New("Invalid seed value or encoding provided")
	}
	c.privKey = SeedToEncKey(bSeed)
	c.from, err = noun.Patp2bn(from)
	if err != nil {
		return err
	}

	// query to addr on eth
	ethRes, err := Lookup(to)
	if err != nil {
		return err
	}

	// create local listener with random port
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		return err
	}
	c.conn = conn

	// handle all incoming packets
	go c.handleConn()

	// ping target

	c.pubKey = noun.B(0)
	c.pubKey.SetString(ethRes.EncryptionKey, 16)

	var pubKeyArr [32]byte
	copy(pubKeyArr[:], noun.BigToLittle(c.pubKey))

	c.symKey = urcrypt.UrcryptEdShar(pubKeyArr, c.privKey)

	c.to, err = noun.Patp2bn(to)
	if err != nil {
		return err
	}

	fromLife := int64(1)
	toLife := int64(1)

	// set bone to 1
	c.bone = 1
	c.num = 1

	// breach moon before connecting
	// this prevents bone and message num conflicts
	if c.breach == true {
		breachBone := 9 // reserved for breaching
		breachFrom, err := noun.Patp2bn(from)

		if err != nil {
			return err
		}

		// breach with helm-moon-breach
		pkt, err := CreatePacket(
			[]string{"ge", "hood"},
			"helm-moon-breach",
			noun.MakeNoun(breachFrom),
			c.num,
			breachBone,
			c.symKey,
			c.from,
			c.to,
			fromLife,
			toLife,
		)
		if err != nil {
			return err
		}
		_, err = c.SendPacket(pkt)

		if err != nil {
			return err
		}
	}

	// ping zod after breach
	err = c.initZod()
	if err != nil {
		return err
	}
	// wait until zod responds as connected
	for c.connected == false {
		time.Sleep(1 * time.Second)
	}

	// now ping target
	_, err = c.Request(
		[]string{"ge", "hood"},
		"helm-hi",
		noun.MakeNoun("urbit-go"),
	)

	return err
}

func (c *Connection) initZod() error {
	zodPatp, err := noun.Patp2bn("~zod")
	if err != nil {
		return err
	}

	// query zod on eth
	ethZodRes, err := Lookup("~zod")
	if err != nil {
		return err
	}

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
		1, // zod bone (1)
		zodSymKey,
		c.from,
		zodPatp,
		int64(1), // zod from life (1)
		int64(6), // zod to life (6)
	)
	if err != nil {
		return err
	}

	_, err = c.SendPacket(pkt)
	if err != nil {
		return err
	}

	// ping zod every 25 seconds
	ticker := time.NewTicker(25 * time.Second)
	count := 1
	go func() {
		for {
			select {
			case <-ticker.C:
				_, err = c.SendPacket(pkt)
				count++
				fmt.Println("zod ", count)
			}
		}
	}()

	return err
}

// Request sends a mark and noun to a connected ship
func (c *Connection) Request(path []string, mark string, data noun.Noun) ([]byte, error) {
	c.mut.Lock()
	defer c.mut.Unlock()
	var pubKeyArr [32]byte
	copy(pubKeyArr[:], noun.BigToLittle(c.pubKey))

	fromLife := int64(1)
	toLife := int64(1)

	pkt, err := CreatePacket(path, mark, data, c.num, c.bone, c.symKey, c.from, c.to, fromLife, toLife)
	c.num++
	if err != nil {
		return nil, err
	}
	_, err = c.SendPacket(pkt)
	return pkt, err
}

func (c *Connection) handleConn() {
	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 512)
	for {
		a, err := bufio.NewReader(c.conn).Read(tmp)
		fmt.Println("for read", tmp, err)
		// TODO check for response from zod
		if err != nil {
			fmt.Println("close:", err)
			if err != io.EOF {
				fmt.Println("conn error:", err)
			}
			break
		}
		buf = append(buf, tmp[:a]...)
		fmt.Println(buf)
	}
	fmt.Println(buf)
}

func (c *Connection) SendPacket(pkt []byte) (int, error) {
	return c.conn.WriteToUDP(pkt, c.rAddr)
}
