package database

type (
	Field struct {
		name  string
		value interface{}
	}
	InsertCommand struct {
		tableName string
		fields    []*Field
	}
)

func NewInsertCommand(tableName string) *InsertCommand {
	return &InsertCommand{tableName: tableName}
}

func (i *InsertCommand) Fields() []*Field {
	return i.fields
}

func (i *InsertCommand) TableName() string {
	return i.tableName
}

func NewField(name string, value interface{}) *Field {
	return &Field{name, value}
}

func (f *Field) Value() interface{} {
	return f.value
}

func (f *Field) Name() string {
	return f.name
}

func (i *InsertCommand) WithFields(fields ...*Field) *InsertCommand {
	i.fields = append(i.fields, fields...)
	return i
}
