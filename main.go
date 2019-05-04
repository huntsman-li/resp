//author: liyan

package resp

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var endOfLine = []byte{'\r', '\n'}

var (
	typeErr     = reflect.TypeOf(errors.New(""))
	typeMessage = reflect.TypeOf(Message{})
)

const (
	// Bulk Strings 用于表示长度最大为512MB，最大单个二进制字符串
	bulkMessageMaxLength = 512 * 1024
)

func byteToTypeName(c byte) string {
	switch c {
	case StringHeader:
		return `status`
	case ErrorHeader:
		return `error`
	case IntegerHeader:
		return `integer`
	case BulkHeader:
		return `bulk`
	case ArrayHeader:
		return `array`
	}
	return `unknown`
}

// Marshal 返回 v 的 RESP 编码。 适用于以下类型
// string, int, []byte, nil 和 []interface{} 类型。
func Marshal(v interface{}) ([]byte, error) {
	//断言v的类型
	switch t := v.(type) {
	case string:
		// 如果传参的是一个字符串，我们将其转换为byte以使其成为二进制流
		v = []byte(t)
	}

	e := NewEncoder(nil)

	if err := e.Encode(v); err != nil {
		return nil, err
	}

	return e.buf, nil
}

// Unmarshal 解析 RESP 编码的数据并将结果存储
// v 适用于 string， int， []byte and
// []interface{} types。
func Unmarshal(data []byte, v interface{}) error {
	var err error

	if v == nil {
		return ErrExpectingPointer
	}

	r := bytes.NewReader(data)
	d := NewDecoder(r)

	if err = d.Decode(v); err != nil {
		return err
	}

	return nil
}

func redisMessageToType(dst reflect.Value, out *Message) error {

	if dst.Type() == typeMessage {
		dst.Set(reflect.ValueOf(*out))
		return nil
	}

	if out.IsNil {
		dst.Set(reflect.Zero(dst.Type()))
		return ErrMessageIsNil
	}

	dstKind := dst.Type().Kind()

	switch out.Type {
	case StringHeader:
		switch dstKind {
		// string -> string.
		case reflect.String:
			dst.Set(reflect.ValueOf(out.Status))
			return nil
		case reflect.Interface:
			dst.Set(reflect.ValueOf(out))
			return nil
		}
	case ErrorHeader:
		switch dstKind {
		// error -> string
		case reflect.String:
			dst.Set(reflect.ValueOf(out.Error.Error()))
			return nil
		// error -> serror
		case typeErr.Kind():
			dst.Set(reflect.ValueOf(out.Error))
			return nil
		case reflect.Interface:
			dst.Set(reflect.ValueOf(out))
			return nil
		}
	case IntegerHeader:
		switch dstKind {
		case reflect.Int:
			// integer -> integer.
			dst.Set(reflect.ValueOf(int(out.Integer)))
			return nil
		case reflect.Int64:
			// integer -> integer64.
			dst.Set(reflect.ValueOf(out.Integer))
			return nil
		case reflect.String:
			// integer -> string.
			dst.Set(reflect.ValueOf(strconv.FormatInt(out.Integer, 10)))
			return nil
		case reflect.Bool:
			// integer -> bool.
			if out.Integer == 0 {
				dst.Set(reflect.ValueOf(false))
			} else {
				dst.Set(reflect.ValueOf(true))
			}
			return nil
		case reflect.Interface:
			dst.Set(reflect.ValueOf(out))
			return nil
		}
	case BulkHeader:
		switch dstKind {
		case reflect.String:
			// []byte -> string
			dst.Set(reflect.ValueOf(string(out.Bytes)))
			return nil
		case reflect.Slice:
			// []byte -> []byte
			dst.Set(reflect.ValueOf(out.Bytes))
			return nil
		case reflect.Int:
			// []byte -> int
			n, _ := strconv.Atoi(string(out.Bytes))
			dst.Set(reflect.ValueOf(n))
			return nil
		case reflect.Int64:
			// []byte -> int64
			n, _ := strconv.Atoi(string(out.Bytes))
			dst.Set(reflect.ValueOf(int64(n)))
			return nil
		case reflect.Interface:
			dst.Set(reflect.ValueOf(out))
			return nil
		}
	case ArrayHeader:
		switch dstKind {
		// slice -> interface
		case reflect.Interface:
			var err error
			var elements reflect.Value
			total := len(out.Array)

			elements = reflect.MakeSlice(reflect.TypeOf([]interface{}{}), total, total)

			for i := 0; i < total; i++ {
				if err = redisMessageToType(elements.Index(i), out.Array[i]); err != nil {
					if err != ErrMessageIsNil {
						return err
					}
				}
			}

			dst.Set(elements)

			return nil
		// slice -> slice
		case reflect.Slice:
			var err error
			var elements reflect.Value
			total := len(out.Array)

			elements = reflect.MakeSlice(dst.Type(), total, total)

			for i := 0; i < total; i++ {
				if err = redisMessageToType(elements.Index(i), out.Array[i]); err != nil {
					if err != ErrMessageIsNil {
						return err
					}
				}
			}

			dst.Set(elements)

			return nil
		}
	}

	return fmt.Errorf(ErrUnsupportedConversion.Error(), byteToTypeName(out.Type), dstKind)
}
