package cfgs

import (
	"fmt"
	"pb"
	"strings"
)

//这里虽然传出了*SConfigStore指针 但要谨慎使用他 因为它的数据可能是不完整的
type FileHandler func(*pb.MsgBytesCfg, *SConfigStore)

type fileHandleReg struct {
	flag string
	fn   FileHandler
}

var fileHandleFns = make([]*fileHandleReg, 0)

//flag 为指定路径的文件夹或文件不带后缀名 例:a/txt
func RegisterHandle(flag string, fn FileHandler) {
	fileHandleFns = append(fileHandleFns, &fileHandleReg{
		flag: flag,
		fn:   fn,
	})
}

func RegisterHandleExt(flag string, ext string, fn FileHandler) {
	key := flag
	if ext != "" {
		key = fmt.Sprintf("[%s]%s", ext, flag)
	}
	RegisterHandle(key, fn)
}

func ClearHandle() {
	fileHandleFns = make([]*fileHandleReg, 0)
}

//解析其它配置(不能通用解析的) 需要事先注册解析函数
func ParserOther(source *pb.MsgOriginalCfgs, out *SConfigStore) {
	if len(fileHandleFns) <= 0 {
		return
	}
	if len(source.OtherCfgs) <= 0 {
		return
	}
	for _, file := range source.OtherCfgs {
		key := file.Filename
		if file.Fileext != "" {
			key = fmt.Sprintf("[%s]%s", file.Fileext, file.Filename)
		}
		for _, fns := range fileHandleFns {
			if strings.HasPrefix(key, fns.flag) {
				fns.fn(file, out)
				break
			}
			if strings.HasPrefix(file.Filename, fns.flag) {
				fns.fn(file, out)
				break
			}
		}
	}
}
