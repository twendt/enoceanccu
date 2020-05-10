package discovery

import (
	"context"
	"fmt"
	"net"
	"time"
)

// maxBufferSize specifies the size of the buffers that
// are used to temporarily hold data from the UDP packets
// that we receive.
const maxBufferSize = 1024
const ccuType = "EnoceanCCU"
const version = ">I3.51.6.20200420"
const writeTimeout = time.Duration(1000)

var discoverSeq = []byte{0x02, 0x8F, 0x91, 0xC0, 0x01, 'e', 'Q', '3', 0x2D, 0x2A, 0x00, 0x2A, 0x00, 0x49}

type Server struct {
	response []byte
	address  string
}

func NewServer(serial, addr string) Server {
	response := []byte{0x02, 0x8f, 0x91, 0xc0, 0x01}
	response = append(response, []byte(ccuType)...)
	response = append(response, 0x00)
	response = append(response, []byte(serial)...)
	response = append(response, 0x00)
	response = append(response, []byte(version)...)
	response = append(response, 0x00, 0x00)

	return Server{
		response: response,
		address:  addr,
	}
}

func (s Server) Listen(ctx context.Context) (err error) {

	pc, err := net.ListenPacket("udp4", s.address)
	if err != nil {
		return
	}

	defer pc.Close()

	doneChan := make(chan error, 1)
	buffer := make([]byte, maxBufferSize)

	go func() {
		for {
			n, addr, err := pc.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}

			fmt.Printf("packet-received: bytes=%d from=%s\n",
				n, addr.String())

			//deadline := time.Now().Add(writeTimeout)
			//err = pc.SetWriteDeadline(deadline)
			//if err != nil {
			//	doneChan <- err
			//	return
			//}

			n, err = pc.WriteTo(s.response, addr)
			if err != nil {
				fmt.Println(err)
				doneChan <- err
				return
			}

			fmt.Printf("packet-written: bytes=%d to=%s\n", n, addr.String())
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("cancelled")
		err = ctx.Err()
	case err = <-doneChan:
	}

	return
}
