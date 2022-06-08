// @Time    : 2022/6/8 2:56 下午
// @Author  : HuYuan
// @File    : main.go
// @Email   : huyuan@virtaitech.com

package main

import (
	"flag"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func recursiveDir(dir string) []string {
	var fileList []string
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("read dir %s error: %v \n", dir, err)
		return nil
	}
	for _, f := range fs {
		if f.IsDir() {
			fileList = append(fileList, recursiveDir(path.Join(dir, f.Name()))...)
			continue
		}
		fileList = append(fileList, path.Join(dir, f.Name()))
	}
	return fileList
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func isFile(path string) bool {
	if exists(path) {
		return !isDir(path)
	}
	return false
}

func checkIsYaml(file string) bool {
	if !isFile(file) {
		return false
	}
	if strings.HasSuffix(file, "yaml") || strings.HasSuffix(file, "yml") {
		return true
	}
	return false
}

func encode(body []byte) ([]byte, error) {
	data := make(map[string]interface{})
	if err := yaml.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return yaml.Marshal(&data)
}

func format(file string) error {
	var perm os.FileMode
	body, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	contentList := strings.Split(string(body), "---")
	var newContentList []string
	for _, content := range contentList {
		if content == "" {
			continue
		}
		data, err := encode([]byte(content))
		if err != nil {
			return err
		}
		newContentList = append(newContentList, string(data))
	}

	d := strings.Join(newContentList, "---\n")

	info, err := os.Stat(file)
	if err != nil {
		perm = os.FileMode(420)
	} else {
		perm = info.Mode()
	}

	return ioutil.WriteFile(file, []byte(d), perm)
}

func formatFile(file string) {
	if !checkIsYaml(file) {
		log.Printf("文件 %s 不存在或者不以 .yaml 和 .yml 结尾\n", file)
		return
	}

	if err := format(file); err != nil {
		log.Printf("格式化 YAML 文件 %s 错误: %v\n", file, err)
		return
	}
}

func formatDir(dir string) {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("read dir %s error: %v \n", dir, err)
		return
	}
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		file := path.Join(dir, f.Name())
		formatFile(file)
	}
}

func formatDirRecursion(dir string) {
	fileList := recursiveDir(dir)
	for _, file := range fileList {
		formatFile(file)
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	var dir, recursion, file string
	flag.StringVar(&dir, "d", "", "格式化指定目录下面的所有yaml文件(以.yaml或者.yml结尾的文件)")
	flag.StringVar(&recursion, "r", "", "格式化指定目录及其子目录下面的所有yaml文件(以.yaml或者.yml结尾的文件)")
	flag.StringVar(&file, "f", "", "格式化指定yaml文件")
	flag.Parse()

	if dir == "" && recursion == "" && file == "" {
		log.Fatalln("请指定格式化的目录或者文件")
	}

	if dir != "" {
		if !isDir(dir) {
			log.Printf("%s 不是目录\n", dir)
			return
		}
		formatDir(dir)
		return
	}

	if recursion != "" {
		if !isDir(recursion) {
			log.Printf("%s 不是目录\n", recursion)
			return
		}
		formatDirRecursion(recursion)
		return
	}

	if file != "" {
		formatFile(file)
		return
	}
}
