## Day 17 notes:

果然又卡在了(图论)最短路径，总结一下这次卡住的原因:

> 想试下 A* 发现配上 walk limit distance 太傻了，还是换回 dijkstra 吧

1. Go 的 priority queue(min heap) 实现有点奇怪，实现了 `heap.Interface` 之后，并不能直接使用相关的 heap method , 导致调试了半天原来是 min heap 相关的操作(push/pop)全部没有生效，看了[官方的例子](https://pkg.go.dev/container/heap#example-package-IntHeap)之后，发现应该直接用`heap.Init`, `heap.Push` 和 `heap.Pop` 来操作, 然后把 wrapper 作为参数传进去，应该不是第一次被这个坑了，得把它塞到 [magejiAoC](https://github.com/magejiCoder/magejiAoc/tree/master/queue/priority) 里面去。

2. Loop or Recurse ? 这个其实不关键，关键在于对 heap 的操作。

3. 这次相对于前几年的裸最短路径题，加了 max walk blocks 的限制 (也符合今年的偏难的基调)，对于这个机制的实现应该使用**向量**在矩阵上操作比较快，我的做法是保留**一个map记录还剩下的blocks**: `map[(up: blocks ,down: blocks,left: blocks,right: blocks)]`，每次前进就让对应方向的 block - 1 ，转弯的话就重置其他 block , 转过去的方向 block - 1，不允许回头。

4. 关于**状态**的定义: dijkstra 算法中定义了状态的概念避免重复轮询整个图，一般以每个坐标点(x,y)作为状态，依靠 min heap 的机制，可以做到先pop的一定是最小的，从而避免其它经过此点的路径。但是这个**状态**并不适用于这个题目，这里的 **状态** 需要考虑 limit blocks , 比如虽然在一个坐标点的值是其他经过此点坐标的路径中值最小的，但是，因为这个点可能是被 limit block 机制 "逼" 过去的，所以可能存在其他的路径比现在的更优。 BUT ! 一旦我们把 limit block 加到状态中去考虑，就会得到新的唯一约束: `在当前点如果拥有相同的 block map ` 和 `此时的前进方向` ，然后我们的状态值就从 `struct{x,y}`  变成了 `struct{x,y,blockMap,direction}` !

5. (P2) (我是没想到P2我也能卡这么久的，最近真是状态烂爆了QmQ): P2 对于 walk 的机制做了改动，**先走4格**，再决定**转向**或者是**继续想前走，直到10格**，依旧可以忽略回头的情况。相对于 P1 其实只加了 edge case 的判断:

	* **一定要走完4格才能决定是否要转向** (<- 卡在这里，看栗子上是在第四格转弯了，然后栗子都过了，答案就是一直WA 😫, 后来仔细看了下，发现转弯口算是第一个点 🤡)，也就是说如果最终的终点如果是在前4格(包括转向)上，那么这个终点是不成立的。
	* **10格之后不能再往前走**
	* **4 ~ 10 格之间可以选择随时转向**
	* **依旧需要按原计划构造 heap，帮助整个图的遍历**

6. (题外话) 这种图论题很难 debug , 特别是在小数据集上表现正常，而在稍大数据集上发生错误(WA), 此时，需要按照基本法对大数据集进行拆分，用二分来 debug 会快很多。
