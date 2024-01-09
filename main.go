package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Block struct {
	Pos       int
	Data      BookCheckOut
	TimeStamp string
	Hash      string
	PrevHash  string
}

type BookCheckOut struct {
	BookID       string `json:"book_id"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis    bool   `json:"is_genesis"`
}

type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	PublishDate string `json:"publish_date"`
	ISBN        string `json:"isbn"`
}

type BlockChain struct {
	blocks []*Block
}

var Blockchain *BlockChain

func (b *Block) GenerateHash() {
	bytes, _ := json.Marshal(b.Data)
	data := string(b.Pos) + b.TimeStamp + string(bytes) + b.PrevHash
	hash := sha256.New()
	hash.Write([]byte(data))
	b.Hash = hex.EncodeToString(hash.Sum(nil))
}

func CreateBlock(prevBlock *Block, checkoutItem BookCheckOut) *Block {
	block := &Block{}
	block.Pos = prevBlock.Pos + 1
	block.TimeStamp = time.Now().String()
	block.PrevHash = prevBlock.Hash
	block.GenerateHash()

	return block
}

func (b *BlockChain) AddBlock(data BookCheckOut) {
	prevBlock := b.blocks[len(b.blocks)-1]

	block := CreateBlock(prevBlock, data)

	if validBlock(block, prevBlock) {
		b.blocks = append(b.blocks, block)
	}
}

func validBlock(block, prevBlock *Block) bool {
	if prevBlock.Hash != block.PrevHash {
		return false
	}
	if !block.ValidateHash(block.Hash) {
		return false
	}
	if prevBlock.Pos+1 != block.Pos {
		return false
	}

	return true
}

func (b *Block) ValidateHash(hash string) bool {
	b.GenerateHash()
	if b.Hash != hash {
		return false
	}

	return true
}

func writeBlock(c *gin.Context) {
	var checkOutItem BookCheckOut
	if err := c.ShouldBindJSON(&checkOutItem); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	Blockchain.AddBlock(checkOutItem)
	c.JSON(200, "new block added sucessfully")
}

func newBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	h := md5.New()
	io.WriteString(h, book.ISBN+book.PublishDate)
	book.ID = fmt.Sprintf("%x", h.Sum(nil))

	c.JSON(200, gin.H{
		"book": book,
	})

}

func GenesisBlock() *Block {
	return CreateBlock(&Block{}, BookCheckOut{IsGenesis: true})
}

func NewBlockChain() *BlockChain {
	return &BlockChain{
		[]*Block{GenesisBlock()},
	}
}

func getBlockChain(c *gin.Context) {

	c.JSON(200, gin.H{
		"blocks": Blockchain.blocks,
	})
}

func main() {

	Blockchain = NewBlockChain()
	r := gin.Default()
	r.GET("/", getBlockChain)
	r.POST("/", writeBlock)
	r.POST("/new", newBook)

	go func() {
		for _, block := range Blockchain.blocks {
			fmt.Printf("Prev.hash : %x\n", block.PrevHash)
			bytes, _ := json.MarshalIndent(block.Data, "", " ")
			fmt.Printf("Data : %v\n", string(bytes))
			fmt.Printf("Hash : %x\n", block.Hash)
			fmt.Println()
		}
	}()

	log.Fatal(r.Run())
}
