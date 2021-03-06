package lc

import (
	"fmt"

	dmr "github.com/polkabana/go-dmr"
)

// Data Format
// ref: ETSI TS 102 361-2 7.2.18
const (
	Format7Bit uint8 = iota
	FormatISO8Bit
	FormatUTF8
	FormatUTF16BE
)

// DataFormatName is a map of data format to string.
var DataFormatName = map[uint8]string{
	Format7Bit:    "7 bit",
	FormatISO8Bit: "ISO 8 bit",
	FormatUTF8:    "unicode utf-8",
	FormatUTF16BE: "unicode utf-16be",
}

// TalkerAliasHeaderPDU Conforms to ETSI TS 102 361-2 7.1.1.4
type TalkerAliasHeaderPDU struct {
	DataFormat uint8
	Length     uint8
	Data       []byte
}

// TalkerAliasBlockPDU Conforms to ETSI TS 102 361-2 7.1.1.5
type TalkerAliasBlockPDU struct {
	Data []byte
}

func movebit(src []byte, srcByte int, srcBit int, dst []byte, dstByte int, dstBit int) {
	bit := (src[srcByte] >> uint8(srcBit)) & dmr.B00000001
	if bit >= 1 {
		dst[dstByte] |= 1 << uint8(dstBit)
	} else {
		dst[dstByte] &= 0 << uint8(dstBit)
	}
}

// ParseTalkerAliasHeaderPDU parses TalkerAliasHeader PDU from bytes
func ParseTalkerAliasHeaderPDU(data []byte) (*TalkerAliasHeaderPDU, error) {
	if len(data) != 7 {
		return nil, fmt.Errorf("dmr/lc/talkeralias: expected 7 bytes, got %d", len(data))
	}

	dataFormat := (data[0] & dmr.B11000000) >> 6

	var out []byte
	if dataFormat == Format7Bit {
		// it will reorganize the bits in the array and return []byte with 7bit chars
		// in each position
		out = make([]byte, 7)
		for i := 7; i < 56; i++ {
			movebit(data, i/8, (7 - (i % 8)), out, (i-7)/7, 6-(i%7))
		}
	} else {
		out = data[1:6]
	}

	return &TalkerAliasHeaderPDU{
		DataFormat: dataFormat,
		Length:     (data[0] & dmr.B00111110) >> 1,
		Data:       out,
	}, nil
}

// Bytes returns object as bytes
func (t *TalkerAliasHeaderPDU) Bytes() []byte {
	var (
		out   []byte
		bit49 byte
	)

	if t.DataFormat == Format7Bit {
		out = make([]byte, 6)
		for i := 7; i < 56; i++ {
			movebit(t.Data, (i-7)/7, 6-(i%7), out, (i/8)-1, (7 - (i % 8)))
		}
	} else {
		out = t.Data
	}

	return []byte{
		((t.DataFormat << 6) & dmr.B11000000) | ((t.Length << 1) & dmr.B00111110) | (bit49 & dmr.B00000001),
		out[0],
		out[1],
		out[2],
		out[3],
		out[4],
		out[5],
	}
}

// DataAsString Returns data part of PDU encoded as string
func (t *TalkerAliasHeaderPDU) DataAsString() string {
	return string(t.Data)
}

func (t *TalkerAliasHeaderPDU) String() string {
	return fmt.Sprintf("TalkerAliasHeader: [ format: %s, length: %d, data: \"%s\" ]",
		DataFormatName[t.DataFormat], t.Length, t.DataAsString())
}

// ParseTalkerAliasBlockPDU parse talker alias block pdu
func ParseTalkerAliasBlockPDU(data []byte) (*TalkerAliasBlockPDU, error) {
	if len(data) != 7 {
		return nil, fmt.Errorf("dmr/lc/talkeralias: expected 7 bytes, got %d", len(data))
	}

	return &TalkerAliasBlockPDU{
		Data: data[0:6],
	}, nil
}

// Bytes returns object as bytes
func (t *TalkerAliasBlockPDU) Bytes() []byte {
	return t.Data
}

// DataAsString Returns data part of PDU encoded as string
func (t *TalkerAliasBlockPDU) DataAsString() string {
	return string(t.Data)
}

func (t *TalkerAliasBlockPDU) String() string {
	return fmt.Sprintf("TalkerAliasBlock: [ data: \"%s\" ]", t.DataAsString())
}
