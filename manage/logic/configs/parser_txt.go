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
		Parser:   parser_txt,
		Param:    ".txt",
	}
	RegisterParser(pk)
}

func parser_txt(name, path string) error {
	log.Printf("parser_txt: %s", path)
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Open File %v error with : %v", path, err)
		return err
	}
	defer file.Close()

	cfg, nname, err := catchCfg(name)
	if err != nil {
		log.Fatalf("Open File %v error with : %v", path, err)
		return err
	}
	if cfg.TableCfgs == nil {
		cfg.TableCfgs = make([]*pb.MsgTableCfg, 0)
	}
	data := &pb.MsgTableCfg{
		Filename: nname,
	}
	cfg.TableCfgs = append(cfg.TableCfgs, data)
	buf := bufio.NewReader(file)
	//读取标题
	for {
		l, err := buf.ReadString('\n')
		if err == io.EOF {
			if len(l) == 0 {
				log.Fatalln(errorTxtNoTitle, name)
				return errorTxtNoTitle
			}
		} else if err != nil {
			log.Fatalln(err)
			return err
		}
		if !strings.HasPrefix(l, "^") {
			continue
		}
		data.Title = l
		break
	}
	//读取内容
	for {
		l, err := buf.ReadString('\n')
		if err == io.EOF {
			if len(l) == 0 {
				break
			}
		} else if err != nil {
			log.Printf("read file error: %v", err)
			return err
		}
		if !strings.HasPrefix(l, "#") {
			continue
		}
		if data.Lines == nil {
			data.Lines = make([]string, 0)
		}
		data.Lines = append(data.Lines, l)
	}
	return nil
}
