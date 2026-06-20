/**
 * Definition for singly-linked list.
 * public class ListNode {
 *     int val;
 *     ListNode next;
 *     ListNode() {}
 *     ListNode(int val) { this.val = val; }
 *     ListNode(int val, ListNode next) { this.val = val; this.next = next; }
 * }
 */
class Solution {
    public ListNode mergeKLists(ListNode[] lists) {
        PriorityQueue<ListNode> pq = new PriorityQueue<>((a, b) -> Integer.compare(a.val, b.val));
        for (ListNode list : lists) {
            if (list != null )
                pq.offer(list);
        }
        ListNode dummy = new ListNode();
        ListNode current = dummy; 
        while (!pq.isEmpty()) {
            ListNode node = pq.poll();
            current.next = node; 
            node = node.next; 
            if (node != null) {
                pq.offer(node);
            } 
            current = current.next;
        }
        return dummy.next; 
    }
}