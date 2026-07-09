class Solution {
    public List<List<Integer>> combinationSum(int[] candidates, int target) {
        List<List<Integer>> result = new ArrayList<>();
        dfs (candidates, target, result, new ArrayList<>(), 0, 0);
        return result; 
    }

    private void dfs(int[] candidates, int target, List<List<Integer>> result, List<Integer> current, int runningSum, int start) {
        if (runningSum == target) {
            result.add(new ArrayList<>(current));
            return; 
        }
        if (runningSum > target) {
            return; 
        }

        for (int i=start; i<candidates.length; i++) {
            current.add(candidates[i]);
            dfs(candidates, target, result, current, runningSum + candidates[i], i);
            current.remove(current.size() - 1);
        }
    }
}