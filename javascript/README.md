# 🚀 JavaScript LeetCode Master Cheat Sheet (ES6+ + Interview Ready)

---

# 🧱 1. Arrays

```javascript
const arr = new Array(n);
const arr = [1, 2, 3];

arr.length;

arr.push(x);
arr.pop();

arr.shift();
arr.unshift(x);

arr.slice(start, end);
arr.splice(index, count);

arr.fill(value);
arr.sort((a, b) => a - b);

const copy = [...arr];
const copy2 = Array.from(arr);
```

---

# 📋 2. Array Methods (MOST IMPORTANT)

```javascript
arr.map(x => x * 2);

arr.filter(x => x > 0);

arr.reduce((sum, x) => sum + x, 0);

arr.forEach(x => console.log(x));

arr.find(x => x === target);

arr.some(x => x > 10);

arr.every(x => x >= 0);

arr.includes(x);

arr.indexOf(x);
```

---

# 🔗 3. Stack (Use Array)

```javascript
const stack = [];

stack.push(x);

stack.pop();

stack.at(-1);
// or
stack[stack.length - 1];
```

---

# 🚀 4. Queue (Simple)

```javascript
const queue = [];

queue.push(x);

const front = queue.shift();
```

---

## ⚡ Efficient Queue (Recommended)

```javascript
const queue = [];
let head = 0;

queue.push(x);

const front = queue[head++];
```

---

# 🔄 5. Deque (Not Built-in)

For interviews, use:

```javascript
const deque = [];

deque.push(x);
deque.pop();

deque.unshift(x);
deque.shift();
```

For performance-critical problems, implement a custom deque.

---

# 🚀 6. Map (MOST IMPORTANT)

```javascript
const map = new Map();

map.set(key, value);

map.get(key);

map.has(key);

map.delete(key);

map.size;

map.clear();
```

---

## 🔥 Frequency Count

```javascript
map.set(x, (map.get(x) || 0) + 1);
```

---

## Grouping

```javascript
if (!map.has(key)) {
    map.set(key, []);
}

map.get(key).push(value);
```

---

## Iteration

```javascript
for (const [k, v] of map) {

}
```

```javascript
map.forEach((value, key) => {

});
```

---

# 🧺 7. Set

```javascript
const set = new Set();

set.add(x);

set.has(x);

set.delete(x);

set.size;
```

---

# 🏔 8. Priority Queue (Heap)

JavaScript has **no built-in heap**.

LeetCode usually provides:

```javascript
const pq = new MinPriorityQueue();

pq.enqueue(value, priority);

pq.dequeue().element;

pq.front().element;

pq.size();
```

Or implement your own Binary Heap.

---

# 🔢 9. Sorting

## Numbers

```javascript
arr.sort((a, b) => a - b);
```

---

## Descending

```javascript
arr.sort((a, b) => b - a);
```

---

## Objects

```javascript
arr.sort((a, b) => a.cost - b.cost);
```

---

## Multiple Conditions

```javascript
arr.sort((a, b) => {
    if (a[0] === b[0]) {
        return a[1] - b[1];
    }
    return a[0] - b[0];
});
```

---

# 🧵 10. Strings

```javascript
const s = "hello";

s.length;

s.slice(1, 4);

s.substring(1, 4);

s.split("");

chars.join("");

s.includes("abc");

s.startsWith("ab");

s.endsWith("yz");

s.repeat(3);
```

---

# 🔣 11. Character Tricks

```javascript
c.charCodeAt(0);

'a'.charCodeAt(0);

String.fromCharCode(97);

const idx = c.charCodeAt(0) - 97;
```

---

# 🧮 12. Math Utilities

```javascript
Math.max(a, b);

Math.min(a, b);

Math.abs(x);

Math.sqrt(x);

Math.floor(x);

Math.ceil(x);

Math.random();
```

---

# 📊 13. Prefix Sum

```javascript
const prefix = new Array(n + 1).fill(0);

for (let i = 0; i < n; i++) {
    prefix[i + 1] = prefix[i] + nums[i];
}
```

---

# 🔍 14. Binary Search Template

```javascript
let left = 0;
let right = nums.length - 1;

while (left <= right) {

    const mid = Math.floor((left + right) / 2);

    if (condition) {
        right = mid - 1;
    } else {
        left = mid + 1;
    }
}
```

---

# 🌐 15. BFS

```javascript
const queue = [[i, j]];
let head = 0;

const visited = new Set();

while (head < queue.length) {

    const [r, c] = queue[head++];

}
```

---

# 🔁 16. DFS

```javascript
function dfs(node) {

    visited.add(node);

    for (const nei of graph[node]) {

        if (!visited.has(nei)) {
            dfs(nei);
        }

    }

}
```

---

# 🔄 17. Backtracking

```javascript
const result = [];

function backtrack(path) {

    if (baseCase) {
        result.push([...path]);
        return;
    }

    for (const choice of choices) {

        path.push(choice);

        backtrack(path);

        path.pop();

    }

}
```

---

# 🧵 18. Arrow Functions

## Basic

```javascript
const add = (a, b) => a + b;
```

---

## Sort

```javascript
(a, b) => a - b

(a, b) => b - a
```

---

## Multiple Statements

```javascript
(a, b) => {

    if (a.cost === b.cost) {
        return a.id - b.id;
    }

    return a.cost - b.cost;

}
```

---

# 🧵 19. Regular Expressions

## Match

```javascript
/[a-z]+/.test(s);
```

---

## Find

```javascript
const matches = s.match(/\d+/g);
```

---

## Replace

```javascript
s.replace(/\d+/g, "#");
```

---

# 📁 20. File I/O (Node.js)

## Read

```javascript
const fs = require("fs");

const input = fs.readFileSync(0, "utf8");
```

---

## Write

```javascript
console.log(answer);
```

---

# 🚀 21. HIGH-YIELD ADVANCED PATTERNS

## Frequency Array

```javascript
const freq = new Array(26).fill(0);

freq[c.charCodeAt(0) - 97]++;
```

---

## Object Frequency

```javascript
const freq = {};

freq[x] = (freq[x] || 0) + 1;
```

---

## Map Frequency

```javascript
const freq = new Map();

freq.set(x, (freq.get(x) || 0) + 1);
```

---

## Clone Array

```javascript
const copy = [...arr];
```

---

## Clone 2D Array

```javascript
const copy = arr.map(row => [...row]);
```

---

## Swap

```javascript
[arr[i], arr[j]] = [arr[j], arr[i]];
```

---

## Destructuring

```javascript
const [a, b] = pair;

const { x, y } = point;
```

---

## Nullish Coalescing

```javascript
const value = map.get(key) ?? 0;
```

---

## Optional Chaining

```javascript
obj?.user?.name;
```

---

# 🧠 FINAL SUMMARY (WHAT YOU SHOULD REMEMBER)

## MUST KNOW (Core 80%)

* Arrays & Array methods (`map`, `filter`, `reduce`)
* Strings
* `Map`
* `Set`
* Stack (`push` / `pop`)
* Queue (head pointer, avoid repeated `shift()`)
* Sorting with comparators
* BFS / DFS / Backtracking
* Binary Search

## VERY IMPORTANT (Interview Edge)

* Arrow functions
* Destructuring
* Spread operator (`...`)
* Frequency counting with `Map`
* Prefix sums
* Multiple-condition sorting
* `??` (Nullish Coalescing)
* `?.` (Optional Chaining)

## OPTIONAL BUT POWERFUL

* Priority Queue / Binary Heap implementation
* Regular Expressions
* Node.js File I/O
* Generators & Iterators
* Typed Arrays (`Uint8Array`, `Int32Array`) for performance