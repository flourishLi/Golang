package consts

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Config struct {
	filepath string
	conflist []map[string]map[string]string
}

//Create an empty configuration file
func SetConfig(filepath string) *Config {
	c := new(Config)
	c.filepath = filepath

	return c
}

//根据key获取配置信息
func (c *Config) GetValue(section, name string) string {
	c.ReadList()
	conf := c.ReadList()
	for _, v := range conf {
		for key, value := range v {
			if key == section {
				return value[name]
			}
		}
	}
	return ""
}

//修改配置文件
func (c *Config) SetValue(section, key, value string) bool {
	c.ReadList()
	data := c.conflist
	var ok bool
	var index = make(map[int]bool)
	var conf = make(map[string]map[string]string)
	for i, v := range data {
		_, ok = v[section]
		index[i] = ok
	}

	i, ok := func(m map[int]bool) (i int, v bool) {
		for i, v := range m {
			if v == true {
				return i, true
			}
		}
		return 0, false
	}(index)

	if ok {
		c.conflist[i][section][key] = value
		return true
	} else {
		conf[section] = make(map[string]string)
		conf[section][key] = value
		c.conflist = append(c.conflist, conf)
		return true
	}

	return false
}

//根据Key删除配置信息
func (c *Config) DeleteValue(section, name string) bool {
	c.ReadList()
	data := c.conflist
	for i, v := range data {
		for key, _ := range v {
			if key == section {
				delete(c.conflist[i][key], name)
				return true
			}
		}
	}
	return false
}

//List all the configuration file
func (c *Config) ReadList() []map[string]map[string]string {

	file, err := os.Open(c.filepath)
	if err != nil {
		CheckErr(err)
	}
	defer file.Close()
	var data map[string]map[string]string
	var section string
	buf := bufio.NewReader(file)
	for {
		l, err := buf.ReadString('\n')
		line := strings.TrimSpace(l)
		if err != nil {
			if err != io.EOF {
				CheckErr(err)
			}
			if len(line) == 0 {
				break
			}
		}
		switch {
		case len(line) == 0:
		case line[0] == '[' && line[len(line)-1] == ']':
			section = strings.TrimSpace(line[1 : len(line)-1])
			data = make(map[string]map[string]string)
			data[section] = make(map[string]string)
		default:
			i := strings.IndexAny(line, "=")
			value := strings.TrimSpace(line[i+1 : len(line)])
			data[section][strings.TrimSpace(line[0:i])] = value
			if c.uniquappend(section) == true {
				c.conflist = append(c.conflist, data)
			}
		}

	}

	return c.conflist
}

func CheckErr(err error) string {
	if err != nil {
		return fmt.Sprintf("Error is :'%s'", err.Error())
	}
	return "Notfound this error"
}

//Ban repeated appended to the slice method
func (c *Config) uniquappend(conf string) bool {
	for _, v := range c.conflist {
		for k, _ := range v {
			if k == conf {
				return false
			}
		}
	}
	return true
}
