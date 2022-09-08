package frame

import "github.com/yomorun/y3"

// StreamFrame is a frame for stream.
type StreamFrame struct {
	dataTags []byte
	metadata []byte
}

// NewStreamFrame creates a new StreamFrame.
func NewStreamFrame(dataTags []byte) *StreamFrame {
	return &StreamFrame{
		dataTags: dataTags,
	}
}

// Type gets the type of StreamFrame.
func (f *StreamFrame) Type() Type {
	return TagOfStreamFrame
}

// DataTags gets the data tags of StreamFrame.
func (f *StreamFrame) DataTags() []byte {
	return f.dataTags
}

// Metadata gets the metadata of StreamFrame.
func (f *StreamFrame) Metadata() []byte {
	return f.metadata
}

// SetMetadata sets the metadata of StreamFrame.
func (f *StreamFrame) SetMetadata(metadata []byte) {
	f.metadata = metadata
}

// Encode to Y3 encoded bytes.
func (f *StreamFrame) Encode() []byte {
	stream := y3.NewNodePacketEncoder(byte(f.Type()))
	// tags
	tagsBlock := y3.NewPrimitivePacketEncoder(byte(TagOfStreamDataTags))
	tagsBlock.SetBytesValue(f.dataTags)
	stream.AddPrimitivePacket(tagsBlock)
	// metadata
	if f.metadata != nil {
		metaBlock := y3.NewPrimitivePacketEncoder(byte(TagOfStreamMetadata))
		metaBlock.SetBytesValue(f.metadata)
		stream.AddPrimitivePacket(metaBlock)
	}

	return stream.Encode()
}

// DecodeToStreamFrame decodes Y3 encoded bytes to StreamFrame.
func DecodeToStreamFrame(buf []byte) (*StreamFrame, error) {
	nodeBlock := y3.NodePacket{}
	_, err := y3.DecodeToNodePacket(buf, &nodeBlock)
	if err != nil {
		return nil, err
	}
	stream := &StreamFrame{}
	for k, v := range nodeBlock.PrimitivePackets {
		switch k {
		case byte(TagOfStreamDataTags):
			stream.dataTags = v.ToBytes()
			break
		case byte(TagOfStreamMetadata):
			stream.metadata = v.ToBytes()
			break
		}
	}
	return stream, nil
}
