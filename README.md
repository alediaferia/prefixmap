# PrefixMap
PrefixMap is a prefix-enhanced map that eases the retrieval of values based on key prefixes.

Creating a PrefixMap
---
```go
// creates the map object
prefixMap := prefixmap.NewMap()
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

Iterate over prefixes
---

`PrefixMap` exposes an `EachPrefix` method that executes a callback function
against every prefix in the map. The prefixes are iterated over using a [Depth First Search](https://en.wikipedia.org/wiki/Depth-first_search)
algorithm. At each iteration the given callback is invoked. The callback allows you to skip a branch
iteration altogether if you're not satisfied with what you're looking for.
Check out `PrefixCallback` documentation for more information.

```go
prefixMap.EachPrefix(func(prefix Prefix) (bool, bool) {
    
    // do something with the current prefix
    doSomething(prefix.Key)
    
    // keep iterating
    return false, false
})
```