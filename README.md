# Search_novel 小说搜索下载器
> 在几个固定的网站搜索小说，并下载成txt文件


- 修改set.ini文件的 小说名称，然后运行程序，自动在下面三个配置好的小说网站中搜对应的小说，找到就会生成小说名字的TXT文件。
- 在DOS下运行的时候，可以给定参数：

## set.ini

```
    search.exe -name 飘邈之旅 -list y 
    上述命令行的意思是， 搜索下载“飘邈之旅”的小说，-list y的意思是要单独生成章节文件，默认是不生成的。
```  


```
    [setall]
xiaoshuo1 = 重生完美时代
小说名称 = 恋上青春期
allsearchnum = 3

[search1]
小说网站地址 = http://www.biquge.tw
小说网站中文名 = 笔趣阁
搜索url = http://zhannei.baidu.com/cse/search?s=16829369641378287696&q=
目录页去前缀 = <div id=\"list\">
目录页去后缀 = <div id=\"footer\" name=\"footer\">
章节页去前缀 = <div id=\"content\">
章节页去后缀 = <div class=\"bottem2\">

[search2]
小说网站地址 = http://www.23us.la
小说网站中文名 = 顶点小说
搜索url = http://zhannei.baidu.com/cse/search?s=6791755209793365737&entry=1&q=
目录页去前缀 = <dl class=\"chapterlist\">
目录页去后缀 = <div class=\"clr\">
章节页去前缀 = <div id=\"content\">
章节页去后缀 = <div class=\"link\">

[search3]
小说网站地址 = http://www.xxbiquge.com
小说网站中文名 = 新笔趣阁
搜索url = http://zhannei.baidu.com/cse/search?s=8823758711381329060&ie=utf-8&q=
目录页去前缀 = 正文
目录页去后缀 = <div id=\"footer\" name=\"footer\">
章节页去前缀 = <div id=\"content\">
章节页去后缀 = 上一章
```
