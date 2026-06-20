package golang

import "sort"
func threeSum(nums []int) [][]int {
    sort.Ints(nums)
    if len(nums) < 3 {
        return nil
    }
    result := [][]int{}
    for i:= 0; i<len(nums) - 2; i++ {
        if i > 0 && nums[i] == nums[i-1] {
            continue
        }
        l := i + 1
        r := len(nums) - 1
        for l < r {
            if nums[i] + nums[l] + nums[r] == 0 {
                result = append(result, []int{nums[i], nums[l], nums[r]})
                l++
                r--
                for l < r && nums[l] == nums[l-1] {
                    l++
                }
                for l < r && nums[r] == nums[r+1] {
                    r--
                }
            } else if (nums[i] + nums[l] + nums[r] < 0) {
                l++
            } else {
                r--
            }
        }
    }
    return result
}