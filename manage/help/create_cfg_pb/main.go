package main

import (
	"bufio"
	"bytes"
	"com/util"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	serverCfgPath = "../../configs"
	pbPath        = "../../../../../Public/"
	end_chat      = "\r\n"
	end_begin     = "\r\n\t"
	pb_array      = "[]"
)

var had_pbs map[string]bool
var wfile *os.File

func ReloadPath() {
	walk_files(serverCfgPath, "")
}

func walk_files(path, ex string) {
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
			walk_files(npath, nex)
		} else {
			nex := ex
			ext := filepath.Ext(file.Name())
			fi := strings.TrimSuffix(file.Name(), ext)
			if nex == "" {
				nex = fi
			} else {
				nex = ex + "/" + fi
			}

			parser(nex, npath, ext)
		}
	}
}

func parser(name, path, ext string) {
	if ext != ".txt" {
		return //只处理txt文件
	}
	log.Printf("parser_txt: %s", path)
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Open File %v error with : %v", path, err)
		return
	}
	defer file.Close()

	buf := bufio.NewReader(file)
	for {
		l, err := buf.ReadString('\n')
		if err == io.EOF {
			if len(l) == 0 {
				return //没有标题配置
			}
		} else if err != nil {
			log.Fatalln(err)
			return //读取出错
		}
		if !strings.HasPrefix(l, "^") {
			continue
		}
		handleTitle(name, l)
		break
	}
}

func handleTitle(pth, line string) {
	strs := strings.FieldsFunc(line, util.SplitRule)
	if len(strs) <= 0 {
		return //不是正确的标题
	}
	if strs[0] == "^" {
		return //无需解析的标题
	}
	pbname := strings.TrimPrefix(strs[0], "^")
	_, find := had_pbs[pbname]
	if find {
		return //已经生成的pb
	}

	name := bytes.NewBuffer([]byte{})
	name.WriteString("//")
	name.WriteString(pth)
	name.WriteString(end_chat)
	name.WriteString("message ")
	name.WriteString(pbname)
	name.WriteString("{")

	had_pbs[pbname] = true

	size := len(strs)
	j := 1
	for i := 1; i < size; i++ {
		ss := strings.Split(strings.TrimSpace(strs[i]), ":")
		if len(ss) == 2 {
			if ss[0] == "" {
				continue //不需要解析的字段
			} else if strings.HasPrefix(ss[0], ".") {
				continue //特殊字段
			}
			name.WriteString(end_begin)
			if strings.HasPrefix(ss[0], pb_array) {
				//数组
				name.WriteString("repeated ")
				name.WriteString(strings.TrimPrefix(ss[0], pb_array))
			} else {
				name.WriteString(ss[0])
			}

			idx := strings.Index(ss[1], ".")
			fn := ss[1]
			if idx > 0 {
				fn = fn[:idx]
			}
			name.WriteString(" ")
			name.WriteString(fn)

			name.WriteString(" = ")
			name.WriteString(strconv.Itoa(j))
			name.WriteString(";")
			j++
		}
	}
	name.WriteString(end_chat)
	name.WriteString("}")
	name.WriteString(end_chat)
	name.WriteString(end_chat)

	wfile.WriteString(name.String())
}

func endChat() {
	wfile.WriteString(end_chat)
}

func writeHead() {
	wfile.WriteString(`syntax = "proto3";`)
	endChat()
	endChat()
	wfile.WriteString(`option optimize_for = LITE_RUNTIME;`)
	endChat()
	endChat()
	wfile.WriteString(`import "msg_public.proto";`)
	endChat()
	wfile.WriteString(`import "msg.proto"; `)
	endChat()
	wfile.WriteString(`import "msg_config.proto"; `)
	endChat()
	endChat()
	wfile.WriteString(`//---------------说明---------------------------`)
	endChat()
	wfile.WriteString(`//自动生成的pb 配置文件`)
	endChat()
	wfile.WriteString(`//----------------------------------------------`)
	endChat()
	endChat()
	wfile.WriteString(`package pb;`)
	endChat()
	endChat()
}

func main() {
	had_pbs = make(map[string]bool)
	var err error
	wfile, err = os.OpenFile(pbPath+"msg_cfg_auto.proto", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("Open write File %v error with : %v", pbPath, err)
	}
	defer wfile.Close()

	writeHead()
	ReloadPath()
}
