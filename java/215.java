class Solution {
    public int findKthLargest(int[] nums, int k) {
        PriorityQueue<Integer> h = new PriorityQueue<>();
        for (int num : nums) {
            h.offer(num);
            if (h.size() > k) {
                h.poll();
            }
        }
        return h.peek();
    }
}