# 获取Gid的方法

- go官方使用案例：**src/net/http/h2_bundle.go**，方法签名：func http2curGoroutineID() uint64

```golang
package http

var http2goroutineSpace = []byte("goroutine ")

func http2curGoroutineID() uint64 {
	bp := http2littleBuf.Get().(*[]byte)
	defer http2littleBuf.Put(bp)
	b := *bp
	b = b[:runtime.Stack(b, false)]
	// Parse the 4707 out of "goroutine 4707 ["
	b = bytes.TrimPrefix(b, http2goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		panic(fmt.Sprintf("No space found in %q", b))
	}
	b = b[:i]
	n, err := http2parseUintBytes(b, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse goroutine ID out of %q: %v", b, err))
	}
	return n
}

var http2littleBuf = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 64)
		return &buf
	},
} 
```

## 实现原理

- 使用runtime.Stack()获取当前Goroutine的栈信息，从栈信息中获取go id(goid)
- 读取的栈关键信息格式如下：goroutine id号，示例 goroutine 1，携程id为1
- bytes.TrimPrefix 去掉goroutine 前缀，得到以goid起始的字符串
- bytes.IndexByte 找到字符串的第一个空格位置，截取字符串得到goid字符串
- 字符串goid转为int64数字goid