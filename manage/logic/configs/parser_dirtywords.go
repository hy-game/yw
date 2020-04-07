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
		Category: PARSER_CG_PREFIX,
		Parser:   parser_dirtywords,
		Param:    "game/dirtywords",
		Ext:      ".txt",
	}
	RegisterParser(pk)
}

func parser_dirtywords(name, path string) error {
	log.Printf("parser_txt: %s", path)
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Open File %v error with : %v", path, err)
		return err
	}
	defer file.Close()

	cfg, nnane, err := catchCfg(name)
	if err != nil {
		log.Fatalf("Open File %v error with : %v", path, err)
		return err
	}
	if cfg.DirtyWords == nil {
		cfg.DirtyWords = make([]*pb.MsgStrKeyValueArray, 0)
	}
	buf := bufio.NewReader(file)
	data := &pb.MsgStrKeyValueArray{
		Key: nnane,
	}
	cfg.DirtyWords = append(cfg.DirtyWords, data)
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
		if data.Value == nil {
			data.Value = make([]string, 0)
		}
		data.Value = append(data.Value, l)
	}
	return nil
}
