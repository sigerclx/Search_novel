// z_html

package z_html

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func HttpGet(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("URL地址获取错误: " + url + "\r\n")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(url + "URL数据读取错误！")
		return ""
	}

	return string(body)
}

func GetHref(src string) []string { //获取url和A标签的名字

	re, _ := regexp.Compile("(h|H)(r|R)(e|E)(f|F)=\"(\\S*)\">([^<]+</a>)")
	return re.FindAllString(src, -1)
	//return re.FindSubmatch([]byte(src))
}

func Removefrist(src string, str string) string {
	//将src前面的都删除
	re, _ := regexp.Compile("[\\S\\s]*" + str)
	src = re.ReplaceAllString(src, "")
	return src
}

func RemoveLast(src string, str string) string {
	//将src后面的都删除
	re, _ := regexp.Compile(str + "[\\s\\S]*")
	src = re.ReplaceAllString(src, "")
	return src
}

func ChangeHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]*?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	return src
}

func RemoveStyle(src string) string {
	//去除STYLE
	re, _ := regexp.Compile("\\<style[\\S\\s]*?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	return src
}

func RemoveScript(src string) string {
	//去除SCRIPT
	re, _ := regexp.Compile("\\<script[\\S\\s]*?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	return src
}

func RemoveRem(src string) string {
	//去除注释行
	re, _ := regexp.Compile("\\<!--[\\S\\s]*?--\\>")
	src = re.ReplaceAllString(src, "")
	return src
}

func RemoveHtml(src string) string {
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ := regexp.Compile("\\<[\\S\\s]*?\\>")
	src = re.ReplaceAllString(src, "\n")
	return src
}

func RemoveReturn(src string) string {
	//去除连续的换行符，替换为空，单空格不管
	re, _ := regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\r\n\r\n")
	return src
}
