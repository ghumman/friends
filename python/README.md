# Python Arrays (Lists) Quick Guide

## Create

```python
arr = [1, 2, 3]
```

## Access

```python
print(arr[0])    # first item
print(arr[-1])   # last item
```

## Change

```python
arr[1] = 10
```

## Add

```python
arr.append(4)        # add at end
arr.insert(1, 99)    # add at index
```

## Delete

```python
arr.remove(10)   # remove by value
del arr[1]       # remove by index
arr.pop()        # remove last item
```

## Length

```python
len(arr)
```

## Loop

```python
for item in arr:
    print(item)
```

## Check Exists

```python
if 2 in arr:
    print("Found")
```

## Slice

```python
arr[1:3]
```

## Sort / Reverse

```python
arr.sort()
arr.reverse()
```

## Combine

```python
c = a + b
```

## Useful Functions

```python
min(arr)
max(arr)
sum(arr)
```

## 2D Array

```python
matrix = [[1,2],[3,4]]

print(matrix[0][1])
```

## Dictionary / Map
location_dict = {}
i = location_dict[complemnt]
location_dict[num] = i

## For Loop
1. for i, num in enumerate(nums):
2. for i in range(len(s)):
3. for ch in s: // s is a string

## While Loop
```
i = 0

while i < 5:
    print(i)
    i += 1
```

## Queues and Stacks in Python
# Stack (LIFO = Last In First Out)

Use a list.

## Create Stack

```python id="xiwtnc"
stack = []
```

## Push

```python id="p6e17d"
stack.append(10)
stack.append(20)
```

## Pop

```python id="6e11g2"
item = stack.pop()
print(item)   # 20
```

## Peek Top

```python id="n2yq1d"
print(stack[-1])
```

---

# Queue (FIFO = First In First Out)

Use `deque`.

```python id="j7j0ko"
from collections import deque

queue = deque()
```

## Enqueue

```python id="3b5kqb"
queue.append(10)
queue.append(20)
```

## Dequeue

```python id="9wuksm"
item = queue.popleft()

print(item)   # 10
```

## Front Item

```python id="4lx4i6"
print(queue[0])
```

---

# Common Operations

```python id="ucd4su"
len(stack)
len(queue)

if not stack:
    print("Empty")
```


## Sets
# Python Sets

Sets store **unique values** and are unordered.

## Create

```python id="b8u9c2"
s = {1, 2, 3}

empty = set()
```

## Add

```python id="5f4b2i"
s.add(4)
```

## Remove

```python id="j6n2w1"
s.remove(2)     # error if not found
s.discard(5)    # no error
```

## Check Exists

```python id="o7k3l9"
if 3 in s:
    print("Found")
```

## Loop

```python id="d4m8q0"
for item in s:
    print(item)
```

## Length

```python id="x1p7v6"
len(s)
```

## Set Operations

```python id="q9r2y5"
a = {1, 2, 3}
b = {3, 4, 5}

print(a | b)   # union
print(a & b)   # intersection
print(a - b)   # difference
```

## Remove Duplicates

```python id="m3k8n1"
nums = [1, 1, 2, 3, 3]

unique = set(nums)

print(unique)
# {1, 2, 3}
```
