package pg

import "fmt"

const (
	BtreeIndex = "btree"
	HashIndex  = "hash"
	GinIndex   = "gin"
)

type Index struct {
	TableName string
	Type      string
	Field     string
	Parameter string
}

func (idx Index) BuildSql() string {
	if len(idx.Type) == 0 {
		idx.Type = BtreeIndex
	}
	if len(idx.Parameter) == 0 {
		idx.Parameter = idx.Field
	}
	return fmt.Sprintf("create index if not exists %s_%s_%s_idx on %s using %s(%s)", idx.TableName, idx.Field, idx.Type, idx.TableName, idx.Type, idx.Parameter)
}
