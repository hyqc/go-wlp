package dirtywords

type DirtyFilter struct {
	datasource []string
	tire       *Trie
}

func NewDirtyFilter(data []string) *DirtyFilter {
	res := &DirtyFilter{datasource: data, tire: NewTrie()}
	res.tire.InsertAll(res.datasource)
	return res
}

// Search 查询词word是否是敏感词
func (d *DirtyFilter) Search(word string) bool {
	return d.tire.Search(word)
}

// Replace 用占位符替换文本中的敏感词，并返回替换后的文本
func (d *DirtyFilter) Replace(text string, placeholder rune) string {
	return d.tire.Replace(text, placeholder)
}

// QueryAll 查询出文本text中全部的敏感词
func (d *DirtyFilter) QueryAll(text string) ([]string, bool) {
	return d.tire.QueryAll(text)
}

// Query 查询文本text中是否存在敏感词
func (d *DirtyFilter) Query(text string) bool {
	return d.tire.Query(text)
}
