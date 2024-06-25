func (m *default{{.upperStartCamelObject}}Model) FindOne(ctx context.Context, id int64,session ...sqlx.Session) (*{{.upperStartCamelObject}}, error) {

	{{if .withCache}}{{.cacheKey}}

    var (
        resp {{.upperStartCamelObject}}
        err  error
    )
    query :=  fmt.Sprintf("select %s from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}} limit 1", {{.lowerStartCamelObject}}Rows, m.table)

    if len(session) > 0 && session[0] != nil {
        // 使用传入的会话执行查询
        err = session[0].QueryRowCtx(ctx, &resp, query, id)
    } else {
        err = m.conn.QueryRowCtx(ctx, &resp, query, id)
    }

	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, xerr.ResourceNotFound
	default:
		return nil, err
	}{{else}}query := fmt.Sprintf("select %s from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}} limit 1", {{.lowerStartCamelObject}}Rows, m.table)

	var resp {{.upperStartCamelObject}}

	err := m.conn.QueryRowCtx(ctx, &resp, query, {{.lowerStartCamelPrimaryKey}})

	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, xerr.ResourceNotFound
	default:
		return nil, err
	}{{end}}
}


func (m *default{{.upperStartCamelObject}}Model) ConditionFindOne(ctx context.Context, rows string, condition squirrel.Sqlizer) (*{{.upperStartCamelObject}}, error) {

	var resp {{.upperStartCamelObject}}

	query, args, _ := squirrel.Select(rows).From(m.table).Where(condition).Limit(1).ToSql()

	err := m.conn.QueryRowPartialCtx(ctx, &resp, query, args...)

	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, xerr.ResourceNotFound
	default:
		return nil, err
	}
}

func (m *default{{.upperStartCamelObject}}Model) Count(ctx context.Context, condition squirrel.Eq) (int64, error) {

	var (
		total int64
		query string
		args  []interface{}
	)
	if condition != nil {
		query, args, _ = squirrel.Select("COUNT(id)").From(m.table).Where(condition).ToSql()
	} else {
		query, args, _ = squirrel.Select("COUNT(id)").From(m.table).ToSql()
	}

	if err := m.conn.QueryRowCtx(ctx, &total, query, args...); err != nil {

		return 0, err
	}
	return total, nil
}

func (m *default{{.upperStartCamelObject}}Model) IsExist(ctx context.Context, condition squirrel.Sqlizer) (bool, error) {
	query, args, _ := squirrel.Select("count(id)").From(m.table).Where(condition).ToSql()
	var total int64

	if err := m.conn.QueryRowCtx(ctx, &total, query, args...); err != nil {
		return false, err
	}
	return total > 0, nil
}

// FindList 查询多条数据
// @params ctx 上下文
// @params rows 查询字段
// @params condition 查询条件
// @params sort 排序字段
func (m *default{{.upperStartCamelObject}}Model) FindList(ctx context.Context, rows string, condition map[string]interface{}, sort ...string) ([]*{{.upperStartCamelObject}}, error) {

	sql := m.rowBuilder(rows, sort, condition)
	query, args, _ := sql.ToSql()

	var resp []*{{.upperStartCamelObject}}
	err := m.conn.QueryRowsPartialCtx(ctx, &resp, query, args...)

	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, xerr.ResourceNotFound
	default:
		return nil, err
	}
}

// FindListWithPage 分页查询多条数据
// @params ctx 上下文
// @params rows 查询字段
// @params condition 查询条件
// @params sort 排序字段
// @params page 页码
// @params pageSize 每页数量
func (m *default{{.upperStartCamelObject}}Model) FindListWithPage(ctx context.Context, rows string, condition map[string]interface{}, page, pageSize uint64, sort ...string) (int64, []*{{.upperStartCamelObject}}, error) {

	query, args, _ := m.rowPageBuilder(rows, sort, condition, page, pageSize).ToSql()
	var resp []*{{.upperStartCamelObject}}
	if err := m.conn.QueryRowsPartialCtx(ctx, &resp, query, args...); err != nil {
		return 0, nil, err
	}
	total, err := m.Count(ctx, condition)
	if err != nil {
		return 0, nil, err
	}
	return total, resp, nil
}

// rowBuilder 构建查询语句
// @params rows 查询字段
// @params rows 排序字段
// @params condition 查询条件
func (m *default{{.upperStartCamelObject}}Model) rowBuilder(rows string, sort []string, condition map[string]interface{}) squirrel.SelectBuilder {
	selectRows, orderBy := m.handleParams(rows, sort)
	sqlBuilder := squirrel.Select(selectRows).From(m.table)
	if len(orderBy) > 0 {
		sqlBuilder = sqlBuilder.OrderBy(orderBy)
	}
	if len(condition) > 0 {
		sqlBuilder = sqlBuilder.Where(condition)
	}
	return sqlBuilder
}

// rowPageBuilder 构建分页查询语句
// @params rows 查询字段
// @params rows 排序字段
// @params condition 查询条件
// @params page 页码
// @params pageSize 每页数量
func (m *default{{.upperStartCamelObject}}Model) rowPageBuilder(rows string, sort []string, condition map[string]interface{}, page, pageSize uint64) squirrel.SelectBuilder {

	return m.rowBuilder(rows, sort, condition).Limit(pageSize).Offset((page - 1) * pageSize)
}

// 处理默认参数
func (m *default{{.upperStartCamelObject}}Model) handleDefaultParams(rows string, sort []string) (string, string) {
	if rows == "*" {
		rows = {{.lowerStartCamelObject}}Rows
	}
	sortStr := strings.Join(sort, ", ")
	return rows, sortStr
}