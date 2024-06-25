// Insert 方法
func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, data *{{.upperStartCamelObject}}, session ...sqlx.Session) (sql.Result, error) {

    query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)

    var (
        result sql.Result
        err    error
    )
    if len(session) > 0 && session[0] != nil {
        result, err = session[0].ExecCtx(ctx, query, {{.expressionValues}})
    } else {
        result, err = m.conn.ExecCtx(ctx, query, {{.expressionValues}})
    }

    return result, err
}

// BatchInsert 方法
func (m *default{{.upperStartCamelObject}}Model) BatchInsert(ctx context.Context, data []*{{.upperStartCamelObject}}, session ...sqlx.Session) (sql.Result, error) {

    valueStrings := make([]string, 0, len(data))
    valueArgs := make([]interface{}, 0, len(data))

    for _, data := range data {
        valueStrings = append(valueStrings, fmt.Sprintf("({{.expression}})"))
        valueArgs = append(valueArgs, {{.expressionValues}})
    }

    query := fmt.Sprintf("insert into %s (%s) values %s", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet, strings.Join(valueStrings, ","))

    var (
        result sql.Result
        err    error
    )
    if len(session) > 0 && session[0] != nil {
        result, err = session[0].ExecCtx(ctx, query, valueArgs...)
    } else {
        result, err = m.conn.ExecCtx(ctx, query, valueArgs...)
    }

    if err != nil {
        return nil, err
    }

    return result, err
}
