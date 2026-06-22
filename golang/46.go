package golang
func permute(nums []int) [][]int {
    result := &[][]int{}
    used := make([]bool, len(nums))
    dfs(nums, result, []int{}, used)
    return *result
}

func dfs(nums []int, result *[][]int, current []int, used []bool) {
    if len(current) == len(nums) {
        *result = append(*result, append([]int{}, current...))
        return
    }
    for i:=0; i<len(nums); i++ {
        if used[i] {
            continue
        }
        used[i] = true
        current = append(current, nums[i])
        dfs(nums, result, current, used)
        current = current[:len(current) - 1]
        used[i] = false
    }
}