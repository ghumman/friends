# 🚀 Java LeetCode Master Cheat Sheet (Modern + Interview Ready)

---

# 🧱 1. Arrays

```java id="a1"
int[] arr = new int[n];
int[] arr = {1, 2, 3};

arr.length;

Arrays.sort(arr);
Arrays.fill(arr, value);
Arrays.copyOf(arr, n);
```

---

# 📋 2. ArrayList

```java id="a2"
List<Integer> list = new ArrayList<>();

list.add(x);
list.get(i);
list.set(i, x);
list.remove(i);
list.size();
```

---

# 🔗 3. LinkedList (used as Deque/Queue)

```java id="a3"
Deque<Integer> dq = new LinkedList<>();
```

(Prefer `ArrayDeque` instead — see below)

---

# ⚡ 4. DEQUE (REPLACES STACK — IMPORTANT)

### ✅ Use this instead of Stack

```java id="a4"
Deque<Integer> stack = new ArrayDeque<>();
```

### Stack operations

```java id="a5"
stack.push(x);     // addFirst
stack.pop();       // removeFirst
stack.peek();      // peekFirst
```

### Queue operations

```java id="a6"
queue.offer(x);
queue.poll();
queue.peek();
```

### Deque operations (both ends)

```java id="a7"
dq.addFirst(x);
dq.addLast(x);

dq.removeFirst();
dq.removeLast();

dq.peekFirst();
dq.peekLast();
```

---

# 🚀 5. HashMap (MOST IMPORTANT DS)

```java id="a8"
Map<Integer, Integer> map = new HashMap<>();

map.put(k, v);
map.get(k);
map.getOrDefault(k, 0);
map.containsKey(k);
map.remove(k);
```

---

## 🔥 Modern Map patterns

### Frequency count

```java id="a9"
map.put(k, map.getOrDefault(k, 0) + 1);
```

### Cleaner Java 8+

```java id="a10"
map.merge(k, 1, Integer::sum);
```

### computeIfAbsent (VERY IMPORTANT)

```java id="a11"
map.computeIfAbsent(k, x -> new ArrayList<>()).add(v);
```

---

## 🔁 Map iteration

```java id="a12"
for (Map.Entry<Integer, Integer> e : map.entrySet()) {
    int k = e.getKey();
    int v = e.getValue();
}
```

```java id="a13"
map.forEach((k, v) -> {
    System.out.println(k + " " + v);
});
```

---

# 🧺 6. HashSet

```java id="a14"
Set<Integer> set = new HashSet<>();

set.add(x);
set.contains(x);
set.remove(x);
```

---

# 🏔 7. PriorityQueue (Heap)

## Min heap

```java id="a15"
PriorityQueue<Integer> pq = new PriorityQueue<>();
```

## Max heap

```java id="a16"
PriorityQueue<Integer> pq =
    new PriorityQueue<>(Collections.reverseOrder());
```

---

## 🔥 Custom comparator (LAMBDA — VERY IMPORTANT)

### Sort by first element

```java id="a17"
PriorityQueue<int[]> pq =
    new PriorityQueue<>((a, b) -> a[0] - b[0]);
```

### Sort by second element descending

```java id="a18"
PriorityQueue<int[]> pq =
    new PriorityQueue<>((a, b) -> b[1] - a[1]);
```

---

## Object-based heap

```java id="a19"
PriorityQueue<Node> pq =
    new PriorityQueue<>((a, b) -> a.cost - b.cost);
```

---

# 🔢 8. Sorting (LAMBDA HEAVY)

## Basic

```java id="a20"
Arrays.sort(arr);
```

## 2D array sorting

```java id="a21"
Arrays.sort(arr, (a, b) -> a[0] - b[0]);
```

## Reverse sort

```java id="a22"
Arrays.sort(arr, (a, b) -> b - a);
```

## List sorting

```java id="a23"
list.sort((a, b) -> a - b);
```

---

# 🧵 9. String + StringBuilder

```java id="a24"
StringBuilder sb = new StringBuilder();

sb.append(x);
sb.toString();
sb.reverse();
sb.deleteCharAt(i);
```

---

# 🔣 10. Character tricks

```java id="a25"
c - 'a'     // index mapping
Character.isDigit(c);
Character.isLetter(c);
```

---

# 🧮 11. Math utilities

```java id="a26"
Math.max(a, b);
Math.min(a, b);
Math.abs(x);
Math.sqrt(x);
```

---

# 📊 12. Prefix Sum

```java id="a27"
int[] prefix = new int[n + 1];

prefix[i + 1] = prefix[i] + arr[i];
```

---

# 🔍 13. Binary Search Template

```java id="a28"
int l = 0, r = n - 1;

while (l <= r) {
    int mid = l + (r - l) / 2;

    if (condition) {
        r = mid - 1;
    } else {
        l = mid + 1;
    }
}
```

---

# 🌐 14. BFS (Queue)

```java id="a29"
Queue<int[]> q = new LinkedList<>();
boolean[][] visited = new boolean[m][n];

q.offer(new int[]{i, j});
visited[i][j] = true;

while (!q.isEmpty()) {
    int[] cur = q.poll();
}
```

---

# 🔁 15. DFS

```java id="a30"
void dfs(int node, boolean[] visited) {
    visited[node] = true;

    for (int nei : graph[node]) {
        if (!visited[nei]) {
            dfs(nei, visited);
        }
    }
}
```

---

# 🔄 16. Backtracking

```java id="a31"
void backtrack(List<Integer> path) {
    if (baseCase) {
        result.add(new ArrayList<>(path));
        return;
    }

    for (choice : choices) {
        path.add(choice);
        backtrack(path);
        path.remove(path.size() - 1);
    }
}
```

---

# 🧵 17. LAMBDA BASICS (IMPORTANT ADDITION)

## Comparator

```java id="a32"
(a, b) -> a - b
(a, b) -> b - a
```

## Multiple conditions

```java id="a33"
(a, b) -> {
    if (a[0] == b[0]) return a[1] - b[1];
    return a[0] - b[0];
}
```

---

## Map lambda iteration

```java id="a34"
map.forEach((k, v) -> {
    System.out.println(k + " " + v);
});
```

---

# 🧵 18. Regex (INTERVIEW-USEFUL)

## Match

```java id="a35"
s.matches("[a-z]+");
```

## Pattern + Matcher

```java id="a36"
Pattern p = Pattern.compile("\\d+");
Matcher m = p.matcher(s);

while (m.find()) {
    System.out.println(m.group());
}
```

## Replace

```java id="a37"
s.replaceAll("\\d+", "#");
```

---

# 📁 19. FILE I/O (RARE BUT COMPLETE)

## Read

```java id="a38"
BufferedReader br = new BufferedReader(new FileReader("input.txt"));

String line;
while ((line = br.readLine()) != null) {
    System.out.println(line);
}
br.close();
```

## Write

```java id="a39"
BufferedWriter bw = new BufferedWriter(new FileWriter("output.txt"));

bw.write("hello");
bw.newLine();
bw.close();
```

## Modern NIO

```java id="a40"
List<String> lines = Files.readAllLines(Paths.get("input.txt"));

Files.write(Paths.get("output.txt"), lines);
```

---

# 🚀 20. HIGH-YIELD ADVANCED PATTERNS

## Frequency array (faster than map)

```java id="a41"
int[] freq = new int[26];
freq[c - 'a']++;
```

---

## Map with default list

```java id="a42"
map.computeIfAbsent(k, x -> new ArrayList<>()).add(v);
```

---

## Heap with lambda object

```java id="a43"
PriorityQueue<int[]> pq =
    new PriorityQueue<>((a, b) -> Integer.compare(a[1], b[1]));
```

---

# 🧠 FINAL SUMMARY (WHAT YOU SHOULD REMEMBER)

## MUST KNOW (core 80%)

* Array / String / StringBuilder
* HashMap / HashSet
* Deque (stack + queue replacement)
* PriorityQueue (heap)
* BFS / DFS / Backtracking
* Binary Search

## VERY IMPORTANT (interview edge)

* Lambdas (`(a,b) -> a-b`)
* `computeIfAbsent`
* `merge`
* custom comparators
* prefix sum

## OPTIONAL BUT POWERFUL

* Regex
* File I/O
* NIO utilities

