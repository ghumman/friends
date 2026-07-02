class Solution {
    public boolean isValid(String s) {
        Map<Character, Character> map = new HashMap<>();

        map.put('(', ')');
        map.put('{', '}');
        map.put('[', ']');

        Deque<Character> stack = new ArrayDeque<>();

        for(char c : s.toCharArray()) {
            if (map.containsKey(c)) {
                stack.push(c);
            } else {
                if (stack.isEmpty()) {
                    return false; 
                }
                char complement = stack.pop();
                if (map.get(complement) != c) {
                    return false; 
                }
            }
        }
        return stack.isEmpty();
    }
}