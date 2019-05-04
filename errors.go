//author: liyan

package resp

import (
	"errors"
)

var (
	// ErrInvalidInput 在解码时返回任意错误的。
	ErrInvalidInput = errors.New(`resp: Invalid input`)

	// ErrMessageIsTooLarge Message 创建时，缓冲区过大
	ErrMessageIsTooLarge = errors.New(`resp: Message is too large`)

	// ErrMissingMessageHeader 返回没有 header 的
	ErrMissingMessageHeader = errors.New(`resp: Missing message header`)

	// ErrExpectingPointer
	// parameter. 非指针类型的返回错误
	ErrExpectingPointer = errors.New(`resp: Expecting pointer value`)

	// ErrUnsupportedConversion 转换出现不兼容的
	ErrUnsupportedConversion = errors.New(`resp: Unsupported conversion: %s to %s`)

	// ErrMessageIsNil 如果编码消息为空时，返回
	ErrMessageIsNil = errors.New(`resp: Message is nil`)

	// ErrMissingReader is 未定义的数据，返回
	ErrMissingReader = errors.New(`resp: Ran out of buffered data and a reader was not defined`)

	// ErrExpectingDestination 编码返回空值， 返回
	ErrExpectingDestination = errors.New(`resp: Expecting a valid destination, but a nil value was provided`)
)
