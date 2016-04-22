package stringmap

import (
	"testing"
)

// testing basic insert functionality
func TestInsert(t *testing.T) {
	m := NewMap()
	expectedValues := []string{ "Diaferia", "alediaferia", "adiaferia" }
	m.Insert("Alessandro", expectedValues...)

	node := m.nodeForKey("Alessandro", false)
	if testEq(expectedValues, node.data) != true {
		t.Errorf("Unexpected value for node 'Alessandro': expected (%v), got (%v)", expectedValues, node.data)
	}
}

// testing if insert appends values to the node
func TestInsertAppends(t *testing.T) {
	m := NewMap()
	expectedValues := []string{ "Diaferia", "alediaferia", "adiaferia" }
	m.Insert("Alessandro", expectedValues[:2]...) // first two values
	m.Insert("Alessandro", expectedValues[2:]...) // rest

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
