package models

import "gorm.io/gorm"

type Book struct{ 
	ID			uint 		`gorm:"primary key;autoIncrement" json:"id"`
	Author		*string		`json:"author"`
	Title		*string		`json:"title"`
	Publisher	*string		`json:"publisher"`
}

func MigrateBooks(db *gorm.DB) error {
	if err := db.AutoMigrate(&Book{}); err != nil {
		return err
	}
	return nil
}
