class Solution {
    /*
    numCourses: 2
    prerequisites: [1, 0]

    requirements[0, 1]
    graph: [0: [1]]

    queue: [0]

    if all 0s in requirements: return true otherwise return false
    */
    public boolean canFinish(int numCourses, int[][] prerequisites) {
        int[] requirements = new int[numCourses];
        Map<Integer, List<Integer>> graph = new HashMap<>();

        Queue<Integer> queue = new LinkedList<>();

        for (int[] p : prerequisites) {
            requirements[p[0]]++;
            List<Integer> l = graph.getOrDefault(p[1], new ArrayList<>());
            l.add(p[0]);
            graph.put(p[1], l);
        }

        for (int i=0; i<requirements.length; i++) {
            if (requirements[i] == 0) {
                queue.offer(i);
            }
        }

        while (!queue.isEmpty()) {
            int n = queue.size();
            for (int i=0; i<n; i++) {
                int course = queue.poll();
                List<Integer> l = graph.getOrDefault(course, new ArrayList<>());
                for ( Integer okCourse : l) {
                    requirements[okCourse]--;
                    if (requirements[okCourse] == 0) {
                        queue.offer(okCourse);
                    }
                }
            }
        }

        for (int r : requirements) {
            if (r != 0) {
                return false; 
            }
        }

        return true; 
    }
}