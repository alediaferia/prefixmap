package prefixmap

// little queue implementation

const q_PAGE_SIZE = 4096 // common page size

type queue struct {
    q                []*Node
    pages            [][]*Node
    h, t, page_index int
}

func (q *queue) enqueue(node *Node) {
    if q.t == cap(q.q) {
        // moving to the next page
        q.page_index += 1

        // incrementing pages slice
        // if no empty pages are available
        if q.page_index == len(q.pages) {
            page := make([]*Node, q_PAGE_SIZE)
            q.pages = append(q.pages, page)
        }
        q.q = q.pages[q.page_index]

        // resetting indexes
        q.t = 0
        q.h = 0
    }
    q.q[q.t] = node
    q.t += 1
}

func (q *queue) isEmpty() bool {
    return q.h == q.t
}

func (q *queue) dequeue() (node *Node) {
    if q.h == q.t {
        if q.page_index > 0 {
            q.page_index -= 1
            q.q = q.pages[q.page_index]
            q.h = 0
            q.t = len(q.q)
        } else {
            node = nil
            return
        }
    }

    node = q.q[q.h]
    q.h += 1
    return
}

func (q *queue) clear() {
    q.h = 0
    q.t = 0
}

func newQueue() *queue {
    q := new(queue)
    q.q = make([]*Node, q_PAGE_SIZE)
    q.pages = [][]*Node{q.q}
    q.h = 0
    q.t = 0
    q.page_index = 0
    return q
}
