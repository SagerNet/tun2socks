package tunnel

import (
	"github.com/xjasonlyu/tun2socks/common/pool"
	M "github.com/xjasonlyu/tun2socks/constant"
	"github.com/xjasonlyu/tun2socks/core"
	"github.com/xjasonlyu/tun2socks/log"
	"github.com/xjasonlyu/tun2socks/proxy"
	"github.com/xjasonlyu/tun2socks/tunnel/statistic"
	"io"
	"net"
	"sync"
)

func newTCPTracker(conn net.Conn, metadata *M.Metadata) net.Conn {
	return statistic.NewTCPTracker(conn, metadata, statistic.DefaultManager)
}

func handleTCP(localConn core.TCPConn) {
	defer localConn.Close()

	id := localConn.ID()
	metadata := &M.Metadata{
		Net:     M.TCP,
		SrcIP:   net.IP(id.RemoteAddress),
		SrcPort: id.RemotePort,
		DstIP:   net.IP(id.LocalAddress),
		DstPort: id.LocalPort,
	}

	targetConn, err := proxy.Dial(metadata)
	if err != nil {
		log.Warnf("[TCP] dial %s error: %v", metadata.DestinationAddress(), err)
		return
	}

	if dialerAddr, ok := targetConn.LocalAddr().(*net.TCPAddr); ok {
		metadata.MidIP = dialerAddr.IP
		metadata.MidPort = uint16(dialerAddr.Port)
	} else { /* fallback */
		metadata.MidIP, metadata.MidPort = parseAddr(targetConn.LocalAddr().String())
	}

	targetConn = newTCPTracker(targetConn, metadata)
	defer targetConn.Close()

	log.Infof("[TCP] %s <-> %s", metadata.SourceAddress(), metadata.DestinationAddress())
	relay(localConn, targetConn) /* relay connections */
}

// relay copies between left and right bidirectionally.
func relay(left, right net.Conn) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		_ = copyBuffer(right, left) /* ignore error */
	}()

	go func() {
		defer wg.Done()
		_ = copyBuffer(left, right) /* ignore error */
	}()

	wg.Wait()
}

func copyBuffer(dst io.Writer, src io.Reader) error {
	buf := pool.Get(pool.RelayBufferSize)
	defer pool.Put(buf)

	_, err := io.CopyBuffer(dst, src, buf)
	return err
}
