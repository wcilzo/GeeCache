package geecache

import "container/list"

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	// 允许使用的最大内存
	maxBytes int64
	// 当前已经使用的内存
	nbytes int64
	// Go 标准库实现的双向链表 list.List
	ll *list.List
	// map 内存映射
	cache map[string]*list.Element
	// optional and executed when an entry is purged.
	// 某条记录被移除时的回调函数 可以为 nil
	OnEvicted func(key string, value Value)
}

// 双向链表节点的数据类型，保存key可以在删除队首节点的时候，直接用key删除
type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
// 用于返回值所占用的内存大小
type Value interface {
	Len() int
}

// 实例化
// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get look ups a key's value
// 查找
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		// 将链表中的节点ele 移动到队尾，约定front为队尾
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest time
// 删除
func (c *Cache) RemoveOldest() {
	// 获取队首节点
	ele := c.ll.Back()
	// 如果存在
	if ele != nil {
		// 从队列移除元素
		c.ll.Remove(ele)
		// 获取元素
		kv := ele.Value.(*entry)
		// 从字典c.cache删除该节点的映射关系
		delete(c.cache, kv.key)
		// 更新当前占用内粗
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		// 如果回调函数不为 nil ，则调用回调函数
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add adds a value to the cache.
// 新增/修改
func (c *Cache) Add(key string, value Value) {
	// 获取缓存中的元素
	// 如果存在，那么更新元素
	if ele, ok := c.cache[key]; ok {
		// 当前元素被利用，移动到队尾
		c.ll.MoveToFront(ele)
		// 获取元素值
		kv := ele.Value.(*entry)
		// 计算更新元素后当前占用的内存
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		// 替换旧值，更新为新值
		kv.value = value
	} else {
		// 如果不存在，新增
		ele := c.ll.PushFront(&entry{key, value})
		// 插入到内存map
		c.cache[key] = ele
		// 重新计算当前占用内存 key + value
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	// 如果超出了最大值，则删除队首元素
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
