* [![Build Status](https://travis-ci.org/alediaferia/prefixmap.svg?branch=master)](https://travis-ci.org/alediaferia/prefixmap)
* [![Coverage Status](https://coveralls.io/repos/github/alediaferia/prefixmap/badge.svg?branch=master)](https://coveralls.io/github/alediaferia/prefixmap?branch=master)
* [![GoDoc](https://godoc.org/github.com/typeflow/typeflow-go/web?status.png)](https://godoc.org/github.com/typeflow/typeflow-go)

# PrefixMap 
PrefixMap is a prefix-enhanced map that eases the retrieval of values based on key prefixes.

Quick Start
===

Creating a PrefixMap
---
```go
// creates the map object
prefixMap := prefixmap.New()
```

Inserting a value
---
```go
 // inserts values 1, "value 2" and false for key 'someKey'
prefixMap.Insert("someKey", 1, "value 2", false)

// map now contains
//
// 'someKey' => [1, "value 2", false]
```

Replace values for key
---
```go
prefixMap.Insert("key", "hello")

// map contents:
//
// 'key' => ["hello"]

prefixMap.Insert("key", "world")

// map contents:
//
// 'key' => ["hello", "world"]

// now replacing the contents for key
prefixMap.Replace("key", "new value")

// map contents:
//
// 'key' => ["new value"]
```

Checking if a key exists
---
```go
prefixMap.Insert("key", "hello")

prefixMap.Contains("k") // #=> false
prefixMap.Contains("key") // #=> true
prefixMap.ContainsPrefix("k") // #=> true
```

Getting by key
---
```go
prefixMap.Insert("foo", "bar", "baz", "quz")

data := prefixMap.Get("foo") // #=> [bar, baz, quz]
```

Getting by keys prefix
---
```go
prefixMap.Insert("prefix1", "prefix1")
prefixMap.Insert("prefix2", "prefix2")
prefixMap.Insert("prefix3", "prefix3")

data := prefixMap.GetByPrefix("prefix") // #=> [prefix1, prefix2, prefix3]
```

Iterate over prefixes
---

[PrefixMap](https://godoc.org/github.com/typeflow/prefixmap) exposes an [EachPrefix](https://godoc.org/github.com/typeflow/prefixmap#PrefixMap.EachPrefix) 
method that executes a callback function against every prefix in the map. 
The prefixes are iterated over using a [Depth First Search](https://en.wikipedia.org/wiki/Depth-first_search)
algorithm. At each iteration the given callback is invoked. The callback allows you to skip a branch
iteration altogether if you're not satisfied with what you're looking for.
Check out [PrefixCallback](https://godoc.org/github.com/typeflow/prefixmap#PrefixCallback) documentation for more information.

```go
prefixMap.EachPrefix(func(prefix Prefix) (bool, bool) {
    
    // do something with the current prefix
    doSomething(prefix.Key)
    
    // keep iterating
    return false, false
})
```

License
===

The code contained in this repository is provided as is under the terms of the MIT license as specified [here](/LICENSE).
