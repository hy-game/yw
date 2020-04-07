package configs

//将屏蔽字替换为指定字符
func Replace(name string, old string, ch rune) string {
	cfg := Config()
	if cfg != nil {
		return cfg.Replace(name, old, ch)
	}
	return old
}

//删除掉string中的所有屏蔽字
func Filter(name string, old string) string {
	cfg := Config()
	if cfg != nil {
		return cfg.Filter(name, old)
	}
	return old
}

//如果包含屏蔽字 返回第一个 和 false
func Validate(name string, old string) (string, bool) {
	cfg := Config()
	if cfg != nil {
		return cfg.Validate(name, old)
	}
	return old, false
}

//默认配置
func ReplaceDef(old string) string {
	return Replace("dirtywords/forName", old, '*')
}

//默认配置
func FilterDef(old string) string {
	return Filter("dirtywords/forName", old)
}

//默认配置
func ValidateDef(old string) (string, bool) {
	return Validate("dirtywords/forName", old)
}
