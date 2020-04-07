package cfgs

import (
	"github.com/importcjj/sensitive"
	"pb"
)

//解析屏蔽字配置
func ParserDw(source *pb.MsgOriginalCfgs, out *SConfigStore) error {
	for _, v := range source.DirtyWords {
		filter := sensitive.New()
		filter.AddWord(v.Value...)
		out.dirty_words[v.Key] = filter
	}
	return nil
}

//将屏蔽字替换为指定字符
func (this *SConfigStore) Replace(name string, old string, ch rune) string {
	filter, ok := this.dirty_words[name]
	if !ok {
		return old
	}
	return filter.Replace(old, ch)
}

//删除掉string中的所有屏蔽字
func (this *SConfigStore) Filter(name string, old string) string {
	filter, ok := this.dirty_words[name]
	if !ok {
		return old
	}
	return filter.Filter(old)
}

//如果包含屏蔽字 返回第一个 和 false
func (this *SConfigStore) Validate(name string, old string) (dw string, ok bool) {
	filter, ok := this.dirty_words[name]
	if !ok {
		return "", true
	}
	ok, dw = filter.Validate(old)
	return
}
