package golang
func maxArea(height []int) int {
    l := 0
    r := len(height) - 1
    maxArea := 0

    for l < r {
        minHeight := height[l]
        if minHeight > height[r] {
            minHeight = height[r]
        }
        currentArea := (r - l) * minHeight
        if currentArea > maxArea {
            maxArea = currentArea
        }
        if height[l] < height[r] {
            l++
        } else {
            r--
        }
    }
    return maxArea
}