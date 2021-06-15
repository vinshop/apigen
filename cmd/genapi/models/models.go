package models

type DbModel struct {
	ModelName  string
	TableName  string
	Attributes []*DbModelAttribute

	NeedImport     bool
	NeedImportTime bool
	NeedImportGorm bool
}

type DbModelAttribute struct {
	FieldName    string
	FieldType    string
	ColumnName   string
	IsPrimaryKey bool
	IsNullable   bool
}
