class Solution {
    public String minWindow(String s, String t) {
        if (s.length() < t.length()) {
            return "";
        }

        Map<Character, Integer> map = new HashMap<>();

        for (char c : t.toCharArray()) {
            map.put(c, map.getOrDefault(c, 0) + 1);
        }

        int l = 0, r = 0; 
        int required = t.length();

        int minLen = Integer.MAX_VALUE;
        int start = 0; 

        while (r < s.length()) {
            char rChar = s.charAt(r);

            if (map.getOrDefault(rChar, 0) > 0) {
                required--;
            }
            map.put(rChar, map.getOrDefault(rChar, 0) - 1);
            r++;

            while ( required == 0) {
                if (r - l < minLen) {
                    minLen = r - l;
                    start = l;
                }

                char lChar = s.charAt(l);
                map.put(lChar, map.getOrDefault(lChar, 0) + 1);

                if (map.getOrDefault(lChar, 0) > 0) {
                    required++;
                }
                l++;
            }
        }

            return minLen == Integer.MAX_VALUE ? "" : s.substring(start, start + minLen); 
        }
}