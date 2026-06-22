package golang

import (
	"regexp"
	"strings"
)

func isPalindrome(s string) bool {
    re := regexp.MustCompile(`[^a-zA-Z0-9]`)
    s = re.ReplaceAllString(s, "")
    s = strings.ToLower(s)
    left, right := 0, len(s) - 1

    for left < right {
        if s[left] != s[right] {
            return false
        }
        left++
        right--
    }
    return true
}