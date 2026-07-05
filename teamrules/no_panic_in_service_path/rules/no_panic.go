package rules

import "github.com/quasilyte/go-ruleguard/dsl"

func noPanicInServicePath(m dsl.Matcher) {
	m.Match(`panic($*_)`).
		Where(m.File().PkgPath.Matches(`/service$`)).
		Report(`service code must return errors instead of panicking`)
}
