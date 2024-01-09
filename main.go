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

// Block represents a block in the blockchain.
type Block struct {
	Pos       int          // Position of the block in the blockchain
	Data      BookCheckOut // Data associated with the block
	TimeStamp string       // Timestamp when the block was created
	Hash      string       // Hash of the block
	PrevHash  string       // Hash of the previous block
}

// BookCheckOut represents information about a book checkout.
type BookCheckOut struct {
	BookID       string `json:"book_id"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis    bool   `json:"is_genesis"`
}

// Book represents information about a book.
type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	PublishDate string `json:"publish_date"`
	ISBN        string `json:"isbn"`
}

// BlockChain represents a blockchain with a collection of blocks.
type BlockChain struct {
	blocks []*Block
}

var Blockchain *BlockChain

// GenerateHash generates the hash for a block based on its data and previous hash.
func (b *Block) GenerateHash() {
	bytes, _ := json.Marshal(b.Data)
	data := fmt.Sprint(b.Pos) + b.TimeStamp + string(bytes) + b.PrevHash
	hash := sha256.New()
	hash.Write([]byte(data))
	b.Hash = hex.EncodeToString(hash.Sum(nil))
}

// CreateBlock creates a new block with the given checkout data and links it to the previous block.
func CreateBlock(prevBlock *Block, checkoutItem BookCheckOut) *Block {
	block := &Block{}
	block.Pos = prevBlock.Pos + 1
	block.TimeStamp = time.Now().String()
	block.PrevHash = prevBlock.Hash
	block.GenerateHash()

	return block
}

// AddBlock adds a new block to the blockchain with the given checkout data.
func (b *BlockChain) AddBlock(data BookCheckOut) {
	prevBlock := b.blocks[len(b.blocks)-1]

	block := CreateBlock(prevBlock, data)

	if validBlock(block, prevBlock) {
		b.blocks = append(b.blocks, block)
	}
}

// validBlock checks if a block is valid by comparing its properties with the previous block.
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

// ValidateHash recalculates the hash for a block and checks if it matches the provided hash.
func (b *Block) ValidateHash(hash string) bool {
	b.GenerateHash()
	if b.Hash != hash {
		return false
	}

	return true
}

// writeBlock handles the endpoint for adding a new block to the blockchain.
func writeBlock(c *gin.Context) {
	var checkOutItem BookCheckOut
	if err := c.ShouldBindJSON(&checkOutItem); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	Blockchain.AddBlock(checkOutItem)
	c.JSON(200, "new block added successfully")
}

// newBook handles the endpoint for adding a new book to the system.
func newBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// Generate a unique ID for the book based on ISBN and publish date.
	h := md5.New()
	io.WriteString(h, book.ISBN+book.PublishDate)
	book.ID = fmt.Sprintf("%x", h.Sum(nil))

	c.JSON(200, gin.H{
		"book": book,
	})
}

// GenesisBlock creates the initial block in the blockchain.
func GenesisBlock() *Block {
	return CreateBlock(&Block{}, BookCheckOut{IsGenesis: true})
}

// NewBlockChain creates a new blockchain with the initial GenesisBlock.
func NewBlockChain() *BlockChain {
	return &BlockChain{
		[]*Block{GenesisBlock()},
	}
}

// getBlockChain handles the endpoint for retrieving the entire blockchain.
func getBlockChain(c *gin.Context) {
	c.JSON(200, gin.H{
		"blocks": Blockchain.blocks,
	})
}

func main() {
	// Initialize the blockchain with the GenesisBlock.
	Blockchain = NewBlockChain()

	// Set up the Gin router.
	r := gin.Default()
	r.GET("/", getBlockChain)
	r.POST("/", writeBlock)
	r.POST("/new", newBook)

	// Print blockchain information in a separate goroutine.
	go func() {
		for _, block := range Blockchain.blocks {
			fmt.Printf("Prev.hash : %x\n", block.PrevHash)
			bytes, _ := json.MarshalIndent(block.Data, "", " ")
			fmt.Printf("Data : %v\n", string(bytes))
			fmt.Printf("Hash : %x\n", block.Hash)
			fmt.Println()
		}
	}()

	// Start the server.
	log.Fatal(r.Run())
}
