package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"text/template"
)

type SPaths struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type SStructs struct {
	Name   string `json:name`
	Field  string `json:field`
	Struct string `json:struct`
	Path   bool   `json:path`
}

type Data struct {
	Paths   []*SPaths   `json:paths`
	Package []string    `json:"packages"`
	Structs []*SStructs `json:structs`
}

func (data *Data) pathJson(path string) error {
	byts, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
		return err
	}
	err = json.Unmarshal(byts, data)
	if err != nil {
		log.Panic(err)
		return err
	}
	return nil
}

var wg sync.WaitGroup

const (
	lfn = `
func {{.FunName}}(index interface{}) (*{{.SName}}, bool) {
	cfg := Config()
	if cfg == nil{
		return nil, false
	}
	f, ok := cfg.GetValue("{{.Path}}", index, "{{.Param}}")
	if !ok {
		return nil, false
	}
	return f.(*{{.SName}}), true
}
`
	lyyn = `
func Y{{.FunName}}(index interface{}) (*{{.SName}}, bool) {
	cfg := YYAct()
	if cfg == nil{
		return nil, false
	}
	f, ok := cfg.GetValue("{{.Path}}", index, "{{.Param}}")
	if !ok {
		return nil, false
	}
	return f.(*{{.SName}}), true
}
`
)

func config() {
	data := &Data{
		Paths:   make([]*SPaths, 0),
		Package: make([]string, 0),
		Structs: make([]*SStructs, 0),
	}
	data.pathJson("build.json")
	//处理
	wfiles := make(map[string]*os.File)
	wBufs := make(map[string]*strings.Builder)
	t := template.Must(template.New("lfn").Parse(lfn))
	endchat := "\r\n"
	//先把文件构建出来
	pkgs := strings.Join(data.Package, endchat)
	pkgstr := "import (" + endchat + pkgs + endchat + ")" + endchat
	for _, paths := range data.Paths {
		fwrite, err := os.OpenFile(paths.Path+"txt.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatalln(err)
		}
		wfiles[paths.Name] = fwrite
		wBufs[paths.Name] = &strings.Builder{}
		//添加一些数据
		fwrite.WriteString("package configs" + endchat)
		fwrite.WriteString(pkgstr)
		fwrite.WriteString("func init() {" + endchat)
	}
	//关闭文件
	defer func() {
		for _, wf := range wfiles {
			wf.Close()
		}
	}()
	//开始生成文件
	for _, field := range data.Structs {
		pos := strings.Index(field.Name, `/`)
		if pos <= 0 {
			log.Panicln("file struct has no server")
		}
		server := field.Name[:pos]
		path := field.Name[pos+1:]
		params := make(map[string]string)
		params["Path"] = path
		fnn := path
		if len(field.Field) > 0 {
			fnn = fnn + "/" + field.Field
		}
		fnn = strings.Title(fnn)
		fnn = strings.ReplaceAll(fnn, "/", "")
		params["FunName"] = fnn
		params["SName"] = field.Struct
		params["Param"] = field.Field

		if server == "util" {
			for _, fw := range wfiles {
				fw.WriteString(`cfgs.RegisterSCreate("` + path + `", "` + field.Field + `", func() interface{} { return &` + field.Struct + `{} })` + endchat)
			}
			if !field.Path {
				for _, b := range wBufs {
					t.Execute(b, params)
				}
			}
		} else {
			fw, ok := wfiles[server]
			if !ok {
				log.Panicf("file struct has no server %v", server)
			}
			fw.WriteString(`cfgs.RegisterSCreate("` + path + `", "` + field.Field + `", func() interface{} { return &` + field.Struct + `{} })` + endchat)
			if !field.Path {
				b, ok := wBufs[server]
				if !ok {
					log.Panicf("file buff create fail %v", server)
				}
				t.Execute(b, params)
			}
		}
	}
	//结尾操作
	for sn, fw := range wfiles {
		fw.WriteString("}" + endchat)
		b, ok := wBufs[sn]
		if !ok {
			log.Panicf("file buff create fail %v", sn)
		}
		fw.WriteString(b.String())
	}

	wg.Done()
}

func yyact() {
	data := &Data{
		Paths:   make([]*SPaths, 0),
		Package: make([]string, 0),
		Structs: make([]*SStructs, 0),
	}
	data.pathJson("buildYYact.json")
	//处理
	wfiles := make(map[string]*os.File)
	wBufs := make(map[string]*strings.Builder)
	t := template.Must(template.New("lyyn").Parse(lyyn))
	endchat := "\r\n"
	//先把文件构建出来
	pkgs := strings.Join(data.Package, endchat)
	pkgstr := "import (" + endchat + pkgs + endchat + ")" + endchat
	for _, paths := range data.Paths {
		fwrite, err := os.OpenFile(paths.Path+"yyact.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatalln(err)
		}
		wfiles[paths.Name] = fwrite
		wBufs[paths.Name] = &strings.Builder{}
		//添加一些数据
		fwrite.WriteString("package configs" + endchat)
		fwrite.WriteString(pkgstr)
		fwrite.WriteString("func init() {" + endchat)
	}
	//关闭文件
	defer func() {
		for _, wf := range wfiles {
			wf.Close()
		}
	}()
	//开始生成文件
	for _, field := range data.Structs {
		params := make(map[string]string)
		params["Path"] = field.Name
		fnn := field.Name
		if len(field.Field) > 0 {
			fnn = fnn + "/" + field.Field
		}
		fnn = strings.Title(fnn)
		fnn = strings.ReplaceAll(fnn, "/", "")
		params["FunName"] = fnn
		params["SName"] = field.Struct
		params["Param"] = field.Field
		for _, fw := range wfiles {
			fw.WriteString(`cfgs.RegisterSCreate("` + field.Name + `", "` + field.Field + `", func() interface{} { return &` + field.Struct + `{} })` + endchat)
		}
		if !field.Path {
			for _, b := range wBufs {
				t.Execute(b, params)
			}
		}
	}
	//结尾操作
	for sn, fw := range wfiles {
		fw.WriteString("}" + endchat)
		b, ok := wBufs[sn]
		if !ok {
			log.Panicf("file buff create fail %v", sn)
		}
		fw.WriteString(b.String())
	}

	wg.Done()
}

func main() {
	wg.Add(2)
	go config()
	go yyact()
	wg.Wait()
}
