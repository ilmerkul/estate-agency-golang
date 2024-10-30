package dto

type RealtorView struct {
	FirstName string `form:"first_name" json:"first_name"`
	LastName  string `form:"last_name" json:"last_name"`
	Phone     string `form:"phone" json:"phone"`
	Email     string `form:"email" json:"email"`
	Rating    int    `form:"rating" json:"rating"`
}
