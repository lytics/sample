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

	assert.Equal(t, 0, Find(prob, 0.16))
	assert.Equal(t, 1, Find(prob, 0.17))
}
