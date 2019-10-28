package model

import (
	"fmt"

	"github.com/IntouchOpec/base-go-echo/model/orm"
	"github.com/jinzhu/gorm"
)

// TemplateSocial
type TemplateSocial struct {
	orm.ModelBase
	Souce                 string         `json:"souce"`
	Type                  uint           `json:"type"`
	ChatChannelID         uint           `json:"chat_channel_id"`
	ChatChannel           ChatChannel    `json:"chat_channel" gorm:"ForeignKey:ChatChannelID;"`
	KeyTemplates          []*KeyTemplate `gorm:"many2many:template_social_key_template;" json:"key_templates"`
	TemplateSocialDetails []*TemplateSocialDetail
}

// TemplateSocialDetail
type TemplateSocialDetail struct {
	orm.ModelBase
	No               int            `json:"no"`
	Souce            string         `json:"souce"`
	BelongsToID      uint           `json:"belongs_id"`
	TemplateSocialID uint           `json:"template_social_id"`
	KeyTemplates     []*KeyTemplate `gorm:"many2many:template_social_detail_key_template;" json:"key_templates"`
	BelongsTo        []*TemplateSocialDetail
	TemplateSocial   TemplateSocial
}

// KeyTemplate
type KeyTemplate struct {
	orm.ModelBase
	Name                  string                  `json:"name"`
	No                    int                     `json:"no"`
	TemplateSocials       []*TemplateSocial       `gorm:"many2many:template_social_key_template;" json:"template_socials"`
	TemplateSocialDetails []*TemplateSocialDetail `gorm:"many2many:template_social_detail_key_template;" json:"template_social_details"`
}

func (templateSoc *TemplateSocial) CreateTample() *TemplateSocial {
	return templateSoc
}

// GetTemplateSocial
func (templateSoc *TemplateSocial) GetTemplateSocial() string {
	if err := DB().Preload("KeyTemplates").Preload("TemplateSocialDetails.KeyTemplates", func(db *gorm.DB) *gorm.DB {
		return db.Order("template_social_detail.no asc")
	}).Find(&templateSoc).Error; err != nil {
		fmt.Println(err)
		return ""
	}
	tempDetails := []*TemplateSocialDetail{}
	tempDs := []*TemplateSocialDetail{}
	var souce map[string]string
	souce = make(map[string]string)
	var arr []string
	subTempDs := []*TemplateSocialDetail{}

	// var parentIDs string
	// var flexMessage string
	arr = append(arr, templateSoc.Souce)
	// flexMessage := fmt.Sprintf(templateSoc.Souce,)

	var make bool = true
	for index := 0; index < len(templateSoc.TemplateSocialDetails); index++ {
		make = true
		DB().Preload("KeyTemplates").Preload("TemplateSocialDetails", "template_social_id = ?", templateSoc.TemplateSocialDetails[index].ID).Find(&tempDetails)
		arr = append(arr, templateSoc.TemplateSocialDetails[index].Souce)

		for i := 0; i < len(tempDetails); i++ {
			DB().Preload("KeyTemplates").Preload("BelongsTo", "belongs_id = ?", tempDetails[i].ID, func(db *gorm.DB) *gorm.DB {
				return db.Order("template_social_detail.no asc")
			}).Find(&tempDs)
			arr = append(arr, tempDetails[i].Souce)

			for make {

				DB().Preload("KeyTemplates").Preload("BelongsTo", "belongs_id = ?", tempDetails[i].ID, func(db *gorm.DB) *gorm.DB {
					return db.Order("template_social_detail.no asc")
				}).Find(&subTempDs)
				// arr = append(arr, subTempDs[x].Souce)

				for x := 0; x < len(subTempDs); x++ {
					arr = append(arr, subTempDs[x].Souce)
					DB().Preload("KeyTemplates").Preload("BelongsTo", "belongs_id = ?", tempDetails[i].ID, func(db *gorm.DB) *gorm.DB {
						return db.Order("template_social_detail.no asc")
					}).Find(&subTempDs)
				}

				if len(subTempDs) == 0 {
					make = false
				}
			}
			souce["souce"] = tempDetails[i].Souce

		}
	}
	fmt.Println(souce)
	return ""
}
