package cfgs

import (
	clog "com/log"
	"com/util"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"pb"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

//解析txt配置
func ParserTxt(source *pb.MsgOriginalCfgs, out *SConfigStore) {
	for _, file := range source.TableCfgs {
		table := &ConfTable{
			name: file.Filename,
			//column: map[string]string{},
			rows: map[string]map[string]interface{}{},
		}
		if ok := table.read(file); ok {
			out.conf_tables[table.name] = table
		}
	}
}

func (table *ConfTable) read(file *pb.MsgTableCfg) bool {
	//解析标题列
	columns := make(map[int]cfgcolumn)
	pbname := table.parserColumn(file.Title, &columns)
	//读取内容
	for _, line := range file.Lines {
		table.parserRow(line, &columns, pbname)
	}
	return true
}

func (table *ConfTable) parserColumn(l string, column *map[int]cfgcolumn) (pb string) {
	//l = strings.TrimPrefix(l, "^")
	l = strings.TrimSpace(l)
	strs := strings.FieldsFunc(l, util.SplitRule)
	i := 0
	field := strings.TrimSpace(strs[0])
	head := true //这个标识标志是否用配置的field
	if field == "^" {
		head = false
	} else {
		pb = fmt.Sprintf("pb.%s", strings.TrimPrefix(field, "^"))
		//pb = strings.TrimPrefix(field, "^")
	}
	strs = strs[1:]
	for _, v := range strs {
		ss := strings.Split(v, ":")
		if len(ss) == 2 {
			nme := ss[0]
			if head && !strings.HasPrefix(nme, ".") {
				nme = "" //不用配置的field字段 因为field字段是
			}
			col := cfgcolumn{
				name:  nme,
				vtype: ss[1],
			}
			(*column)[i] = col
		} else {
			col := cfgcolumn{}
			(*column)[i] = col
		}
		i++
	}
	return
}

func (table *ConfTable) parserRow(l string, columns *map[int]cfgcolumn, pbname string) {
	l = strings.TrimPrefix(l, "#")
	l = strings.TrimSpace(l)
	strs := strings.FieldsFunc(l, util.SplitRule)

	key := strs[0]
	multi := ""
	if len(*columns) > len(strs) {
		if col, ok := (*columns)[len(strs)]; ok {
			if col.name == ".key" {
				key1 := ""
				ss := strings.Split(col.vtype, ",")
				for _, si := range ss {
					pos, err := strconv.Atoi(si)
					if err == nil {
						key1 += strs[pos] + ","
					} else {
						log.Fatalf("Parser Configs Key Error: %v", err)
					}
				}
				if key1 != "" {
					key = key1[:len(key1)-1]
				}
			} else if col.name == ".multi" {
				multi = col.vtype
			}
		}
	}

	var row map[string]interface{}
	ok := false
	if multi != "" {
		row, ok = table.rows[key]
	}
	if !ok {
		row = make(map[string]interface{})
	}

	for i, value := range strs {
		col, ok := (*columns)[i]
		if !ok {
			continue //没有该index字段
		}
		if col.vtype == "" {
			continue //该字段被忽略
		}
		var v interface{}
		switch col.vtype {
		case "int":
			j, _ := strconv.Atoi(value)
			v = j
		case "string":
			v = value
		case "float":
			f, _ := strconv.ParseFloat(value, 32)
			v = float32(f)
		case "bool":
			b, _ := strconv.ParseBool(value)
			v = b
		default:
			{
				v, ok = row[col.name]
				if !ok {
					fn, ok := getSCreate(table.name, col.name)
					if ok {
						v = fn()
					} else if pbname != "" {
						v = reflect.New(proto.MessageType(pbname).Elem()).Interface()
						if v == nil {
							continue
						}
					} else {
						continue //没有获取到构造函数
					}
				}

				if multi != "" {
					initMulti(v, multi)
					multi = "" //只增加一行
				}

				if fullData(v, col.vtype, value) == false {
					continue //填充数据失败
				}
			}
		}

		row[col.name] = v
	}

	table.rows[key] = row
}

//获取总行数
func (this *SConfigStore) GetLength(table_name string) (int, bool) {
	t, ok := this.conf_tables[table_name]
	if ok {
		return len(t.rows), true
	}
	return 0, false
}

//获取一行的配置数据
func (this *SConfigStore) Find(table_name string, index interface{}) (map[string]interface{}, bool) {
	t, ok := this.conf_tables[table_name]
	if ok {
		i := util.ToString(index)
		r, ok := t.rows[i]
		return r, ok
	}
	return nil, false
}

//获取特定行 特定列的配置数据
func (this *SConfigStore) GetValue(table_name string, index interface{}, column string) (interface{}, bool) {
	r, ok := this.Find(table_name, index)
	if ok {
		c, ok := r[column]
		return c, ok
	} else {
		return nil, false
	}
}

//获取txt配置的int值 该函数只能用于配置中映射为int的值
func (this *SConfigStore) GetValueInt(table_name string, index interface{}, column string) (int, bool) {
	v, ok := this.GetValue(table_name, index, column)
	if ok {
		return v.(int), true
	} else {
		return 0, false
	}
}

//获取txt配置的string值 该函数只能用于配置中映射为string的值
func (this *SConfigStore) GetValueString(table_name string, index interface{}, column string) (string, bool) {
	v, ok := this.GetValue(table_name, index, column)
	if ok {
		return v.(string), true
	} else {
		return "", false
	}
}

//获取txt配置的float32值 该函数只能用于配置中映射为float32的值
func (this *SConfigStore) GetValueFloat(table_name string, index interface{}, column string) (float32, bool) {
	v, ok := this.GetValue(table_name, index, column)
	if ok {
		return v.(float32), true
	} else {
		return 0.0, false
	}
}

//获取txt配置的bool值 该函数只能用于配置中映射为bool的值
func (this *SConfigStore) GetValueBool(table_name string, index interface{}, column string) (bool, bool) {
	v, ok := this.GetValue(table_name, index, column)
	if ok {
		return v.(bool), true
	} else {
		return false, false
	}
}

//该函数遍历指定配置文件每一行数据 传入处理函数可以对其中数据做特殊处理
func (this *SConfigStore) MapTable(table_name string, field string, fn func(interface{}) bool) bool {
	t, ok := this.conf_tables[table_name]
	if !ok {
		clog.Warnf("table %v not found", table_name)
		return false
	}
	for _, v := range t.rows {
		i, ok := v[field]
		if !ok {
			clog.Warnf("table %v field %v not found", table_name, field)
			return false
		}
		ok = fn(i)
		if !ok {
			clog.Warnf("table %v field %v handle error", table_name, field)
			return false
		}
	}
	return true
}

func (this *SConfigStore) Keys(table_name string) sort.StringSlice {
	t, ok := this.conf_tables[table_name]
	if !ok {
		clog.Warnf("table %v not found", table_name)
		return nil
	}
	srt := make(sort.StringSlice, 0, len(t.rows))

	for k, _ := range t.rows {
		srt = append(srt, k)
	}
	srt.Sort()
	return srt
}

func (this *SConfigStore) KeysInt(table_name string) sort.IntSlice {
	t, ok := this.conf_tables[table_name]
	if !ok {
		clog.Warnf("table %v not found", table_name)
		return nil
	}
	srt := make(sort.IntSlice, 0, len(t.rows))

	for k, _ := range t.rows {
		ik, err := strconv.Atoi(k)
		if err != nil {
			return nil
		}
		srt = append(srt, ik)
	}
	srt.Sort()
	return srt
}

func (this *SConfigStore) KeysFloat64(table_name string) sort.Float64Slice {
	t, ok := this.conf_tables[table_name]
	if !ok {
		clog.Warnf("table %v not found", table_name)
		return nil
	}
	srt := make(sort.Float64Slice, 0, len(t.rows))

	for k, _ := range t.rows {
		fk, err := strconv.ParseFloat(k, 64)
		if err != nil {
			return nil
		}
		srt = append(srt, fk)
	}
	srt.Sort()
	return srt
}
