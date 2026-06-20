/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
package golang

import "container/heap"
type MinHeap []*ListNode

func (h MinHeap) Len()int {
    return len(h)
}

func (h MinHeap) Less(i, j int)bool {
    return h[i].Val < h[j].Val
}

func (h MinHeap) Swap(i, j int) {
    h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(n interface{}) {
    *h = append(*h, n.(*ListNode))
}

func (h *MinHeap) Pop()interface{} {
    old := *h
    length := len(old)
    val := old[length - 1]
    *h = old[:length - 1]
    return val
}

func mergeKLists(lists []*ListNode) *ListNode {
    h := &MinHeap{}
    heap.Init(h)

    for _, l := range(lists) {
        if l != nil {
            heap.Push(h, l)
        }
    }

    dummy := &ListNode{}
    current := dummy

    for h.Len() > 0 {
        node := heap.Pop(h).(*ListNode)
        current.Next = node
        if node.Next != nil {
            heap.Push(h, node.Next)
        }
        current = current.Next
    }
    return dummy.Next
}