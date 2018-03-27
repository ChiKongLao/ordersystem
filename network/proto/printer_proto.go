package proto

import (
	"bufio"
	"encoding/binary"
	"io"
	"sync"
	"github.com/henrylee2cn/teleport/socket"
	"github.com/henrylee2cn/teleport/utils"
)

var NewStringProtoFunc = func(rw io.ReadWriter) socket.Proto {
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
	return &stringproto{
		id:   'j',
		name: "json",
		r:    bufio.NewReaderSize(rw, readBufioSize),
		w:    rw,
	}
}

type stringproto struct {
	id   byte
	name string
	r    *bufio.Reader
	w    io.Writer
	rMu  sync.Mutex
}

// Version returns the protocol's id and name.
func (j *stringproto) Version() (byte, string) {
	return j.id, j.name
}

const stringprotoFormat = `{"seq":%d,"ptype":%d,"uri":%q,"meta":%q,"body_codec":%d,"body":"%s","xfer_pipe":%s}`

// Pack writes the Packet into the connection.
// Note: Make sure to write only once or there will be package contamination!
func (j *stringproto) Pack(p *socket.Packet) error {
	// marshal body
	bodyBytes, err := p.MarshalBody()
	if err != nil {
		return err
	}
	// do transfer pipe
	bodyBytes, err = p.XferPipe().OnPack(bodyBytes)
	//if err != nil {
	//	return err
	//}
	//
	//// marshal transfer pipe ids
	//var xferPipeIds = make([]int, p.XferPipe().Len())
	//for i, id := range p.XferPipe().Ids() {
	//	xferPipeIds[i] = int(id)
	//}
	//xferPipeIdsBytes, _ := json.Marshal(xferPipeIds)
	//// marshal whole
	//var s = fmt.Sprintf(format,
	//	p.Seq(),
	//	p.Ptype(),
	//	p.Uri(),
	//	p.Meta().QueryString(),
	//	p.BodyCodec(),
	//	bytes.Replace(bodyBytes, []byte{'"'}, []byte{'\\', '"'}, -1),
	//	xferPipeIdsBytes,
	//)
	//// set size
	//p.SetSize(uint32(len(s)))
	//var all = make([]byte, p.Size()+4)
	//binary.BigEndian.PutUint32(all, p.Size())
	//copy(all[4:], s)
	//_, err = j.w.Write(all)


	_, err = j.w.Write(bodyBytes)

	return err
}

// Unpack reads bytes from the connection to the Packet.
// Note: Concurrent unsafe!
func (j *stringproto) Unpack(p *socket.Packet) error {
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
	//if err != nil {
	//	return err
	//}
	//s := string(bb.B)
	//
	//// read transfer pipe
	//xferPipe := gjson.Get(s, "xfer_pipe")
	//for _, r := range xferPipe.Array() {
	//	p.XferPipe().Append(byte(r.Int()))
	//}
	//
	//// read body
	//p.SetBodyCodec(byte(gjson.Get(s, "body_codec").Int()))
	//body := gjson.Get(s, "body").String()
	//bodyBytes, err := p.XferPipe().OnUnpack([]byte(body))
	//if err != nil {
	//	return err
	//}
	//
	//// read other
	//p.SetSeq(uint64(gjson.Get(s, "seq").Int()))
	//p.SetPtype(byte(gjson.Get(s, "ptype").Int()))
	//p.SetUri(gjson.Get(s, "uri").String())
	//meta := gjson.Get(s, "meta").String()
	//p.Meta().ParseBytes(goutil.StringToBytes(meta))
	//
	//// unmarshal new body
	//err = p.UnmarshalBody(bodyBytes)








	return err
}
