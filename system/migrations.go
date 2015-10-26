package system

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/Sirupsen/logrus"
	"github.com/rubenv/sql-migrate"
)

//migrateNew creates new blank migration
func migrateNew(box *rice.Box) {
	if len(flag.Args()) == 0 {
		logrus.Error("Migrations name not specified")
		return
	}
	name := path.Join(box.Name(), fmt.Sprintf("%d_%s.sql", time.Now().Unix(), flag.Arg(0)))
	file, err := os.Create(name)
	if err != nil {
		logrus.Error(err)
		return
	}
	fmt.Fprintf(file, "-- +migrate Up\n")
	fmt.Fprintf(file, "-- SQL in section 'Up' is executed when this migration is applied\n\n\n")
	fmt.Fprintf(file, "-- +migrate Down\n")
	fmt.Fprintf(file, "-- SQL in section 'Down' is executed when this migration is rolled back\n\n\n")
	err = file.Close()
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Infof("File %s has been successfully created\n", name)
	}
}

//migrateUp applies db migrations
func migrateUp(DB *sql.DB, box *rice.Box) {
	migrations := getRiceMigrations(box)
	n, err := migrate.Exec(DB, "postgres", migrations, migrate.Up)
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Infof("%d migrations applied", n)
	}
}

//migrateDown rolls back db migrations
func migrateDown(DB *sql.DB, box *rice.Box) {
	migrations := getRiceMigrations(box)
	n, err := migrate.Exec(DB, "postgres", migrations, migrate.Down)
	if err != nil {
		logrus.Error(err)
	} else {
		logrus.Infof("%d migrations rolled back", n)
	}
}

//getRiceMigrations builds migration source from go.rice storage
func getRiceMigrations(box *rice.Box) *migrate.MemoryMigrationSource {
	source := &migrate.MemoryMigrationSource{}
	fn := func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			migFile, err := box.Open(path)
			if err != nil {
				return err
			}
			mig, err := migrate.ParseMigration(path, migFile)
			migFile.Close()
			if err != nil {
				return err
			}
			source.Migrations = append(source.Migrations, mig)
		}
		return nil
	}
	err := box.Walk("", fn)
	if err != nil {
		logrus.Fatal(err)
		return nil
	}
	return source
}
