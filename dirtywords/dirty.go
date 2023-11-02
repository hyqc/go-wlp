package dirtywords

type TrieNode struct {
	isEnd    bool
	children map[rune]*TrieNode
}

type Trie struct {
	root *TrieNode
}

type ITrie interface {
	Insert(word string)                           // 插入一个敏感词
	InsertAll(words []string)                     // 插入一批敏感词
	Search(word string) bool                      // 传入的词是否是敏感词
	Replace(text string, placeholder rune) string // 将文本中的敏感词替换词指定的字符
	QueryAll(text string) ([]string, bool)        // 查出文本中的全部敏感词
	Query(text string) bool                       // 查出文本中是否存在任意敏感词
}

func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{
			isEnd:    false,
			children: make(map[rune]*TrieNode),
		},
	}
}

func (t *Trie) Insert(word string) {
	cur := t.root
	for _, w := range word {
		if _, ok := cur.children[w]; !ok {
			cur.children[w] = &TrieNode{
				isEnd:    false,
				children: make(map[rune]*TrieNode),
			}
		}
		cur = cur.children[w]
	}
	cur.isEnd = true
}

func (t *Trie) InsertAll(words []string) {
	for _, w := range words {
		t.Insert(w)
	}
}

// Search 词word是否是敏感词
func (t *Trie) Search(word string) bool {
	cur := t.root
	for _, w := range word {
		if _, ok := cur.children[w]; !ok {
			return false
		}
		cur = cur.children[w]
	}
	return cur.isEnd
}

// Replace 用替换符替换文本text中全部的敏感词，并返回替换后的文本
func (t *Trie) Replace(text string, placeholder rune) string {
	node := t.root
	words := []rune(text)
	res := make([]int, 0, len(words))
	for i, w := range words {
		if _, ok := node.children[w]; ok {
			if tmp, is := t.searchDirtyWordIndexes(words[i:], i); is {
				res = append(res, tmp...)
			}
		}
	}
	for _, i := range res {
		words[i] = placeholder
	}
	return string(words)
}

func (t *Trie) searchDirtyWordIndexes(text []rune, start int) ([]int, bool) {
	cur := t.root
	res := make([]int, 0, len(text)/4)
	for i, w := range text {
		if _, ok := cur.children[w]; !ok {
			return nil, false
		}
		res = append(res, start+i)
		cur = cur.children[w]
		if cur.isEnd {
			break
		}
	}
	if !cur.isEnd {
		return nil, false
	}
	return res, cur.isEnd
}

// QueryAll 查询文本text中全部的敏感词列表
// []string 文本text中存在的全部敏感词列表
// bool 文本text中存在敏感词时，返回true，否则返回false
func (t *Trie) QueryAll(text string) ([]string, bool) {
	words := []rune(text)
	node := t.root
	res := make([]*replaceIndexes, 0, len(text))
	for i, w := range words {
		if _, ok := node.children[w]; ok {
			if tmp, is := t.searchDirtyWords(words[i:], i); is {
				res = append(res, tmp)
			}
		}
	}
	if len(res) == 0 {
		return nil, false
	}
	result := make([]string, 0, len(res))
	for _, item := range res {
		result = append(result, string(words[item.start:item.end+1]))
	}
	return result, true
}

// Query 查询文本text中是否存在敏感词
func (t *Trie) Query(text string) bool {
	words := []rune(text)
	node := t.root
	for i, w := range words {
		if _, ok := node.children[w]; ok {
			if _, is := t.searchDirtyWords(words[i:], i); is {
				return true
			}
		}
	}
	return false
}

type replaceIndexes struct {
	start int
	end   int
}

func (t *Trie) searchDirtyWords(text []rune, start int) (*replaceIndexes, bool) {
	cur := t.root
	res := &replaceIndexes{
		start: start,
		end:   -1,
	}
	for i, w := range text {
		if _, ok := cur.children[w]; !ok {
			return nil, false
		}
		cur = cur.children[w]
		if cur.isEnd {
			res.end = start + i
			break
		}
	}
	return res, cur.isEnd
}

func buildPatternTable(pattern string) []int {
	table := make([]int, 0, len(pattern))
	table[0] = 0
	length := 0
	i := 1
	for i < len(pattern) {
		if pattern[i] == pattern[length] {
			length++
			table[i] = length
			i++
		} else {
			if length != 0 {
				length = table[length-1]
			} else {
				table[i] = 0
				i++
			}
		}
	}
	return table
}
