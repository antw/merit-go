package merit

import (
	"sort"
	"testing"
)

func TestDispatchableListSort(t *testing.T) {
	d1 := Dispatchable{Cost: 2.0}
	d2 := Dispatchable{Cost: 2.1}
	d3 := Dispatchable{Cost: 1.0}

	list := DispatchableList{&d1, &d2, &d3}
	sort.Sort(list)

	expected := []Dispatchable{d3, d1, d2}

	for i, disp := range list {
		if expected[i] != *disp {
			t.Errorf("Sorted DispatchableList[%d] = {Cost: %f}, "+
				"want {Cost: %f}", i, disp.Cost, expected[i].Cost)
		}
	}
}
