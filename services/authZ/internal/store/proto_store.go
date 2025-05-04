package store

import (
	"errors"
	"fmt"
	"time"

	"github.com/klauspost/compress/zstd"
	"google.golang.org/protobuf/proto"
)

type ProtoStore struct {
	db      InMemoryStore
	encoder *zstd.Encoder
	decoder *zstd.Decoder
}

func NewProtoStore(db InMemoryStore) (*ProtoStore, error) {
	encoder, err := zstd.NewWriter(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create zstd encoder: %w", err)
	}
	decoder, err := zstd.NewReader(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create zstd decoder: %w", err)
	}
	return &ProtoStore{db: db, encoder: encoder, decoder: decoder}, nil
}

// SetProto stores a compressed protobuf message with TTL
func (p *ProtoStore) SetStoreWithTTL(key string, msg proto.Message, ttlSeconds int) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal proto: %w", err)
	}
	compressed := p.encoder.EncodeAll(data, make([]byte, 0, len(data)))
	return p.db.SetWithTTL(key, compressed, toDuration(ttlSeconds))
}

// SetProto stores a compressed protobuf message
func (p *ProtoStore) SetStore(key string, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal proto: %w", err)
	}
	compressed := p.encoder.EncodeAll(data, make([]byte, 0, len(data)))
	return p.db.Set(key, compressed)
}

// GetProto retrieves a protobuf message by key and unmarshals into the target
func (p *ProtoStore) GetProto(key string, out proto.Message) error {
	raw, err := p.db.Get(key)
	if err != nil {
		return err
	}
	compressed, ok := raw.([]byte)
	if !ok {
		return errors.New("invalid cache value type")
	}
	decompressed, err := p.decoder.DecodeAll(compressed, nil)
	if err != nil {
		return fmt.Errorf("zstd decompression failed: %w", err)
	}
	return proto.Unmarshal(decompressed, out)
}

func toDuration(ttlSeconds int) time.Duration {
	return time.Duration(ttlSeconds) * time.Second
}
