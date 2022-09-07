package frame

import "github.com/yomorun/y3"

// StreamFrame is a frame for stream.
type StreamFrame struct {
	sinkTags []byte
	metadata []byte
}

// NewStreamFrame creates a new StreamFrame.
func NewStreamFrame(sinkTags []byte) *StreamFrame {
	return &StreamFrame{
		sinkTags: sinkTags,
	}
}

// Type gets the type of StreamFrame.
func (f *StreamFrame) Type() Type {
	return TagOfStreamFrame
}

// SinkTags gets the sink tags of StreamFrame.
func (f *StreamFrame) SinkTags() []byte {
	return f.sinkTags
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
	// sink
	sinkBlock := y3.NewPrimitivePacketEncoder(byte(TagOfStreamSinkTags))
	sinkBlock.SetBytesValue(f.sinkTags)
	stream.AddPrimitivePacket(sinkBlock)
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
		case byte(TagOfStreamSinkTags):
			stream.sinkTags = v.ToBytes()
			break
		case byte(TagOfStreamMetadata):
			stream.metadata = v.ToBytes()
			break
		}
	}
	return stream, nil
}
