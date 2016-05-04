package stringmap

import (
    _ "strings"
)

// constants
const (
    // default allocation space for keys
    k_DEFAULT_KEY_ALLOC_SIZE = 10
)

type Node struct {
    IsLeaf   bool
    Parent   *Node
    Children []*Node

    // private
    key    []byte
    isRoot bool
    //depth  int64
    data   []string
}

func newNode() (m *Node) {
    m = new(Node)

    m.IsLeaf = false
    m.Parent = nil
    
    return
}

// NewMap returns a new empty map
func NewMap() (m *Node) {
    m = newNode()
    m.isRoot = true
    //m.depth = 0

    return
}

func (m *Node) Depth() int {
    depth := 0
    parent := m.Parent
    for parent != nil {
        depth++
        parent = parent.Parent
    }
    
    return depth
}

// This method traverses the map to find an appropriate node
// for the given key. Optionally, if no node is found, one is created.
//
// Algorithm: BFS
func (m *Node) nodeForKey(key string, createIfMissing bool) *Node {
    var last_node = m
    var current_node = m

    key_ := []byte(key) // we need to edit the key

    // holds the next children to explore
    var children []*Node

    var lcpI int // last lcp index

    for current_node != nil && len(key_) > 0 {
        // root is special case for us
        // since it doesn't hold any information
        if current_node.isRoot {
            if len(current_node.Children) == 0 {
                break
            }
            children = current_node.Children
            if len(children) > 0 {
                current_node, children = children[0], children[1:]
                continue
            }
            break
        }

        lcpI = lcp(key_, current_node.key)

        // current node is not the one
        if lcpI == -1 {
            if len(children) > 0 {
                current_node, children = children[0], children[1:]
                continue
            }
            break
        }
               
        // key matches current node: returning it
        if lcpI == len(key_) - 1 && lcpI == len(current_node.key) - 1 {
            return current_node
        }
        key_ = key_[lcpI+1:]

        // in this case the key we are looking for is a substring
        // of the current node key so we need to split the node
        if len(key_) == 0 {
            if createIfMissing == true {
                current_node.split(lcpI+1)
            }
            return current_node
        }
        
        // current node key is a substring of the requested
        // key so we go deep in the tree from here
        if lcpI == len(current_node.key) - 1 {
            last_node = current_node
            children = current_node.Children
            if len(children) == 0 {
                break
            }
            current_node, children = children[0], children[1:]
            continue
        }

        if createIfMissing == true {
            current_node.split(lcpI+1)
            last_node = current_node
            break
        }
        
        // third case: given key partially matches with
        // current node key
        // this means we have to split the existing node
        // into two nodes and append the new content accordingly
        //
        // e.g.
        // Key to be inserted: 'string'
        // Node found: 'stringmap'
        // => we need to split 'stringmap' into 'string' and 'map'
        //    in order to be able to set a value for the key 'string'
        //    and still maintain the value(s) associated with 'stringmap'
        //    in the new 'map' node
        //
        // States can be represented as following:
        //
        // State 1:
        //
        //         o (root)
        //         |
        //         o (stringmap) = (some values)
        //
        // State 2 after inserting key 'string' into the map:
        //
        //        o (root)
        //        |
        //        o (string) = (some values associated with 'string' key)
        //        |
        //        o (map)    = (some values associated with 'stringmap' key)
    }

    if createIfMissing == true {
        newNode := newNodeWithKey(key_)
        return last_node.appendNode(newNode)
    }

    return last_node
}

func (m *Node) split(index int) {
    rightKey := m.key[index:]
    leftKey  := m.key[:index]
    subNode := m.copyNode()
    subNode.key = rightKey
    subNode.Parent = m
    subNode.IsLeaf = true

    // adjusting children parent
    for _, child := range subNode.Children {
        child.Parent = subNode
    }

    m.key = []byte(leftKey)
    m.Children = []*Node{ subNode }
    m.data = []string{}
    m.IsLeaf = false
}

func (m *Node) copyNode() (*Node) {
    n := &Node{}
    *n = *m
    return n
}

func newNodeWithKey(key []byte) *Node {
    n := newNode()
    n.key = key
    return n
}

func (m *Node) appendNode(n *Node) *Node {
    //n.depth = m.depth + 1
    m.Children = append(m.Children, n)
    n.IsLeaf = true
    n.Parent = m
    return n
}

// Insert inserts a new value in the map for the specified key
// If the key is already present in the map, the value is appended
// to the values list associated with the given key
func (m *Node) Insert(key string, values ...string) {
    n := m.nodeForKey(key, true)
    n.data = append(n.data, values...)
}

func (m *Node) Contains(key string) bool {
    return m.nodeForKey(key, false) != nil
}

// func (m *Node) Replace(key string, values ...string) {
// }

// -------- auxiliary functions -------- //

// LCP: Longest Common Prefix
// Implementation freely inspired from:
// https://rosettacode.org/wiki/Longest_common_prefix#Go
//
// returns the lcp and the index of the last
// character matching
//
func lcp(strs ...[]byte) int {
    if len(strs) < 2 {
        return -1
    }
    // Special cases first
    switch len(strs) {
    case 0:
        return -1
    case 1:
        return 0
    }
    // LCP of min and max (lexigraphically)
    // is the LCP of the whole set.
    min, max := strs[0], strs[0]
    part := strs[1:]
    for i := 0; i < len(part); i++ {
        s := part[i]
        switch {
        case len(s) < len(min):
            min = s
        case len(s) > len(max):
            max = s
        }
    }
    for i := 0; i < len(min) && i < len(max); i++ {
        if min[i] != max[i] {
            return i - 1
        }
    }
    // In the case where lengths are not equal but all bytes
    // are equal, min is the answer ("foo" < "foobar").
    return len(min) - 1
}
