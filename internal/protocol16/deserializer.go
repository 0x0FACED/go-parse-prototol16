package protocol16

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"sync"
	"unsafe"

	"github.com/0x0FACED/go-parse-protocol16/internal/protocol16/photon"
	"go.uber.org/zap"
)

// deserializer intfc
type Deserializer interface {
	DeserializeNoCode(input protocol16Stream) (any, error)
	Deserialize(input protocol16Stream, code byte) any
}

type deserializer struct {
	bufferPool sync.Pool
	logger     *zap.Logger
}

var (
	ErrGetBufferFromPool = errors.New("cannot get buffer from pool")
)

func NewDeserializer() *deserializer {
	log, err := zap.NewProduction()
	if err != nil {
		return nil
	}
	return &deserializer{
		bufferPool: sync.Pool{
			New: func() any {
				return make([]byte, 8)
			},
		},
		logger: log,
	}
}

// func (d *deserializer) getByteBuffer() []byte {
// 	return d.bufferPool.Get().([]byte)
// }

// func (d *deserializer) putByteBuffer(buffer []byte) {
// 	d.bufferPool.Put(&buffer)
// }

func (d *deserializer) DeserializeNoCode(input protocol16Stream) (any, error) {
	b, err := input.ReadByte()
	if err != nil {
		return nil, errors.New("cant read byte DeserializeNoCode()")
	}
	res, err := d.Deserialize(input, b)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *deserializer) Deserialize(input protocol16Stream, code byte) (any, error) {
	switch Protocol16Type(code) {
	case Protocol16Type(Unknown):
	case Protocol16Type(Null):
		return nil, nil
	case Protocol16Type(Dictionary):
		return d.deserializeDictionary(input)
	case Protocol16Type(StringArray):
		return d.deserializeStringArray(input)
	case Protocol16Type(Byte):
		return d.deserializeByte(input)
	case Protocol16Type(EventData):
		return d.deserializeEventData(input)
	case Protocol16Type(Double):
		return d.deserializeDouble(input)
	case Protocol16Type(Float):
		return d.deserializeFloat(input)
	case Protocol16Type(Integer):
		return d.deserializeInt(input)
	case Protocol16Type(Hashtable):
		return d.deserializeHashtable(input)
	case Protocol16Type(Short):
		return d.deserializeShort(input)
	case Protocol16Type(Long):
		return d.deserializeLong(input)
	case Protocol16Type(IntegerArray):
		return d.deserializeIntArray(input)
	case Protocol16Type(Boolean):
		return d.deserializeBoolean(input)
	case Protocol16Type(OperationResponse):
		return d.deserializeOperationResp(input)
	case Protocol16Type(OperationRequest):
		return d.deserializeOperationReq(input)
	case Protocol16Type(String):
		return d.deserializeString(input)
	case Protocol16Type(ByteArray):
		return d.deserializeByteArray(input)
	case Protocol16Type(Array):
		return d.deserializeArray(input)
	case Protocol16Type(ObjectArray):
		return d.deserializeObjectArray(input)
	default:
		return nil, fmt.Errorf("unknown code: %v", code)
	}

	return nil, nil
}

func (d *deserializer) deserializeOperationReq(input protocol16Stream) (*photon.OperationRequest, error) {
	code, err := d.deserializeByte(input)
	if err != nil {
		return nil, err
	}
	params, err := d.deserializeParamTable(input)
	if err != nil {
		return nil, err
	}

	return &photon.OperationRequest{Code: code, Params: params}, nil
}

func (d *deserializer) deserializeOperationResp(input protocol16Stream) (*photon.OperationResponse, error) {
	code, err := d.deserializeByte(input)
	if err != nil {
		return nil, err
	}
	respCode, err := d.deserializeShort(input)
	if err != nil {
		return nil, err
	}

	// MAYBE ERR
	debug, err := d.deserializeString(input)
	if err != nil {
		return nil, err
	}
	debugMsg := string(debug)
	params, err := d.deserializeParamTable(input)
	if err != nil {
		return nil, err
	}

	return &photon.OperationResponse{Code: code, Params: params, DebugMsg: debugMsg, RespCode: photon.Short(respCode)}, nil
}

func (d *deserializer) deserializeEventData(input protocol16Stream) (*photon.EventData, error) {
	code, err := d.deserializeByte(input)
	if err != nil {
		return nil, err
	}

	params, err := d.deserializeParamTable(input)
	if err != nil {
		return nil, err
	}

	return &photon.EventData{Code: code, Params: params}, nil
}

func (d *deserializer) deserializeByte(input protocol16Stream) (byte, error) {
	b, err := input.ReadByte()
	if err != nil {
		return 0, err
	}
	return b, nil
}

func (d *deserializer) deserializeBoolean(input protocol16Stream) (bool, error) {
	b, err := input.ReadByte()
	if err != nil {
		return false, err
	}
	return b != 0, nil
}

func (d *deserializer) deserializeShort(input protocol16Stream) (int16, error) {
	buf := d.bufferPool.Get()
	defer d.bufferPool.Put(buf)

	byteBuf, ok := byteBuf(buf)
	if !ok {
		return -1, ErrGetBufferFromPool
	}
	// size of short = 2
	_, err := input.Read(byteBuf, 0, 2)
	if err != nil {
		return 0, err
	}

	return int16(byteBuf[0]<<8) | int16(byteBuf[1]), nil
}

func (d *deserializer) deserializeInt(input protocol16Stream) (int, error) {
	buf := d.bufferPool.Get()
	defer d.bufferPool.Put(buf)

	byteBuf, ok := byteBuf(buf)
	if !ok {
		return -1, ErrGetBufferFromPool
	}
	// size of int = 4
	_, err := input.Read(byteBuf, 0, 4)
	if err != nil {
		return -1, err
	}

	return int(byteBuf[0]<<24) | int(byteBuf[1]<<16) | int(byteBuf[2]<<8) | int(byteBuf[3]), nil
}

func (d *deserializer) deserializeHashtable(input protocol16Stream) (int16, error) {
	panic("not impl")
}

func (d *deserializer) deserializeLong(input protocol16Stream) (int64, error) {
	buf := d.bufferPool.Get()
	defer d.bufferPool.Put(buf)

	byteBuf, ok := byteBuf(buf)
	if !ok {
		return -1, ErrGetBufferFromPool
	}

	_, err := input.Read(byteBuf, 0, 8)
	if err != nil {
		return -1, err
	}

	if isLittleEndian() {
		return int64(byteBuf[0])<<56 | int64(byteBuf[1])<<48 | int64(byteBuf[2])<<40 | int64(byteBuf[3])<<32 |
			int64(byteBuf[4])<<24 | int64(byteBuf[5])<<16 | int64(byteBuf[6])<<8 | int64(byteBuf[7]), nil
	}

	return int64(binary.BigEndian.Uint64(byteBuf[:8])), nil
}

func (d *deserializer) deserializeFloat(input protocol16Stream) (float32, error) {
	buf := d.bufferPool.Get()
	defer d.bufferPool.Put(buf)

	byteBuf, ok := byteBuf(buf)
	if !ok {
		return -1, ErrGetBufferFromPool
	}

	_, err := input.Read(byteBuf, 0, 4)
	if err != nil {
		return -1, err
	}

	if isLittleEndian() {
		byteBuf[0], byteBuf[3] = byteBuf[3], byteBuf[0]
		byteBuf[1], byteBuf[2] = byteBuf[2], byteBuf[1]
	}

	bits := binary.BigEndian.Uint32(byteBuf)
	return math.Float32frombits(bits), nil
}

func (d *deserializer) deserializeDouble(input protocol16Stream) (float64, error) {
	buf := d.bufferPool.Get()
	defer d.bufferPool.Put(buf)

	byteBuf, ok := byteBuf(buf)
	if !ok {
		return -1, ErrGetBufferFromPool
	}

	_, err := input.Read(byteBuf, 0, 8)
	if err != nil {
		return -1, err
	}

	if isLittleEndian() {
		byteBuf[0], byteBuf[7] = byteBuf[7], byteBuf[0]
		byteBuf[1], byteBuf[6] = byteBuf[6], byteBuf[1]
		byteBuf[2], byteBuf[5] = byteBuf[5], byteBuf[2]
		byteBuf[3], byteBuf[4] = byteBuf[4], byteBuf[3]
	}

	bits := binary.BigEndian.Uint64(byteBuf)
	return math.Float64frombits(bits), nil
}

func (d *deserializer) deserializeString(input protocol16Stream) (string, error) {
	size, err := d.deserializeShort(input)
	if err != nil {
		return "", err
	}

	if size == 0 {
		return "", nil
	}

	buffer := make([]byte, size)
	_, err = input.Read(buffer, 0, int(size))
	if err != nil {
		return "", err
	}

	return string(buffer), nil
}

func (d *deserializer) deserializeByteArray(input protocol16Stream) ([]byte, error) {
	size, err := d.deserializeInt(input)
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, size)
	_, err = input.Read(buffer, 0, size)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func (d *deserializer) deserializeIntArray(input protocol16Stream) ([]int, error) {
	size, err := d.deserializeInt(input)
	if err != nil {
		return nil, err
	}

	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i], err = d.deserializeInt(input)
		if err != nil {
			return nil, err
		}
	}

	return arr, nil
}

func (d *deserializer) deserializeStringArray(input protocol16Stream) ([]string, error) {
	size, err := d.deserializeShort(input)
	if err != nil {
		return nil, err
	}

	arr := make([]string, int(size))
	for i := 0; i < int(size); i++ {
		arr[i], err = d.deserializeString(input)
		if err != nil {
			return nil, err
		}
	}

	return arr, nil
}

func (d *deserializer) deserializeObjectArray(input protocol16Stream) ([]any, error) {
	size, err := d.deserializeShort(input)
	if err != nil {
		return nil, err
	}

	arr := make([]any, int(size))
	for i := 0; i < int(size); i++ {
		code, err := input.ReadByte()
		if err != nil {
			return nil, err
		}
		arr[i], err = d.Deserialize(input, code)
		if err != nil {
			return nil, err
		}
	}

	return arr, nil
}

func (d *deserializer) deserializeDictionary(input protocol16Stream) (map[any]any, error) {
	keyTypeCode, _ := input.ReadByte()
	valTypeCode, _ := input.ReadByte()
	dictSize, _ := d.deserializeShort(input)

	output := make(map[any]any)

	err := d.deserializeDictElements(input, output, int(dictSize), byte(keyTypeCode), byte(valTypeCode))
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (d *deserializer) deserializeDictElements(input protocol16Stream, dict map[any]any, size int, keyType, valType byte) error {
	for i := 0; i < size; i++ {
		key, err := d.Deserialize(input, keyType)
		if err != nil {
			return err
		}
		val, err := d.Deserialize(input, valType)
		if err != nil {
			return err
		}

		dict[key] = val
	}

	return nil
}
func (d *deserializer) deserializeArray(input protocol16Stream) ([]any, error) {
	size, err := d.deserializeShort(input)
	if err != nil {
		return nil, err
	}

	code, err := input.ReadByte()
	if err != nil {
		return nil, err
	}

	switch Protocol16Type(code) {
	case Protocol16Type(Array):
		arr, err := d.deserializeArray(input)
		if err != nil {
			return nil, err
		}

		res := make([]any, size)
		res[0] = arr
		for i := 1; i < int(size); i++ {
			arr, err = d.deserializeArray(input)
			if err != nil {
				return nil, err
			}

			res[i] = arr
		}

		return res, nil

	case Protocol16Type(ByteArray):
		byteArr := make([]any, size)
		for i := 0; i < int(size); i++ {
			byteArr[i], err = d.deserializeByteArray(input)
			if err != nil {
				return nil, err
			}
		}
		return byteArr, nil
	case Protocol16Type(Dictionary):
		res, err := d.deserializeDictArray(input, int(size))
		if err != nil {
			return nil, err
		}
		return res, nil
	default:
		res := make([]any, size)

		for i := 0; i < int(size); i++ {
			val, err := d.Deserialize(input, code)
			if err != nil {
				return nil, err
			}
			res[i] = val
		}

		return res, nil
	}
}

func (d *deserializer) deserializeDictArray(input protocol16Stream, size int) ([]any, error) {
	dictType, _ := d.deserializeDictType(input)

	res := make([]any, size)

	for i := 0; i < size; i++ {
		dict := make(map[any]any, 0)

		arrSize, err := d.deserializeShort(input)
		if err != nil {
			return nil, err
		}

		for j := 0; j < int(arrSize); j++ {
			var key any
			if dictType[0] > 0 {
				key, _ = d.Deserialize(input, dictType[0])
			} else {
				nextKeyCode, _ := input.ReadByte()
				key, _ = d.Deserialize(input, nextKeyCode)
			}
			var val any
			if dictType[1] > 0 {
				val, _ = d.Deserialize(input, dictType[1])
			} else {
				nextValCode, _ := input.ReadByte()
				val, _ = d.Deserialize(input, nextValCode)
			}

			dict[key] = val

		}

		res = append(res, dict)
	}

	return res, nil
}

func (d *deserializer) deserializeDictType(input protocol16Stream) ([]byte, error) {
	// dont read errors
	keyCode, _ := input.ReadByte()
	valCode, _ := input.ReadByte()

	return []byte{keyCode, valCode}, nil
}

func (d *deserializer) deserializeParamTable(input protocol16Stream) (map[byte]any, error) {
	panic("not impl")
}

// func getTypeOfCode(code byte) Protocol16Type {
// 	switch Protocol16Type(code) {
// 	case Protocol16Type(Unknown):
// 		return nil
// 	case Protocol16Type(Null):
// 		return nil
// 	case Protocol16Type(Dictionary):
// 	case Protocol16Type(StringArray):
// 	case Protocol16Type(Byte):
// 	case Protocol16Type(EventData):
// 	case Protocol16Type(Double):
// 	case Protocol16Type(Float):
// 	case Protocol16Type(Integer):
// 	case Protocol16Type(Hashtable):
// 	case Protocol16Type(Short):
// 	case Protocol16Type(Long):
// 	case Protocol16Type(IntegerArray):
// 	case Protocol16Type(Boolean):
// 	case Protocol16Type(OperationResponse):
// 	case Protocol16Type(OperationRequest):
// 	case Protocol16Type(String):
// 	case Protocol16Type(ByteArray):
// 	case Protocol16Type(Array):
// 	case Protocol16Type(ObjectArray):
// 	default:
// 		return nil
// 	}
// }

func byteBuf(buf any) ([]byte, bool) {
	byteBuf, ok := buf.([]byte)
	if !ok {
		return nil, false
	}
	return byteBuf, true
}

func isLittleEndian() bool {
	var x int = 1
	return *(*byte)(unsafe.Pointer(&x)) == 1
}
