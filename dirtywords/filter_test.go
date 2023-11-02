package dirtywords

import (
	"encoding/csv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewDirtyFilter(t *testing.T) {
	f, err := os.Open("./dirtywords.csv")
	assert.Nil(t, err, err)
	defer f.Close()
	reader := csv.NewReader(f)
	rows, err := reader.ReadAll()
	assert.Nil(t, err, err)
	data := make([]string, 0, 10000)
	for _, row := range rows {
		for _, val := range row {
			data = append(data, val)
		}
	}
	filter := NewDirtyFilter(data)
	assert.Equal(t, true, filter.Query("你好啊张三个测试个贵大头"), "错误")
	assert.Equal(t, false, filter.Query("你好啊张哥三个测哥试个贵大头"), "错误")
	text := "你好啊，张三个挨揍皮带哥三个挂测试了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个谋杀啊"
	assert.Equal(t, "你好啊，张三个挨揍皮带哥三个挂**了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个**啊", filter.Replace(text, '*'), "错误")
}

func BenchmarkDirtyFilter_Replace2(b *testing.B) {
	f, err := os.Open("./dirtywords.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	reader := csv.NewReader(f)
	rows, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	data := make([]string, 0, 10000)
	for _, row := range rows {
		for _, val := range row {
			data = append(data, val)
		}
	}
	filter := NewDirtyFilter(data)
	text := "你好啊，张三个挨揍皮带哥三个挂测试了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a张三个挨揍皮带哥三个挂了啊李焕英Abc就好你好啊李焕英a"
	for i := 0; i < b.N; i++ {
		filter.Replace(text, '*')
	}
}

func BenchmarkDirtyFilter_Replace(b *testing.B) {
	f, err := os.Open("./dirtywords.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	reader := csv.NewReader(f)
	rows, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	data := make([]string, 0, 10000)
	for _, row := range rows {
		for _, val := range row {
			data = append(data, val)
		}
	}
	filter := NewDirtyFilter(data)
	texts := []string{
		"谋杀是英雄水电费水电费水电费水电费水电费",
		"造反是英雄水电费水电费水电费水电费水电费",
		"赌马是英雄水电费水电费水电费水赌马电费水电费",
		"是英雄水下注电费水电费水电费水电费水电费",
		"是英雄水电费水下注电费水电费水电费水电费",
		"谋杀是英雄水电费水电费水电费水电费水电费",
		"是英雄水电费水电费水押大电费水电费水电费",
		"是英雄水电费水电费杜冷丁水电费水电费水电费",
		"是英雄水电费水电费水老虎机电坐庄费水电费水电费",
		"是英雄水电费水电费水电费水电费水电费",
		"是英雄水电费水电安非他命费水电费水电费中华人民实话实说水电费",
		"是英雄水电费水电费释迦牟尼水电费水电费水电费穆罕默德",
		"是英雄水电费水电费赌球水电费不文小丈夫之银座嬉春水电费水电费",
		"是英雄水电费水电费水电费水电费水电费穆罕默德",
		"是英雄水电费水电费水艺坛照妖镜之96应召名册电费水电费水电费",
		"是英雄水电费水电费水电费水电费水电费",
		"是英雄水电费水电不文小丈夫之银座嬉春费水电费水电费水电费",
		"是英雄水电费水电费水电费水电费水电不文小丈夫之银座嬉春费中华人民实话实说",
		"中华人民实话实说是英雄水电费水电费水电费水电费水电费",
		"是英雄水电费水电费水电费聊斋之欲焰三娘子水电费水电费",
		"偷窥洗澡是英雄水聊斋之欲焰三娘子电费水电费水电费水电费水电费",
		"是英雄水电费水电费水电费水电费水电费偷窥洗澡",
		"是英雄水电费水电费水电偷窥洗澡费水电费水电费",
	}
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(texts); j++ {
			filter.Replace(texts[j], '*')
		}
	}
}
