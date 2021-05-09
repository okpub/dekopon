package packet

type Message interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

//包头(必须大于1个字节)
const PacketHeadSize = 5

//class packet
type Packet struct {
	Header [PacketHeadSize]byte
	Body   []byte
}

func New() *Packet {
	return &Packet{}
}

func With(body []byte) *Packet {
	return &Packet{Body: body}
}

func (p *Packet) Len() int {
	return PacketHeadSize + len(p.Body)
}

func (p *Packet) Marshal() (body []byte, err error) {
	body = make([]byte, p.Len())
	copy(body[:PacketHeadSize], p.Header[:PacketHeadSize])
	copy(body[PacketHeadSize:], p.Body)
	return
}

func (p *Packet) Unmarshal(body []byte) (err error) {
	p.Body = make([]byte, len(body)-PacketHeadSize)
	copy(p.Header[:PacketHeadSize], body[:PacketHeadSize])
	copy(p.Body, body[PacketHeadSize:])
	return
}
