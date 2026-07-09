class Solution {
    public List<List<Integer>> permute(int[] nums) {
        List<List<Integer>> result = new ArrayList<>();
        Set<Integer> set = new HashSet<>();
        dfs(nums, result, set, new ArrayList<>(), 0);
        return result; 
    }

    private void dfs(int[] nums, List<List<Integer>> result, Set<Integer> set, List<Integer> currentList, int index) {
        if (nums.length == currentList.size()) {
            result.add(new ArrayList<>(currentList));
            return; 
        }



        for (int i=0; i<nums.length; i++) {
        if (set.contains(nums[i])) {
            continue;
        }
            set.add(nums[i]);
            currentList.add(nums[i]);
            dfs(nums, result, set, currentList, i);
            currentList.remove(currentList.size() - 1);
            set.remove(nums[i]);
        }
    }
}