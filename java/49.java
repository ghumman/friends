class Solution {
    public List<List<String>> groupAnagrams(String[] strs) {
        Map<String, List<String>> map = new HashMap<>();
        for (String str : strs) {
            char[] word = str.toCharArray();
            Arrays.sort(word);
            String sortedStr = new String(word);
            map.computeIfAbsent(sortedStr, k -> new ArrayList<>()).add(str);

        }

        return new ArrayList<>(map.values());

    }
}