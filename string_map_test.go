package stringmap

import (
	"testing"
)

func TestInsert(t *testing.T) {
	m := NewMap()
	expectedValues := []string{ "Diaferia", "alediaferia", "adiaferia" }
	m.Insert("Alessandro", expectedValues...)

	node := m.nodeForKey("Alessandro", false)
	if testEq(expectedValues, node.data) != true {
		t.Errorf("Unexpected value for node 'Alessandro': expected (%v), got (%v)", expectedValues, node.data)
	}
}

/* utils */
func testEq(a, b []string) bool {

    if a == nil && b == nil { 
        return true; 
    }

    if a == nil || b == nil { 
        return false; 
    }

    if len(a) != len(b) {
        return false
    }

    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }

    return true
}
