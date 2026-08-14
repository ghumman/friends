class Solution {
    public int numIslands(char[][] grid) {
        int count = 0; 
        int n = grid.length; 
        int m = grid[0].length; 
        for (int i=0; i<n; i++) {
            for (int j=0; j<m; j++) {
                if (grid[i][j] == '1') {
                    count++;
                    dfs(grid, i, j);
                }
            }
        }
        return count; 
    }

    private void dfs(char[][] grid, int i, int j) {
        if (i < 0 || j < 0 || i > grid.length - 1 || j > grid[0].length - 1 || grid[i][j] != '1') {
            return; 
        }

        grid[i][j] = '*';
        dfs(grid, i+1, j);
        dfs(grid, i, j+1);
        dfs(grid, i-1, j);
        dfs(grid, i, j-1);
    }
}