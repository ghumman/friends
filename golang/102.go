package golang

/**
 * Definition for a binary tree node.
**/
 type TreeNode struct {
     Val int
     Left *TreeNode
     Right *TreeNode
 }
 
func levelOrder(root *TreeNode) [][]int {
    if root == nil {
        return nil
    }
    queue := []*TreeNode{}
    result := [][]int{}

    queue = append(queue, root)
    for len(queue) > 0 {
        currentLen := len(queue)
        currentElements := make([]int, 0, currentLen)
        for i:=0; i<currentLen; i++ {
            currentElements = append(currentElements, queue[i].Val)
            if queue[i].Left != nil {
                queue = append(queue, queue[i].Left)
            }
            if queue[i].Right != nil {
                queue = append(queue, queue[i].Right)
            }
        }
        result = append(result, currentElements)
        queue = queue[currentLen:]
    }
    return result
}