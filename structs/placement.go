package structs

type Placement struct {
	User   *User
	Points int
}

func NewPlacement(user *User, points int) (Placement, error) {
	return Placement{user, points}, nil
}
