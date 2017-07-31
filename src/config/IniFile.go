package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"strings"
)

func ReadSection(filename string, section string) *ini.Section {
	//LoadOptions 可以读取
	//MySQL's configuration allows a key without value as follows:
	//skip-name-resolve
	cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, filename)
	if err != nil {
		fmt.Println("load ini config file error:", err)
	}

	cfg.BlockMode = false
	section_ret := cfg.Section(section)

	return section_ret
}

func ReadSectionKey(filename string, section string, key string) string {

	section_val := ReadSection(filename, section)
	value := section_val.Key(key).String()

	return strings.TrimSpace(value)
	//cfg.Section("mysqld").NewKey("datadir","/data01/data/")
	//cfg.Section("mysqld").NewKey("user","mysql")
	//
	//cfg.SaveTo(filename)
}

func WriteSection(filename string, Newsection string) bool {
	cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, filename)
	if err != nil {
		fmt.Println("load ini config file error:", err)
		return false
	}

	sectionNames := cfg.SectionStrings()
	//判断新的Section是否存在
	for _, s := range sectionNames {
		if s == Newsection {
			fmt.Println("new section already exist!")
			return true
		}
	}
	_, err = cfg.NewSection(Newsection)
	if err != nil {
		fmt.Println("add new section error:", err)
		return false
	}
	cfg.SaveTo(filename)
	return true
}

func WriteKey(filename string, sectionName string, key string, value string) bool {
	cfg, err := ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, filename)
	if err != nil {
		fmt.Println("load ini config file error:", err)
		return false
	}
	//首先检查写入的section是否存在
	ret := WriteSection(filename, sectionName)
	if ret {
		//如果存在section,那么检查写入的key是否存在
		if cfg.Section(sectionName).HasKey(key) {
			fmt.Println("your key is already exits!")
			return false
		}

		//在对应的section下写入key value
		_, err = cfg.Section(sectionName).NewKey(key, value)
		cfg.SaveTo(filename)
		if err != nil {
			fmt.Println("add key :", key, "error:", err)
			return false
		}
	}

	return true
}
