package server_model

// 基础结构体，包含 UserId 字段
type User struct {
	UserId string `header:"user_id" json:"user_id" form:"user_id" binding:"required,len=13" label:"User"`
}

// Base 结构体实现 SetUserId 方法
func (u *User) SetUserId(userId string) {
	u.UserId = userId
}

type PageReq struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

func (p *PageReq) CheckPage() {
	if p.PageSize < 1 || p.PageSize > 1000 {
		p.PageSize = 10
	}
	if p.Page < 1 {
		p.Page = 1
	}
}
