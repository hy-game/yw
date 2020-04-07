package cfgs

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type SCreateFunc func() interface{}

type SCreateReg struct {
	Path  string
	Field string
	Fn    SCreateFunc
}

var structFuncHandle = make(map[string]*SCreateReg)
var structLock sync.Mutex

//注册配置文件和对应的struct
func RegisterSCreate(space string, fnname string, fn SCreateFunc) {
	structLock.Lock()
	defer structLock.Unlock()
	value := &SCreateReg{
		Path:  space,
		Field: fnname,
		Fn:    fn,
	}
	if space != "" {
		fnname = "[" + space + "]" + fnname
	}
	structFuncHandle[fnname] = value
}

func getSCreate(space string, fnname string) (SCreateFunc, bool) {
	//先进行文件匹配
	key := fnname
	if space != "" {
		key = "[" + space + "]" + fnname
	}
	f, ok := structFuncHandle[key]
	if ok {
		return f.Fn, true //文件匹配成功
	}

	//文件匹配失败 进行路径匹配
	for _, v := range structFuncHandle {
		//这里路径传入空时 就忽略路径的匹配
		if space != "" {
			if !strings.HasPrefix(v.Path, space) {
				continue //路径匹配不上
			}
		}
		if fnname == v.Field {
			return v.Fn, true
		}
	}
	return nil, false
}

func initMulti(data interface{}, field string) {
	t := reflect.ValueOf(data)
	tv := t.Elem().FieldByName(field)
	if tv.IsNil() {
		slic := reflect.MakeSlice(tv.Type(), 0, 0)
		tv.Set(slic)
	}
	//log.Debugf("%v", tv.Type().Elem().Elem())
	vv := reflect.New(tv.Type().Elem().Elem())
	tv.Set(reflect.Append(tv, vv))
}

func fullData(data interface{}, line string, value string) bool {
	if line == "" {
		return false
	}
	t := reflect.ValueOf(data)
	pars := strings.Split(line, ".")
	if len(pars) == 1 {
		field := t.Elem().FieldByName(pars[0])
		fullField(&field, value)
	} else {
		tv := t
		for _, str := range pars {
			if tv.Type().Kind() == reflect.Slice {
				i, err := strconv.Atoi(str)
				if err != nil {
					//多行的情况
					i = tv.Len() - 1
					tv = tv.Index(i)
					tv = tv.Elem().FieldByName(str)
				} else {
					for tv.Len() <= i {
						vvv := reflect.New(tv.Type().Elem().Elem())
						tv.Set(reflect.Append(tv, vvv))
					}
					tv = tv.Index(i)
				}
			} else {
				tv = tv.Elem().FieldByName(str)
				//这里初始化
				if tv.Type().Kind() == reflect.Ptr {
					if tv.Type().Elem().Kind() == reflect.Struct {
						if tv.IsNil() {
							vvv := reflect.New(tv.Type().Elem())
							tv.Set(vvv)
						}
					} else if tv.Type().Elem().Kind() == reflect.Slice {
						if tv.IsNil() {
							slic := reflect.MakeSlice(tv.Type(), 0, 0)
							tv.Set(slic)
						}
					}
				}
			}
		}
		fullField(&tv, value)
	}
	return true
}

func fullField(field *reflect.Value, value string) bool {
	switch field.Type().Kind() {
	case reflect.String,
		reflect.Int32,
		reflect.Int,
		reflect.Int64,
		reflect.Uint32,
		reflect.Uint,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Bool:
		{
			v, ok := parseValue(field.Type().Kind(), value)
			if !ok {
				return false
			}
			field.Set(reflect.ValueOf(v).Convert(field.Type()))
		}
	case reflect.Slice:
		{
			typ := field.Type().Elem().Kind()
			v, ok := parseValue(typ, value)
			if !ok {
				return false
			}
			field.Set(reflect.Append(*field, reflect.ValueOf(v).Convert(field.Type().Elem())))
		}
	}
	return true
}

func parseValue(typ reflect.Kind, value string) (interface{}, bool) {
	switch typ {
	case reflect.String:
		{
			return value, true
		}
	case reflect.Int32:
		{
			f, err := strconv.Atoi(value)
			if err != nil {
				return nil, false
			}
			return int32(f), true
		}
	case reflect.Int,
		reflect.Int64:
		{
			f, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, false
			}
			return f, true
		}
	case reflect.Uint32:
		{
			f, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return nil, false
			}
			return f, true
		}
	case reflect.Uint,
		reflect.Uint64:
		{
			f, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return nil, false
			}
			return f, true
		}
	case reflect.Float32:
		{
			f, err := strconv.ParseFloat(value, 10)
			if err != nil {
				return nil, false
			}
			return float32(f), true
		}
	case reflect.Float64:
		{
			f, err := strconv.ParseFloat(value, 10)
			if err != nil {
				return nil, false
			}
			return f, true
		}
	case reflect.Bool:
		{
			f, err := strconv.ParseBool(value)
			if err != nil {
				return nil, false
			}
			return f, true
		}
	}
	return nil, false
}
