package protocol16

import (
	"errors"
	"io"
)

var (
	ErrOutOfBounds = errors.New("out of bounds")
)

type protocol16Stream struct {
	buffer []byte
	pos    int // current stream position
	len    int
}

func NewStream(size int) *protocol16Stream {
	return &protocol16Stream{
		buffer: make([]byte, size),
		len:    size,
	}
}

func NewStreamFromBuffer(buf []byte) *protocol16Stream {
	return &protocol16Stream{
		buffer: buf,
	}
}

func (p *protocol16Stream) Read(buffer []byte, offset int, count int) (int, error) {
	dif := len(p.buffer) - p.pos
	if dif <= 0 {
		return -1, ErrOutOfBounds
	}
	if count > dif {
		count = dif
	}
	copy(buffer[offset:], p.buffer[p.pos:p.pos+count])
	p.pos += count

	return count, nil
}

func (s *protocol16Stream) ReadByte() (byte, error) {
	if s.pos >= s.len {
		return 0, ErrOutOfBounds
	}

	arr := s.buffer
	res := arr[s.pos]
	s.pos++
	return res, nil
}

func (p *protocol16Stream) Write(buffer []byte, offset int, count int) {
	sum := p.pos + count
	p.expand(sum)
	if sum > len(p.buffer) {
		p.len = sum
	}
	copy(p.buffer[p.pos:p.pos+count], buffer[offset:offset+count])
	p.pos = sum
}

func (s *protocol16Stream) Seek(offset int64, whence int) (int64, error) {
	var newPos int

	switch whence {
	case io.SeekStart:
		newPos = int(offset)
	case io.SeekCurrent:
		newPos = s.pos + int(offset)
	case io.SeekEnd:
		newPos = len(s.buffer) + int(offset)
	default:
		return 0, errors.New("invalid whence")
	}

	if newPos < 0 {
		return 0, errors.New("negative position")
	}

	if newPos > s.Length() {
		return 0, errors.New("seek after end")
	}

	s.pos = newPos
	return int64(s.pos), nil
}

func (s *protocol16Stream) Bytes() []byte {
	return append([]byte(nil), s.buffer...)
}

func (s *protocol16Stream) SetPosition(pos int) error {
	if pos < 0 || pos > len(s.buffer) {
		return errors.New("position out of bounds")
	}
	s.pos = pos
	return nil
}
func (s *protocol16Stream) Length() int {
	return s.len
}

func (s *protocol16Stream) Position() int {
	return s.pos
}

func (s *protocol16Stream) expand(requiredSize int) {
	if requiredSize <= len(s.buffer) {
		return
	}
	newBuffer := make([]byte, requiredSize)
	copy(newBuffer, s.buffer)
	s.buffer = newBuffer
}
