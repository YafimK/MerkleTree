package main

import (
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/Yafimk/MerkleTree/MerkleTree"
	"log"
	"math/rand"
	"strings"
	"time"
)

func main() {
	treeSize := flag.Int("values", 1, "number of random values to generate")
	treeNodeChildrenNumber := flag.Int("children", 1, "number of children per node in the tree")
	flag.Parse()

	hashFunc := func(data []byte) []byte {
		result := sha256.Sum256(data)
		return result[:]
	}
	tree := MerkleTree.NewMerkleTree(hashFunc, *treeNodeChildrenNumber)
	var children [][]byte
	fmt.Println("The randomized dates for insertion are - ")
	for i := 0; i < *treeSize; i++ {
		randomDate := GenerateRandomDate()
		fmt.Println(time.Unix(randomDate, 0))
		children = append(children, convertDateToByteSlice(randomDate))
	}

	if err := tree.InsertNodes(children); err != nil {
		log.Fatalf("Encountered error inserting children to the tree: %v", err)
	}

	fmt.Printf("\nCreated a tree with %v random values, \neach node has %v children and the tree has in total %v nodes\n", *treeSize, *treeNodeChildrenNumber, tree.Size())

	fmt.Printf("\n\nThe tree:\n")
	PrintTree(tree)
	return
}

func convertDateToByteSlice(date int64) []byte {
	result := make([]byte, 8)
	binary.LittleEndian.PutUint64(result, uint64(date))
	return result
}

func convertByteSliceToDate(date []byte) uint64 {
	return binary.LittleEndian.Uint64(date)
}

func GenerateRandomDate() int64 {
	min := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2018, 12, 31, 23, 59, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return sec
}

func formatNodeValue(value []byte) string {
	return fmt.Sprintf("%v", time.Unix(int64(convertByteSliceToDate(value)), 0))
}

func printNode(mn *MerkleTree.MerkleNode, pre string, sb *strings.Builder) {
	nodeChildren := mn.GetChildren()
	if len(nodeChildren) == 0 {
		sb.WriteString(fmt.Sprintf("╴%v\n", "leaf"))
		return
	}
	sb.WriteString(fmt.Sprintf("┐%v\n", "node"))
	last := len(nodeChildren) - 1
	for _, ch := range nodeChildren[:last] {
		sb.WriteString(fmt.Sprintf("%v%v", pre, "├─"))
		printNode(ch, pre+"│ ", sb)
	}
	sb.WriteString(fmt.Sprintf("%v%v", pre, "└─"))
	printNode(nodeChildren[last], pre+"  ", sb)
}

func PrintTree(mt *MerkleTree.MerkleTree) {
	sb := strings.Builder{}
	if mt.GetRoot() == nil {
		fmt.Println("<empty>")
		return
	}
	printNode(mt.GetRoot(), "", &sb)
	fmt.Println(sb.String())
}
