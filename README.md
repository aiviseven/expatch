# expatch
java项目的增量补丁打包工具，学习go语言编写的小工具

go get github.com/bsh100220/expatch

```
expatch --help
 -conf string
        idea模块配置文件路径，当type=idea时候有效  
 -out string
        输出目录，默认为当前目录下的【patch+当前时间】目录
 -svn string
        需要对比的两个版本，用英文冒号(:)分隔，例：100:95
 -type string
        项目类型，可选(idea,eclipse)，默认为idea (default "idea")
```

目前只支持打包svn版本控制的项目，git的还未补充

build后的可执行程序需要放到环境变量中去，该工具不带有java编译功能，打包前最好编译一下项目，不然可能导致打包出来的补丁包中包含的文件不是最新的。
