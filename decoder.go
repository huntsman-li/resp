//author: liyan

package resp

import (
	"errors"
	"io"
	"reflect"
	"strconv"
)

// Decoder 从输入流中读取和解码 RESP 对象。
type Decoder struct {
	r *Reader
}

// NewDecoder 创建并返回 Decoder.
func NewDecoder(r io.Reader) *Decoder {
	d := &Decoder{
		r: NewReader(r),
	}
	return d
}

// 尝试 decode 下一条Message.
func (d *Decoder) next(out *Message) (err error) {
	// 期望以 \r\n 结尾的消息
	var line []byte
	if out.Type, line, err = d.r.ReadLine(); err != nil {
		return err
	}

	switch out.Type {

	case StringHeader:
		out.Status = string(line)
		return

	case ErrorHeader:
		out.Error = errors.New(string(line))
		return

	case IntegerHeader:
		if out.Integer, err = strconv.ParseInt(string(line), 10, 64); err != nil {
			return err
		}
		return

	case BulkHeader:
		//获取字符串长度.
		var msgLen int

		if msgLen, err = strconv.Atoi(string(line)); err != nil {
			return
		}

		if msgLen > bulkMessageMaxLength {
			err = ErrMessageIsTooLarge
			return
		}

		if msgLen < 0 {
			out.IsNil = true
			return
		}

		if out.Bytes, err = d.r.ReadMessageBytes(msgLen); err != nil {
			return
		}

		return
	case ArrayHeader:
		// 获取字符串长度.
		var arrLen int

		if arrLen, err = strconv.Atoi(string(line)); err != nil {
			return
		}

		if arrLen < 0 {
			// The concept of Null Array exists as well, and is an alternative way to
			// 指定一个 Null 值 (通常使用 Null Bulk String , 但是由于
			// 历史原因，有两种格式).
			out.IsNil = true
			return
		}

		out.Array = make([]*Message, arrLen)

		for i := 0; i < arrLen; i++ {
			out.Array[i] = new(Message)
			if err = d.next(out.Array[i]); err != nil {
				return err
			}
		}

		return
	}

	return ErrInvalidInput
}

// Decode 尝试在缓冲区解码整个消息
func (d *Decoder) Decode(v interface{}) (err error) {
	out := new(Message)

	if err = d.next(out); err != nil {
		return err
	}

	if v == nil {
		return ErrExpectingDestination
	}

	dst := reflect.ValueOf(v)

	if dst.Kind() != reflect.Ptr || dst.IsNil() {
		return ErrExpectingPointer
	}

	if err = redisMessageToType(dst.Elem(), out); err != nil {
		if out.Type == ErrorHeader {
			return errors.New(out.Error.Error())
		}
	}

	return err
}
