package valve

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	// use mysql driver
	_ "github.com/go-sql-driver/mysql"
)

var db *sqlx.DB

func init() {
	d, err := connect(
		"root",
		"",
		"127.0.0.1",
		"3306",
		"",
	)
	if err != nil {
		log.Fatalln(err)
	}
	tx := d.MustBegin()

	if _, err := tx.Exec("CREATE DATABASE IF NOT EXISTS `mq` DEFAULT CHARACTER SET utf8mb4"); err != nil {
		log.Fatalln(err)
	}
	if _, err := tx.Exec("use mq"); err != nil {
		log.Fatalln(err)
	}
	if _, err := tx.Exec("CREATE TABLE IF NOT EXISTS `message` (`id` int(11) unsigned NOT NULL AUTO_INCREMENT,`body` text NOT NULL,`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,`expire` timestamp NULL DEFAULT NULL,`flag` int(11) NOT NULL DEFAULT '0',`hash` varchar(512) NOT NULL DEFAULT '',PRIMARY KEY (`id`),KEY `hash` (`hash`)) ENGINE=InnoDB AUTO_INCREMENT=6145 DEFAULT CHARSET=utf8mb4;"); err != nil {
		log.Fatalln(err)
	}

	tx.Commit()
	d.Close()

	db, _ = connect(
		"root",
		"",
		"127.0.0.1",
		"3306",
		"mq",
	)
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
