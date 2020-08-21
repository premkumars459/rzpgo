package ut

type Urltable struct {
	Id               int    `gorm:"AUTO_INCREMENT"`
	Url              string `gorm:"unique;not null"`
	CrawlTimeout     int    `gorm:"not null"`
	Frequency        int    `gorm:"not null"`
	FailureThreshold int    `gorm:"not null"`
	Status           string `gorm:"not null"`
	FailureCount     int    `gorm:"not null"`
}
