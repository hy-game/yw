package configs

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"pb"
)

func init() {
	//此处注册需要特殊处理的配置文件
	registerBytes(PARSER_CG_PREFIX, "fight/RegionLogic", ".xml")
}

func registerBytes(category uint8, param string, ext string) {
	pk := &ParserKandler{
		Category: category,
		Parser:   parser_bytes,
		Param:    param,
		Ext:      ext,
	}
	RegisterParser(pk)
}

func parser_bytes(name, path string) error {
	log.Printf("parser_bytes: %s", path)
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
	if cfg.OtherCfgs == nil {
		cfg.OtherCfgs = make([]*pb.MsgBytesCfg, 0)
	}
	data := &pb.MsgBytesCfg{
		Filename: nname,
		Fileext:  filepath.Ext(path),
	}
	data.Filebody, err = ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Read other file %v error %v", path, err)
	}
	cfg.OtherCfgs = append(cfg.OtherCfgs, data)
	return nil
}
