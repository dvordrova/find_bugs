package rules

import "github.com/quasilyte/go-ruleguard/dsl"

func noWallClockInDomain(m dsl.Matcher) {
	m.Match(`time.Now()`).
		Where(m.File().PkgPath.Matches(`/domain$`)).
		Report(`domain logic must not call time.Now directly`)
}
