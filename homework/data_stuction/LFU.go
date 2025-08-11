package data

type LFUCache struct {
	have     map[int]*Node
	capacity int
	tail     *Node
	head     *Node
}

type Node struct {
	key   int
	f     int
	value int
	next  *Node
	pre   *Node
}

func Constructor(capacity int) LFUCache {
	head := &Node{value: -1}
	tail := &Node{value: -1}
	head.next = tail
	tail.pre = head
	return LFUCache{
		have:     make(map[int]*Node, capacity),
		capacity: capacity,
		head:     head,
		tail:     tail,
	}
}

func (this *LFUCache) Get(key int) int {
	val, ok := this.have[key]
	if !ok {
		return -1
	}
	val.f++
	this.Order(val)
	return val.value
}

func (this *LFUCache) Put(key int, value int) {
	if this.capacity == 0 {
		return
	}
	if val, ok := this.have[key]; ok {
		val.f++
		val.value = value
		this.Order(val)
		return
	}
	if len(this.have) == this.capacity {
		this.Delete()
	}
	node := &Node{
		key:   key,
		value: value,
		f:     0,
	}
	this.Insert(node)
	this.have[key] = node
}

func (this *LFUCache) Order(val *Node) {
	// 先从链表中摘除当前节点
	val.pre.next = val.next
	val.next.pre = val.pre

	// 找到合适位置，频率降序排列
	temp := this.head.next
	for temp != this.tail && temp.f > val.f {
		temp = temp.next
	}

	// 插入到 temp 前面
	val.next = temp
	val.pre = temp.pre
	temp.pre.next = val
	temp.pre = val
}

func (this *LFUCache) Insert(val *Node) {
	val.f = 0
	temp := this.tail.pre
	for temp != this.head && temp.f == 0 {
		temp = temp.pre
	}

	// 插入到 temp 后面
	val.next = temp.next
	val.pre = temp
	temp.next.pre = val
	temp.next = val
}

func (this *LFUCache) Delete() {
	if this.tail.pre == this.head {
		return
	}
	temp := this.tail.pre
	temp.pre.next = this.tail
	this.tail.pre = temp.pre
	delete(this.have, temp.key)
}

func main() {

}
