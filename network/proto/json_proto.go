package proto

import (
"bufio"
"bytes"
"encoding/binary"
"fmt"
"io"
"sync"

"github.com/henrylee2cn/teleport/socket"
"github.com/henrylee2cn/teleport/utils"
	"github.com/tidwall/gjson"
)

// NewJsonProtoFunc is creation function of JSON socket protocol.
//  Packet data demo: `83{"seq":%d,"ptype":%d,"uri":%q,"meta":%q,"body_codec":%d,"body":"%s","xfer_pipe":%s}`
var NewJsonProtoFunc2 = func(rw io.ReadWriter) socket.Proto {
	var (
		readBufioSize             int
		readBufferSize, isDefault = socket.ReadBuffer()
	)
	if isDefault {
		readBufioSize = 1024 * 4
	} else if readBufferSize == 0 {
		readBufioSize = 1024 * 35
	} else {
		readBufioSize = readBufferSize / 2
	}
	return &jsonproto2{
		id:   'j',
		name: "json",
		r:    bufio.NewReaderSize(rw, readBufioSize),
		w:    rw,
	}
}

type jsonproto2 struct {
	id   byte
	name string
	r    *bufio.Reader
	w    io.Writer
	rMu  sync.Mutex
}

// Version returns the protocol's id and name.
func (j *jsonproto2) Version() (byte, string) {
	return j.id, j.name
}

//const format = `{"seq":%d,"ptype":%d,"uri":%q,"meta":%q,"body_codec":%d,"body":"%s","xfer_pipe":%s}`
const format = `{"ptype":3,"body":"%s"}`

// Pack writes the Packet into the connection.
// Note: Make sure to write only once or there will be package contamination!
func (j *jsonproto2) Pack(p *socket.Packet) error {
	// marshal body
	bodyBytes, err := p.MarshalBody()
	if err != nil {
		return err
	}
	//// do transfer pipe
	//bodyBytes, err = p.XferPipe().OnPack(bodyBytes)
	//if err != nil {
	//	return err
	//}

	// marshal whole
	var s = fmt.Sprintf(format,
		bytes.Replace(bodyBytes, []byte{'"'}, []byte{'\\', '"'}, -1),
	)
	// set size
	p.SetSize(uint32(len(s)))
	var all = make([]byte, p.Size()+4)
	binary.BigEndian.PutUint32(all, p.Size())
	copy(all[4:], s)
	_, err = j.w.Write(all)
	return err
}

// Unpack reads bytes from the connection to the Packet.
// Note: Concurrent unsafe!
func (j *jsonproto2) Unpack(p *socket.Packet) error {
	j.rMu.Lock()
	defer j.rMu.Unlock()
	var size uint32
	err := binary.Read(j.r, binary.BigEndian, &size)
	if err != nil {
		return err
	}
	if err = p.SetSize(size); err != nil {
		return err
	}
	if p.Size() == 0 {
		return nil
	}
	bb := utils.AcquireByteBuffer()
	defer utils.ReleaseByteBuffer(bb)
	bb.ChangeLen(int(p.Size()))
	_, err = io.ReadFull(j.r, bb.B)
	if err != nil {
		return err
	}
	s := string(bb.B)
	//println("unpack: ",s)

	p.SetPtype(byte(gjson.Get(s, "ptype").Int()))
	//body := gjson.Get(s, "body").String()
	//bodyBytes, err := p.XferPipe().OnUnpack([]byte(body))
	//if err != nil {
	//	return err
	//}

	//// unmarshal new body
	//err = p.UnmarshalBody([]byte(body))
	return err
}
