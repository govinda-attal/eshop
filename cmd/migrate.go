package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/spf13/cobra"
)

const (
	up   = "up"
	down = "down"
)

var migSrc string

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "reads migrations from sources and applies them to the database",
	Run:   migrateDB,
}

func migrateDB(cmd *cobra.Command, args []string) {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	action := strings.Join(args, "")
	srcURL := fmt.Sprintf("file://%s/%s/", dir, migSrc)
	if strings.HasPrefix(migSrc, "/") {
		srcURL = fmt.Sprintf("file://%s", migSrc)
	}

	m, err := migrate.New(srcURL, cfg.DB.MigrateUrl)
	if err != nil {
		log.Fatal(err)
	}

	switch action {
	case up:
		err = m.Up()
	case down:
		err = m.Down()
	default:
		log.Fatal("migration action not supported: ", action, "!")
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println("eshop migrate", action, "complete")
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.PersistentFlags().StringVar(&migSrc, "src", "scripts/db/migrations", "migrations source directory (migrations)")
}
