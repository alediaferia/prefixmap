package prefixmap

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "testing"
)

// testing basic insert functionality
func TestInsert(t *testing.T) {
    m := New()
    n := (*Node)(m)
    expectedValues := []interface{}{"bar", "baz", "quz"}
    m.Insert("foo", expectedValues...)

    node, _ := n.nodeForKey("foo", false)
    if testEq(expectedValues, node.data) != true {
        t.Errorf("Unexpected value for node 'foo': expected (%v), got (%v)", expectedValues, node.data)
    }
}

func TestInsertWithLengtheningKeys(t *testing.T) {
    m := New()
    n := (*Node)(m)

    testCases := []struct {
        key1, key2 string
    }{
        {
            "a", "aa",
        },
    }

    for _, tc := range testCases {
        m.Insert(tc.key1, tc.key1)
        m.Insert(tc.key2, tc.key2)

        node, _ := n.nodeForKey(tc.key1, false)
        if k := node.Key(); k != tc.key1 {
            t.Errorf("Unexpected key: got '%s' (%v), expected '%s' (%v); values = %v", k, []byte(k), tc.key1, []byte(tc.key1), node.data)
        }

        node, _ = n.nodeForKey(tc.key2, false)
        if k := node.Key(); k != tc.key2 {
            t.Errorf("Unexpected key: got '%s' (%v), expected '%s' (%v); values = %v", k, []byte(k), tc.key2, []byte(tc.key2), node.data)
            n.print(-1)
        }
    }
}

// testing if insert appends values to the node
func TestInsertAppends(t *testing.T) {
    m := New()
    n := (*Node)(m)
    expectedValues := []interface{}{"bar", "baz", "quz"}
    m.Insert("foo", expectedValues[:2]...) // first two values
    m.Insert("foo", expectedValues[2:]...) // rest

    node, _ := n.nodeForKey("foo", false)
    if testEq(expectedValues, node.data) != true {
        t.Errorf("Unexpected value for node 'foo': expected (%v), got (%v)", expectedValues, node.data)
    }
}

func TestReplaceReplaces(t *testing.T) {
    m := New()
    n := (*Node)(m)
    expectedValues := []interface{}{"bar", "baz", "quz"}
    m.Replace("foo", expectedValues[:2]...) // first two values
    m.Replace("foo", expectedValues[2:]...) // rest

    node, _ := n.nodeForKey("foo", false)
    if testEq(expectedValues[2:], node.data) != true {
        t.Errorf("Unexpected value for node 'foo': expected (%v), got (%v)", expectedValues, node.data)
    }
}

func TestInsertSubstringKeys(t *testing.T) {
    m := New()
    n := (*Node)(m)
    expectedValues := []interface{}{"Diaferia", "alediaferia", "adiaferia"}
    m.Insert("stringmap", expectedValues[:2]...) // first two values
    m.Insert("string", expectedValues[2:]...)    // rest

    node, _ := n.nodeForKey("stringmap", false)
    if testEq(expectedValues[:2], node.data) != true {
        t.Errorf("Unexpected value for node 'stringmap': expected (%v), got (%v)", expectedValues[:2], node.data)
    }
    node, _ = n.nodeForKey("string", false)
    if testEq(expectedValues[2:], node.data) != true {
        t.Errorf("Unexpected value for node 'string': expected (%v), got (%v)", expectedValues[2:], node.data)
    }

    if count := n.countNodes(); count != 3 {
        t.Errorf("Unexpected node count: got %d, expected %d", count, 3)
    }
}

func TestGet(t *testing.T) {
    testCases := []struct {
        getKey string
        values            []interface{}
    }{
        {
            getKey:    "string",
            values:    []interface{}{"a", "b", "c"},
        },
    }
    
    for _, testCase := range testCases {
        m := New()
        m.Insert(testCase.getKey, testCase.values...)
        if data := m.Get(testCase.getKey); testEq(data, testCase.values) != true {
            t.Errorf("Unexpected value for key '%s': expected (%v), got (%v)", testCase.getKey, testCase.values, data)
        }
    }
}

func TestPrefixAsKey(t *testing.T) {
    testCases := []struct {
        insertKey, getKey string
        values            []interface{}
    }{
        {
            insertKey: "stringmap",
            getKey:    "string",
            values:    []interface{}{"a", "b", "c"},
        },
    }

    for _, testCase := range testCases {
        m := New()
        n := (*Node)(m)
        m.Insert(testCase.insertKey, testCase.values...)
        if node, _ := n.nodeForKey(testCase.getKey, false); testEq(node.data, testCase.values) != true {
            t.Errorf("Unexpected value for node '%s': expected (%v), got (%v)", testCase.getKey, testCase.values, node.data)
        }
    }
}

func TestLcp(t *testing.T) {
    lcpTestCases := []struct {
        source, destination, expected string
        index                         int // expected index
    }{{
            "string", "stringmap", "string",
            5,
       },
        {
            "romane", "romanus", "roman",
            4,
        },
        {
            "r", "a", "",
            -1,
        },
        {
            "foobar", "bar", "",
            -1,
        },
    }

    for _, testCase := range lcpTestCases {
        index := lcpIndex(testCase.source, testCase.destination)
        if index != testCase.index {
            t.Errorf("Unexpected lcp index: got %d, expected %d", index, testCase.index)
        }

        if testCase.source[:testCase.index+1] != testCase.expected {
            t.Errorf("Unexpected substring index: got %s, expected %s", testCase.source[:testCase.index+1], testCase.expected)
        }
    }
}

func TestSplit(t *testing.T) {
    m := New()
    n := (*Node)(m)
    m.Insert("stringmap", "a", "b", "c")

    node, _ := n.nodeForKey("stringmap", false)
    node.split(6)
    if len(node.Children) != 1 {
        t.Errorf("'stringmap' node should only have 1 child")
    }

    if string(node.key) != "string" {
        t.Errorf("Node is expected to have key 'string'")
    }

    if string(node.Children[0].key) != "map" {
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
        "A",
    }, 15,
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
        m := New()
        n := (*Node)(m)
        // appending nodes
        for _, w := range v.words {
            m.Insert(w, w)
            node, _ := n.nodeForKey(w, false)
            if node == nil {
                t.Errorf("Cannot retrieve node for key: %s", w)
            }
            if node.Key() != w {
                t.Errorf("Unexpected key for node: got %s, expected %s", node.Key(), w)
                n.print(-1)
                t.FailNow()
            }
        }

        count := n.countNodes()
        if count != v.nodes {
            n.print(-1)
            t.Errorf("Unexpected node count: got %d, expected %d", count, v.nodes)
        }
    }
}

func TestPrefixIteration(t *testing.T) {
    testCases := []struct {
        keys []string
        expectedPrefixes []interface{}
    }{
        {
            keys:             []string{"benchmark", "bench", "bob", "blueray", "bluetooth"},
            expectedPrefixes: []interface{}{"b", "blue", "bluetooth", "blueray", "bob", "bench", "benchmark"},
        },
    }

    for _, tc := range testCases {
        m := New()
        for _, key := range tc.keys {
            m.Insert(key, key)
        }

        foundPrefixes := []interface{}{}
        m.EachPrefix(func(prefix Prefix) (bool, bool) {
            foundPrefixes = append(foundPrefixes, prefix.Key)
            return false, false
        })
        if testEq(foundPrefixes, tc.expectedPrefixes) != true {
            t.Errorf("Unexpected prefixes list: got %v, expected %v", foundPrefixes, tc.expectedPrefixes)
        }
    }
}

func TestContains(t *testing.T) {
    type expectedResult struct {
            key    string
            result bool
    }
    testCases := []struct{
        keys []string
        expectedResults []expectedResult
    }{
        {
            keys: []string{"a","b","c"},
            expectedResults: []expectedResult{
                {
                    "a", true,
                },
                {
                    "d", false,
                },
            },
        },
    }
    
    for _, tc := range testCases {
        m := New()
        for _, key := range tc.keys {
            m.Insert(key, key)
        }
        for _, er := range tc.expectedResults {
            if got := m.Contains(er.key); got != er.result {
                t.Errorf("Unexpected result for key %s: got %v, expected %v", er.key, got, er.result)
            }
        }
    }
}

func TestContainsPrefix(t *testing.T) {
        type expectedResult struct {
            key    string
            result bool
    }
    testCases := []struct{
        keys []string
        expectedResults []expectedResult
    }{
        {
            keys: []string{"foobar","golang"},
            expectedResults: []expectedResult{
                {
                    "foo", true,
                },
                {
                    "bar", false,
                },
                {
                    "go", true,
                },
                {
                    "lang", false,
                },
            },
        },
    }
    
    for _, tc := range testCases {
        m := New()
        for _, key := range tc.keys {
            m.Insert(key, key)
        }
        for _, er := range tc.expectedResults {
            if got := m.ContainsPrefix(er.key); got != er.result {
                t.Errorf("Unexpected result for key %s: got %v, expected %v", er.key, got, er.result)
                (*Node)(m).print(-1)
            }
        }
    }
}

func BenchmarkInsertAllocations(b *testing.B) {
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

    m := New()

    b.ResetTimer()
    b.StartTimer()
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        word := words[i%len(words)]
        m.Insert(word, word)
    }
}

func BenchmarkInsertNativeMapAllocations(b *testing.B) {
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

    m := make(map[string][]string)

    b.ResetTimer()
    b.StartTimer()
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        word := words[i%len(words)]
        if v, ok := m[word]; ok == true {
            v = append(v, word)
            m[word] = v
        } else {
            m[word] = v
        }
    }
}

func BenchmarkContains(b *testing.B) {
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
    m := New()

    for i := 0; i < b.N; i++ {
        word := words[i%len(words)]
        m.Insert(word, word)
    }

    b.ResetTimer()
    b.StartTimer()

    var found = false
    for i := 0; i < b.N; i++ {
        word := words[i%len(words)]
        found = m.Contains(word)
        if found != true {
            b.Errorf("Unexpected: couldn't find word '%s'", word)
            n := (*Node)(m)
            if node, _ := n.nodeForKey(word, false); node != nil {
                k := n.Key()
                b.Logf("Node has unexpected key: %s (%v)", k, []byte(k))
                n.print(2)
            }
            b.FailNow()
        }
    }
}

func BenchmarkContainsNativeMap(b *testing.B) {
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
    m := make(map[string][]string)

    for i := 0; i < b.N; i++ {
        word := words[i%len(words)]

        if v, ok := m[word]; ok == true {
            v = append(v, word)
            m[word] = v
        } else {
            m[word] = v
        }
    }

    b.ResetTimer()
    b.StartTimer()

    for i := 0; i < b.N; i++ {
        word := words[i%len(words)]
        _, found := m[word]
        if found != true {
            b.Fatalf("Unexpected: couldn't find word '%s'", word)
        }
    }
}

/* utils */
func testEq(a, b []interface{}) bool {

    if a == nil && b == nil {
        return true
    }

    if a == nil || b == nil {
        return false
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

func (m *Node) print(maxDepth int) {
    q := newQueue()
    last_depth := m.Depth()

    fmt.Print("Map: \n")
    fmt.Print("---------\n")
    q.enqueue(m)
    for !q.isEmpty() {
        n := q.dequeue()
        if maxDepth > -1 && n.Depth() > maxDepth {
            break
        }
        if n.isRoot {
            fmt.Print("/ (root)")
        } else {
            if n.Depth() > last_depth {
                fmt.Println("\n|")
            }
            fmt.Printf("- (%s,depth=%d) => %v;\t", string(n.key), n.Depth(), n.data)
        }
        for _, c := range n.Children {
            q.enqueue(c)
        }
        last_depth = n.Depth()
    }
    fmt.Print("\n---------\n\n")
}
