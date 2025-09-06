package models

import (
	"gorm.io/gorm"
	"time"
)

// News struct corresponds to News model in Django project
// Adding translation fields as well

type News struct {
	gorm.Model
	Title       string `gorm:"size:250;unique"`
	Header      string `gorm:"type:text"`
	Description string `gorm:"type:text"`
	Features    string `gorm:"type:text"`
	PostedOn    time.Time
	// Translation fields for Polish language
	TitlePL       string `gorm:"size:250;column:title_pl"`
	DescriptionPL string `gorm:"type:text;column:description_pl"`
	HeaderPL      string `gorm:"type:text;column:header_pl"`
	FeaturesPL    string `gorm:"type:text;column:features_pl"`
	//Translation fields for English language
	TitleEN       string `gorm:"size:250;column:title_en"`
	DescriptionEN string `gorm:"type:text;column:description_en"`
	HeaderEN      string `gorm:"type:text;column:header_en"`
	FeaturesEN    string `gorm:"type:text;column:features_en"`
	//Translation fields for Ukrainian language
	TitleUK       string `gorm:"size:250;column:title_uk"`
	DescriptionUK string `gorm:"type:text;column:description_uk"`
	HeaderUK      string `gorm:"type:text;column:header_uk"`
	FeaturesUK    string `gorm:"type:text;column:features_uk"`
	// Images URLs and date when news were posted
	ImageLeft  string `gorm:"type:text"`
	ImageRight string `gorm:"type:text"`
}
