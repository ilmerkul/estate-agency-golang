package entity

type Realtor struct {
	ID         int    `form:"id" json:"id"`
	FirstName  string `form:"first_name" json:"first_name" validate:"len=0|min=2,max=20"`
	LastName   string `form:"last_name" json:"last_name" validate:"len=0|min=2,max=20"`
	Phone      string `form:"phone" json:"phone" validate:"len=0|e164"`
	Email      string `form:"email" json:"email" validate:"len=0|email"`
	Rating     int    `form:"rating" json:"rating" validate:"gte=0,lte=50"`
	Experience int    `form:"experience" json:"experience" validate:"gte=0"`
}
