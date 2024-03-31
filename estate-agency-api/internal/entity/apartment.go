package entity

type Apartment struct {
	ID         int    `form:"id" json:"id"`
	Title      string `form:"title" json:"title" binding:"max=100"`
	Price      int    `form:"price" json:"price"`
	City       string `form:"city" json:"city" binding:"max=20"`
	Rooms      int    `form:"rooms" json:"rooms"`
	Address    string `form:"address" json:"address" binding:"max=100"`
	Square     int    `form:"square" json:"square"`
	IDRealtor  int    `form:"id_realtor" json:"id_realtor"`
	UpdateTime string `json:"update_time"`
	CreateTime string `json:"create_time"`
}
