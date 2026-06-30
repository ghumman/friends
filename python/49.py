class Solution:
    def groupAnagrams(self, strs: List[str]) -> List[List[str]]:
        groups = {}
        for str in strs:
            count = [0] * 26
            for ch in str:
                count[ord(ch) - ord('a')] += 1
            key = tuple(count)

            if key not in groups:
                groups[key] = []
            groups[key].append(str)
        return list(groups.values())