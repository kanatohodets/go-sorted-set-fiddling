package main

import (
	"container/heap"
	"sort"
)

type Doc uint64

type DocSlice []Doc

func (a DocSlice) Len() int           { return len(a) }
func (a DocSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DocSlice) Less(i, j int) bool { return a[i] < a[j] }

type DocSetsHeap [][]Doc

func (h DocSetsHeap) Len() int           { return len(h) }
func (h DocSetsHeap) Less(i, j int) bool { return h[i][0] < h[j][0] }
func (h DocSetsHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *DocSetsHeap) Push(x interface{}) {
	t := x.([]Doc)
	*h = append(*h, t)
}

func (h *DocSetsHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func HeapUnion(docSets [][]Doc) []Doc {
	h := DocSetsHeap(docSets)
	heap.Init(&h)
	set := []Doc{}
	for h.Len() > 0 {
		cur := h[0]
		metric := cur[0]
		if len(set) == 0 || set[len(set)-1] != metric {
			set = append(set, metric)
		}
		if len(cur) == 1 {
			heap.Pop(&h)
		} else {
			h[0] = cur[1:]
			heap.Fix(&h, 0)
		}
	}
	return set
}

func HeapIntersect(docSets [][]Doc) []Doc {
	if len(docSets) == 0 {
		return []Doc{}
	}

	for _, list := range docSets {
		// any empty set --> empty intersection
		if len(list) == 0 {
			return []Doc{}
		}
	}

	h := DocSetsHeap(docSets)
	heap.Init(&h)
	set := []Doc{}
	for {
		cur := h[0]
		smallestDoc := cur[0]
		present := 0
		for _, candidate := range h {
			if candidate[0] == smallestDoc {
				present++
			} else {
				// any further matches will be purged by the fixup loop
				break
			}
		}

		// found something in every subset
		if present == len(docSets) {
			if len(set) == 0 || set[len(set)-1] != smallestDoc {
				set = append(set, smallestDoc)
			}
		}

		for h[0][0] == smallestDoc {
			list := h[0]
			if len(list) == 1 {
				return set
			}

			h[0] = list[1:]
			heap.Fix(&h, 0)
		}
	}
}

func RepeatedPairwiseIntersect(docSets [][]Doc) []Doc {
	if len(docSets) == 0 {
		return []Doc{}
	}
	for _, set := range docSets {
		// any empty set -> no intersection
		if len(set) == 0 {
			return set
		}
	}
	sort.Sort(DocSetsHeap(docSets))
	// result can contain at most the number of items in the smallest set
	a := docSets[0]
	result := make([]Doc, 0, len(docSets[0]))
	var ridx int
	for i := 1; i < len(docSets); i++ {
		b := docSets[i]
		var aidx, bidx int

	scan:
		for aidx < len(a) && bidx < len(b) {
			if a[aidx] == b[bidx] {
				//log.Printf("i: %v, ridx: %v, aidx: %v, bidx: %v, value %v, res: %v", i, ridx, aidx, bidx, a[aidx], result)
				if len(result) == 0 {
					result = append(result, a[aidx])
				} else {
					if result[ridx] != a[aidx] {
						result = append(result, a[aidx])
					}
					ridx++
				}
				aidx++
				bidx++
				if aidx == len(a) || bidx == len(b) {
					break scan
				}
			}

			for a[aidx] < b[bidx] {
				aidx++
				if aidx == len(a) {
					break scan
				}
			}

			for a[aidx] > b[bidx] {
				bidx++
				if bidx == len(b) {
					break scan
				}
			}
		}
		a = result
		ridx = 0
	}

	return a
}
