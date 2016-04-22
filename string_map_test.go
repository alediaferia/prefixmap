package stringmap

import (
	"testing"
)

func TestInsert(t *testing.T) {
	m := NewMap()
	m.Insert("Alessandro", "Diaferia", "alediaferia")

	t.Log(m.nodeForKey("Alessandro", false).data)
}
