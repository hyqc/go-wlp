package dirtywords

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTrie(t *testing.T) {
	n := NewTrie()
	n.Insert("你好啊李焕英")
	n.Insert("李焕英")
	n.Insert("王三")
	n.Insert("张三")
	n.Insert("李四")
	n.Insert("三个")
	n.InsertAll([]string{"王八蛋", "八嘎", "SB", "大B"})
	assert.Equal(t, true, n.Search("张三"), "错误")
	assert.Equal(t, true, n.Search("李焕英"), "错误")
	assert.Equal(t, false, n.Search("张三三"), "错误")
	text := "张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a"
	res, ok := n.searchDirtyWordIndexes([]rune(text), 0)
	assert.Equal(t, true, ok, "错误")
	assert.Equal(t, []int{0, 1}, res, "错误")

	text = "你好啊，张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a"
	res, ok = n.searchDirtyWordIndexes([]rune(text), 0)
	assert.Equal(t, false, ok, "错误")
	assert.Nil(t, res, "错误")

	ress, ok := n.QueryAll(text)
	assert.Equal(t, true, ok, "错误")
	assert.Equal(t, []string{"张三", "三个", "三个", "李焕英", "你好啊李焕英", "李焕英"}, ress, "错误")

	assert.Equal(t, true, n.Query(text), "错误")
	assert.Equal(t, false, n.Query("你好歹啊的大多数地方"), "错误")

	assert.Equal(t, "你好啊，***挨揍皮带哥**挂了啊***Abc就好******a", n.Replace(text, '*'), "错误")
}

func BenchmarkTrie_Query(b *testing.B) {
	n := NewTrie()
	n.Insert("你好啊李焕英")
	n.Insert("李焕英")
	n.Insert("王三")
	n.Insert("张三")
	n.Insert("李四")
	n.Insert("三个")
	n.InsertAll([]string{"王八蛋", "八嘎", "SB", "大B"})
	text := "你好啊，张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a"
	for i := 0; i < b.N; i++ {
		n.Replace(text, '*')
	}
}
