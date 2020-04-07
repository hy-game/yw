package cfgs

import (
	"github.com/importcjj/sensitive"
	"pb"
)

type cfgcolumn struct {
	name  string
	vtype string
}

//这里存储了txt配置文件的数据
type ConfTable struct {
	name string //配置文件名称
	//column map[string]string                 //配置文件中列字段 和 列字段的数据类型
	rows map[string]map[string]interface{} //配置文件数据
}

type SConfigStore struct {
	conf_tables map[string]*ConfTable
	ini_map     map[string]string
	dirty_words map[string]*sensitive.Filter
	XmlMap      map[string]interface{}
}

//这里特殊处理配置加载完成后需要特殊处理的配置项
type AfterFunc func(*SConfigStore)

var afterFunc []AfterFunc = make([]AfterFunc, 0)

//注册对配置数据的特殊处理函数
func RegisterAfterFunc(fn AfterFunc) {
	if fn != nil {
		afterFunc = append(afterFunc, fn)
	}
}

//清除所有对配置的特殊处理 该函数应该没用
func ClearAfterFunc() {
	afterFunc = make([]AfterFunc, 0)
}

//获取一个容器存储配置数据
func NewCfgStore() *SConfigStore {
	return &SConfigStore{
		conf_tables: make(map[string]*ConfTable),
		ini_map:     make(map[string]string),
		dirty_words: make(map[string]*sensitive.Filter),
		XmlMap:      make(map[string]interface{}),
	}
}

//收到manage配置数据后 解析配置数据
func (this *SConfigStore) Reload(in *pb.MsgOriginalCfgs) {
	ParserIni(in, this)
	ParserTxt(in, this)
	ParserDw(in, this)
	ParserOther(in, this)

	this.afterLoad()
}

func (this *SConfigStore) afterLoad() {
	for _, fn := range afterFunc {
		if fn != nil {
			fn(this)
		}
	}
}
