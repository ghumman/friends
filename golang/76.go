package golang

import "math"
func minWindow(s string, t string) string {
    if len(s) == 0 || len(t) == 0 || len(s) < len(t) {
        return ""
    }

    dictT := make([]int, 128)
    for i:=0; i<len(t); i++ {
        dictT[t[i]]++
    }

    required := 0
    for _, count := range dictT {
        if count > 0 {
            required++
        }
    }

    windowCounts := make([]int, 128)

    left, right := 0, 0
    formed := 0


    minLen := math.MaxInt32
    startIdx := 0

    for right < len(s) {
        charRight := s[right]
        windowCounts[charRight]++

        if dictT[charRight] > 0 && windowCounts[charRight] == dictT[charRight] {
            formed++
        }

        for left <= right && formed == required {
            charLeft := s[left]

            currentWindowLength := right - left + 1
            if currentWindowLength < minLen {
                minLen = currentWindowLength
                startIdx = left
            }

            windowCounts[charLeft]--
            if dictT[charLeft] > 0 && windowCounts[charLeft] < dictT[charLeft] {
                formed--
            }
            left++
        }
        right++
    }
    if minLen == math.MaxInt32 {
        return ""
    }
    return s[startIdx : startIdx + minLen]
}