package configs

import (
	"bufio"
	"io"
	"log"
	"os"
	"pb"
	"strings"
)

func init() {
	pk := &ParserKandler{
		Category: PARSER_CG_EXT,
		Parser:   parser_ini,
		Param:    ".ini",
	}
	RegisterParser(pk)
}

func parser_ini(name, path string) error {
	log.Printf("parser_txt: %s", path)
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Open File %v error with : %v", path, err)
		return err
	}
	defer file.Close()

	cfg, _, err := catchCfg(name)
	if err != nil {
		log.Fatalf("Open File %v error with : %v", path, err)
		return err
	}
	if cfg.Inicfg == nil {
		cfg.Inicfg = make([]*pb.MsgStrKeyValue, 0)
	}

	buf := bufio.NewReader(file)
	title := ""
	for {
		l, err := buf.ReadString('\n')
		if err == io.EOF {
			if len(l) == 0 {
				break
			}
		} else if err != nil {
			log.Fatalf("Read File %v error with error %v", path, err)
			return err
		}
		l = strings.TrimSpace(l)
		if len(l) == 0 { //空行
			continue
		} else if strings.HasPrefix(l, "#") || strings.HasPrefix(l, ";") { //注释
			continue
		} else if strings.HasPrefix(l, "[") { //分组
			title = l
			j := strings.Index(title, ";") //后面的为注释
			if j > 0 {
				title = strings.TrimSpace(string([]byte(l)[:j]))
			}
			title = strings.TrimSpace(title)
			continue
		}

		//配置行 去除注释
		i := strings.Index(l, ";")
		if i > 0 {
			l = strings.TrimSpace(l[:i])
		}
		//解析key value
		i = strings.Index(l, "=")
		if i <= 0 {
			continue
		}
		key := strings.TrimSpace(l[:i])
		value := strings.TrimSpace(l[i+1:])
		value = strings.TrimPrefix(value, `"`)
		value = strings.TrimSuffix(value, `"`)
		if title != "" {
			key = title + key
		}
		//添加进slice
		data := &pb.MsgStrKeyValue{
			Key:   key,
			Value: value,
		}
		cfg.Inicfg = append(cfg.Inicfg, data)
	}
	return nil
}
