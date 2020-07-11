package entity

type User struct {
	Email              string `bson:"email" json:"email"`
	Password           string `bson:"password" json:"password"`
	Name               string `bson:"name" json:"name"`
}

func (u *User) New() *User {
	return &User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}
