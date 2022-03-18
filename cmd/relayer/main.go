package main

import (
	"bytes"
	"crypto/rand"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/google"
	"math/big"
)

func mustHashMessage(message []byte) *bn256.G1 {
	reader := bytes.NewBuffer(crypto.Keccak256(message))
	_, HM, err := bn256.RandomG1(reader)
	if err != nil {
		panic(err)
	}
	return HM
}

var G2Base *bn256.G2

func init() {
	G2Base = new(bn256.G2)
	exp := big.NewInt(1)
	G2Base.ScalarBaseMult(exp)
}

func hack(sig *bn256.G1, msg []byte) {
	reader := bytes.NewBuffer(crypto.Keccak256(msg))
	k, _, _ := bn256.RandomG1(reader)
	k.ModInverse(k, bn256.Order)
	sig.ScalarMult(sig, k)
	println(hexutil.Encode(sig.Marshal()))
}

func main() {
	msg1 := mustHashMessage([]byte("Hui"))
	msg2 := mustHashMessage([]byte("Pidr"))

	sk1, pk1, _ := bn256.RandomG2(rand.Reader)
	println(hexutil.Encode(sk1.Bytes()))
	sk2, pk2, _ := bn256.RandomG2(rand.Reader)
	println(hexutil.Encode(sk2.Bytes()))

	//publicKeysHash := crypto.Keccak256(append(append([]byte{}, pk1.Marshal()...), pk2.Marshal()...))

	sig11 := new(bn256.G1).ScalarMult(msg1, sk1)
	sig12 := new(bn256.G1).ScalarMult(msg1, sk2)
	sig21 := new(bn256.G1).ScalarMult(msg2, sk1)
	sig22 := new(bn256.G1).ScalarMult(msg2, sk2)

	hack(sig11, []byte("Hui"))

	pkAgg := new(bn256.G2).Add(pk1, pk2)
	sigAgg := new(bn256.G1).Add(sig11, sig12)
	sigAgg = sigAgg.Add(sigAgg, sig21)
	sigAgg = sigAgg.Add(sigAgg, sig22)

	left := bn256.Pair(new(bn256.G1).Add(msg1, msg2), pkAgg).Marshal()
	right := bn256.Pair(sigAgg, G2Base).Marshal()

	if bytes.Equal(left, right) {
		println("OK")
	} else {
		println("NOT OK")
	}

}
