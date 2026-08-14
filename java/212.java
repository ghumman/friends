class Solution {

    class TrieNode {
        TrieNode[] childern = new TrieNode[26];
        String word; 

    }

    List<String> result = new ArrayList<>();
    public List<String> findWords(char[][] board, String[] words) {
        TrieNode root = new TrieNode();

        for (String word : words) {
            TrieNode current = root; 
            for (char c : word.toCharArray()) {
                int index = c - 'a';
                if (current.childern[index] == null)
                    current.childern[index] = new TrieNode();
                current = current.childern[index];
            }
            current.word = word; 
        }
        int m = board.length; 
        int n = board[0].length; 
        for (int i=0; i<m; i++) {
            for (int j=0; j<n; j++) {
                dfs(board, i, j, root);
            }
        }
        return result; 
    }

    private void dfs(char[][] board, int i, int j, TrieNode current) {
        if (i < 0 || j < 0 || i >= board.length || j >= board[0].length) {
            return;
        }
        char c = board[i][j];
        if (c == '#') {
            return;
        }

        TrieNode next = current.childern[c - 'a'];

        if (next == null) {
            return;
        }

        if (next.word != null) {
            result.add(next.word);
            next.word = null; 
        }

        board[i][j] = '#';

        dfs(board, i+1, j, next);
        dfs(board, i, j+1, next);
        dfs(board, i-1, j, next);
        dfs(board, i, j-1, next);

        board[i][j] = c;
    }
}