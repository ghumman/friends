class Solution:
    def combinationSum(self, candidates: List[int], target: int) -> List[List[int]]:

        result = []
        currentList = []
        self.dfs(result, currentList, candidates, target, 0, 0)
        return result
    
    def dfs(self, result:List[List[int]], currentList: List[int], candidates: List[int], target: int, runningSum: int, currentIndex: int) :
        if runningSum == target:
            result.append(list(currentList))
            return
        if runningSum > target:
            return
        
        for i in range(currentIndex, len(candidates)):
            val = candidates[i]
            currentList.append(val)
            self.dfs(result, currentList, candidates, target, runningSum + val, i)
            currentList.pop()
        

    
        