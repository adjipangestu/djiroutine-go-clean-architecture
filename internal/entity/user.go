package entity

type User struct {
	ID        int     `gorm:"primaryKey;column:id"`
	Username  string  `gorm:"column:username"`
	Email     string  `gorm:"column:email"`
	FirstName *string `gorm:"column:first_name"`
	LastName  *string `gorm:"column:last_name"`
}

type UserResponse struct {
	ID        int     `gorm:"primaryKey;column:id" json:"id"`
	Username  string  `gorm:"column:username" json:"username"`
	Email     string  `gorm:"column:email" json:"email"`
	FirstName *string `gorm:"column:first_name" json:"first_name"`
	Lastname  *string `gorm:"column:last_name" json:"last_name"`
}

func (User) TableName() string {
	return "auth_user"
}
