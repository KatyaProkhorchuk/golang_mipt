package consistenthash

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
)

type node string

func (n node) ID() string { return string(n) }

func TestHash_SingleNode(t *testing.T) {
	h := New[node]()

	n1 := node("1")
	h.AddNode(&n1)

	require.Equal(t, &n1, h.GetNode("key0"))
}

func TestHash_TwoNodes(t *testing.T) {
	h := New[node]()

	n1 := node("1")
	h.AddNode(&n1)

	n2 := node("2")
	h.AddNode(&n2)

	n := h.GetNode("key0")
	require.True(t, n == &n1 || n == &n2)
	for i := 0; i < 32; i++ {
		require.Equal(t, n, h.GetNode("key0"))
	}

	differs := false
	for i := 0; i < 32; i++ {
		other := h.GetNode(fmt.Sprintf("key%d", i))
		if other != n {
			differs = true
		}
	}
	require.True(t, differs)
}

func TestHash_EvenDistribution(t *testing.T) {
	h := New[node]()

	const K = 32
	for i := 0; i < K; i++ {
		n := node(fmt.Sprint(i))
		h.AddNode(&n)
	}

	counts := map[*node]int{}
	const N = 1 << 16
	for i := 0; i < N; i++ {
		counts[h.GetNode(fmt.Sprintf("key%d", i))] += 1
	}

	const P = 1 / float64(K)
	const variance = N * (P) * (1 - P)
	stddev := math.Sqrt(variance)
	t.Logf("P = %v, var = %v, stddev = %v", P, variance, stddev)
	t.Logf("counts = %v", maps.Values(counts))

	for _, count := range counts {
		require.Greater(t, count, int(N/K-stddev*10))
		require.Less(t, count, int(N/K+stddev*10))
	}
}

func TestHash_ConsistentDistribution(t *testing.T) {
	h := New[node]()

	const K = 32
	for i := 0; i < K; i++ {
		n := node(fmt.Sprint(i))
		h.AddNode(&n)
	}

	nodes := map[string]*node{}

	const N = 1 << 16
	for i := 0; i < N; i++ {
		key := fmt.Sprintf("key%d", i)
		nodes[key] = h.GetNode(key)
	}

	newNode := node("new_node")
	h.AddNode(&newNode)

	changed := 0
	movedToNewNode := 0

	for key, oldNode := range nodes {
		n := h.GetNode(key)
		if n != oldNode {
			changed++
		}

		if n == &newNode {
			movedToNewNode++
		}
	}

	t.Logf("changed = %d, movedToNewNode = %d", changed, movedToNewNode)
	assert.Less(t, changed, N/K*2)
	assert.Equal(t, movedToNewNode, changed)
}

func BenchmarkHashSpeed(b *testing.B) {
	for _, K := range []int{32, 1024, 4096} {
		h := New[node]()

		for i := 0; i < K; i++ {
			n := node(fmt.Sprint(i))
			h.AddNode(&n)
		}

		b.Run(fmt.Sprintf("K=%d", K), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = h.GetNode(fmt.Sprintf("key%d", i))
			}
		})
	}
}
