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

var ZOD = "~zod"

type OnPacket func(c *Connection, pkt Packet)

type Ames struct {
	PrivateKey [32]byte
	breach     bool
	Ship       *big.Int
	Life       int64
	RAddr      *net.UDPAddr
	conn       *net.UDPConn
	Peers      map[string]*Peer
	connected  bool
	OnPacket
}

type Connection struct {
	ames      *Ames
	mut       sync.Mutex
	bone, num int
	Peer      *Peer
}

type Peer struct {
	ship        *big.Int
	pubKey      *big.Int
	symKey      []byte
	life        int64
	Connections map[int]*Connection
	nextBone    int
}

type Packet struct {
	Path []string
	Mark string
	Data noun.Noun
	Num  int
}

func NewAmes(seed string, onPacket OnPacket) (*Ames, error) {
	bSeed, ok := hexSeedToBig(seed)
	if !ok {
		return &Ames{}, errors.New("Invalid seed value or encoding provided")
	}
	shp, life, privKey, err := ParseSeed(bSeed)
	ames := &Ames{
		breach:     true,
		Ship:       shp,
		Life:       life.Int64(),
		PrivateKey: privKey,
		Peers:      make(map[string]*Peer),
		OnPacket:   onPacket,
	}
	raddr, err := net.ResolveUDPAddr("udp", zodAddr)
	if err != nil {
		return ames, err
	}

	ames.RAddr = raddr

	// create local listener with random port
	conn, err := net.ListenUDP("udp", nil)
	conn.SetReadBuffer(4096 * 16)
	if err != nil {
		return ames, err
	}
	ames.conn = conn

	// handle all incoming packets
	go ames.handleConn()

	// breach parent
	parent := noun.B(0).Mod(ames.Ship, noun.Bex(noun.B(32)))

	patp, err := noun.BN2patp(parent)
	fmt.Println("connecting to:", patp)
	if err != nil {
		return ames, err
	}

	c, err := ames.Connect(patp)
	if err != nil {
		return ames, err
	}

	// breach moon before connecting
	// this prevents bone and message num conflicts
	if c.ames.breach == true {
		// breach with helm-moon-breach
		_, err := c.Request(
			[]string{"ge", "hood"},
			"helm-moon-breach",
			noun.MakeNoun(c.ames.Ship),
		)
		if err != nil {
			return ames, err
		}
	}
	// delay for zod to catch up
	time.Sleep(2 * time.Second)

	// ping zod after breach
	err = ames.initZod()
	if err != nil {
		return ames, err
	}
	// wait until zod responds as connected
	for ames.connected == false {
		time.Sleep(1 * time.Second)
	}
	return ames, err
}

// initZod initiates the ames vane through zod
func (a *Ames) initZod() error {
	c, err := a.Connect(ZOD)

	pkt, err := c.CreatePacket(
		[]string{"ge", "hood"},
		"helm-hi",
		noun.MakeNoun("urbit-go"),
	)
	if err != nil {
		return err
	}

	_, err = a.SendPacket(pkt)
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
				_, err = a.SendPacket(pkt)
				count++
			}
		}
	}()

	return err
}

func (a *Ames) newPeer(name *big.Int) (Peer, error) {
	peer := Peer{
		ship:        name,
		nextBone:    1,
		Connections: make(map[int]*Connection),
	}

	bnp, err := noun.BN2patp(name)
	if err != nil {
		return peer, err
	}
	// query to addr on eth
	ethRes, err := Lookup(bnp)
	if err != nil {
		return peer, err
	}
	symKey := a.GenerateSymKey(ethRes.EncryptionKey)

	peer.symKey = symKey
	peer.life = ethRes.Life
	return peer, nil
}

func (a *Ames) GetPeer(name *big.Int) (*Peer, error) {
	n, err := noun.BN2patp(name)
	if err != nil {
		return &Peer{}, err
	}
	peer, ok := a.Peers[n]
	if ok {
		return peer, nil
	}
	p, err := a.newPeer(name)
	if err != nil {
		return &Peer{}, err
	}
	a.Peers[n] = &p
	return &p, nil
}

func (a *Ames) Connect(name string) (*Connection, error) {
	p, err := noun.Patp2bn(name)
	if err != nil {
		return &Connection{}, err
	}
	peer, err := a.GetPeer(p)
	return a.GetConnection(p, peer.nextBone)
}

// GetConnection retrieves or creates a Connection
func (a *Ames) GetConnection(p *big.Int, bone int) (*Connection, error) {
	peer, err := a.GetPeer(p)
	if err != nil {
		return &Connection{}, err
	}
	cn, ok := peer.Connections[bone]
	if ok {
		return cn, nil
	}
	c := &Connection{
		Peer: peer,
		bone: bone,
		ames: a,
		num:  1,
	}
	peer.Connections[bone] = c
	peer.nextBone += 4
	return c, nil
}

func (a *Ames) GenerateSymKey(encryptionKey string) []byte {
	theirPubkey := noun.B(0)
	theirPubkey.SetString(encryptionKey, 16)
	var pubTheirKeyArr [32]byte
	copy(pubTheirKeyArr[:], noun.BigToLittle(theirPubkey))
	return urcrypt.UrcryptEdShar(pubTheirKeyArr, a.PrivateKey)
}

// Request sends a mark and data (noun) to a connected ship
func (c *Connection) Request(path []string, mark string, data noun.Noun) ([]byte, error) {
	c.mut.Lock()
	defer c.mut.Unlock()

	pkt, err := c.CreatePacket(path, mark, data)
	c.num++
	if err != nil {
		return nil, err
	}
	_, err = c.ames.SendPacket(pkt)
	return pkt, err
}

func (c *Connection) CreatePacket(path []string, mark string, data noun.Noun) ([]byte, error) {
	poke := ConstructPoke(path, mark, data)
	msg := SplitMessage(c.num, poke)
	// TODO: create each packet vs msg[0]
	pat := FragmentToShutPacket(msg[0], c.bone)
	pack, err := EncodeShutPacket(pat, c.Peer.symKey, c.ames.Ship, c.Peer.ship, c.ames.Life, c.Peer.life)
	if err != nil {
		return []byte{}, err
	}
	packet := EncodePacket(pack)

	return packet, nil
}

func (a *Ames) handleConn() {
	for {
		buf := make([]byte, 0)
		tmp := make([]byte, 4096)

		for {
			ln, err := bufio.NewReader(a.conn).Read(tmp)
			if err != nil {
				fmt.Println("close:", err)
				if err != io.EOF {
					fmt.Println("conn error:", err)
				}
				break
			}
			buf = append(buf, tmp[:ln]...)
			if ln < 4096 {
				break
			}
		}
		packet, c, err := a.ParsePacket(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		// if res is from zod
		if c.Peer.ship.Cmp(noun.B(0)) == 0 && !a.connected {
			// we are now connected
			a.connected = true
		}

		// messages
		if a.OnPacket != nil {
			a.OnPacket(c, packet)
		}
	}
}

// ParsePacket is the reverse of CreatePacket
func (a *Ames) ParsePacket(pkt []byte) (Packet, *Connection, error) {
	from, to, fromTick, toTick, content, err := DecodePacket(pkt)
	if err != nil {
		return Packet{}, &Connection{}, err
	}

	peer, err := a.GetPeer(from)

	if err != nil {
		return Packet{}, &Connection{}, err
	}

	pat, err := DecodeShutPacket(content, peer.symKey, from, to, fromTick, toTick, peer.life, a.Life)
	if err != nil {
		return Packet{}, &Connection{}, err
	}
	msg1, bone, num, isFrag, err := ShutPacketToFragment(pat)
	if err != nil {
		return Packet{}, &Connection{}, err
	}

	conn, err := a.GetConnection(from, bone)
	if isFrag {
		_, msg, err := JoinMessage([]noun.Noun{msg1})
		if err != nil {
			return Packet{}, &Connection{}, err
		}
		path, mark, data, err := DestructPoke(msg)
		packet := Packet{
			Path: path,
			Mark: mark,
			Data: data,
			Num:  num,
		}

		return packet, conn, err
	}
	ack := Packet{
		Path: []string{""},
		Mark: "ack",
		Data: noun.MakeNoun(0),
		Num:  num,
	}
	return ack, conn, nil
}

// SendPacket writes the packet input to the connected target
func (a *Ames) SendPacket(pkt []byte) (int, error) {
	return a.conn.WriteToUDP(pkt, a.RAddr)
}
