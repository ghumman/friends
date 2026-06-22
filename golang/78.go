package golang
func subsets(nums []int) [][]int {
    result := &[][]int{}
    dps(nums, result, []int{}, 0)
    return *result
}

func dps (nums []int, result *[][]int, current []int, start int) {
    *result = append(*result, append([]int{}, current...))
    for i:= start; i<len(nums); i++ {
        current = append(current, nums[i])
        dps(nums, result, current, i + 1)
        current = current[:len(current) - 1]
    }
}