package geecache

// A ByteView holds in immutable view of bytes.
type ByteView struct {
	// 存储真实的缓存值
	b []byte
}

// Len returns the view's length.
// 实现Len 方法 [被缓存对象需要实现]
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice returns a copy of the data as a byte slice.
// 返回一个拷贝，防止缓存被外部程序修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String returns the data as a string, making a copy if necessary.
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
