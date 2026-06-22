package golang

/**
 * Definition for a Node.
 * type Node struct {
 *     Val int
 *     Neighbors []*Node
 * }
 */

type Node struct {
    Val int
    Neighbors []*Node
}

func cloneGraph(node *Node) *Node {
    seen := make(map[*Node]*Node)
    return dfs(seen, node)
}

func dfs(seen map[*Node]*Node, node *Node) *Node {
    if node == nil {
        return nil
    }
    
    if n, ok := seen[node]; ok{
        return n
    }
    clone := &Node{Val : node.Val}
    seen[node] = clone
    for _, nei := range node.Neighbors {
        clone.Neighbors = append(clone.Neighbors, dfs(seen, nei))
    }
    return clone
}