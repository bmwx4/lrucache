package lrucache

import (
	"errors"
)

//Element 节点元素描述信息
type Element struct {
	Key int
	Val int
}

//ListNode 用于描述操作数据(node)信息，并且应用双向链表实现快速的添加和删除 node
type ListNode struct {
	*Element
	Next *ListNode
	Pre  *ListNode
}

//LinkedList 双向链表，用于快速更新节点的位置
type LinkedList struct {
	Head *ListNode //当有新数据被访问，或新增数据时，可以快速存放该数据(node)到链表的最上层
	Last *ListNode //当超出 cache 容量时，可以快速删除最近最少使用一个，也就是链表的最后一个
}

// LRUCache  构建 LRUCache 结构体，使用hash map 来实现O(1)时间复杂度查询key是否存在，value 是该entry所属的指针
type LRUCache struct {
	M   map[int]*ListNode
	Cap int // cache 容量
	Len int // cache 中当前key数量
	L   LinkedList
}

// NewLRUCache 创建一个 LRUCache 实例
func NewLRUCache(cap int) *LRUCache {
	return &LRUCache{
		M:   make(map[int]*ListNode, cap),
		Cap: cap,
		Len: 0,
		L: LinkedList{
			Head: nil,
			Last: nil,
		},
	}
}

//Get 查询 cache,
//cache 为空，返回 -1,没有命中,返回 -1,并返回相应的错误信息
//命中返回对应的 value,并尝试把 node 移动到链表的表头
func (c *LRUCache) Get(key int) (int, error) {
	if c.Len == 0 || c.Cap == 0 {
		return -1, errors.New("Cache is empty or capcity is 0") // cache 为空，返回 -1
	}

	res, ok := c.M[key]
	if !ok { // 如果找不到直接返回 -1
		return -1, errors.New("Not Found")
	}

	if c.L.Head == res { // 查找到key，并且该node在链表表头，直接返回其值
		return res.Val, nil
	}

	newHead := &ListNode{ // 如果该node不在表头的位置, 则需要进行一次移动，先定义一个 头节点
		Element: &Element{
			Val: res.Val,
			Key: key,
		},
		Pre: nil,
	}

	c.move(res, newHead, key) // 进行分析和挪动
	return res.Val, nil
}

//move 操作双向链表添加和移除节点
func (c *LRUCache) move(res *ListNode, newHead *ListNode, key int) {
	// 查询节点在在末尾
	if res == c.L.Last {
		// 把数放到第一位,NEXT指针指向之前的HEAD节点,
		// 长度因为没有超出,所以把该数的上一位设置成最后一个节点,
		oldHead := c.L.Head        //保存旧的 HEAD,防止被覆盖
		oldHead.Pre = newHead      // 新节点将作为 旧HEAD节点的Pre节点
		newHead.Next = oldHead     // 旧节点将作为 新节点的Next节点
		res.Pre.Next = nil         // 使用原尾部节点的Pre节点进行重置尾部节点,需要将其NEXT设置为nil
		c.M[res.Pre.Key] = res.Pre //重置索引信息
		c.L.Last = res.Pre
		c.M[key] = newHead
		c.L.Head = newHead
		return
	}

	// 查询的节点在链表的中间位置,则首先把该node放到HEAD位置,然后把该node之前的Next节点指向它的Next节点
	oldHead := c.L.Head
	oldHead.Pre = newHead
	newHead.Next = oldHead
	c.M[key].Pre.Next = c.M[key].Next      //把该node之前的Next节点指向它的Next节点
	c.M[key].Next.Pre = c.M[key].Pre       //把该node之前的Pre节点指向它的Next节点
	c.M[c.M[key].Next.Key] = c.M[key].Next //重置索引信息,保存节点指针到hash map
	c.M[oldHead.Key] = oldHead
	c.M[c.M[key].Pre.Key] = c.M[key].Pre
	c.M[key] = newHead
	c.L.Head = newHead
	return
}

//Put 添加新节点到 list，并保存node指针到hash map, true 代表添加成功，false 代表失败，并反回相应错误信息
// 如果cache 容量为0，则直接返回错误,
// 如果当前 cache 中节点数量为0，则直接存入节点,
// 如果新放入的key==>value在节点中存在，则更新value和node的值,
// 如果当前的key没有找到,那么将此值放到链表头位置,更新map，并判断链表总长度,如果超出容量,就删掉最后一个
func (c *LRUCache) Put(key int, value int) (bool, error) {
	if c.Cap <= 0 {
		return false, errors.New("Cache capcity is 0")
	}
	if c.Len == 0 {
		newHead := &ListNode{
			Element: &Element{
				Val: value,
				Key: key,
			},
			Pre:  nil,
			Next: nil,
		}
		c.L.Head = newHead
		c.L.Last = newHead
		c.M[key] = newHead
		c.Len++
	}

	if res, ok := c.M[key]; ok {
		if c.L.Head == res { //如果在 表头的位置，则直接更新value，否则进行移动
			res.Val = value
			c.M[key] = res
			return true, nil
		}
		newHead := &ListNode{
			Element: &Element{
				Val: value,
				Key: key,
			},
			Pre: nil,
		}
		c.move(res, newHead, key)
		return true, nil
	}

	newHead := &ListNode{
		Element: &Element{
			Val: value,
			Key: key,
		},
		Pre: nil,
	}
	oldHead := c.L.Head   // 将新节点放到头位置
	oldHead.Pre = newHead //与之前的头结点交换位置
	newHead.Next = oldHead
	c.L.Head = newHead
	c.M[oldHead.Key] = oldHead //更新map
	c.M[key] = newHead

	c.Len++ // 如果当前链表长度超过总容量了,就删掉最后一个节点,之前最后一个节点的上一个节点设置成最后一个节点
	if c.Len > c.Cap {
		c.L.Last.Pre.Next = nil
		delete(c.M, c.L.Last.Key)
		c.L.Last = c.L.Last.Pre
		c.Len--
	}
	return true, nil
}

//DumpKeys 查询所有的key
func (c *LRUCache) DumpKeys() []Element {
	s := make([]Element, c.Len)
	i := 0
	for k, v := range c.M {
		s[i] = Element{
			Key: k,
			Val: v.Val,
		}
		i++
	}
	return s
}

//CleanLRUCache 清空lrucache
func (c *LRUCache) CleanLRUCache() {
	for k := range c.M {
		c.M[k].Next = nil //防止内存泄漏
		c.M[k].Pre = nil
		c.M[k].Element = nil
		delete(c.M, k)
	}
}
