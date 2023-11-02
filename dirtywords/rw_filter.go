package dirtywords

import "sync"

type RWDirtyFilter struct {
	datasource []string
	tire       *Trie
	rw         *sync.RWMutex
}

func NewRWDirtyFilter(data []string) *RWDirtyFilter {
	res := &RWDirtyFilter{datasource: data, tire: NewTrie(), rw: &sync.RWMutex{}}
	res.tire.InsertAll(res.datasource)
	return res
}

func (d *RWDirtyFilter) Insert(word string) {
	d.rw.Lock()
	d.datasource = append(d.datasource, word)
	d.tire.Insert(word)
	d.rw.Unlock()
}

func (d *RWDirtyFilter) InsertBatch(words []string) {
	d.rw.Lock()
	d.datasource = append(d.datasource, words...)
	d.tire.InsertAll(words)
	d.rw.Unlock()
}

// Search 查询词word是否是敏感词
func (d *RWDirtyFilter) Search(word string) bool {
	d.rw.RLock()
	defer d.rw.RUnlock()
	return d.tire.Search(word)
}

// Replace 用占位符替换文本中的敏感词，并返回替换后的文本
func (d *RWDirtyFilter) Replace(text string, placeholder rune) string {
	d.rw.RLock()
	defer d.rw.RUnlock()
	return d.tire.Replace(text, placeholder)
}

// QueryAll 查询出文本text中全部的敏感词
func (d *RWDirtyFilter) QueryAll(text string) ([]string, bool) {
	d.rw.RLock()
	defer d.rw.RUnlock()
	return d.tire.QueryAll(text)
}

// Query 查询文本text中是否存在敏感词
func (d *RWDirtyFilter) Query(text string) bool {
	d.rw.RLock()
	defer d.rw.RUnlock()
	return d.tire.Query(text)
}
