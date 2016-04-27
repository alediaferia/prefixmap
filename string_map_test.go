package stringmap

import (
    "testing"
    "bufio"
    "os"
    "io"
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

func TestSubstringKeys(t *testing.T) {
    m := NewMap()
    expectedValues := []string{ "Diaferia", "alediaferia", "adiaferia" }
    m.Insert("stringmap", expectedValues[:2]...) // first two values
    m.Insert("string", expectedValues[2:]...) // rest

    node := m.nodeForKey("stringmap", false)
    if testEq(expectedValues[:2], node.data) != true {
        t.Errorf("Unexpected value for node 'stringmap': expected (%v), got (%v)", expectedValues[:2], node.data)
    }
    node = m.nodeForKey("string", false)
    if testEq(expectedValues[2:], node.data) != true {
        t.Errorf("Unexpected value for node 'string': expected (%v), got (%v)", expectedValues[2:], node.data)
    }

    if count := m.countNodes(); count != 3 {
        t.Errorf("Unexpected node count: got %d, expected %d", count, 3)
    }
}

func TestLcp(t *testing.T) {
    lcpTestCases := []struct {
        source, destination, expected string
        index int // expected index
    } {{
        "string", "stringmap", "string",
        len("string") - 1,
    }}

    for _, testCase := range lcpTestCases {
        value, index := lcp(testCase.source, testCase.destination)
        if value != testCase.expected {
            t.Errorf("Unexpected lcp value: got '%s', expected '%s'", value, testCase.expected)
        }

        if index != testCase.index {
            t.Errorf("Unexpected lcp index: got %d, expected %d", index, testCase.index)
        }
    }
}

func TestSplit(t *testing.T) {
    m := NewMap()
    m.Insert("stringmap", "a", "b", "c")

    node := m.nodeForKey("stringmap", false)
    node.split("string")
    if len(node.Children) != 1 {
        t.Errorf("'stringmap' node should only have 1 child")
    }

    if (string(node.key) != "string") {
        t.Errorf("Node is expected to have key 'string'")
    }

    if (string(node.Children[0].key) != "map") {
        t.Errorf("Node is expected to have key 'map'")
    }
}

type nodeTest struct {
    words []string
    nodes int
}

var nodeTests = []nodeTest{
    {[]string{
        "romane",
        "romanus",
        "romulus",
        "rubens",
        "ruber",
        "rubicon",
        "rubicundus",
    }, 14,
    },
    {[]string{
        "arma",
        "armatura",
        "armento",
    }, 5,
    },
}

func TestNodeCount(t *testing.T) {
    for _, v := range nodeTests {
        m := NewMap()
        // appending nodes
        for _, w := range v.words {
            m.Insert(w, w)
        }

        count := m.countNodes()
        if count != v.nodes {
            t.Errorf("Unexpected node count: got %d, expected %d", count, v.nodes)
        }
    }
}

func Benchmark_nodesAllocation(b *testing.B) {
    b.StopTimer()

    // building country name
    // source from file
    file, err := os.Open("/usr/share/dict/words")
    words := make([]string, 0)
    if err != nil {
        b.Log("Cannot open expected file /usr/share/dict/words. Skipping this benchmark.")
        b.SkipNow()
        return
    }
    reader := bufio.NewReader(file)
    for {
        line, err := reader.ReadString('\n')
        if err == io.EOF {
            break
        }
        words = append(words, line[:len(line)-1])
    }

    b.Logf("Inserting %d words in the trie", b.N)

    m := NewMap()

    b.ResetTimer()
    b.StartTimer()
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        word := words[i % len(words)]
        m.Insert(word, word)
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

func (m *Node) countNodes() int {
    queue := newQueue()

    queue.enqueue(m)
    count := 0
    for !queue.isEmpty() {
        node := queue.dequeue()
        count++

        for _, c := range node.Children {
            queue.enqueue(c)
        }
    }

    return count
}
