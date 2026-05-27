class Solution:
    def isValid(self, s: str) -> bool:
        complements = {
            ')' : '(',
            '}' : '{',
            ']' : '['
        }
        stack = []

        for ch in s:
            if ch in complements:
                if not stack:
                    return False
                if stack.pop() != complements[ch]:
                    return False
            else:
                stack.append(ch)
        return not stack

        