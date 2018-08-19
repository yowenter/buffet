package lib

import "net/url"

// Item describe a info card
type Item struct {
	Metadata    map[string]interface{}
	Tags        []string
	Subject     string
	Description string
	Link        url.URL
	References  []string
}

//  let me see see .. 好代码的基础有一半是好的数据结构。
// so 如何设计这个数据结构呢？
// 作为一个资深图书馆员， 信息存储的目的是为了信息使用，so， 这个信息一定要有好的可供索引的字段。 Tags （或者 Keywrods 是必须的）
// 描述一个东西， 基本的 subject/title 是需要有的。
// 为了使信息可追溯， 所以要表明来源 -）
// 当然， 因为所有的信息都是息息相关的，所以要加一个 参考文档
// Metadata 这东西，为了扩展，先扔个 hash 表。 万一以后要加什么字段也方便。
