package main

import (
	"time"
	"strconv"
	"bytes"
	"crypto/sha256"
	"fmt"
)

type Block struct {
	Timestamp    int64
	Data         []byte
	PreviousHash []byte
	Hash         []byte
}

type BlockChain struct {
	blocks []*Block
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
	block := &Block{time.Now().Unix(), []byte(data), PreviousHash, []byte{}}
	block.setHash()

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
