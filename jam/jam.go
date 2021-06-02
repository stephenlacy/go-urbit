package jam

import (
	"fmt"
	"math/big"
)

type MatTupl [2]*big.Int

type nounMap map[string]int64
type cueNounMap map[int64]Noun

type Noun interface {
	isNoun()
	String() string
}

type Atom struct {
	Value *big.Int
}

type Cell struct {
	Head Noun
	Tail Noun
}

func B(i int64) *big.Int {
	return big.NewInt(i)
}

func (Atom) isNoun() {}
func (a Atom) String() string {
	return a.Value.Text(10)
}
func (Cell) isNoun() {}
func (a Cell) String() string {
	return "[" + a.innerString() + "]"
}
func (a Cell) innerString() string {
	switch t := a.Tail.(type) {
	case Cell:
		return a.Head.String() + " " + t.innerString()
	default:
		return a.Head.String() + " " + a.Tail.String()
	}
}

func cut(start, run int64, b *big.Int) *big.Int {
	b1 := B(0).Rsh(b, uint(start))
	c1 := B(0).Mod(b1, B(0).Lsh(B(1), uint(run)))
	return c1
}

// Mat is Jam on atoms
func Mat(arg *big.Int) MatTupl {
	if arg.Cmp(B(0)) == 0 {
		return MatTupl{B(1), B(1)}
	}
	b := int64(arg.BitLen())
	c := int64(len(fmt.Sprintf("%b", b)))
	tup1 := B(b + c + c)

	d1 := 1 << c // 2 ** c
	var d2 int64 = b % (1 << (c - 1))
	d3 := B(0).Lsh(arg, uint(c-1))
	d4 := B(0).Xor(d3, B(d2))
	d5 := B(0).Lsh(d4, uint(len(fmt.Sprintf("%b", d1))))
	tup2 := B(0).Add(d5, B(int64(d1)))

	return MatTupl{tup1, tup2}
}

func Rub(index int64, b *big.Int) (int64, Atom) {
	var c int64 = 0

	for ; b.Bit(int(index+c)) == 0; c++ {
	}

	if c == 0 {
		return 1, Atom{Value: B(0)}
	}

	d := index + c + 1
	d1 := cut(d, c-1, b)
	e := B(0).Add(d1, B(0).Lsh(B(1), uint(c-1)))
	return c + c + e.Int64(), Atom{Value: cut(d+c-1, e.Int64(), b)}
}

func bConcat(a, b *big.Int) *big.Int {
	l := uint(a.BitLen())
	b2 := B(0).Lsh(b, l)
	a2 := B(0).Xor(a, b2)
	return a2
}

func MakeNoun(arg interface{}) Noun {
	switch t := arg.(type) {
	case int:
		{
			return Atom{Value: B(int64(t))}
		}
	case *big.Int:
		{
			return Atom{Value: t}
		}
	case []interface{}:
		{
			l := len(t)
			if l == 0 {
				return Atom{Value: B(0)}
			}
			if l == 1 {
				return MakeNoun(t[0])
			}
			c := Cell{
				Head: MakeNoun(t[l-2]),
				Tail: MakeNoun(t[l-1]),
			}

			for k := range t[:l-2] {
				c = Cell{
					Head: MakeNoun(t[l-k-3]),
					Tail: c,
				}
			}
			return c
		}
	default:
		return Atom{Value: B(0)}
	}
}

func jamIn(nmap nounMap, n Noun, index int64) (int64, *big.Int) {
	if pIndex, ok := nmap[n.String()]; ok {
		switch t := n.(type) {
		case Atom:
			{
				if t.Value.BitLen() > int(pIndex) {
					d := Mat(t.Value)
					return 1 + d[0].Int64(), B(0).Lsh(d[1], 1)
				}
			}
		}

		d1 := Mat(B(pIndex))
		d2 := B(0).Lsh(d1[1], 2)
		d3 := B(0).Xor(d2, B(3))
		return 2 + d1[0].Int64(), d3
	}

	nmap[n.String()] = index

	switch t := n.(type) {
	case Atom:
		{
			d := Mat(t.Value)
			return 1 + d[0].Int64(), d[1].Lsh(d[1], 1)
		}
	case Cell:
		{
			index += 2
			hidx, d1 := jamIn(nmap, t.Head, index)
			index += hidx
			tidx, d2 := jamIn(nmap, t.Tail, index)
			index += tidx
			d3 := bConcat(d1, d2)
			d4 := B(0).Lsh(d3, 2)
			d5 := B(0).Xor(d4, B(1))
			return index, d5
		}
	}
	return index, B(0)
}

func Jam(n Noun) *big.Int {
	var nmap nounMap = make(nounMap)
	var index int64 = 0

	_, q1 := jamIn(nmap, n, index)
	return q1
}

func cueIn(nmap cueNounMap, b *big.Int, index int64) (int64, Noun) {
	a := b.Bit(int(index))
	// a == 0 > a is an atom
	index1 := index + 1
	if a == 0 {
		i, a1 := Rub(index1, b)
		nmap[index] = a1
		return i + 1, a1
	}

	index2 := index + 2
	a2 := b.Bit(int(index1))
	// when it is a Cell
	if a2 == 0 {
		i1, n1 := cueIn(nmap, b, index2)
		i2, n2 := cueIn(nmap, b, index2+i1)
		cell := Cell{
			Head: n1,
			Tail: n2,
		}
		nmap[index] = cell
		return i1 + i2 + 2, cell
	}

	// when it is a pointer, not atom or cell
	i3, a3 := Rub(index2, b)
	n3 := nmap[a3.Value.Int64()]

	return i3 + 2, n3
}

func Cue(b *big.Int) Noun {
	var nmap cueNounMap = make(cueNounMap)
	var index int64 = 0

	_, q1 := cueIn(nmap, b, index)
	return q1
}
