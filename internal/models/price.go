package models

import "gorm.io/gorm"

// Price struct corresponds to Django model called Price.
// We will add translation fields to match database schema.

type Price struct {
	gorm.Model

	ItemName string `gorm:"size:150;unique"`
	Price    float32
	// Translation for Polish
	ItemNamePL string `gorm:"size:150;column:item_name_pl"`
	// Translation for English
	ItemNameEN string `gorm:"size:150;column:item_name_en"`
	// Translation for Ukrainian
	ItemNameUK string `gorm:"size:150;column:item_name_uk"`
	Category   string `gorm:"size:2"`
}
