package go_ethernet_ip

import (
	"bytes"

	"github.com/teldio-operations/go-ethernet-ip/typedef"
)

type MessageRouterRequest struct {
	Service         typedef.Usint
	RequestPathSize typedef.Usint
	RequestPath     []byte
	RequestData     []byte
}

func (m *MessageRouterRequest) Encode() []byte {
	if m.RequestPathSize == 0 {
		m.RequestPathSize = typedef.Usint(len(m.RequestPath) / 2)
	}

	buffer := new(bytes.Buffer)
	WriteByte(buffer, m.Service)
	WriteByte(buffer, m.RequestPathSize)
	WriteByte(buffer, m.RequestPath)
	WriteByte(buffer, m.RequestData)

	return buffer.Bytes()
}

func (m *MessageRouterRequest) New(service typedef.Usint, path []byte, data []byte) {
	m.Service = service
	m.RequestPathSize = typedef.Usint(len(path) / 2)
	m.RequestPath = path
	m.RequestData = data
}

type MessageRouterResponse struct {
	ReplyService           typedef.Usint
	Reserved               typedef.Usint
	GeneralStatus          typedef.Usint
	SizeOfAdditionalStatus typedef.Usint
	AdditionalStatus       []byte
	ResponseData           []byte
}

type FragmentedReadResponse struct {
	ReplyService           typedef.Usint
	Reserved               typedef.Usint
	GeneralStatus          typedef.Usint
	SizeOfAdditionalStatus typedef.Usint
	AdditionalStatus       []byte
	TagTypeValue           uint16
	ResponseData           []byte
}

func (m *FragmentedReadResponse) Decode(data []byte, isString bool) {
	dataReader := bytes.NewReader(data)
	ReadByte(dataReader, &m.ReplyService)
	ReadByte(dataReader, &m.Reserved)
	ReadByte(dataReader, &m.GeneralStatus)
	ReadByte(dataReader, &m.SizeOfAdditionalStatus)
	m.AdditionalStatus = make([]byte, m.SizeOfAdditionalStatus*2)
	ReadByte(dataReader, &m.AdditionalStatus)
	ReadByte(dataReader, &m.TagTypeValue)
	// Strings have an extra 2 bytes for some reason
	if isString {
		var _t uint16
		ReadByte(dataReader, &_t)
	}
	m.ResponseData = make([]byte, dataReader.Len())
	ReadByte(dataReader, &m.ResponseData)
}

func (m *MessageRouterResponse) Decode(data []byte) {
	dataReader := bytes.NewReader(data)
	ReadByte(dataReader, &m.ReplyService)
	ReadByte(dataReader, &m.Reserved)
	ReadByte(dataReader, &m.GeneralStatus)
	ReadByte(dataReader, &m.SizeOfAdditionalStatus)
	m.AdditionalStatus = make([]byte, m.SizeOfAdditionalStatus*2)
	ReadByte(dataReader, &m.AdditionalStatus)
	m.ResponseData = make([]byte, dataReader.Len())
	ReadByte(dataReader, &m.ResponseData)
}
