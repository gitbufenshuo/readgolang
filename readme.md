# 阅读和对golang源码进行实验，打印日志，以了解一些细节

# 预备工作,准备可以随时污染(加日志)的另一份go源码

- 首先要弄一份golang源码，去github clone就好了,为什么要另搞一份源码，是为了不污染你本身的环境
- 假如golang 源码的 路径是 ~/go 那么 
- cd ~/go/src
- sh all.bash
- 以上两部执行完了之后，在你的 ~/go/bin 就会生成一个新的 go 文件，就是go程序了
- 用这份go源码编译你的文件可以这样
- GOROOT=~/go ~/go/bin/go build yourcode.go
- 那么你在 ~/go/src 这里改了东西的话，会反应出来了，比如说加了日志。

# 本项目每一个目录，都主要对一个点进行实验。

# 比如说 begin 项目就是一个简单的，我们在 ~/go/src/runtime/proc.go 的main 函数里，第一句，加了个打印

```
// The main goroutine.
func main() {
	printstring("main go routine\n")
	g := getg()
```

然后，执行 
- GOROOT=~/go ~/go/bin/go build -o a.out main.go
- ./a.out
结果就是

```
main go routine
```

# 看源码，要尽情的使用 panic，以便看出调用路径。。。