module github.com/xjasonlyu/tun2socks

go 1.17

require (
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/sys v0.0.0-20210909193231-528a39cd75f3
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac
	gvisor.dev/gvisor v0.0.0
)

require github.com/google/btree v1.0.1 // indirect

replace gvisor.dev/gvisor v0.0.0 => github.com/sagernet/gvisor v0.0.0-20210909160323-ce37d6df1e92
