# Definition for singly-linked list.
# class ListNode:
#     def __init__(self, val=0, next=None):
#         self.val = val
#         self.next = next
import heapq
class Solution:
    def mergeKLists(self, lists: List[Optional[ListNode]]) -> Optional[ListNode]:
        heap = []
        for node in lists:
            if node:
                heapq.heappush(heap, (node.val, id(node), node))
        
        dummy = ListNode()
        current = dummy
        while heap:
            _, _, curr = heapq.heappop(heap)
            current.next = curr
            current = current.next
            curr = curr.next
            if curr:
                heapq.heappush(heap, (curr.val, id(curr), curr))
        
        return dummy.next