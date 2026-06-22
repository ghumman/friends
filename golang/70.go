package golang
func climbStairs(n int) int {
    if n < 3 {
        return n
    }
    prev2, prev1  := 1, 2
    curr := 0
    for i:=3; i<=n; i++ {
        curr = prev1 + prev2
        prev2 = prev1
        prev1 = curr
    }
    return prev1
}