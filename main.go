package main

import (
	"time"
	"strconv"
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
	"math"
)

type Block struct {
	Timestamp    int64
	Data         []byte
	PreviousHash []byte
	Hash         []byte
	Nonce        int
}

type BlockChain struct {
	blocks []*Block
}

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

const targetBits = 24

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	return &ProofOfWork{block, target}
}

func (proofOfWork *ProofOfWork) prepareData(nonce int) []byte {
	block := proofOfWork.block
	return bytes.Join([][]byte{block.PreviousHash,
		block.Data,
		IntToHex(block.Timestamp),
		IntToHex(int64(targetBits)),
		IntToHex(int64(nonce))}, []byte{})
}

func (proofOfWork *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := proofOfWork.prepareData(proofOfWork.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	return hashInt.Cmp(proofOfWork.target) == -1
}

func (proofOfWork *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	maxNonce := math.MaxInt64
	fmt.Printf("Executing proof of work for %s\n", proofOfWork.block.Data)

	for nonce < maxNonce {
		data := proofOfWork.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(proofOfWork.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Printf("Hash %x\n", hash)
	fmt.Printf("End of proof of work\n")
	return nonce, hash[:]

}

func IntToHex(i int64) []byte {
	return []byte(strconv.FormatInt(i, 16))
}

func (blockChain *BlockChain) AddBlock(data string) {
	previousBlock := blockChain.blocks[len(blockChain.blocks)-1]
	newBlock := NewBlock(data, previousBlock.Hash)
	blockChain.blocks = append(blockChain.blocks, newBlock)
}

func (b *Block) setHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PreviousHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func NewBlock(data string, PreviousHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), PreviousHash, []byte{}, 0}
	proofOfWork := NewProofOfWork(block)
	nonce, hash := proofOfWork.Run()
	block.Hash = hash
	block.Nonce = nonce
	fmt.Printf("PoW: %s\n", strconv.FormatBool(proofOfWork.Validate()))

	return block
}

func NewBlockChain() *BlockChain {
	genesisBlock := NewBlock("Genesis", []byte{})
	return &BlockChain{[]*Block{genesisBlock}}
}

func main() {
	blockchain := NewBlockChain()
	blockchain.AddBlock("First Block")
	blockchain.AddBlock("Second Block")

	for _, block := range blockchain.blocks {
		fmt.Printf("Previous hash: %x\n", block.PreviousHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Data: %s\n", block.Data)
	}
}
