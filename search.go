package main

import (
	"flag"
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	//"z_goini"
	"./z_file" //读写文件
	"./z_goini"
	"./z_html" //去除HTML标签
)

var wg sync.WaitGroup //定义一个同步等待的组
type MyMap struct {
	xiaoshuoData map[int]string
	sync.RWMutex
}

type XiaoShuoSet struct { //小说搜索配置信息
	WebName             string
	WebHome             string
	SearchUrl           string
	ZhangNum            int    //总体章节数量
	FristRemove         string //目录页去除的前缀文字
	LastRemove          string //目录页去除的后缀文字
	zhangjieFristRemove string //章节页去除的前缀文字
	zhangjieLastRemove  string //章节页去除的前缀文字

	AllSearch int //预置了多少个可以用于搜小说的网站
}

var struct_xiaoshuo MyMap
var Current_XiaoShuoSet XiaoShuoSet

var novel = flag.String("name", "", "小说名称")        //""为默认值
var list = flag.String("list", "n", "章节文件开关(n/y)") //""为默认值
var XiaoShuoName = ""

func main() {
	flag.Parse()
	fmt.Println(*novel)
	struct_xiaoshuo.xiaoshuoData = make(map[int]string)
	maxProcs := runtime.NumCPU() //获set取cpu个数
	runtime.GOMAXPROCS(maxProcs) //限制同时运行的goroutines数量

	conf := goini.SetConfig("set.ini")
	//取set.ini中登记小说网站的数量

	Current_XiaoShuoSet.AllSearch, _ = strconv.Atoi(conf.GetValue("setall", "allsearchnum"))
	XiaoShuoName = *novel //如果参数有小说名字，就按参数中的小说名字搜，否则按INI中事先设定的搜
	if len(XiaoShuoName) == 0 {
		XiaoShuoName = conf.GetValue("setall", "小说名称")
	}
	searchResult1 := get_all_search(XiaoShuoName)

	for k, URL := range searchResult1 {
		get_ini(k)
		if len(URL) > 0 {
			fmt.Println("1. 找到小说链接：", URL)
			getXiaoshuo(k, URL)
			wg.Wait() //阻塞等待所有组内成员都执行完毕退栈

			fmt.Println("All DONE!!!")
			subject := XiaoShuoName + " From " + Current_XiaoShuoSet.WebHome + "\r\n\r\n"
			z_file.Writefile(subject, XiaoShuoName+"_"+Current_XiaoShuoSet.WebName+".txt")
			for a := 1; a <= Current_XiaoShuoSet.ZhangNum; a++ {
				z_file.AppendToFile(struct_xiaoshuo.xiaoshuoData[a], XiaoShuoName+"_"+Current_XiaoShuoSet.WebName+".txt")
			}
		} else {
			fmt.Println(Current_XiaoShuoSet.WebName + " " + Current_XiaoShuoSet.WebHome + " 没有搜到：" + XiaoShuoName)
		}

	}

}

func get_all_search(xiaoshuoname string) map[int]string {

	var searchResult = make(map[int]string) //返回值，存储搜索到的小说链接url

	conf := goini.SetConfig("set.ini")
	URL := ""

	//取set.ini中登记小说网站的数量
	Current_XiaoShuoSet.AllSearch, _ = strconv.Atoi(conf.GetValue("setall", "allsearchnum"))

	//在所有INI中预先定义的小说网站中搜索
	//根据ini的配置顺序，获取搜索小说事先配置项
	for i := 1; i <= Current_XiaoShuoSet.AllSearch; i++ {
		get_ini(i)
		URL = search_a_web()

		searchResult[i] = URL
		//fmt.Println(searchResult[i])
		if len(URL) == 0 {
			//fmt.Println(Current_XiaoShuoSet.WebHome, " 没有搜到对应小说！")
		}
	}
	return searchResult
}

func getXiaoshuo(k int, URL string) {
	get_ini(k)
	//传入URL是小说的具体地址，然后获取小说所有章节文字
	body := z_html.HttpGet(URL)
	body = z_html.Removefrist(body, Current_XiaoShuoSet.FristRemove)
	body = z_html.RemoveLast(body, Current_XiaoShuoSet.LastRemove)
	//fmt.Println(body)
	//z_file.Writefile(body, "zhangjie.txt")
	body = mulu_dowith(body, Current_XiaoShuoSet.WebHome)

	//把小说的目录章节写入文件
	fmt.Println("3. 生成章节目录......")
	if *list == "y" { //本程序的第二个命令行参数，y的话就把章节单独写文件
		z_file.Writefile(body, XiaoShuoName+"_"+Current_XiaoShuoSet.WebName+"_LIST.txt")
	}
}

func get_ini(num int) {
	//根据ini的配置顺序，获取搜索小说事先配置项
	conf := goini.SetConfig("set.ini")
	//取set.ini中登记小说网站的数量

	search_set := "search" + strconv.Itoa(num)
	Current_XiaoShuoSet.WebHome = conf.GetValue(search_set, "小说网站地址")
	Current_XiaoShuoSet.WebName = conf.GetValue(search_set, "小说网站中文名")
	Current_XiaoShuoSet.SearchUrl = conf.GetValue(search_set, "搜索url")
	Current_XiaoShuoSet.FristRemove = conf.GetValue(search_set, "目录页去前缀")
	Current_XiaoShuoSet.LastRemove = conf.GetValue(search_set, "目录页去后缀")
	Current_XiaoShuoSet.zhangjieFristRemove = conf.GetValue(search_set, "章节页去前缀")
	Current_XiaoShuoSet.zhangjieLastRemove = conf.GetValue(search_set, "章节页去后缀")

}

func search_a_web() string {
	// 搜索一个小说网站，成功返回小说的具体URL，失败返回空
	body := z_html.HttpGet(Current_XiaoShuoSet.SearchUrl + XiaoShuoName)

	zhengze := "(h|H)(r|R)(e|E)(f|F)=\"(\\S*)\"\\s+title=\"" + XiaoShuoName + "\""
	re, _ := regexp.Compile(zhengze)
	xiashuourl := re.FindSubmatch([]byte(body))

	if len(xiashuourl) == 0 {
		//搜索不到对应的小说
		return ""
	}
	x := string(xiashuourl[5]) //第5个匹配()就是url
	return x

}

func Set(key int, value string) {
	//MAP写入数据的锁
	struct_xiaoshuo.Lock()
	struct_xiaoshuo.xiaoshuoData[key] = value
	struct_xiaoshuo.Unlock()
}

func Get(key int) string {
	//MAP获取数据的锁
	struct_xiaoshuo.RLock()
	defer struct_xiaoshuo.RUnlock()
	return struct_xiaoshuo.xiaoshuoData[key]
}

func getandwrite(thisURL string, myLink string, zhang1 int) {
	//多线程模块：根据具体章节url获取该章节内容
	time1 := time.Now() //记录开始时间
	mystr := z_html.HttpGet(thisURL)
	fmt.Println(myLink)
	//z_file.Writefile(mystr, myLink+"1.txt")
	mystr = z_html.Removefrist(mystr, Current_XiaoShuoSet.zhangjieFristRemove)
	//z_file.Writefile(mystr, myLink+"2.txt")
	mystr = z_html.RemoveLast(mystr, Current_XiaoShuoSet.zhangjieLastRemove)
	//z_file.Writefile(mystr, myLink+".txt")
	mystr = z_html.RemoveStyle(mystr)
	mystr = z_html.RemoveScript(mystr)
	mystr = z_html.RemoveHtml(mystr)
	mystr = z_html.RemoveReturn(mystr)
	mystr = strings.Replace(mystr, "&nbsp;", " ", -1)
	fmt.Println("正在下载 ：" + myLink)
	mystr = myLink + "\r\n\r\n" + mystr
	//z_file.Writefile(mystr, strconv.Itoa(zhang1)+".txt") //把每章写入文件
	Set(zhang1, mystr) //Map 变量的写入，锁住为了防止同时写入MAP发生崩溃
	fmt.Println(time.Since(time1))
	defer wg.Done()
}

func mulu_dowith(src string, webHome string) string {
	//找出href连接的字符串
	fmt.Println("2. 分析章节目录......")
	var mystr string //返回值
	var zhangjiename string
	var thisURL string
	Current_XiaoShuoSet.ZhangNum = 0
	mystr1 := z_html.GetHref(src) //获取A标签里href的值

	re, _ := regexp.Compile("(h|H)(r|R)(e|E)(f|F)=\"(\\S*)\">([^<]+)</a>")
	var xiashuourl [][]byte
	mystr = XiaoShuoName + " From : " + webHome + "\r\n"
	for _, x := range mystr1 {
		Current_XiaoShuoSet.ZhangNum++
		xiashuourl = re.FindSubmatch([]byte(x))
		if len(xiashuourl) == 0 {
			fmt.Println("章节节点匹配出错!!!")
			return x
		}
		thisURL = string(xiashuourl[5]) //第5匹配()就是url
		zhangjiename = string(xiashuourl[6])
		//fmt.Println(x, zhangjiename)

		thisURL = webHome + thisURL
		mystr = mystr + strconv.Itoa(Current_XiaoShuoSet.ZhangNum) + " : " + zhangjiename + " , " + thisURL + "\r\n"

		//fmt.Println(webHome + "  ： " + thisURL)
		wg.Add(1) //为同步等待组增加一个成员 ,加入的成员视cpu个数，进行排队
		go getandwrite(thisURL, zhangjiename, Current_XiaoShuoSet.ZhangNum)

	}
	return mystr
}
