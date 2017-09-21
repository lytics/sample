package sample

import (
	"fmt"
	"testing"
	"time"

	"github.com/bmizerany/assert"
)

var (
	testSampler = NewSampler(time.Now().UnixNano())
)

func TestSampleFloats(t *testing.T) {
	x := []float64{1.0, 2.0, 3.0}
	weights := []float64{0.2, 0.2, 9.6}

	sample, err := testSampler.SampleFloats(x, 30, true, weights)
	assert.Equal(t, nil, err)

	samplemap := make(map[string]int)
	for _, val := range sample {
		samplemap[fmt.Sprintf("%.f", val)]++
	}
	assert.Tf(t, samplemap["3"] > 1, "3.0 should have been sampled more than once")

	sample, err = testSampler.SampleFloats(x, 3, false, weights)
	samplemap = make(map[string]int)
	for _, val := range sample {
		samplemap[fmt.Sprintf("%.f", val)]++
	}

	assert.Equal(t, 3, len(samplemap), "samplemap didn't sample without replacement")
	for _, count := range samplemap {
		assert.Equal(t, 1, count, "samplemap didn't sample without replacement")
	}
}

func TestRemoveIndex(t *testing.T) {
	testIndex := index{1, 2, 3, 4}

	// remove the 2th element
	testIndex = testIndex.Remove(2)
	assert.Equal(t, index{1, 2, 4}, testIndex)

	// don't remove anything (removal out of bounds)
	testIndex = testIndex.Remove(20)
	assert.Equal(t, index{1, 2, 4}, testIndex)
}

func TestFind(t *testing.T) {
	w := vector{1, 2, 3}
	prob := w.CumProb()
	assert.Equal(t, 0, find(prob, -0.01))
	assert.Equal(t, 0, find(prob, 0.0))
	assert.Equal(t, 0, find(prob, 0.16))
	assert.Equal(t, 1, find(prob, 0.17))
	assert.Equal(t, 2, find(prob, 0.99))
	assert.Equal(t, 2, find(prob, 1.00))
	assert.Equal(t, 2, find(prob, 1.01))

	w2 := vector{0.1, 0.3, 0.4}
	prob2 := w2.CumProb()
	assert.Equal(t, 0, find(prob2, -0.1))
	assert.Equal(t, 0, find(prob2, 0.12))
	assert.Equal(t, 1, find(prob2, 0.13))
	assert.Equal(t, 1, find(prob2, 0.49))
	assert.Equal(t, 2, find(prob2, 0.50))
}

func TestRemove(t *testing.T) {
	x := vector{0, 10, 200, 3000}

	x = x.Remove(uint(0))
	assert.Equal(t, vector{10, 200, 3000}, x)

	x = x.Remove(uint(1))
	assert.Equal(t, vector{10, 3000}, x)

	x = x.Remove(uint(1))
	assert.Equal(t, vector{10}, x)

	x = x.Remove(uint(0))
	assert.Equal(t, vector{}, x)
}
