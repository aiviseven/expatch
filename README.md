# expatch
java项目的增量补丁打包工具，学习go语言编写的小工具

go build -ldflags="-w -s"

```
expatch --help
 -conf string
        idea模块配置文件路径，当type=idea时候有效
 -ignore string
        忽略文件的配置文件路径，需要忽略的文件路径按行填写 (default ".expatch_ignore")
 -out string
        输出目录，默认为当前目录下的【patch+当前时间】目录
 -svn string
        需要对比的两个版本，版本号用英文冒号(:)分隔，带有文件路径则用空格分割，例：'100:95'、'100 .'，为空时则对比当前目录与svn最新版本
 -type string
        项目类型，可选(idea,eclipse)，默认为idea (default "idea")
```

目前只支持打包svn版本控制的项目，git的还未补充

build后的可执行程序需要放到环境变量中去，该工具不带有java编译功能，打包前最好编译一下项目，不然可能导致打包出来的补丁包中包含的文件不是最新的。
