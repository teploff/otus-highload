package database

type QueryFactoryFunc func() Query

type Query interface {
	From(table string) Query
	Select(expr string) Query
	AndWhere(expr string, args ...interface{}) Query
	OrWhere(expr string, args ...interface{}) Query
	GroupBy(columns ...string) Query
	Having(expr string, args ...interface{}) Query
	OrderBy(columns ...string) Query
	Limit(l int) Query
	Offset(o int) Query
	ToSql() (string, []interface{}, error)
	Exec() ([]map[string]interface{}, error)
}
