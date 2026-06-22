package golang
func trap(height []int) int {
    left, right := 0, len(height) - 1
    water := 0
    leftMax := height[left]
    rightMax := height[right]

    for left < right {
        if height[left] < height[right] {
            if height[left] > leftMax {
                leftMax = height[left]
            } else {
                water += leftMax - height[left]
            }
            left++
        } else {
            if height[right] > rightMax {
                rightMax = height[right]
            } else {
                water += rightMax - height[right]
            }
            right--
        }
    }
    return water
}