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
    key    []rune
    isRoot bool
    depth  int64
    data   []string
}

func newNode() (m *Node) {
    m = new(Node)
    
    m.IsLeaf = false
    m.Parent = nil
    
    m.key = make([]rune, 0, k_DEFAULT_KEY_ALLOC_SIZE)
    m.data = make([]string, 0, k_DEFAULT_KEY_ALLOC_SIZE)

    return
}

// NewMap returns a new empty map
func NewMap() (m *Node) {
    m = newNode()
    
    m.IsLeaf = false
    
    m.isRoot = true
    m.depth = 0

    return
}

func (m *Node) increaseDepth() {
    last_depth := m.depth
    q := newQueue()
    q.enqueue(m)

    for !q.isEmpty() {
        n := q.dequeue()
        if n.depth > last_depth {
            last_depth = n.depth
        }

        n.depth++
        for _, c := range n.Children {
            q.enqueue(c)
        }
    }
}

// This method traverses the map to find an appropriate node
// for the given key. Optionally, if no node is found, one is created.
//
// Algorithm: BFS
func (m *Node) nodeForKey(key string, createIfMissing bool) *Node {
    var last_node *Node
    var current_node = m

    key_ := []rune(key) // we need to edit the key

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
            goto next_child
        }

        // key matches current node: returning it
        if string(key_) == string(current_node.key) {
            return current_node
        }

        _, lcpI = lcp(string(key_), string(current_node.key))

        // current node is not the one
        if lcpI == -1 {
            goto next_child
        }
        
        // in this case the key we are looking for is a substring
        // of the current node key so we need to split the node
        if lcpI == (len(key_) - 1) && lcpI < (len(current_node.key) - 1) {
            current_node.split(string(key_))
            return current_node
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
        
        key_ = key_[lcpI+1:]
        last_node = current_node
        children = current_node.Children

    next_child:
        if len(children) > 0 {
            current_node, children = children[0], children[1:]
            continue
        }
        break
    }

    if createIfMissing && last_node == nil {
        last_node = nodeWithKey(key_)
        m.appendNode(last_node)
    }

    return last_node
}

func (m *Node) split(subkey string) {
    _, lcpIndex := lcp(string(m.key), subkey)
    rightKey := m.key[lcpIndex+1:]
    leftKey  := subkey
    subNode := m.copyNode()
    subNode.key = rightKey
    subNode.increaseDepth()
    subNode.Parent = m
    subNode.IsLeaf = true
    
    // adjusting children parent
    for _, child := range subNode.Children {
        child.Parent = subNode
    }
    
    m.key = []rune(leftKey)
    m.Children = []*Node{ subNode }
    m.data = []string{}
}

func (m *Node) copyNode() (*Node) {
    n := newNode()
    *n = *m
    return n
}

func nodeWithKey(key []rune) *Node {
    n := newNode()
    n.key = key
    return n
}

func (m *Node) appendNode(n *Node) *Node {
    n.depth = m.depth + 1
    m.Children = append(m.Children, n)
    n.IsLeaf = true
    n.Parent = m
    return n
}

func (m *Node) Insert(key string, values ...string) {
    n := m.nodeForKey(key, true)
    n.data = append(n.data, values...)
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
func lcp(strs ...string) (string, int) {
    // Special cases first
    switch len(strs) {
    case 0:
        return "", -1
    case 1:
        return strs[0], 0
    }
    // LCP of min and max (lexigraphically)
    // is the LCP of the whole set.
    min, max := strs[0], strs[0]
    for _, s := range strs[1:] {
        switch {
        case s < min:
            min = s
        case s > max:
            max = s
        }
    }
    for i := 0; i < len(min) && i < len(max); i++ {
        if min[i] != max[i] {
            return min[:i], i
        }
    }
    // In the case where lengths are not equal but all bytes
    // are equal, min is the answer ("foo" < "foobar").
    return min, (len(min) - 1)
}
