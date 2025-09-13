package models

import "gorm.io/gorm"

// Program struct corresponds to the Django project 'Program' model.
// We are manually adding the translation fields to match the database schema
// used by the django-modeltranslation library.

type Program struct {
	gorm.Model // This automatically includes ID, CreatedAt, UpdatedAt, DeletedAt

	// Original and defaulf language fields
	Title       string `gorm:"size:250;unique" form:"title"`
	Description string `gorm:"type:text" form:"description"`
	Results     string `gorm:"type:text" form:"results"`

	// Translation fields for Polish language
	TitlePL       string `gorm:"size:250;column:title_pl" form:"title_pl"`
	DescriptionPL string `gorm:"type:text;column:description_pl" form:"description_pl"`
	ResultsPL     string `gorm:"type:text;column:results_pl" form:"results_pl"`

	//Translation fields for English language
	TitleEN       string `gorm:"size:250;column:title_en" form:"title_en"`
	DescriptionEN string `gorm:"type:text;column:description_en" form:"description_en"`
	ResultsEN     string `gorm:"type:text;column:results_en" form:"results_en"`

	//Translation fields for Ukrainian language
	TitleUK       string `gorm:"size:250;column:title_uk" form:"title_uk"`
	DescriptionUK string `gorm:"type:text;column:description_uk" form:"description_uk"`
	ResultsUK     string `gorm:"type:text;column:results_uk" form:"results_uk"`

	Category string `gorm:"size:2" form:"category"`
}
