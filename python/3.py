class Solution:
    def lengthOfLongestSubstring(self, s: str) -> int:
        if len(s) == 0:
            return 0
        
        l = 0
        result = 0
        unique_letters = set()

        for r in range(len(s)):
            while s[r] in unique_letters:
                unique_letters.remove(s[l])
                l += 1
            unique_letters.add(s[r])
            result = max(result, r - l + 1)
        return result