package models

import (
	"time"
	"gorm.io/gorm"
)

// News struct corresponds to News model in Django project
// Adding translation fields as well

type News struct {
	gorm.Model
	Title       string `gorm:"size:250;unique" form:"title"`
	Header      string `gorm:"type:text" form:"header"`
	Description string `gorm:"type:text" form:"description"`
	Features    string `gorm:"type:text" form:"features"`
	PostedOn    time.Time
	// Translation fields for Polish language
	TitlePL       string `gorm:"size:250;column:title_pl" form:"title_pl"`
	DescriptionPL string `gorm:"type:text;column:description_pl" form:"description_pl"`
	HeaderPL      string `gorm:"type:text;column:header_pl" form:"header_pl"`
	FeaturesPL    string `gorm:"type:text;column:features_pl" form:"features_pl"`
	//Translation fields for English language
	TitleEN       string `gorm:"size:250;column:title_en" form:"title_en"`
	DescriptionEN string `gorm:"type:text;column:description_en" form:"description_en"`
	HeaderEN      string `gorm:"type:text;column:header_en" form:"header_en"`
	FeaturesEN    string `gorm:"type:text;column:features_en" form:"features_en"`
	//Translation fields for Ukrainian language
	TitleUK       string `gorm:"size:250;column:title_uk" form:"title_uk"`
	DescriptionUK string `gorm:"type:text;column:description_uk" form:"description_uk"`
	HeaderUK      string `gorm:"type:text;column:header_uk" form:"header_uk"`
	FeaturesUK    string `gorm:"type:text;column:features_uk" form:"features_uk"`
	// Images URLs
	ImageLeft  string `gorm:"type:text"`
	ImageRight string `gorm:"type:text"`
}
