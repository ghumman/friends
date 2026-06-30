class Solution:
    def permute(self, nums: List[int]) -> List[List[int]]:
        used = [False] * len(nums)
        result = []
        self.dfs(nums, used, result, [])
        return result

    def dfs(self, nums: List[int], used: list[bool], result: List[List[int]], current: List[int]):
        if len(nums) == len(current):
            result.append(current.copy())
            return

        for i in range(len(nums)):
            if used[i]:
                continue

            used[i] = True
            current.append(nums[i])

            self.dfs(nums, used, result, current)

            current.pop()
            used[i] = False