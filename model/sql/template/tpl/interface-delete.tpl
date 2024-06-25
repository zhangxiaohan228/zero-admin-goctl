Delete(ctx context.Context, id int64 ,session ...sqlx.Session) (int64, error)
ConditionDelete(ctx context.Context, where squirrel.Sqlizer,session ...sqlx.Session) (int64, error)
BatchDelete(ctx context.Context, id []int64,session ...sqlx.Session) (int64, error)