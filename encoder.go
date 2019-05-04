//author: liyan

package resp

import (
	"io"
	"sync"
)

const digitbuflen = 20

var (
	encoderNil = []byte("$-1\r\n")
	digits     = []byte("0123456789")
)

func intToBytes(v int) []byte {
	buf := make([]byte, digitbuflen)

	i := len(buf)

	for v >= 10 {
		i--
		buf[i] = digits[v%10]
		v = v / 10
	}

	i--
	buf[i] = digits[v%10]

	return buf[i:]
}

// Encoder 提供 Encode() 方式，直接把数据编码到io.Writer.
type Encoder struct {
	w   io.Writer
	buf []byte
	mu  *sync.Mutex
}

// NewEncoder 使用初始的io.Writer创建并返回*Encoder结构值
func NewEncoder(w io.Writer) *Encoder {
	e := &Encoder{
		w:   w,
		buf: []byte{},
		mu:  new(sync.Mutex),
	}
	return e
}

// 编码将给定参数序列化到RESP消息中并输出
func (e *Encoder) Encode(v interface{}) error {
	return e.writeEncoded(e.w, v)
}

func (e *Encoder) writeEncoded(w io.Writer, data interface{}) (err error) {

	var b []byte

	switch v := data.(type) {

	case []byte:
		n := intToBytes(len(v))

		b = make([]byte, 0, 1+len(n)+2+len(v)+2)

		b = append(b, BulkHeader)
		b = append(b, n...)
		b = append(b, endOfLine...)
		b = append(b, v...)
		b = append(b, endOfLine...)

	case string:
		q := []byte(v)

		b = make([]byte, 0, 1+len(q)+2)
		b = append(b, StringHeader)
		b = append(b, q...)
		b = append(b, endOfLine...)

	case error:
		q := []byte(v.Error())

		b = make([]byte, 0, 1+len(q)+2)
		b = append(b, ErrorHeader)
		b = append(b, q...)
		b = append(b, endOfLine...)

	case int:
		q := intToBytes(int(v))
		b = make([]byte, 0, 1+len(q)+2)
		b = append(b, IntegerHeader)
		b = append(b, q...)
		b = append(b, endOfLine...)

	case [][]byte:
		n := intToBytes(len(v))

		b = make([]byte, 0, 1+len(n)+2)
		b = append(b, ArrayHeader)
		b = append(b, n...)
		b = append(b, endOfLine...)

		for i := range v {
			q := intToBytes(len(v[i]))

			z := make([]byte, 0, 1+len(q)+2+len(v[i])+2)

			z = append(z, BulkHeader)
			z = append(z, q...)
			z = append(z, endOfLine...)
			z = append(z, v[i]...)
			z = append(z, endOfLine...)

			b = append(b, z...)
		}

	case []string:
		q := intToBytes(len(v))

		b = make([]byte, 0, 1+len(q)+2)
		b = append(b, ArrayHeader)
		b = append(b, q...)
		b = append(b, endOfLine...)

		for i := range v {
			p := []byte(v[i])

			z := make([]byte, 0, 1+len(p)+2)
			z = append(z, StringHeader)
			z = append(z, p...)
			z = append(z, endOfLine...)

			b = append(b, z...)
		}

	case []int:
		n := intToBytes(len(v))

		b = make([]byte, 0, 1+len(n)+2)
		b = append(b, ArrayHeader)
		b = append(b, n...)
		b = append(b, endOfLine...)

		for i := range v {
			m := intToBytes(v[i])

			z := make([]byte, 0, 1+len(m)+2)
			z = append(z, IntegerHeader)
			z = append(z, m...)
			z = append(z, endOfLine...)

			b = append(b, z...)
		}

	case []interface{}:
		q := intToBytes(len(v))

		b = make([]byte, 0, 1+len(q)+2)
		b = append(b, ArrayHeader)
		b = append(b, q...)
		b = append(b, endOfLine...)

		e.buf = append(e.buf, b...)

		if w != nil {
			e.mu.Lock()
			w.Write(e.buf)
			e.buf = []byte{}
			e.mu.Unlock()
		}

		for i := range v {
			if err = e.writeEncoded(w, v[i]); err != nil {
				return err
			}
		}

		return nil

	case *Message:
		switch v.Type {
		case ErrorHeader:
			return e.writeEncoded(w, v.Error)
		case IntegerHeader:
			return e.writeEncoded(w, int(v.Integer))
		case BulkHeader:
			return e.writeEncoded(w, v.Bytes)
		case StringHeader:
			return e.writeEncoded(w, v.Status)
		case ArrayHeader:
			return e.writeEncoded(w, v.Array)
		default:
			return ErrMissingMessageHeader
		}

	case nil:
		b = make([]byte, 0, len(encoderNil))
		b = append(b, encoderNil...)

	default:
		return ErrInvalidInput
	}

	e.buf = append(e.buf, b...)

	if w != nil {
		e.mu.Lock()
		w.Write(e.buf)
		e.buf = []byte{}
		e.mu.Unlock()
	}

	return nil
}
