package MerkleTree_test

import (
	"github.com/Yafimk/MerkleTree/MerkleTree"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestMerkleTree(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MerkleTRee Suite")
}

func MockHashFunction(value []byte) []byte {
	return value
}

var _ = Describe("MerkleTree", func() {
	BeforeEach(func() {})
	It("create empty merkle tree", func() {
		tree := MerkleTree.NewMerkleTree(func(value []byte) []byte {
			return value
		}, 3)
		Expect(tree.Size()).To(Equal(0))
	})
	It("check insertion to merkle tree", func() {
		tree := MerkleTree.NewMerkleTree(func(value []byte) []byte {
			return value
		}, 3)
		Expect(tree.Size()).To(Equal(0))
		value1 := []byte("porky")
		value2 := []byte("dooky")
		value3 := []byte("pooky")
		values := [][]byte{value1, value2, value3}
		err := tree.InsertNodes(values)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(tree.Size()).To(Equal(4))
		Expect(len(tree.GetRoot().GetChildren())).To(Equal(3))

		var children []byte
		for _, item := range values {
			children = append(children, MockHashFunction(item)...)
		}
		expectedChildrennHash := MockHashFunction(children)
		Expect(tree.GetRoot().GetValue()).To(BeEquivalentTo(expectedChildrennHash))
	})
	It("check insertion to merkle tree that require small padding", func() {
		tree := MerkleTree.NewMerkleTree(func(value []byte) []byte {
			return value
		}, 3)
		Expect(tree.Size()).To(Equal(0))
		value1 := []byte("porky")
		value2 := []byte("dooky")
		values := [][]byte{value1, value2, value2}
		err := tree.InsertNodes(values)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(tree.Size()).To(Equal(4))
		Expect(len(tree.GetRoot().GetChildren())).To(Equal(3))

		var children []byte
		for _, item := range values {
			children = append(children, MockHashFunction(item)...)
		}
		expectedChildrenHash := MockHashFunction(children)
		Expect(tree.GetRoot().GetValue()).To(Equal(expectedChildrenHash))
	})
	It("check insertion to merkle tree that require large padding", func() {
		tree := MerkleTree.NewMerkleTree(func(value []byte) []byte {
			return value
		}, 3)
		Expect(tree.Size()).To(Equal(0))
		value1 := []byte("porky")
		value2 := []byte("dooky")
		value3 := []byte("pooky")
		value4 := []byte("tooky")
		value5 := []byte("choky")
		values := [][]byte{value1, value2, value3, value4, value5, value5, value5, value5, value5}
		err := tree.InsertNodes(values)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(tree.Size()).To(Equal(13))
		Expect(len(tree.GetRoot().GetChildren())).To(Equal(3))

		var children []byte
		for _, item := range values {
			children = append(children, MockHashFunction(item)...)
		}
		expectedChildrenHash := MockHashFunction(children)
		Expect(tree.GetRoot().GetValue()).To(Equal(expectedChildrenHash))
	})
	It("check insertion to merkle tree that require large padding, for binary tree", func() {
		tree := MerkleTree.NewMerkleTree(func(value []byte) []byte {
			return value
		}, 2)
		Expect(tree.Size()).To(Equal(0))
		value1 := []byte("porky")
		value2 := []byte("dooky")
		value3 := []byte("pooky")
		value4 := []byte("tooky")
		value5 := []byte("choky")
		values := [][]byte{value1, value2, value3, value4, value5, value5, value5, value5}
		err := tree.InsertNodes(values)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(tree.Size()).To(Equal(15))
		Expect(len(tree.GetRoot().GetChildren())).To(Equal(2))

		var children []byte
		for _, item := range values {
			children = append(children, MockHashFunction(item)...)
		}
		expectedChildrenHash := MockHashFunction(children)
		Expect(tree.GetRoot().GetValue()).To(Equal(expectedChildrenHash))
	})
})
