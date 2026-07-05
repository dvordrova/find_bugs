package rules

import "github.com/quasilyte/go-ruleguard/dsl"

func transactionBoundary(m dsl.Matcher) {
	m.ImportAs("database/sql", "sql")

	m.Match(`$db.BeginTx($*_)`).
		Where(
			m["db"].Type.Is(`*sql.DB`) &&
				!m.File().PkgPath.Matches(`/transaction$`),
		).
		Report(`transactions belong in transaction manager packages`)

	m.Match(
		`$tx.Commit()`,
		`$tx.Rollback()`,
	).
		Where(
			m["tx"].Type.Is(`*sql.Tx`) &&
				!m.File().PkgPath.Matches(`/transaction$`),
		).
		Report(`transactions belong in transaction manager packages`)
}
