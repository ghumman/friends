package golang
func combinationSum(candidates []int, target int) [][]int {
    result := [][]int{}
    dfs(candidates, target, &result, []int{}, 0, 0)
    return result; 
}

func dfs(candidates []int, target int, result *[][]int, current []int, runningSum int, index int) {
    if (runningSum == target) {
        *result = append(*result, append([]int{}, current...))
        return
    }
    if (runningSum > target) {
        return
    }

    for i:=index; i<len(candidates); i++ {
        current = append(current, candidates[i])
        dfs(candidates, target, result, current, runningSum + candidates[i], i)
        current = current[:len(current) - 1]
    }
}