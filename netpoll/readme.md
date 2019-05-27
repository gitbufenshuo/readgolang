# 先看看 golang 的 internal 机制，因为 netpoll 代码里有

- 一句话，一个包 apackage ，一个internal包 bpackage，如果 apackage 要 import internal/bpackage，
那么 apackage 和 bpackage 必须要有同一根目录

- 所以 net 包 可以 import internal/poll, 因为他们都属于 src 目录

# 这个项目主要看看 netpoll 的运行机制，所以在 runtime/netpoll.go 开始就说明了一下，

- netpoll.go 实现了跨平台 poll，netpoll.go这份代码是跟平台无关的。
- 具体的平台，具体的写自己的代码 举个例子 runtime/netpoll_epoll.go 就是我们最熟悉的 linux epoll 模式
- 所以你要是mac，runtime/netpoll_kqueue.go , 你就需要看这里的代码
- 还是最好搞一个 linux 环境吧，，，
- GOOS=linux GOARCH=amd64 GOROOT=/Users/newone/code/go /Users/newone/code/go/bin/go build -o a.out main.go
- 把 a.out 拷贝到 linux 系统里，就能 ./a.out 运行了，我们可以在关键地方，panic 一下，比如说 runtime/netpoll_epoll.go , netpollinit 直接 panic

```
func netpollinit() {
	panic("netpollinit")

```
- 哈哈，这样我们就可以看见，netpollinit 的调用轨迹了。

- 首先 普通的 socket 系统调用 生成 了 fd 3
- 然后 打开 /proc/sys/net/core/somaxconn 文件， 生成了 fd 4
- 这个 fd init 的时候，顺带 netpollinit
- 这里面好几个 fd 的结构体，抽象，
- 首先是 net.netFD 结构 代表网络 fd
- 然后 net.netFD 里包含了 一个 internal/poll.FD ,这个 FD 是os包和net包都可以用的一种抽象，因为一切皆文件吗，这相当于更底层的一种FD
- 然后 internal/poll.FD 里包含了 一个 pollDesc 专门用来 io poll 的。
- 一层一层的 init， 到了 pollDesc 的init时候，就会 调用 netpollinit 这函数 ，这函数只会被调用一次， sync.once 做了这个事，在linux下，这函数就简单了，
就是 epollcreate
- pollDesc.init 还会做一件事，你可能已经想到了，就是把自己的 fd 加到 刚才的 epoll里去。
- 到这为止，算是搞清楚了，哦，golang 在 linux 上，也是用了 epoll，这没什么，重点在后面
- 而这一切，都是由 net.Listen 函数触发的。我们还没走到 l.Accept() 函数,而这个函数，就要接受外界的连接了。

- l.Accept() 函数 最终会调用系统调用 accept，如果 出现了 erragain ，说明什么，再来一次吧，说明还没有连接过来，这个网络 io 还没有东西就绪，
- 一步一步往下看，就发现了，最后，到了 netpollblock 这里，而且不断掉循环调用这个函数
- 这个函数就是什么呢，如果发现 fd 是 pdREADY的，就返回true，可以开始读了，那什么时候，谁来标记这个 fd 状态呢。
- 首先，如果发现没有 pdREADY ， 就讲 fd 变成 pdWAIT， 然后把自己挂起，对，就是那个大名鼎鼎的 gopark ，让自己这个goroutine 挂起，变成WAIT状态
- 注意，这个函数还没完呢，这就是goroutine，挂起之后，没完
- 然后，另一个系统goroutine，在那里调用 netpoll_epoll.go 里的 netpoll，这就是普通的linux那一套，最关键就是，扫到可以用的fd之后，应该就是通过这些
- fd，去找找，有哪些 goroutine 卡在这些fd上了，然后把他们弄醒，继续跑，读数据。