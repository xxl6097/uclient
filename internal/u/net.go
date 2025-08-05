package u

import (
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

	_ = conn.SetDeadline(time.Now().Add(2 * time.Second))
	_, err := conn.Write(data)
	if err != nil {
		return false
	}

	reply := make([]byte, 1500)
	_, err = conn.Read(reply)
	return err == nil
}
