class TrieNode {
    TrieNode[] childern; 
    boolean isWord;

    public TrieNode() {
        childern = new TrieNode[26];
        isWord = false; 
    }
}

class Trie {
    TrieNode root;
    public Trie() {
        root = new TrieNode();
    }
    
    public void insert(String word) {
        TrieNode current = root; 
        for (char c : word.toCharArray()) {
            int index = c - 'a';
            if (current.childern[index] == null) {
                current.childern[index] = new TrieNode();
            }
            current = current.childern[index];
        }
        current.isWord = true; 
    }
    
    public boolean search(String word) {
        TrieNode current = root; 
        for (char c : word.toCharArray()) {
            int index = c - 'a';
            if (current.childern[index] == null) {
                return false; 
            }
            current = current.childern[index];
        }
        return current.isWord;
    }
    
    public boolean startsWith(String prefix) {
        TrieNode current = root; 
        for (char c : prefix.toCharArray()) {
            int index = c - 'a';
            if (current.childern[index] == null) {
                return false; 
            }
            current = current.childern[index];
        }
        return true; 
    }
}

/**
 * Your Trie object will be instantiated and called as such:
 * Trie obj = new Trie();
 * obj.insert(word);
 * boolean param_2 = obj.search(word);
 * boolean param_3 = obj.startsWith(prefix);
 */