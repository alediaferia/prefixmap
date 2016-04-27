package stringmap

import (
    "testing"
)

type queue_test struct {
    w string
}

var queue_tests = []queue_test{
    {"alessandro"},
    {"typeflow"},
    {"triego"},
}

func Test_queue(t *testing.T) {
    for _, v := range queue_tests {
        q := newQueue()
        word := make([]rune, 0, len(v.w))
        for _, c := range v.w {
            t_ := newNode()
            t_.key = append(t_.key, c)
            q.enqueue(t_)
        }

        for !q.isEmpty() {
            t_ := q.dequeue()
            word = append(word, t_.key[0])
        }

        if string(word) != v.w {
            t.Errorf("Unexpected characters dequeued: got '%s', expected '%s'", string(word), v.w)
        }
    }
}

func Benchmark_queue_enqueue(b *testing.B) {
    b.ReportAllocs()
    q := newQueue()
    t := newNode()
    b.StopTimer()
    b.ResetTimer()
    b.StartTimer()
    for i := 0; i < b.N; i++ {
        q.enqueue(t)
    }
}

func Benchmark_queue_dequeue(b *testing.B) {
    b.ReportAllocs()
    q := newQueue()
    t := newNode()
    b.StopTimer()
    for i := 0; i < b.N; i++ {
        q.enqueue(t)
    }
    b.ResetTimer()
    b.StartTimer()
    for i := 0; i < b.N; i++ {
        q.dequeue()
    }
}
