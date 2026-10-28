// matchSet provides an O(1) bitset-backed membership test for matched item indices.
type matchSet struct {
	bitset  []uint64
	offset  int32
	hashSet map[int32]struct{}
}

// matched checks whether the item index is in the match set with strict bounds safety.
func (ms *matchSet) matched(index int32) bool {
	if ms == nil {
		return false
	}
	if ms.hashSet != nil {
		_, ok := ms.hashSet[index]
		return ok
	}
	idx := index - ms.offset
	if idx >= 0 && int(idx) < len(ms.bitset)*64 {
		return (ms.bitset[idx>>6] & (uint64(1) << (idx & 63))) != 0
	}
	return false
}

// MatchSet builds an O(1) matchSet by scanning lists directly, completely bypassing
// mergedGet sorting and tournament comparisons.
func (mg *Merger) MatchSet() *matchSet {
	if mg == nil || mg.count == 0 {
		return &matchSet{}
	}

	minIndex := mg.minIndex
	maxIndex := mg.maxIndex

	// Normal case: indices fit within 20 million items (under 2.5 MB bitset).
	if maxIndex > minIndex && (maxIndex-minIndex) <= 20000000 {
		size := int(maxIndex - minIndex)
		words := make([]uint64, (size+63)/64)
		for _, list := range mg.lists {
			for i := range list {
				idx := list[i].item.Index() - minIndex
				if idx >= 0 && int(idx) < size {
					words[idx>>6] |= uint64(1) << (idx & 63)
				}
			}
		}
		return &matchSet{bitset: words, offset: minIndex}
	}

	// Safe fallback for wrapped indices or massive datasets (> 20M items).
	hashSet := make(map[int32]struct{}, mg.count)
	for _, list := range mg.lists {
		for i := range list {
			hashSet[list[i].item.Index()] = struct{}{}
		}
	}
	return &matchSet{hashSet: hashSet}
}

// ToMap converts merger results to a map. Iterates lists directly in O(K)
// without triggering O(K * M) mergedGet comparisons.
func (mg *Merger) ToMap() map[int32]Result {
	ret := make(map[int32]Result, mg.count)
	for _, list := range mg.lists {
		for i := range list {
			ret[list[i].item.Index()] = list[i]
		}
	}
	return ret
}
