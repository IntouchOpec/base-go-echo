package model

import "github.com/hb-go/gorm"

// TemplateSocial
type TemplateSocial struct {
	gorm.Model
	Souce                 string         `json:"souce"`
	Type                  uint           `json:"type"`
	ChatChannelID         uint           `json:"chat_channel_id"`
	ChatChannel           ChatChannel    `json:"chat_channel" gorm:"foreignkey:ChatChannelID;"`
	KeyTemplates          []*KeyTemplate `gorm:"many2many:template_social_key_template;" json:"key_templates"`
	TemplateSocialDetails []*TemplateSocialDetail
}

// TemplateSocialDetail
type TemplateSocialDetail struct {
	gorm.Model
	No           int            `json:"no"`
	Souce        string         `json:"souce"`
	ParentID     uint           `json:"parent_id"`
	KeyTemplates []*KeyTemplate `gorm:"many2many:template_social_detail_key_template;" json:"key_templates"`
}

// KeyTemplate
type KeyTemplate struct {
	gorm.Model
	Name                  string                  `json:"name"`
	No                    int                     `json:"no"`
	TemplateSocials       []*TemplateSocial       `gorm:"many2many:template_social_key_template;" json:"template_socials"`
	TemplateSocialDetails []*TemplateSocialDetail `gorm:"many2many:template_social_detail_key_template;" json:"template_social_details"`
}
