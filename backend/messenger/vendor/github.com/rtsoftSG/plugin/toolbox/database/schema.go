package database

type DataType int

const (
	_ DataType = iota
	TypeInt8
	TypeInt32
	TypeInt64
	TypeFloat64
	TypeString
	TypeArrayInt64
	TypeArrayFloat64

	TypeNullInt8
	TypeNullInt64
	TypeNullFloat64
	TypeNullString
)

type (
	Column struct {
		name     string
		dataType DataType
	}
	SchemaBuilder interface {
		SetTableName(string) SchemaBuilder
		SetColumns(...Column) SchemaBuilder
		Build() error
	}
)

func NewColumn(name string, colType DataType) Column {
	return Column{name, colType}
}

func (c *Column) DataType() DataType {
	return c.dataType
}

func (c *Column) Name() string {
	return c.name
}
