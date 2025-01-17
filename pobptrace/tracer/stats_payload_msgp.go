

package tracer

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *groupedStats) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Service":
			z.Service, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Name":
			z.Name, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Resource":
			z.Resource, err = dc.ReadString()
			if err != nil {
				return
			}
		case "HTTPStatusCode":
			z.HTTPStatusCode, err = dc.ReadUint32()
			if err != nil {
				return
			}
		case "Type":
			z.Type, err = dc.ReadString()
			if err != nil {
				return
			}
		case "DBType":
			z.DBType, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Hits":
			z.Hits, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Errors":
			z.Errors, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Duration":
			z.Duration, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "OkSummary":
			z.OkSummary, err = dc.ReadBytes(z.OkSummary)
			if err != nil {
				return
			}
		case "ErrorSummary":
			z.ErrorSummary, err = dc.ReadBytes(z.ErrorSummary)
			if err != nil {
				return
			}
		case "Synthetics":
			z.Synthetics, err = dc.ReadBool()
			if err != nil {
				return
			}
		case "TopLevelHits":
			z.TopLevelHits, err = dc.ReadUint64()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *groupedStats) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 13
	// write "Service"
	err = en.Append(0x8d, 0xa7, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Service)
	if err != nil {
		return
	}
	// write "Name"
	err = en.Append(0xa4, 0x4e, 0x61, 0x6d, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Name)
	if err != nil {
		return
	}
	// write "Resource"
	err = en.Append(0xa8, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Resource)
	if err != nil {
		return
	}
	// write "HTTPStatusCode"
	err = en.Append(0xae, 0x48, 0x54, 0x54, 0x50, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.HTTPStatusCode)
	if err != nil {
		return
	}
	// write "Type"
	err = en.Append(0xa4, 0x54, 0x79, 0x70, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Type)
	if err != nil {
		return
	}
	// write "DBType"
	err = en.Append(0xa6, 0x44, 0x42, 0x54, 0x79, 0x70, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.DBType)
	if err != nil {
		return
	}
	// write "Hits"
	err = en.Append(0xa4, 0x48, 0x69, 0x74, 0x73)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Hits)
	if err != nil {
		return
	}
	// write "Errors"
	err = en.Append(0xa6, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x73)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Errors)
	if err != nil {
		return
	}
	// write "Duration"
	err = en.Append(0xa8, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Duration)
	if err != nil {
		return
	}
	// write "OkSummary"
	err = en.Append(0xa9, 0x4f, 0x6b, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.OkSummary)
	if err != nil {
		return
	}
	// write "ErrorSummary"
	err = en.Append(0xac, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79)
	if err != nil {
		return
	}
	err = en.WriteBytes(z.ErrorSummary)
	if err != nil {
		return
	}
	// write "Synthetics"
	err = en.Append(0xaa, 0x53, 0x79, 0x6e, 0x74, 0x68, 0x65, 0x74, 0x69, 0x63, 0x73)
	if err != nil {
		return
	}
	err = en.WriteBool(z.Synthetics)
	if err != nil {
		return
	}
	// write "TopLevelHits"
	err = en.Append(0xac, 0x54, 0x6f, 0x70, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x48, 0x69, 0x74, 0x73)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.TopLevelHits)
	if err != nil {
		return
	}
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *groupedStats) Msgsize() (s int) {
	s = 1 + 8 + msgp.StringPrefixSize + len(z.Service) + 5 + msgp.StringPrefixSize + len(z.Name) + 9 + msgp.StringPrefixSize + len(z.Resource) + 15 + msgp.Uint32Size + 5 + msgp.StringPrefixSize + len(z.Type) + 7 + msgp.StringPrefixSize + len(z.DBType) + 5 + msgp.Uint64Size + 7 + msgp.Uint64Size + 9 + msgp.Uint64Size + 10 + msgp.BytesPrefixSize + len(z.OkSummary) + 13 + msgp.BytesPrefixSize + len(z.ErrorSummary) + 11 + msgp.BoolSize + 13 + msgp.Uint64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *statsBucket) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Start":
			z.Start, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Duration":
			z.Duration, err = dc.ReadUint64()
			if err != nil {
				return
			}
		case "Stats":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Stats) >= int(zb0002) {
				z.Stats = (z.Stats)[:zb0002]
			} else {
				z.Stats = make([]groupedStats, zb0002)
			}
			for za0001 := range z.Stats {
				err = z.Stats[za0001].DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *statsBucket) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "Start"
	err = en.Append(0x83, 0xa5, 0x53, 0x74, 0x61, 0x72, 0x74)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Start)
	if err != nil {
		return
	}
	// write "Duration"
	err = en.Append(0xa8, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteUint64(z.Duration)
	if err != nil {
		return
	}
	// write "Stats"
	err = en.Append(0xa5, 0x53, 0x74, 0x61, 0x74, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Stats)))
	if err != nil {
		return
	}
	for za0001 := range z.Stats {
		err = z.Stats[za0001].EncodeMsg(en)
		if err != nil {
			return
		}
	}
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *statsBucket) Msgsize() (s int) {
	s = 1 + 6 + msgp.Uint64Size + 9 + msgp.Uint64Size + 6 + msgp.ArrayHeaderSize
	for za0001 := range z.Stats {
		s += z.Stats[za0001].Msgsize()
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *statsPayload) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Hostname":
			z.Hostname, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Env":
			z.Env, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Version":
			z.Version, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Stats":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Stats) >= int(zb0002) {
				z.Stats = (z.Stats)[:zb0002]
			} else {
				z.Stats = make([]statsBucket, zb0002)
			}
			for za0001 := range z.Stats {
				err = z.Stats[za0001].DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *statsPayload) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "Hostname"
	err = en.Append(0x84, 0xa8, 0x48, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Hostname)
	if err != nil {
		return
	}
	// write "Env"
	err = en.Append(0xa3, 0x45, 0x6e, 0x76)
	if err != nil {
		return
	}
	err = en.WriteString(z.Env)
	if err != nil {
		return
	}
	// write "Version"
	err = en.Append(0xa7, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteString(z.Version)
	if err != nil {
		return
	}
	// write "Stats"
	err = en.Append(0xa5, 0x53, 0x74, 0x61, 0x74, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Stats)))
	if err != nil {
		return
	}
	for za0001 := range z.Stats {
		err = z.Stats[za0001].EncodeMsg(en)
		if err != nil {
			return
		}
	}
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *statsPayload) Msgsize() (s int) {
	s = 1 + 9 + msgp.StringPrefixSize + len(z.Hostname) + 4 + msgp.StringPrefixSize + len(z.Env) + 8 + msgp.StringPrefixSize + len(z.Version) + 6 + msgp.ArrayHeaderSize
	for za0001 := range z.Stats {
		s += z.Stats[za0001].Msgsize()
	}
	return
}
