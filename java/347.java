class Solution {
    public int[] topKFrequent(int[] nums, int k) {
        Map<Integer, Integer> map = new HashMap<>();
        for (int num : nums) {
            map.put(num, map.getOrDefault(num, 0) + 1);
        }

        PriorityQueue<Map.Entry<Integer, Integer>> pq = new PriorityQueue<>((a, b) -> a.getValue() - b.getValue());


        for (Map.Entry<Integer, Integer> e : map.entrySet()) {
            pq.add(e);

            if (pq.size() > k) {
                pq.poll();
            }
        }
        int[] result = new int[k];
        for (int i=0; i<k; i++) {
            result[i] = pq.remove().getKey();
        }
        return result; 
    }
}