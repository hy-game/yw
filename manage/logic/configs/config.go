package configs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"pb"
	"strings"
	"sync"
)

var (
	errorNoServer   = errors.New("Config had not have server")
	errorTxtNoTitle = errors.New("Config had not title reflect")
)

const (
	PARSER_CG_EXT    uint8 = iota //后缀名
	PARSER_CG_PREFIX              //前缀字符串
	PARSER_CG_FILE                //文件名匹配
)

const (
	serverCfgPath = "./configs"
	topicYYAct    = "yyact"
	topicUtil     = "util"
)

//读取文件回调函数
type PaserFunc func(name, path string) error

//注册解释函数 和 筛选方式
type ParserKandler struct {
	Category uint8
	Parser   PaserFunc
	Param    string
	Ext      string //在以文件或文件夹注册时 可指定后缀名
}

//定义解析器类型
type parsers struct {
	ext_parser    map[string]PaserFunc
	prefix_parser map[string]PaserFunc
	file_parser   map[string]PaserFunc
	lock          sync.Mutex
}

//定义解析器
var handles = parsers{
	ext_parser:    make(map[string]PaserFunc),
	prefix_parser: make(map[string]PaserFunc),
	file_parser:   make(map[string]PaserFunc),
}

func (p *parsers) isRegister(pk *ParserKandler) bool {
	switch pk.Category {
	case PARSER_CG_EXT:
		{
			_, ok := handles.ext_parser[pk.Param]
			return ok
		}
	case PARSER_CG_PREFIX:
		{
			key := pk.Param
			if pk.Ext != "" {
				key = fmt.Sprintf("[%s]%s", pk.Ext, pk.Param)
			}
			_, ok := handles.prefix_parser[key]
			return ok
		}
	case PARSER_CG_FILE:
		{
			key := pk.Param
			if pk.Ext != "" {
				key = fmt.Sprintf("[%s]%s", pk.Ext, pk.Param)
			}
			_, ok := handles.file_parser[key]
			return ok
		}
	}
	return false
}

func (p *parsers) matchFunc(name, ext string) (PaserFunc, bool) {
	//优先匹配文件名
	key := fmt.Sprintf("[%s]%s", ext, name)
	if f, ok := p.file_parser[key]; ok {
		return f, true
	}
	if f, ok := p.file_parser[name]; ok {
		return f, true
	}
	//其次匹配前缀
	for k, v := range p.prefix_parser {
		if ok := strings.HasPrefix(key, k); ok {
			return v, true
		}
		if ok := strings.HasPrefix(name, k); ok {
			return v, true
		}
	}
	//最后匹配后缀名
	if r, ok := p.ext_parser[ext]; ok {
		return r, true
	}
	return nil, false
}

//注册配置的解析函数
func RegisterParser(pk *ParserKandler) bool {
	handles.lock.Lock()
	defer handles.lock.Unlock()

	if ok := handles.isRegister(pk); ok {
		return false
	}
	switch pk.Category {
	case PARSER_CG_EXT:
		{
			if handles.ext_parser == nil {
				handles.ext_parser = make(map[string]PaserFunc)
			}
			handles.ext_parser[pk.Param] = pk.Parser
		}
	case PARSER_CG_PREFIX:
		{
			if handles.prefix_parser == nil {
				handles.prefix_parser = make(map[string]PaserFunc)
			}
			key := pk.Param
			if pk.Ext != "" {
				key = fmt.Sprintf("[%s]%s", pk.Ext, pk.Param)
			}
			handles.prefix_parser[key] = pk.Parser
		}
	case PARSER_CG_FILE:
		{
			if handles.file_parser == nil {
				handles.file_parser = make(map[string]PaserFunc)
			}
			key := pk.Param
			if pk.Ext != "" {
				key = fmt.Sprintf("[%s]%s", pk.Ext, pk.Param)
			}
			handles.file_parser[key] = pk.Parser
		}
	}
	return true
}

//重读配置 起服时执行一次 运行中一般由GM后台调用
func Reload(yyact bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Run time panic: %v", err)
		}
	}()

	reloadServer()
	//广播
	if yyact {
		broadCastYYAct()
	} else {
		broadCastCfg()
	}
}

func reloadServer() {
	//生成新的实例
	loadingcfg = make(mapCfgs)

	var group sync.WaitGroup
	walk_files(serverCfgPath, "", &group)
	group.Wait()

	//替换
	origincfg.Store(loadingcfg)
}

func walk_files(path, ex string, group *sync.WaitGroup) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalln("Read Configs error")
	}
	for _, file := range files {
		npath := filepath.Join(path, file.Name())
		if file.IsDir() {
			nex := ex
			if nex == "" {
				nex = file.Name()
			} else {
				nex = ex + "/" + file.Name()
			}
			walk_files(npath, nex, group)
		} else {
			nex := ex
			ext := filepath.Ext(file.Name())
			fi := strings.TrimSuffix(file.Name(), ext)
			if nex == "" {
				nex = fi
			} else {
				nex = ex + "/" + fi
			}
			if fn, ok := handles.matchFunc(nex, ext); ok {
				group.Add(1)
				go func() {
					fn(nex, npath)
					group.Done()
					fmt.Printf("%v : %v \n", nex, npath)
				}()
			}
		}
	}
}

func subPathServer(path string) (string, string) {
	pos := strings.Index(path, `/`)
	if pos > 0 {
		return path[:pos], path[pos+1:]
	}
	return "", path
}

func catchCfg(name string) (*pb.MsgOriginalCfgs, string, error) {
	server, path := subPathServer(name)
	if server == "" {
		log.Fatalln(errorNoServer, name)
		return nil, name, errorNoServer
	}

	cfg, ok := loadingcfg[server]
	if !ok {
		//还没创建 就先创建
		cfg = &pb.MsgOriginalCfgs{}
		loadingcfg[server] = cfg
	}
	return cfg, path, nil
}
