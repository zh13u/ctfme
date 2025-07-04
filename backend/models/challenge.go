package models

import "gorm.io/gorm"

type Challenge struct {
	gorm.Model

	/*
	Challenge
	Tên trường	Kiểu dữ liệu	Ý nghĩa
	id		int		Khóa chính
	title		string		Tên challenge
	description	text		Mô tả challenge
	category	string		Thể loại (web, pwn, ...)
	points		int		Số điểm
	flag		string		Flag đúng
	file_url	string		Link file (nếu có)
	visible		bool		Hiển thị/ẩn
	created_at	datetime	Ngày tạo
	*/

	Title       string `gorm:"not null"`
	Description string `gorm:"type:text"`
	Category    string `gorm:"not null"`
	Points      int    `gorm:"not null"`
	Flag        string `gorm:"not null"`
	FileURL     string
	Visible     bool   `gorm:"default:true"`
	SolvedBy    []User `gorm:"many2many:challenge_solvers;"`
}
