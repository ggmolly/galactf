package orm

import "time"

type Attachment struct {
	ID          uint64    `json:"id" gorm:"primaryKey"`
	ChallengeID uint64    `json:"-" gorm:"index"`
	Type        string    `json:"type" gorm:"type:varchar(32)"` // url, file
	URL         string    `json:"url" gorm:"type:varchar(255)"`
	Title       string    `json:"filename" gorm:"type:varchar(255)"`
	Size        uint64    `json:"size" gorm:"type:bigint"`
	CreatedAt   time.Time `json:"-" gorm:"autoCreateTime"`

	Challenge Challenge `json:"-" gorm:"foreignKey:ChallengeID"`
}
