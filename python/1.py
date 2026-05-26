from typing import List

class Solution:
    def twoSum(self, nums: List[int], target: int) -> List[int]:
        location_dict = {}
        for i, num in enumerate(nums):
            complement = target - num
            if complement in location_dict:
                return [location_dict[complement], i]
            location_dict[num] = i
        return []
        