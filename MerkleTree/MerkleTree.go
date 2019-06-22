package MerkleTree

import (
	"bytes"
	"container/list"
	"log"
	"math"
	"reflect"
)

type TreeHashFunc func(data []byte) []byte

type MerkleTree struct {
	root             *MerkleNode
	hashFunction     TreeHashFunc
	treeChildrenSize int
	size             int
}

func (mt *MerkleTree) GetRoot() *MerkleNode {
	return mt.root
}

func (mt *MerkleTree) Size() int {
	return mt.size
}

func (mt *MerkleTree) InsertNodes(values [][]byte) error {
	values = padInsertedValues(values, mt.treeChildrenSize)
	hashQueue := list.New()
	for _, value := range values {
		hashValue := mt.hashFunction(value)
		hashQueue.PushBack(NewMerkleNode(hashValue, nil))
	}

	var newNode *MerkleNode
	for hashQueue.Len() > 1 {
		var childrenSum []byte
		newNode = NewMerkleNode(nil, nil)
		for count := 0; count < mt.treeChildrenSize; count++ {
			node := hashQueue.Front()
			merkleNode, isCastOk := node.Value.(*MerkleNode)
			if !isCastOk {
				log.Fatalf("merkle tree walk queue encountred unrelated node of type [%v] ", reflect.TypeOf(node.Value))
			}
			childrenSum = append(childrenSum, merkleNode.GetValue()...)
			hashQueue.Remove(node)
			newNode.AddChild(merkleNode)
		}
		childrenHashSum := mt.hashFunction(childrenSum)
		newNode.SetValue(childrenHashSum)
		hashQueue.PushBack(newNode)
		mt.size += mt.treeChildrenSize
	}
	mt.size += 1
	mt.SetRoot(newNode)

	return nil
}

func (mt *MerkleTree) SetRoot(root *MerkleNode) {
	mt.root = root
}

func lookup(rootNode *MerkleNode, lookupValue []byte) (*list.List, bool) {
	var queue = list.New()
	var nodeTraverser func(mn *MerkleNode, lookupValue []byte) bool
	nodeTraverser = func(mn *MerkleNode, lookupValue []byte) bool {
		if len(mn.GetChildren()) == 0 {
			if bytes.Compare(lookupValue, mn.GetValue()) == 0 {
				return true
			}
		} else {
			junctionElement := queue.PushFront(mn)
			for _, child := range mn.GetChildren() {
				traversalResult := nodeTraverser(child, lookupValue)
				if traversalResult {
					return true
				}
			}
			queue.Remove(junctionElement)
		}
		return false
	}
	nodeTraversalResult := nodeTraverser(rootNode, lookupValue)
	if nodeTraversalResult {
		return queue, true
	}
	return nil, false
}

//GetLookupValueProofPath Looks for the hashed values and returnes the values required for merkle tree proof
func (mt *MerkleTree) GetLookupValueProofPath(lookupValue []byte) ([][]byte, bool) {
	nodeQueue, isContained := lookup(mt.GetRoot(), lookupValue)
	if !isContained {
		return nil, false
	}
	var proofPath [][]byte
	lastNodeValue := lookupValue
	for node := nodeQueue.Front(); node != nil; node = node.Next() {
		castNode := node.Value.(*MerkleNode)
		for _, child := range castNode.GetChildren() {
			if bytes.Compare(child.value, lastNodeValue) != 0 {
				proofPath = append(proofPath, castNode.GetValue())
			}
		}
	}
	return proofPath, true
}

func padInsertedValues(values [][]byte, childrenSize int) [][]byte {
	requiredValueNumber := 0
	for i := 0; requiredValueNumber < len(values); i++ {
		requiredValueNumber = int(math.Pow(float64(childrenSize), float64(i)))
	}
	requiredValueNumber -= len(values)
	for i := 0; i < requiredValueNumber; i++ {
		// we don't have enough to create perfect merkle tree, so we fill last items till we get
		// to a power of tree children size
		values = append(values, []byte("0"))
	}
	return values
}

//NewMerkleTree receives function used to hash the inserted values and number of children per node - i.e. if set to 3, it's ternery tree
func NewMerkleTree(hashFunction TreeHashFunc, childrenPerNode int) *MerkleTree {
	return &MerkleTree{hashFunction: hashFunction, treeChildrenSize: childrenPerNode}
}

type MerkleNode struct {
	value    []byte
	children []*MerkleNode
}

func (mn *MerkleNode) String() string {
	return string(mn.value)
}

func NewMerkleNode(value []byte, children []*MerkleNode) *MerkleNode {
	return &MerkleNode{value: value, children: children}
}

func (mn *MerkleNode) SetValue(value []byte) {
	mn.value = value
}

func (mn *MerkleNode) AddChild(child *MerkleNode) {
	mn.children = append(mn.children, child)
}

func (mn *MerkleNode) GetValue() []byte {
	return mn.value
}

func (mn *MerkleNode) GetChildren() []*MerkleNode {
	return mn.children
}
