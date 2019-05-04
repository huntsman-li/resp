//author: liyan

package resp

import (
	"bufio"
	"bytes"
	"io"
)

// Reader 从输入流中读取Redis Token
type Reader struct {
	br *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	d := &Reader{
		br: bufio.NewReader(r),
	}
	return d
}

// 读取输入流及类型
func (r *Reader) ReadLine() (lineType byte, line []byte, err error) {
	buf := bytes.NewBuffer(nil)
	end := endOfLine[len(endOfLine)-1]
	for !bytes.HasSuffix(buf.Bytes(), endOfLine) {
		if tmp, err := r.br.ReadBytes(end); err != nil {
			return 0, nil, err
		} else {
			buf.Write(tmp)
		}
	}
	// Line must be at least 1 byte + EOL marker
	if buf.Len() < (1 + len(endOfLine)) {
		return 0, nil, ErrInvalidInput
	}

	if lineType, err = buf.ReadByte(); err != nil {
		return 0, nil, err
	}
	buf.Truncate(buf.Len() - len(endOfLine))
	line = buf.Bytes()
	return lineType, line, nil
}

// Read a message from Redis of length n bytes (not including EOL marker)
func (r *Reader) ReadMessageBytes(n int) (buf []byte, err error) {
	bytesRemaining := n + len(endOfLine)
	buf = make([]byte, bytesRemaining)

	for {
		readStart := len(buf) - bytesRemaining
		var bytesRead int
		if bytesRead, err = r.br.Read(buf[readStart:]); err != nil {
			return nil, err
		}
		if bytesRead == bytesRemaining {
			break
		} else {
			bytesRemaining -= bytesRead
		}
	}
	// 消息必须以EOL标记终止
	if !bytes.HasSuffix(buf, endOfLine) {
		return nil, ErrInvalidInput
	}

	// 从返回缓冲区截断EOL标记
	buf = buf[:n]
	return buf, nil
}
