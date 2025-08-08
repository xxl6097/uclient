package u

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"time"
)

func Ping(target string) bool {
	conn, _ := net.Dial("ip4:icmp", target)
	defer conn.Close()

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Body: &icmp.Echo{ID: 0, Seq: 1, Data: []byte("")},
	}
	data, _ := msg.Marshal(nil)

	_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
	_, err := conn.Write(data)
	if err != nil {
		return false
	}

	reply := make([]byte, 1500)
	_, err = conn.Read(reply)
	return err == nil
}

func ByteCountSpeed(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B/s", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %ciB/s", float64(b)/float64(div), "KMGTPE"[exp])
}
