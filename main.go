package main

import (
	. "MercleTree/MerkleTree"
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
)

func main() {
	treeSize := flag.Int("treeValuesNumber", 0, "number of values stored in the tree")
	treeNodeChildrenNumber := flag.Int("treeNodeChildNumber", 1, "number of children per node in the tree")
	flag.Parse()

	hashFunc := func(data []byte) []byte {
		result := sha256.Sum256(data)
		return result[:]
	}
	tree := NewMerkleTree(hashFunc, *treeNodeChildrenNumber)
	var children [][]byte
	for i := 0; i < *treeSize; i++ {
		randomDate := GenerateRandomDate()
		children = append(children, convertDateToByteSlice(randomDate))
	}
	if err := tree.InsertNodes(children); err != nil {
		log.Fatalf("Encountered error inserting children to the tree: %v", err)
	}

	fmt.Printf("Youv'e requested to create a tree with %v random values, where each node has %v children, the tree has in total %v nodes", *treeSize, *treeNodeChildrenNumber, tree.Size())
	return
}

func convertDateToByteSlice(date int64) []byte {
	result := make([]byte, 8)
	binary.LittleEndian.PutUint64(result, uint64(date))
	return result
}

func GenerateRandomDate() int64 {
	min := time.Date(2018, 1, 1, 0, 0, 0, 0, time.Local).Unix()
	max := time.Date(2018, 12, 31, 23, 59, 0, 0, time.Local).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return sec
}
