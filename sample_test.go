package sample

import (
	"fmt"
	"math"
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

func testSampleOrder(t *testing.T, x []int, weights vector, replace bool, results []vector) {
	assert.Equalf(t, len(x), len(results), "length of results must equal length of x")

	// test sampling order
	orders := make([][]int, len(x))
	for i := range x {
		orders[i] = make([]int, len(x))
	}

	niter := int(1e5)
	for iter := 0; iter < niter; iter++ {
		randorder, err := testSampler.SampleInts(x, len(x), replace, weights)
		assert.Equal(t, nil, err)
		for i, v := range randorder {
			orders[i][v]++
		}
	}

	probs := make([]vector, len(x))
	for i := range x {
		probs[i] = make(vector, len(x))
	}

	for i, order := range orders {
		for j, count := range order {
			probs[i][j] = float64(count) / float64(niter)
		}
	}

	fmt.Printf("Testing:\n")
	fmt.Printf("x:       %v\n", x)
	fmt.Printf("weights: %v\n", weights)
	fmt.Printf("replace: %v\n", replace)
	for i, prob := range probs {
		fmt.Printf("\ti(%d): %v\n", i, prob)
	}

	for i, prob := range probs {
		assert.Tf(t, approxequal(prob, results[i], 0.05), "(%d) %v != %v", i, prob, results[i])
	}
}

// sampler.summary(sampler(c(0.5, 0.5), replace = FALSE))
func TestEqualSampleWeighting(t *testing.T) {
	testSampleOrder(t, []int{0, 1}, []float64{0.5, 0.5}, false, []vector{
		vector{0.5, 0.5},
		vector{0.5, 0.5},
	})
}

// sampler.summary(sampler(c(1/3, 1/2, 1/4, 1/2), replace = FALSE))
func TestWeightedSamplingWOReplacement(t *testing.T) {
	x := []int{0, 1, 2, 3}
	weights := vector{1 / 3., 1 / 2., 1 / 4., 1 / 2.}
	testSampleOrder(t, x, weights, false, []vector{
		vector{0.211, 0.316, 0.158, 0.316},
		vector{0.234, 0.289, 0.188, 0.289},
		vector{0.272, 0.241, 0.246, 0.241},
		vector{0.283, 0.154, 0.409, 0.154},
	})
}

// sampler.summary(sampler(rep(0.25, 4), replace = FALSE))
func TestUnweightedSamplingWOReplacement(t *testing.T) {
	x := []int{0, 1, 2, 3}
	testSampleOrder(t, x, nil, false, []vector{
		vector{0.25, 0.25, 0.25, 0.25},
		vector{0.25, 0.25, 0.25, 0.25},
		vector{0.25, 0.25, 0.25, 0.25},
		vector{0.25, 0.25, 0.25, 0.25},
	})
}

// sampler.summary(sampler(c(1/3, 1/2, 1/4, 1/2), replace = TRUE))
func TestWeightedSamplingWReplacement(t *testing.T) {
	x := []int{0, 1, 2, 3}
	weights := vector{1 / 3., 1 / 2., 1 / 4., 1 / 2.}
	scaled := weights.Scale()
	testSampleOrder(t, x, weights, true, []vector{
		scaled, scaled, scaled, scaled,
	})
}

// sampler.summary(sampler(rep(0.25, 4), replace = TRUE))
func TestUnweightedSamplingWReplacement(t *testing.T) {
	x := []int{0, 1, 2, 3}
	weights := vector{0.25, 0.25, 0.25, 0.25}
	testSampleOrder(t, x, nil, true, []vector{
		weights, weights, weights, weights,
	})
}

func TestIndexMethods(t *testing.T) {
	// remove elements
	assert.Equal(t, index{2, 3, 4}, index{1, 2, 3, 4}.Remove(0))
	assert.Equal(t, index{1, 3, 4}, index{1, 2, 3, 4}.Remove(1))
	assert.Equal(t, index{1, 2, 4}, index{1, 2, 3, 4}.Remove(2))
	assert.Equal(t, index{1, 2, 3}, index{1, 2, 3, 4}.Remove(3))
	assert.Equal(t, index{1, 2, 3, 4}, index{1, 2, 3, 4}.Remove(4))
}

func TestVectorMethods(t *testing.T) {
	v := vector{0.1, 0.2, 0.3, 0.4}
	assert.T(t, approxequal(vector{0.1, 0.3, 0.6, 1.0}, v.CumProb(), 0.001))
	v2 := v.Copy()
	v2 = v2.Remove(2).Remove(200)
	assert.T(t, approxequal(vector{1 / 7., 3 / 7., 7 / 7.}, v2.CumProb(), 0.001))
}

func TestErrors(t *testing.T) {
	_, err := testSampler.SampleInts([]int{1, 9, 3}, 4, false, nil)
	assert.NotEqual(t, nil, err)

	_, err = testSampler.SampleInts([]int{1, 9, 3}, 2, true, vector{0.1, 0.1})
	assert.NotEqual(t, nil, err)

	_, err = testSampler.SampleFloats([]float64{0.1, 0.9, 0.3}, 4, false, nil)
	assert.NotEqual(t, nil, err)

	_, err = testSampler.SampleFloats([]float64{0.1, 0.9, 0.3}, 2, true, vector{0.1, 0.1})
	assert.NotEqual(t, nil, err)
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

func approxequal(f1, f2 vector, precision float64) bool {
	if len(f1) != len(f2) {
		return false
	}
	for i := range f1 {
		if math.Abs(f1[i]-f2[i]) > precision {
			return false
		}
	}
	return true
}
