package mysql

import (
	"time"
)

type BaseModel struct {
	ID        uint       `gorm:"primary_key" json:"id" title:"ID"`
	CreatedAt time.Time  `json:"-" title:"创建时间"`
	UpdatedAt time.Time  `json:"-" title:"更新时间"`
	DeletedAt *time.Time `json:"-" title:"删除时间"`
}
