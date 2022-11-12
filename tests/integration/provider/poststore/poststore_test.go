package poststore

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/yerassyldanay/makala/pkg/convx"
	"github.com/yerassyldanay/makala/provider/poststore"

	"github.com/yerassyldanay/makala/datastore/postgres_store"
	"github.com/yerassyldanay/makala/pkg/configx"
	"github.com/yerassyldanay/makala/pkg/errorx"
)

func getLogger() *log.Logger {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile)
	return logger
}

func dbDropDatabase(t *testing.T, db *sql.DB) {
	_, err := db.Exec(`SELECT pg_terminate_backend(pid) 
	FROM pg_stat_activity 
	WHERE pid <> pg_backend_pid() AND datname = 'test_db';`)
	require.NoError(t, err, "failed to revoke database from public")
	_, err = db.Exec("drop database if exists test_db")
	require.NoError(t, err, "failed to empty test playground")
}

func dbCreateDatabase(t *testing.T, db *sql.DB) {
	_, err := db.Exec("create database test_db")
	require.NoError(t, err, "failed to empty test playground")
	_, err = db.Exec("GRANT CONNECT ON DATABASE test_db TO public")
	require.NoError(t, err, "failed to grant database to public")
}

func getDbConnection(t *testing.T) *sql.DB {
	conf, err := configx.NewConfiguration()
	require.NoErrorf(t, err, "failed to parse configurations")

	{
		// connect to datastore & prepare a playground for tests
		db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=%s",
			conf.PostgresHostname, conf.PostgresPort, conf.PostgresUsername, conf.PostgresPassword, conf.PostgresSSLMode))
		require.NoErrorf(t, err, "failed to establish connection with database")
		require.NoErrorf(t, db.Ping(), "failed to ping datastore")

		dbDropDatabase(t, db)
		dbCreateDatabase(t, db)
	}

	conf.PostgresDBName = "test_db"
	db, err := postgres_store.NewDB(*conf)
	require.NoErrorf(t, err, "failed to establish connection with database")

	dir, err := os.Getwd()
	require.NoErrorf(t, err, "failed to Getwd")

	filePath := filepath.Join(dir, "/../../../../datastore/postgres_store/migration")

	// migrate
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	errorx.PanicIfError(err)

	migrateInstance, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", filePath),
		"postgres", driver)
	errorx.PanicIfError(err)
	_ = migrateInstance.Up()

	return db
}

func TestPostStore(t *testing.T) {
	db := getDbConnection(t)
	defer func() {
		fmt.Printf("[DB] closing the connection with datastore. err: %v \n", db.Close())
	}()

	postQuerier := poststore.New(db)

	var postsCreated []poststore.FeedPost
	{
		for i := 1; i <= 10; i++ {
			postCreated, err := postQuerier.Create(context.Background(), poststore.CreateParams{
				Title:     fmt.Sprintf("title %d", i),
				Author:    "t2_author12",
				Link:      convx.StrToPtr("https://makala.com"),
				Submakala: "submakala",
				Score:     float64(i) + 0.5,
			})
			require.NoErrorf(t, err, "failed to create a post")

			postsCreated = append(postsCreated, postCreated)
		}
	}

	{
		var ids []int64
		for _, each := range postsCreated {
			ids = append(ids, each.ID)
		}
		postsFetched, err := postQuerier.GetByIds(context.Background(), ids)
		require.NoErrorf(t, err, "failed to fetch posts by id")

		postListContainer := poststore.PostSortContainer(postsFetched)
		sort.Sort(postListContainer)

		postsCreatedContainer := poststore.PostSortContainer(postsCreated)
		sort.Sort(postsCreatedContainer)

		require.Equal(t, postsCreated, postsFetched, "posts must be the same")

		for i, eachPost := range postsFetched {
			expectedPost := poststore.FeedPost{
				ID:        eachPost.ID,
				Title:     fmt.Sprintf("title %d", 10-i),
				Author:    "t2_author12",
				Link:      convx.StrToPtr("https://makala.com"),
				Submakala: "submakala",
				Score:     float64(10-i) + 0.5,
			}
			require.Equalf(t, expectedPost, eachPost, "posts are not the same")
		}
	}
}

func TestDatastoreConnection(t *testing.T) {
	db := getDbConnection(t)
	defer func() {
		fmt.Printf("[DB] closing the connection with datastore. err: %v \n", db.Close())
	}()
}
