class Solution {
    public boolean isAnagram(String s, String t) {
        Map<Character, Integer> map = new HashMap<>();

        for (char ch : s.toCharArray()) {
            map.put(ch, map.getOrDefault(ch, 0) + 1);
        }

        for (char ch : t.toCharArray()) {
            if (map.containsKey(ch)) { 
                int count = map.get(ch);
                if (count == 1){
                    map.remove(ch);
                } else {
                    map.put(ch, count - 1);
                }
            } else {
                return false; 
            }
        }
        return map.isEmpty();
    }
}