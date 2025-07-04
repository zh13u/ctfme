package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	/*
	User
	Tên trường	Kiểu dữ liệu	Ý nghĩa
	id		int		Khóa chính
	username	string		Tên đăng nhập
	email		string		Email
	password_hash	string		Mật khẩu đã mã hóa
	team_id		int/null	Tham chiếu team (nếu có)
	is_admin	bool		Quyền admin
	created_at	datetime	Ngày tạo
	*/

	Username	string	`gorm:"unique;not null"`
	Email		string	`gorm:"unique;not null"`
	Password	string	`gorm:"not null"`
	TeamID		*uint
	IsAdmin		bool	`gorm:"default:false"`
}