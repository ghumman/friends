package golang

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func maxPathSum(root *TreeNode) int {
    if root == nil {
        return 0
    }
    result := root.Val
    var dfs func (*TreeNode)int
    dfs = func(root *TreeNode) int {
        if root == nil {
            return 0
        }
        leftGain := dfs(root.Left)
        if leftGain < 0 {
            leftGain = 0
        }
        rightGain := dfs(root.Right)
        if rightGain < 0 {
            rightGain = 0
        }

        currentMaxSum := leftGain + root.Val + rightGain
        if currentMaxSum > result {
            result = currentMaxSum
        }

        if root.Val + leftGain > root.Val + rightGain {
            return root.Val + leftGain
        }
        return root.Val + rightGain
    }
    dfs(root)
    return result
}