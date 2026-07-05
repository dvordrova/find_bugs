package rules

import "github.com/quasilyte/go-ruleguard/dsl"

func forceSQLCQueryLayer(m dsl.Matcher) {
	m.ImportAs("database/sql", "sql")

	m.Match(
		`$db.ExecContext($*_)`,
		`$db.QueryContext($*_)`,
		`$db.QueryRowContext($*_)`,
	).
		Where(
			m["db"].Type.Is(`*sql.DB`) &&
				!m.File().PkgPath.Matches(`/sqlc$`),
		).
		Report(`database/sql calls must go through generated sqlc packages`)

	m.Match(
		`$tx.ExecContext($*_)`,
		`$tx.QueryContext($*_)`,
		`$tx.QueryRowContext($*_)`,
	).
		Where(
			m["tx"].Type.Is(`*sql.Tx`) &&
				!m.File().PkgPath.Matches(`/sqlc$`),
		).
		Report(`database/sql calls must go through generated sqlc packages`)
}
