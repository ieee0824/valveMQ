package valve

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func DBInit(cfg *Config) error {
	d, err := connect(
		cfg.DbUser,
		cfg.DbPass,
		cfg.DbHost,
		cfg.DbPort,
		"",
	)
	if err != nil {
		return err
	}
	tx := d.MustBegin()

	if _, err := tx.Exec("CREATE DATABASE IF NOT EXISTS `mq` DEFAULT CHARACTER SET utf8mb4"); err != nil {
		return err
	}
	if _, err := tx.Exec("use mq"); err != nil {
		return err
	}
	if _, err := tx.Exec("CREATE TABLE IF NOT EXISTS `message` (`id` int(11) unsigned NOT NULL AUTO_INCREMENT,`body` text NOT NULL,`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,`expire` timestamp NULL DEFAULT NULL,`flag` int(11) NOT NULL DEFAULT '0',`hash` varchar(512) NOT NULL DEFAULT '',PRIMARY KEY (`id`),KEY `hash` (`hash`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;"); err != nil {
		return err
	}

	if _, err := tx.Exec("CREATE TABLE IF NOT EXISTS `log` (`id` int(11) unsigned NOT NULL AUTO_INCREMENT, `last_dequeue_time` timestamp(6) NOT NULL DEFAULT '1980-01-01 00:00:00.000000', `hash` varchar(512) NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;"); err != nil {
		return err
	}

	tx.Commit()
	d.Close()

	db, err = connect(
		cfg.DbUser,
		cfg.DbPass,
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbName,
	)
	if err != nil {
		return err
	}
	return nil
}

func connect(
	dbUser, dbPass, dbHost, dbPort, dbName string,
) (*sqlx.DB, error) {
	return sqlx.Connect(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true",
			dbUser,
			dbPass,
			dbHost,
			dbPort,
			dbName,
		),
	)
}
