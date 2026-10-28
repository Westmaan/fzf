package main // adjust package name if needed

import "testing"

func TestMatchSet(t *testing.T) {
	// 1. Empty merger
	r := revision{}
	emptyMg := EmptyMerger(r)
	emptySet := emptyMg.MatchSet()
	if emptySet.matched(0) || emptySet.matched(10) || emptySet.matched(-1) {
		t.Error("Empty merger should match nothing")
	}

	// 2. Multi-list merger
	lists, _ := buildLists(true)
	mg := NewMerger(nil, lists, true, false, r, 0, 100000)
	ms := mg.MatchSet()

	// Verify all items in lists are marked as matched
	expectedCount := 0
	for _, list := range lists {
		for _, item := range list {
			expectedCount++
			if !ms.matched(item.item.Index()) {
				t.Errorf("Expected item index %d to be matched", item.item.Index())
			}
		}
	}

	// Verify an unindexed number outside the lists does not match
	if ms.matched(9999999) {
		t.Error("Unmatched index reported as matched")
	}
	if ms.matched(-5) {
		t.Error("Negative index reported as matched")
	}

	// 3. Fallback path test (force hashSet fallback via large index delta)
	mgLarge := NewMerger(nil, lists, true, false, r, 0, 30000000)
	msLarge := mgLarge.MatchSet()
	if msLarge.hashSet == nil {
		t.Error("Expected fallback to hashSet for large index span")
	}
	for _, list := range lists {
		for _, item := range list {
			if !msLarge.matched(item.item.Index()) {
				t.Errorf("Fallback hashSet: expected item %d to match", item.item.Index())
			}
		}
	}
}
