Insert(ctx context.Context, data *{{.upperStartCamelObject}}, session ...sqlx.Session) (sql.Result,error)
BatchInsert(ctx context.Context, data []*{{.upperStartCamelObject}}, session ...sqlx.Session) (sql.Result, error)
