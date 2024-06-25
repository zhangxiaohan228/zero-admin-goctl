func (m *default{{.upperStartCamelObject}}Model) Delete(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}},session ...sqlx.Session) (int64, error) {

    query := fmt.Sprintf("delete from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table)

    var (
        result sql.Result
        err    error
    )

    if len(session) > 0 && session[0] != nil {
    		// 使用传入的会话执行查询
    		result, err = session[0].ExecCtx(ctx, query, id)
    	} else {
    		result, err = m.conn.ExecCtx(ctx, query, id)
    	}

    if err != nil {
        return 0, err
    }

    affected, _ := result.RowsAffected()

    return affected, nil
}

func (m *default{{.upperStartCamelObject}}Model) BatchDelete(ctx context.Context, id []int64, session ...sqlx.Session) (int64, error) {

	query, args, _ := squirrel.Delete(m.table).Where(squirrel.Eq{"id": id}).ToSql()

    var (
        result sql.Result
        err    error
    )
    if len(session) > 0 && session[0] != nil {
        // 使用传入的会话执行查询
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

func (m *default{{.upperStartCamelObject}}Model) ConditionDelete(ctx context.Context, where squirrel.Sqlizer, session ...sqlx.Session) (int64, error) {

	query, args, _ := squirrel.Delete(m.table).Where(where).ToSql()

    var (
        result sql.Result
        err    error
    )
    if len(session) > 0 && session[0] != nil {
        // 使用传入的会话执行查询
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