package sender

import (
	"bytes"
	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"log"
	"net"
	"tencent_im_reminder"
	"time"
)

var Servers = map[uint32]string{
	1: "10.133.37.216:20000",
	2: "10.133.2.214:20000",
	3: "172.27.204.187:20000",
}

func SendPkg(server_type uint32, pkg *tencent_im_reminder.ReminderPackage) (*tencent_im_reminder.ReminderPackage, int64) {
	buffer, err := proto.Marshal(pkg)
	if err != nil {
		log.Panicln("marshal:", err)
	}

	conn, err := net.Dial("udp", Servers[server_type])
	if err != nil {
		log.Println("dial:", err)
		return nil, 0
	}
	defer conn.Close()

	// encode
	send := new(bytes.Buffer)
	binary.Write(send, binary.BigEndian, byte('('))
	binary.Write(send, binary.BigEndian, uint32(len(buffer)))
	binary.Write(send, binary.BigEndian, buffer)
	binary.Write(send, binary.BigEndian, byte(')'))

	before := time.Now()
	_, err = conn.Write(send.Bytes())
	if err != nil {
		log.Panicln("Write:", err)
	}

	// log.Println("reading...", send.Bytes())
	buffer = make([]byte, 65535)
	recv_len, err := conn.Read(buffer)
	if err != nil {
		log.Panicln("Read:", err)
	}

	after := time.Now()

	var diff int64 = after.UnixNano() - before.UnixNano()

	// decode
	// log.Println(buffer)
	buffer = buffer[5 : recv_len-1]
	// log.Println(buffer)

	resp := new(tencent_im_reminder.ReminderPackage)
	err = proto.Unmarshal(buffer, resp)
	if err != nil {
		log.Println("Unmarshal:", err, ", recv_len:", recv_len)
	}
	return resp, diff
}
