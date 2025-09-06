package models

import "gorm.io/gorm"

// Program struct corresponds to the Django project 'Program' model.
// We are manually adding the translation fields to match the database schema
// used by the django-modeltranslation library.

type Program struct {
	gorm.Model // This automatically includes ID, CreatedAt, UpdatedAt, DeletedAt

	// Original and defaulf language fields
	Title       string `gorm:"size:250;unique"`
	Description string `gorm:"type:text"`
	Results     string `gorm:"type:text"`

	// Translation fields for Polish language
	TitlePL       string `gorm:"size:250;column:title_pl"`
	DescriptionPL string `gorm:"type:text;column:description_pl"`
	ResultsPL     string `gorm:"type:text;column:results_pl"`

	//Translation fields for English language
	TitleEN       string `gorm:"size:250;column:title_en"`
	DescriptionEN string `gorm:"type:text;column:description_en"`
	ResultsEN     string `gorm:"type:text;column:results_en"`

	//Translation fields for Ukrainian language
	TitleUK       string `gorm:"size:250;column:title_uk"`
	DescriptionUK string `gorm:"type:text;column:description_uk"`
	ResultsUK     string `gorm:"type:text;column:results_uk"`

	Category string `gorm:"size:2"`
}

// Category constants for the Program model
const (
	Kosmetologia      = "KS"
	Laseroterapia     = "LS"
	Kosmetyka         = "KT"
	ZabiegiEstetyczne = "ZE"
	Trychologia       = "TR"
	Podologia         = "PD"
	Masaze            = "MS"
)
