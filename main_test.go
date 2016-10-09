package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"testing"
)

var testData [][]Doc
var benchData [][]Doc
var overlap []Doc

var numSets = 30
var setSizeMean float64 = 10000
var setSizeStdDev float64 = 2000

var commonSetSize = 40

func TestMain(m *testing.M) {
	overlap = []Doc{Doc(1), Doc(7777777), Doc(1234567890), Doc(18446111111111111111)}
	testData = make([][]Doc, 0, numSets)
	for i := 0; i < numSets; i++ {
		testData = append(testData, generateDocSlice())
	}

	benchData = generateBenchData(10)
	os.Exit(m.Run())
}

func generateBenchData(n int) [][]Doc {
	sets := make([][]Doc, 0, n)
	for i := 0; i < n; i++ {
		sets = append(sets, generateDocSlice())
	}
	return sets
}

func generateDocSlice() []Doc {
	length := randSetSize()

	docs := map[Doc]bool{}
	for i := 0; i < length; i++ {
		docs[randDoc()] = true
	}

	for _, common := range overlap {
		docs[common] = true
	}

	res := make([]Doc, 0, len(docs))
	for doc, _ := range docs {
		res = append(res, doc)
	}

	sort.Sort(DocSlice(res))

	return res
}

func randDoc() Doc {
	return Doc(rand.Uint32())<<32 + Doc(rand.Uint32())
}

func randSetSize() int {
	sample := rand.NormFloat64()*setSizeStdDev + setSizeMean
	// breaks normal dist, but just a teeny bit
	if sample < 0 {
		sample = 0
	}
	return int(sample)
}

func copySeed(seed [][]Doc) [][]Doc {
	res := make([][]Doc, 0, len(testData))
	for _, set := range seed {
		res = append(res, set[:])
	}
	return res
}

func TestHeapIntersect(t *testing.T) {
	res := HeapIntersect(copySeed(testData))
	if fmt.Sprintf("%v", res) != fmt.Sprintf("%v", overlap) {
		t.Errorf("heap intersect: expected %v, got %v", overlap, res)
	}
}

func TestRepeatedPairwiseIntersect(t *testing.T) {
	res := RepeatedPairwiseIntersect(copySeed(testData))
	if fmt.Sprintf("%v", res) != fmt.Sprintf("%v", overlap) {
		t.Errorf("repeated pairwise intersect: expected %v, got %v", overlap, res)
	}
}

func BenchmarkHeapIntersect(b *testing.B) {
	tests := make([][][]Doc, 0, b.N)
	for i := 0; i < b.N; i++ {
		tests = append(tests, copySeed(benchData))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HeapIntersect(tests[i])
	}
}

func BenchmarkRepeatedPairwiseIntersect(b *testing.B) {
	tests := make([][][]Doc, 0, b.N)
	for i := 0; i < b.N; i++ {
		tests = append(tests, copySeed(benchData))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RepeatedPairwiseIntersect(tests[i])
	}
}
