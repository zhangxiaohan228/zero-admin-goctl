Update(ctx context.Context, {{if .containsIndexCache}}newData{{else}}data{{end}} *{{.upperStartCamelObject}}, session ...sqlx.Session) (int64, error)
ConditionUpdate(ctx context.Context, where squirrel.Sqlizer, value map[string]interface{}, session ...sqlx.Session) (int64, error)