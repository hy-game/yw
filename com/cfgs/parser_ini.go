package cfgs

import (
	"pb"
	"regexp"
	"strconv"
	"strings"
)

//解析ini文件配置
func ParserIni(source *pb.MsgOriginalCfgs, out *SConfigStore) error {
	for _, v := range source.Inicfg {
		out.ini_map[v.Key] = v.Value
	}
	replaceIni(out)
	return nil
}

//ini通配符处理
func replaceIni(out *SConfigStore) {
	r, _ := regexp.Compile("%\\(.*\\)")
	for k, v := range out.ini_map {
		for {
			if r.MatchString(v) {
				v = out.regexpIni(r, v)
			} else {
				out.ini_map[k] = v
				break
			}
		}
	}
}

func (cfg *SConfigStore) regexpIni(r *regexp.Regexp, v string) string {
	ss := r.FindStringSubmatch(v)
	ii := r.FindStringSubmatchIndex(v)
	str := ""
	if ii[0] > 0 {
		str = string([]byte(v)[:ii[0]])
	}
	s := strings.TrimPrefix(ss[0], "%(")
	s = strings.TrimSuffix(s, ")")
	str += cfg.ini_map[s]
	if ii[1] < len(v) {
		str += string([]byte(v)[ii[1]:])
	}
	return str
}

//读取配置文件
func (cfg *SConfigStore) GetIniInt(space string, key string) (int, bool) {
	if space != "" {
		key = "[" + space + "]" + key
	}
	v, ok := cfg.ini_map[key]
	if !ok {
		return 0, false
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return i, true
}

//读取配置文件
func (cfg *SConfigStore) GetIniString(space string, key string) (string, bool) {
	if space != "" {
		key = "[" + space + "]" + key
	}
	v, ok := cfg.ini_map[key]
	return v, ok
}

//读取配置文件
func (cfg *SConfigStore) GetIniFloat(space string, key string) (float32, bool) {
	if space != "" {
		key = "[" + space + "]" + key
	}
	v, ok := cfg.ini_map[key]
	if !ok {
		return 0.0, false
	}
	f, err := strconv.ParseFloat(v, 32)
	if err != nil {
		return 0.0, false
	}
	return float32(f), true
}

//读取配置文件
func (cfg *SConfigStore) GetInitBool(space string, key string) (bool, bool) {
	if space != "" {
		key = "[" + space + "]" + key
	}
	v, ok := cfg.ini_map[key]
	if !ok {
		return false, false
	}
	f, err := strconv.ParseBool(v)
	if err != nil {
		return false, false
	}
	return f, true
}

//读取int类型的ini配置数据
func (cfg *SConfigStore) IniInt(key string) (int, bool) {
	return cfg.GetIniInt("", key)
}

//读取string类型的ini配置数据
func (cfg *SConfigStore) IniString(key string) (string, bool) {
	return cfg.GetIniString("", key)
}

//读取float32类型的ini配置数据
func (cfg *SConfigStore) IniFloat(key string) (float32, bool) {
	return cfg.GetIniFloat("", key)
}

//读取bool类型的ini配置数据
func (cfg *SConfigStore) InitBool(key string) (bool, bool) {
	return cfg.GetInitBool("", key)
}
