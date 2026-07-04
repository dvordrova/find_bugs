package profile

type Profile struct {
	Email string
	Name  string
}

type Repository struct {
	byEmail map[string]*Profile
}

func NewRepository() *Repository {
	return &Repository{
		byEmail: map[string]*Profile{
			"admin@example.com": {
				Email: "admin@example.com",
				Name:  "Admin",
			},
		},
	}
}

func (r *Repository) FindByEmail(email string) (*Profile, error) {
	if p, ok := r.byEmail[email]; ok {
		return p, nil
	}

	return nil, nil
}
