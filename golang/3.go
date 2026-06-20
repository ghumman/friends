package golang
func lengthOfLongestSubstring(s string) int {
    l := 0
    result := 0
    seen := make(map[byte]bool)

    for r:=0; r<len(s); r++ {
        for seen[s[r]] {
            delete(seen, s[l])
            l++
        }
        seen[s[r]] = true

        if r - l + 1 >  result {
            result = r - l + 1
        }
    }
    return result
}