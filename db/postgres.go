package db

import (
	"feedbacks/models"
	"github.com/kr/pretty"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var db *gorm.DB

//setupPostgres initializes the database instance
func setupPostgres() {
	var err error
	db, err = gorm.Open(postgres.Open(models.Conf.DB), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		log.Println("db.SetupPostgres err:", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(100)
	//AutoMigration
	autoMigrate()
	pretty.Logln("Postgres DB successfully connected! ")
}

// closePostgres closes database connection (unnecessary)
func closePostgres() {
	sqlDB, err := db.DB()
	sqlDB.Close()
	if err != nil {
		pretty.Logln("Error on closing the DB: ", err)
	}
	log.Println("Postgres closed")
}

func GetPGSQL() *gorm.DB {
	return db
}

func autoMigrate() {
	for _, model := range []interface{}{
		//(*models.SdpClients)(nil),
	} {
		dbSilent := db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)})
		err := dbSilent.AutoMigrate(model)
		if err != nil {
			log.Fatalf("create model %s : %s", model, err)
		}
	}
}
