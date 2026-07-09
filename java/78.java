class Solution {
    public List<List<Integer>> subsets(int[] nums) {
        List<List<Integer>> result = new ArrayList<>();
        dfs(nums, result, new ArrayList<>(), 0);
        return result; 
    }

    private void dfs(int[] nums, List<List<Integer>> result, List<Integer> currentList, int start) {
        result.add(new ArrayList<>(currentList));
        for (int i=start; i<nums.length; i++) {
            currentList.add(nums[i]);
            dfs(nums, result, currentList, i+1);
            currentList.remove(currentList.size() - 1);
        }
    }
}