package go_ethernet_ip

import (
	"bytes"
	"errors"
	"time"

	"github.com/teldio-operations/go-ethernet-ip/typedef"
)

type ListIdentityItem struct {
	ItemTypeCode                 typedef.Uint
	ItemLength                   typedef.Uint
	EncapsulationProtocolVersion typedef.Uint
	SinFamily                    typedef.Int
	SinPort                      typedef.Uint
	SinAddr                      typedef.Udint
	SinZero                      typedef.Ulint
	VendorID                     typedef.Uint
	DeviceType                   typedef.Uint
	ProductCode                  typedef.Uint
	Major                        typedef.Usint
	Minor                        typedef.Usint
	Status                       typedef.Word
	SerialNumber                 typedef.Udint
	NameLength                   typedef.Usint
	ProductName                  []byte
	State                        typedef.Usint
}

type ListIdentity struct {
	ItemCount typedef.Uint
	Items     []ListIdentityItem
}

func (l *ListIdentity) Decode(data []byte) {
	dataReader := bytes.NewReader(data)
	ReadByte(dataReader, &l.ItemCount)

	for i := typedef.Uint(0); i < l.ItemCount; i++ {
		item := &ListIdentityItem{}
		ReadByte(dataReader, &item.ItemTypeCode)
		ReadByte(dataReader, &item.ItemLength)
		ReadByte(dataReader, &item.EncapsulationProtocolVersion)
		ReadByteBigEndian(dataReader, &item.SinFamily)
		ReadByteBigEndian(dataReader, &item.SinPort)
		ReadByteBigEndian(dataReader, &item.SinAddr)
		ReadByteBigEndian(dataReader, &item.SinZero)
		ReadByte(dataReader, &item.VendorID)
		ReadByte(dataReader, &item.DeviceType)
		ReadByte(dataReader, &item.ProductCode)
		ReadByte(dataReader, &item.Major)
		ReadByte(dataReader, &item.Minor)
		ReadByte(dataReader, &item.Status)
		ReadByte(dataReader, &item.SerialNumber)
		ReadByte(dataReader, &item.NameLength)
		item.ProductName = make([]byte, item.NameLength)
		ReadByte(dataReader, &item.ProductName)
		ReadByte(dataReader, &item.State)

		l.Items = append(l.Items, *item)
	}
}

func NewListIdentity(context typedef.Ulint) *EncapsulationPacket {
	encapsulationPacket := &EncapsulationPacket{}
	encapsulationPacket.Command = EIPCommandListIdentity
	encapsulationPacket.SenderContext = context
	return encapsulationPacket
}

func (e *EIPTCP) ListIdentity() (*ListIdentity, error) {
	ctx := CtxGenerator()
	e.receiverMutex.Lock()
	e.receiver[ctx] = make(chan *EncapsulationPacket)
	e.receiverMutex.Unlock()

	encapsulationPacket := NewListIdentity(ctx)
	b, _ := encapsulationPacket.Encode()

	if e.tcpConn != nil {
		e.sender <- b
	}

	for {
		select {
		case <-time.After(e.config.TCPTimeout):
			return nil, errors.New("tcp timeout")
		case received := <-e.receiver[ctx]:
			return e.ListIdentityDecode(received), nil
		}
	}
}

func (e *EIPUDP) ListIdentity() {
	encapsulationPacket := NewListIdentity(0)
	b, _ := encapsulationPacket.Encode()

	if e.udpConn != nil {
		e.send(b)
	}
}

func (e *EIPTCP) ListIdentityDecode(encapsulationPacket *EncapsulationPacket) *ListIdentity {
	if len(encapsulationPacket.CommandSpecificData) == 0 {
		return nil
	}

	list := &ListIdentity{}
	list.Decode(encapsulationPacket.CommandSpecificData)

	return list
}

func (e *EIPUDP) ListIdentityDecode(encapsulationPacket *EncapsulationPacket) *ListIdentity {
	if len(encapsulationPacket.CommandSpecificData) == 0 {
		return nil
	}

	list := &ListIdentity{}
	list.Decode(encapsulationPacket.CommandSpecificData)

	return list
}
