// Update 修改
func (m *default{{.upperStartCamelObject}}Model) Update(ctx context.Context, data *{{.upperStartCamelObject}}, session ...sqlx.Session) (int64, error) {

	query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)

	var (
        result sql.Result
        err    error
    )

    if len(session) > 0 && session[0] != nil {
  		result, err = session[0].ExecCtx(ctx, query, {{.expressionValues}})
  	} else {
  		result, err = m.conn.ExecCtx(ctx, query, {{.expressionValues}})
  	}

  	if err != nil {
  		return 0, err
  	}
  	affected, _ := result.RowsAffected()

  	return affected, nil
}

func (m *default{{.upperStartCamelObject}}Model) ConditionUpdate(ctx context.Context, where squirrel.Sqlizer, value map[string]interface{}, session ...sqlx.Session) (int64, error) {

	query, args, _ := squirrel.Update(m.table).SetMap(value).Where(where).ToSql()

	var (
        result sql.Result
        err    error
    )
    if len(session) > 0 && session[0] != nil {
        result, err = session[0].ExecCtx(ctx, query, args...)
    } else {
        result, err = m.conn.ExecCtx(ctx, query, args...)
    }

    if err != nil {
        return 0, err
    }
    affected, _ := result.RowsAffected()

    return affected, nil
}