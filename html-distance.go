package main

import (
	"fmt"
	"net/http"
	"bytes"
	"io"
	"github.com/mfonda/simhash"
	"golang.org/x/net/html"

)
func main() {
	fmt.Println("Enter first url: ")
	var url1 string
	fmt.Scanf("%s\n", &url1)
	fmt.Println("Enter second url: ")
	var url2 string
	fmt.Scanf("%s\n", &url2)

	s := 2

	resp1, _ := http.Get(url1)
	fmt.Printf("\nFetching %s, Got %d\n\n", url1, resp1.StatusCode)
	resp2, _ := http.Get(url2)
	fmt.Printf("Fetching %s, Got %d\n\n", url2, resp2.StatusCode)

	f1 := Fingerprint(resp1.Body, s)
	f2 := Fingerprint(resp2.Body, s)
	d  := Distance(f1,f2)
	fmt.Printf("Fingerprint1 %64b\n", f1)
	fmt.Printf("Fingerprint2 %64b\n\n", f2)
	fmt.Printf("Feature Distance is %d.\nShingle factor is %d.\nHTML Similarity is %3.2f%%\n",d, s, (1-float32(d)/64)*100)
}

func Fingerprint(r io.Reader, shingle int) uint64 {
	if shingle < 1 {
		shingle = 1
	}
	// collect the features via this cf channel.
	cf := make(chan string, 1000)
	cs := make(chan uint64, 1000)
	v := simhash.Vector{}

	// Tokenize and then Generate Features. .
	go func() {
		defer close(cf)
		z := html.NewTokenizer(r)
		// TODO - export the max token count as an function argument.
		count := 0
		for tt := z.Next(); count < 5000 && tt != html.ErrorToken; tt = z.Next() {
			t := z.Token()
			count++
			genFeatures(&t, cf)
		}

	}()

	// Collect the features.
	go func() {
		defer close(cs)
		a := make([][]byte, shingle)
		for f := <-cf; f != ""; f = <-cf {
			// shingle: generate the k-gram token as a single feature.
			a = append(a[1:], []byte(f))
			// fmt.Printf("%#v\n", a)
			// fmt.Printf("%s\n", bytes.Join(a, []byte(" ")))
			cs <- simhash.NewFeature(bytes.Join(a, []byte(" "))).Sum()
			// cs <- simhash.NewFeature([]byte(f)).Sum()
		}
	}()

	// from the checksum (of feature), append to vector.
	for s := <-cs; s != 0; s = <-cs {
		for i := uint8(0); i < 64; i++ {
			bit := ((s >> i) & 1)
			if bit == 1 {
				v[i]++
			} else {
				v[i]--
			}
		}
	}

	return simhash.Fingerprint(v)

}

func genFeatures(t *html.Token, cf chan<- string) {

	s := ""

	switch t.Type {
	case html.StartTagToken:
		s = "A:" + t.DataAtom.String()
	case html.EndTagToken:
		s = "B:" + t.DataAtom.String()
	case html.SelfClosingTagToken:
		s = "C:" + t.DataAtom.String()
	case html.DoctypeToken:
		s = "D:" + t.DataAtom.String()
	case html.CommentToken:
		s = "E:" + t.DataAtom.String()
	case html.TextToken:
		s = "F:" + t.DataAtom.String()
	case html.ErrorToken:
		s = "Z:" + t.DataAtom.String()
	}
	// fmt.Println(s)
	cf <- s

	for _, attr := range t.Attr {
		switch attr.Key {
		case "class":
			s = "G:" + t.DataAtom.String() + ":" + attr.Key + ":" + attr.Val
		// case "id":
		// 	s = "G:" + t.DataAtom.String() + ":" + attr.Key + ":" + attr.Val
		case "name":
			s = "G:" + t.DataAtom.String() + ":" + attr.Key + ":" + attr.Val
		case "rel":
			s = "G:" + t.DataAtom.String() + ":" + attr.Key + ":" + attr.Val
		default:
			s = "G:" + t.DataAtom.String() + ":" + attr.Key
		}
		// fmt.Println(s)
		cf <- s
	}

	// fmt.Println(s)

}

type Oracle struct {
	fingerprint uint64      // node value.
	nodes       [65]*Oracle // leaf nodes
}

// NewOracle return an oracle that could tell if the fingerprint has been seen or not.
func NewOracle() *Oracle {
	return newNode(0)
}

func newNode(f uint64) *Oracle {
	return &Oracle{fingerprint: f}
}

// Distance return the similarity distance between two fingerprint.
func Distance(a, b uint64) uint8 {
	return simhash.Compare(a, b)
}

// See asks the oracle to see the fingerprint.
func (n *Oracle) See(f uint64) *Oracle {
	d := Distance(n.fingerprint, f)

	if d == 0 {
		// current node with same fingerprint.
		return n
	}

	// the target node is already set,
	if c := n.nodes[d]; c != nil {
		return c.See(f)
	}

	n.nodes[d] = newNode(f)
	return n.nodes[d]
}

// Seen asks the oracle if anything closed to the fingerprint in a range (r) is seen before.
func (n *Oracle) Seen(f uint64, r uint8) bool {
	d := Distance(n.fingerprint, f)
	if d < r {
		return true
	}

	// TODO - should search from d, d-1, d+1, ... until d-r and d+r, for best performance
	for k := d - r; k <= d+r; k++ {
		if k > 64 {
			break
		}
		if c := n.nodes[k]; c != nil {
			if c.Seen(f, r) == true {
				return true
			}
		}
	}
	return false
}
