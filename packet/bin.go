package packet

// type PacketBin struct {
// 	Timestamp [PACKET_TIMESTAMP_SIZE_BYTES]byte
// 	Type      [PACKET_TYPE_SIZE_BYTES]byte
// 	Payload   [PACKET_PAYLOAD_SIZE_BYTES]byte
// }

// func New(packType PacketType, payload [PACKET_PAYLOAD_SIZE_BYTES]byte) PacketBin {
// 	p := PacketBin{}

// 	switch packType {
// 	case PING_PACKET:
// 		p.Type = [PACKET_TYPE_SIZE_BYTES]byte{'p', 'i', 'n', 'g'}
// 	case PONG_PACKET:
// 		p.Type = [PACKET_TYPE_SIZE_BYTES]byte{'p', 'o', 'n', 'g'}
// 	default:
// 		p.Type = [PACKET_TYPE_SIZE_BYTES]byte{'n', 'o', 'n', 'e'}
// 	}

// 	binary.BigEndian.PutUint64(p.Timestamp[:], uint64(time.Now().UnixNano()))
// 	p.Payload = payload

// 	return p
// }

// func (p PacketBin) Marshal() [MAX_SAFE_UDP_PAYLOAD_SIZE_BYTES]byte {
// 	var buffer bytes.Buffer
// 	enc := gob.NewEncoder(&buffer)
// 	if err := enc.Encode(p); err != nil {
// 		return [MAX_SAFE_UDP_PAYLOAD_SIZE_BYTES]byte{}
// 	}
// 	var data [MAX_SAFE_UDP_PAYLOAD_SIZE_BYTES]byte
// 	copy(data[:], buffer.Bytes()[:MAX_SAFE_UDP_PAYLOAD_SIZE_BYTES])
// 	return data
// }

// func (p PacketBin) String() string {
// 	str := ""

// 	return str
// }

// func (p PacketBin) Print() {
// 	fmt.Println(p.Timestamp, p.Type)
// }

// // func Unmarshal(data [MAX_SAFE_UDP_PAYLOAD_SIZE_BYTES]byte) PacketBin {
// // 	p := PacketBin{}

// // 	var buffer bytes.Buffer
// // 	buffer.Write(data[:])
// // 	if err := binary.Read(&buffer, binary.BigEndian, &p); err != nil {
// // 		log.Fatal(err)
// // 		return PacketBin{}
// // 	}

// // 	return p
// // }
