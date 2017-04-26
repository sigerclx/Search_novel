// z_file.go
package z_file

import (
	"fmt"
	"os"
)

func Writefile(mystr string, userFile string) {
	//userFile := "write.txt"          //文件路径
	fout, err := os.Create(userFile) //根据路径创建File的内存地址
	defer fout.Close()               //延迟关闭资源
	if err != nil {
		fmt.Println(userFile, err)
		return
	}

	//fout.Write(myfile1);  二进制写入文件
	_, err = fout.WriteString(mystr) //写入字符串
	if err != nil {
		fmt.Println(err)
		return
	}
}

// fileName:文件名字(带全路径)
// content: 写入的内容
func AppendToFile(content string, fileName string) error {
	// 以只写的模式，打开文件
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0644)

	if err != nil {
		//fmt.Println("cacheFileList.yml file create failed. err: " + err.Error())
		f, err = os.Create(fileName) //根据路径创建File的内存地址
	}
	// 查找文件末尾的偏移量
	n, _ := f.Seek(0, os.SEEK_END)
	// 从末尾的偏移量开始写入内容
	_, err = f.WriteAt([]byte(content), n)

	defer f.Close()
	return err
}

func Readfile(userFile string) string {
	var mystr string
	//userFile := "read.txt"        //文件路径
	fin, err := os.Open(userFile) //打开文件,返回File的内存地址
	defer fin.Close()             //延迟关闭资源
	if err != nil {
		fmt.Println(userFile, err)
		//		return ""
	}
	buf := make([]byte, 1024) //创建一个初始容量为1024的slice,作为缓冲容器
	for {
		//循环读取文件数据到缓冲容器中,返回读取到的个数
		n, _ := fin.Read(buf)

		if 0 == n {
			break //如果读到个数为0,则读取完毕,跳出循环
		}
		//从缓冲slice中写出数据,从slice下标0到n,通过os.Stdout写出到控制台
		//os.Stdout.Write(buf[:n])
		mystr += string(buf[:n])
	}
	return mystr

}
