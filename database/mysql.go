package database

import (
	"github.com/sirupsen/logrus"
	"github.com/synchthia/nebula-api/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Mysql struct {
	client   *gorm.DB
	database string
}

func NewMysqlClient(mysqlConStr, database string) *Mysql {
	logrus.WithField("connection", mysqlConStr).Infof("[MySQL] Connecting to MySQL...")

	client, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       mysqlConStr,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: logger.NewGorm(),
	})
	if err != nil {
		logrus.Fatalf("[MySQL] Failed to connect to MySQL: %s", err)
		return nil
	}

	m := &Mysql{
		client:   client,
		database: database,
	}

	if err := m.client.AutoMigrate(&Servers{}); err != nil {
		logrus.Fatalf("[MySQL] Failed to migrate: %s", err)
		return nil
	}

	if err := m.client.AutoMigrate(&Bungee{}); err != nil {
		logrus.Fatalf("[MySQL] Failed to migrate: %s", err)
		return nil
	}

	if err := m.client.AutoMigrate(&Players{}); err != nil {
		logrus.Fatalf("[MySQL] Failed to migrate: %s", err)
		return nil
	}

	if err := m.InitBungeeTable(); err != nil {
		return nil
	}

	logrus.Infof("[MySQL] Connected to MySQL")

	return m
}
