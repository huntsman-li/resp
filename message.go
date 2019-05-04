//author: liyan

package resp

const (
	// StringHeader 是用于为单一字符串 (或状态
	// 消息)。 String 消息不是二进制安全的。
	StringHeader = '+'
	// ErrorHeader 是用于为错误消息添加前缀的 Header。
	ErrorHeader = '-'
	// IntegerHeader 是用于为整数添加前缀的 Header。
	IntegerHeader = ':'
	// BulkHeader 是用于为二进制安全消息添加前缀的 Header。
	BulkHeader = '$'
	// ArrayHeader 是用于为消息数组添加前缀的 Header。
	ArrayHeader = '*'
)

// Message 是 RESP 的数据结构
type Message struct {
	Error   error
	Integer int64
	Bytes   []byte
	Status  string
	Array   []*Message
	IsNil   bool
	Type    byte
}

// SetStatus 插入一个数据状态.
func (m *Message) SetStatus(s string) {
	m.Type = StringHeader
	m.Status = s
}

// SetError 插入消息类型错误.
func (m *Message) SetError(e error) {
	m.Type = ErrorHeader
	m.Error = e
}

// SetInteger 设置整数类型的消息.
func (m *Message) SetInteger(i int64) {
	m.Type = IntegerHeader
	m.Integer = i
}

// SetBytes 设置二进制安全消息
func (m *Message) SetBytes(b []byte) {
	m.Type = BulkHeader
	m.Bytes = b
}

// SetArray 插入一个消息数组类型.
func (m *Message) SetArray(a []*Message) {
	m.Type = ArrayHeader
	m.Array = a
}

// SetNil 插入消息等于Null.
func (m *Message) SetNil() {
	m.Type = 0
	m.IsNil = true
}

// Interface 返回消息的当前值，为interface{}
func (m Message) Interface() interface{} {
	switch m.Type {
	case ErrorHeader:
		return m.Error
	case IntegerHeader:
		return m.Integer
	case BulkHeader:
		return m.Bytes
	case StringHeader:
		return m.Status
	case ArrayHeader:
		return m.Array
	}
	return nil
}
