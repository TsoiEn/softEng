package model

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// BlockChain structure contains a slice of blocks.
type BlockChain struct {
	Blocks []Block
}

// Block represents a block in the blockchain.
type Block struct {
	Index     int
	Timestamp string
	Data      []byte
	Hash      []byte
	PrevHash  []byte
}

// Serialize serializes the block into a JSON byte slice.
func (b *Block) Serialize() ([]byte, error) {
	return json.Marshal(b)
}

// DeriveHash generates a hash for the block using its index, timestamp, data, and previous hash.
func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{[]byte(fmt.Sprintf("%d", b.Index)), []byte(b.Timestamp), b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

// SetPrevHash sets the previous hash for the block.
func (b *Block) SetPrevHash(prevHash []byte) {
	if len(prevHash) == 0 {
		genesisBlock := Genesis()
		b.PrevHash = genesisBlock.Hash
	} else {
		b.PrevHash = prevHash
	}
}

// CreateBlock creates a new block with the given data and previous hash.
func CreateBlock(index int, blockData []byte, prevHash []byte) *Block {
	block := &Block{
		Index:     index,
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      blockData,
		PrevHash:  prevHash,
	}
	block.DeriveHash()
	return block
}

// AddBlock adds a new block to the blockchain.
// AddBlock adds a new block to the blockchain.
func (chain *BlockChain) AddBlock(blockData []byte) {
	if len(chain.Blocks) == 0 {
		fmt.Println("Blockchain is empty, adding Genesis block first.")
		genesisBlock := Genesis()
		chain.Blocks = append(chain.Blocks, *genesisBlock)
	}

	prevBlock := chain.Blocks[len(chain.Blocks)-1]

	// Validate the previous block's hash
	prevBlock.DeriveHash()
	if !bytes.Equal(prevBlock.Hash, prevBlock.Hash) {
		fmt.Println("Previous block's hash is invalid.")
		return
	}

	// Create a new block, set the PrevHash properly
	newIndex := prevBlock.Index + 1
	newBlock := CreateBlock(newIndex, blockData, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, *newBlock)
	log.Printf("New Block Added: %+v", newBlock)
}

// Genesis creates the first block in the blockchain.
func Genesis() *Block {
	return CreateBlock(0, []byte("Genesis Block"), []byte{})
}

// NewBlockChain creates a blockchain with the genesis block.
func NewBlockChain() *BlockChain {
	return &BlockChain{Blocks: []Block{*Genesis()}}
}

// FindCredentialByID searches the blockchain for a credential with the given ID.
func (chain *BlockChain) FindCredentialByID(id string) (*Credential, error) {
	for _, block := range chain.Blocks {
		var cred Credential
		if err := json.Unmarshal(block.Data, &cred); err == nil {
			if cred.ID == id {
				return &cred, nil
			}
		}
	}
	return nil, fmt.Errorf("credential with ID %s not found", id)
}

// CredentialChain represents a blockchain specifically for credentials.

// AddCredentialToBlockchain adds a new credential to the blockchain.
func (chain *CredentialChain) AddCredentialToBlockchain(cred *Credential) error {
	// 1. Validate the credential data
	if err := ValidateCredentialData(cred); err != nil {
		return fmt.Errorf("invalid credential data: %w", err)
	}

	// 2. Generate a hash for the credential
	cred.Hash = GenerateCredentialHash(cred)

	// 3. Serialize the credential into JSON format for block data
	blockData, err := json.Marshal(cred)
	if err != nil {
		return fmt.Errorf("failed to serialize credential: %w", err)
	}

	// 4. Check if the blockchain is initialized
	if len(chain.Blocks) == 0 {
		// Initialize the blockchain with a genesis block
		genesisBlock := CreateBlock(0, []byte("Genesis Block"), nil)
		chain.Blocks = append(chain.Blocks, *genesisBlock)
	}

	// 5. Get the last block in the chain
	prevBlock := chain.Blocks[len(chain.Blocks)-1]

	// 6. Create a new block for the credential, set the PrevHash properly
	newBlock := CreateBlock(prevBlock.Index+1, blockData, prevBlock.Hash)

	// 7. Add the new block to the blockchain
	chain.Blocks = append(chain.Blocks, *newBlock)

	log.Printf("Credential added to blockchain: %+v", cred)
	return nil
}
