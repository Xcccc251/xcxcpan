package hashRing

import (
	"crypto/sha256"
	"sort"
	"strconv"
)

// 一致性哈希实现负载均衡
var Hash = NewConsistentHash(5, nil)

// HashFunc 定义哈希函数类型
type HashFunc func(data []byte) uint32

// ConsistentHash 定义一致性哈希结构
type ConsistentHash struct {
	replicas int               // 虚拟节点倍数
	hashFunc HashFunc          // 哈希函数
	ring     []uint32          // 哈希环
	nodes    map[uint32]string // 虚拟节点到物理节点的映射
}

// NewConsistentHash 创建一致性哈希对象
func NewConsistentHash(replicas int, hashFunc HashFunc) *ConsistentHash {
	if hashFunc == nil {
		// 默认使用 SHA256
		hashFunc = func(data []byte) uint32 {
			hash := sha256.Sum256(data)
			return uint32(hash[0])<<24 | uint32(hash[1])<<16 | uint32(hash[2])<<8 | uint32(hash[3])
		}
	}
	return &ConsistentHash{
		replicas: replicas,
		hashFunc: hashFunc,
		ring:     []uint32{},
		nodes:    make(map[uint32]string),
	}
}

// Add 添加节点
func (c *ConsistentHash) Add(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < c.replicas; i++ {
			// 创建虚拟节点的哈希值
			virtualNode := c.hashFunc([]byte(node + strconv.Itoa(i)))
			c.ring = append(c.ring, virtualNode)
			c.nodes[virtualNode] = node
		}
	}
	// 排序哈希环
	sort.Slice(c.ring, func(i, j int) bool {
		return c.ring[i] < c.ring[j]
	})
}

// Remove 移除节点
func (c *ConsistentHash) Remove(node string) {
	for i := 0; i < c.replicas; i++ {
		virtualNode := c.hashFunc([]byte(node + strconv.Itoa(i)))
		index := c.findIndex(virtualNode)
		// 删除对应虚拟节点
		if index != -1 {
			c.ring = append(c.ring[:index], c.ring[index+1:]...)
			delete(c.nodes, virtualNode)
		}
	}
}

// Get 获取数据对应的节点
func (c *ConsistentHash) Get(key string) string {
	if len(c.ring) == 0 {
		return ""
	}
	hash := c.hashFunc([]byte(key))
	index := c.findIndex(hash)
	return c.nodes[c.ring[index]]
}

// findIndex 找到顺时针第一个匹配的节点
func (c *ConsistentHash) findIndex(hash uint32) int {
	// 使用二分查找提升效率
	idx := sort.Search(len(c.ring), func(i int) bool {
		return c.ring[i] >= hash
	})
	if idx == len(c.ring) {
		return 0 // 环形回绕
	}
	return idx
}
