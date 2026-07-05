package projection

type AccountOpened struct {
	AccountID string
}

type CreditReserved struct {
	AccountID string
	Cents     int
}

type Projector struct {
	accounts map[string]bool
	reserved map[string]int
}

func New() *Projector {
	return &Projector{
		accounts: make(map[string]bool),
		reserved: make(map[string]int),
	}
}

func (p *Projector) Apply(event any) {
	switch event := event.(type) {
	case AccountOpened:
		p.accounts[event.AccountID] = true
	case CreditReserved:
		if !p.accounts[event.AccountID] {
			return
		}
		p.reserved[event.AccountID] += event.Cents
	}
}

func (p *Projector) ReservedCents(accountID string) int {
	return p.reserved[accountID]
}
