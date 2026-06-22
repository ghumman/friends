package golang
func longestConsecutive(nums []int) int {
    numSet := make(map[int]bool)
    for _, num := range nums {
        numSet[num] = true
    }
    result := 0
    for num := range numSet {
        if !numSet[num - 1] {
            streak := 0
            for numSet[num] {
                streak++
                num++
            }
            if result < streak {
                result = streak
            }
        }
    }
    return result
}