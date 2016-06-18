package prefixmap

import (
  "gopkg.in/alediaferia/stackgo.v1"
)

// Node is a single node within
// the map
type Node struct {
  // true if this node is a leaf node
  IsLeaf bool

  // the reference to the parent node
  Parent *Node

  // the children nodes
  Children []*Node

  // private
  key    string
  isRoot bool
  data   []interface{}
}

// PrefixMap type
type PrefixMap Node

func newNode() (m *Node) {
  m = new(Node)

  m.IsLeaf = false
  m.Parent = nil

  return
}

// New returns a new empty map
func New() *PrefixMap {
  m := newNode()
  m.isRoot = true

  return (*PrefixMap)(m)
}

// Depth returns the depth of the
// current node within the map
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
// Returns an additional bool indicating if the node key is an exact match
// for the given key parameter or false if it is the closest match found
// Algorithm: BFS
func (m *Node) nodeForKey(key string, createIfMissing bool) (*Node, bool) {
  var lastNode *Node = nil
  var currentNode = m

  // holds the next children to explore
  var children []*Node

  var lcpI int // last lcp index

  for currentNode != nil && len(key) > 0 {
    // root is special case for us
    // since it doesn't hold any information
    if currentNode.isRoot {
      if len(currentNode.Children) == 0 {
        break
      }
      children = currentNode.Children
      if len(children) > 0 {
        currentNode, children = children[0], children[1:]
        continue
      }
      break
    }

    lcpI = lcpIndex(key, currentNode.key)

    // current node is not the one
    if lcpI == -1 {
      if len(children) > 0 {
        currentNode, children = children[0], children[1:]
        continue
      }
      break
    }

    // key matches current node: returning it
    if lcpI == len(key)-1 && lcpI == len(currentNode.key)-1 {
      return currentNode, true
    }
    key = key[lcpI+1:]

    // in this case the key we are looking for is a substring
    // of the current node key so we need to split the node
    if len(key) == 0 {
      if createIfMissing == true {
        currentNode.split(lcpI + 1)
        return currentNode, true
      }
      return currentNode, false
    }

    // current node key is a substring of the requested
    // key so we go deep in the tree from here
    if lcpI == len(currentNode.key)-1 {
      lastNode = currentNode
      children = currentNode.Children
      if len(children) == 0 {
        break
      }
      currentNode, children = children[0], children[1:]
      continue
    }

    if createIfMissing == true {
      currentNode.split(lcpI + 1)
      lastNode = currentNode
      break
    }

    // Important Case: given key partially matches with
    // current node key.
    //
    // This means we have to split the existing node
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
    newNode := newNodeWithKey(key)
    if lastNode == nil {
      return m.appendNode(newNode), true
    }
    return lastNode.appendNode(newNode), true
  }

  return lastNode, false
}

func (m *Node) split(index int) {
  rightKey := m.key[index:]
  leftKey := m.key[:index]
  subNode := m.copyNode()
  subNode.key = rightKey
  subNode.Parent = m
  subNode.IsLeaf = true

  // adjusting children parent
  for _, child := range subNode.Children {
    child.Parent = subNode
  }

  m.key = leftKey
  m.Children = []*Node{subNode}
  m.data = []interface{}{}
  m.IsLeaf = false
}

func (m *Node) copyNode() *Node {
  n := &Node{}
  *n = *m
  return n
}

func newNodeWithKey(key string) *Node {
  n := newNode()
  n.key = key
  return n
}

func (m *Node) appendNode(n *Node) *Node {
  m.Children = append(m.Children, n)
  n.IsLeaf = true
  n.Parent = m
  return n
}

// Insert inserts a new value in the map for the specified key
// If the key is already present in the map, the value is appended
// to the values list associated with the given key
func (m *PrefixMap) Insert(key string, values ...interface{}) {
  mNode := (*Node)(m)
  n, _ := mNode.nodeForKey(key, true)
  n.data = append(n.data, values...)
}

// Replace replaces the value(s) for the given key in the map
// with the give ones. If no such key is present, this method
// behaves the same as Insert
func (m *PrefixMap) Replace(key string, values ...interface{}) {
  mNode := (*Node)(m)
  n, _ := mNode.nodeForKey(key, true)
  n.data = values
}

// Contains checks if the given key is present in the map
// In this case, an exact match case is considered
// If you're interested in prefix-based check: ContainsPrefix
func (m *PrefixMap) Contains(key string) bool {
  mNode := (*Node)(m)
  retrievedNode, exactMatch := mNode.nodeForKey(key, false)
  return retrievedNode != nil && exactMatch
}

// Get returns the data associated with the given key in the map
// or nil if no such key is present in the map
func (m *PrefixMap) Get(key string) []interface{} {
  mNode := (*Node)(m)
  retrievedNode, exactMatch := mNode.nodeForKey(key, false)
  if !exactMatch {
    return nil
  }

  return retrievedNode.data
}

// GetByPrefix returns a flattened collection of values
// associated with the given prefix key
func (m *PrefixMap) GetByPrefix(key string) []interface{} {
  mNode := (*Node)(m)
  retrievedNode, _ := mNode.nodeForKey(key, false)
  if retrievedNode == nil {
    return []interface{}{}
  }

  // now, fetching all the values (DFS)
  stack := stackgo.NewStack()
  values := []interface{}{}
  stack.Push(retrievedNode)
  for stack.Size() > 0 {
    node := stack.Pop().(*Node)
    values = append(values, node.data...)
    for _, c := range node.Children {
      stack.Push(c)
    }
  }

  return values
}

// ContainsPrefix checks if the given prefix is present as key in the map
func (m *PrefixMap) ContainsPrefix(key string) bool {
  mNode := (*Node)(m)
  retrievedNode, _ := mNode.nodeForKey(key, false)
  return retrievedNode != nil
}

// Key Retrieves current node key
// complexity: MAX|O(log(N))| where N
// is the number of nodes in the map.
// Number of nodes in the map cannot exceed
// number of keys + 1.
func (m *Node) Key() string {
  node := m
  k := make([]byte, 0, len(m.key))
  for node != nil && node.isRoot != true {
    key := string(node.key) // triggering a copy here
    k = append([]byte(key), k...)
    node = node.Parent
  }
  return string(k)
}

// PrefixCallback is invoked by EachPrefix for each prefix reached
// by the traversal. The callback has the ability to affect the traversal.
// Returning skipBranch = true will make the traversal skip the current branch
// and jump to the sibling node in the map. Returning halt = true, instead,
// will halt the traversal altogether.
type PrefixCallback func(prefix Prefix) (skipBranch bool, halt bool)

// Prefix holds prefix information
// passed to the PrefixCallback instance by
// the EachPrefifx method.
type Prefix struct {
  node *Node

  // The current prefix string
  Key string

  // The values associated to the current prefix
  Values []interface{}
}

// Depth returns the depth of the corresponding
// node for this prefix in the map.
func (p *Prefix) Depth() int {
  return p.node.Depth()
}

// EachPrefix iterates over the prefixes contained in the
// map using a DFS algorithm. The callback can be used to skip
// a prefix branch altogether or halt the iteration.
func (m *PrefixMap) EachPrefix(callback PrefixCallback) {
  mNode := (*Node)(m)
  stack := stackgo.NewStack()
  prefix := []byte{}

  skipsubtree := false
  halt := false
  addedLengths := stackgo.NewStack()
  lastDepth := mNode.Depth()

  stack.Push(mNode)
  for stack.Size() != 0 {
    node := stack.Pop().(*Node)
    if !node.isRoot {
      // if we are now going up
      // in the radix (e.g. we have
      // finished with the current branch)
      // then we adjust the current prefix
      currentDepth := node.Depth()
      if lastDepth >= node.Depth() {
        var length = 0
        for i := 0; i < (lastDepth-currentDepth)+1; i++ {
          length += addedLengths.Pop().(int)
        }
        prefix = prefix[:len(prefix)-length]
      }
      lastDepth = currentDepth
      prefix = append(prefix, node.key...)
      addedLengths.Push(len(node.key))

      // building the info
      // data to pass to the callback
      info := Prefix{
        node:   node,
        Key:    string(prefix),
        Values: node.data,
      }

      skipsubtree, halt = callback(info)
      if halt {
        return
      }
      if skipsubtree {
        continue
      }
    }
    for i := 0; i < len(node.Children); i++ {
      stack.Push(node.Children[i])
    }
  }
}

// -------- auxiliary functions -------- //

// LCP: Longest Common Prefix
// Implementation freely inspired from:
// https://rosettacode.org/wiki/Longest_common_prefix#Go
//
// returns the lcp and the index of the last
// character matching
//
func lcpIndex(strs ...string) int {
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
    case len(s) >= len(max):
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
