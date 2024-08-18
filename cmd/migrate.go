package cmd

import (
	// "database/sql"

	"database/sql"
	"errors"
	"fmt"
	"log"
	"snap_chat_server/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateDropCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateFreshCmd)
}

func getMigrationInstance() *migrate.Migrate {
	fmt.Println(config.Env.Db)
	dataSource := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.Env.Db.Username,
		config.Env.Db.Password,
		config.Env.Db.Host,
		config.Env.Db.Port,
		config.Env.Db.Name,
	)
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		log.Fatal(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///"+basepath+"/database/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	return m
}

var migrateCmd = &cobra.Command{
	Use:              "migrate [COMMANDS]",
	TraverseChildren: true,
	Short:            "Migrate database schema",
	Long:             `This command will execute all migration file in migrations folder if not already executed.`,
	Run: func(cmd *cobra.Command, args []string) {
		m := getMigrationInstance()
		err := m.Up()
		if errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}

	},
}

var migrateFreshCmd = &cobra.Command{
	Use:   "fresh",
	Short: "Migrate fresh database schema",
	Long:  `This command will drop all tables and execute all migration file in migrations folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		m := getMigrationInstance()
		err := m.Down()
		if err != nil {
			if !errors.Is(err, migrate.ErrNoChange) {
				log.Fatal(err)
			}
		}
		err = m.Up()
		if err != nil {
			if !errors.Is(err, migrate.ErrNoChange) {
				log.Fatal(err)
			}
		}
	},
}

var migrateDropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop database schema",
	Long:  `This command will drop all tables in database.`,
	Run: func(cmd *cobra.Command, args []string) {
		m := getMigrationInstance()
		err := m.Drop()
		if err != nil {
			log.Fatal(err)
		}
	},
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Migrate database schema",
	Long:  `This command will execute all migration file in migrations folder if not already executed.`,
	Run: func(cmd *cobra.Command, args []string) {
		m := getMigrationInstance()
		err := m.Up()
		if err != nil {
			log.Fatal(err)
		}
	},
}
