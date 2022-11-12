package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/yerassyldanay/makala/datastore/postgres_store"
	"github.com/yerassyldanay/makala/datastore/redis_store"
	"github.com/yerassyldanay/makala/pkg/configx"
	"github.com/yerassyldanay/makala/pkg/errorx"
	"github.com/yerassyldanay/makala/provider/adstore"
	"github.com/yerassyldanay/makala/provider/feedstore"
	"github.com/yerassyldanay/makala/provider/poststore"
	"github.com/yerassyldanay/makala/server/rest/handler"
	"github.com/yerassyldanay/makala/server/synchron"
	"github.com/yerassyldanay/makala/service/postfeed"
)

func main() {
	fmt.Println("[SERVICE] version v1.0.0...")

	// at this moment, we will just print
	logger := log.New(os.Stdout, "", log.LstdFlags)

	conf, err := configx.NewConfiguration()
	errorx.PanicIfError(err)

	db, err := postgres_store.NewDB(*conf)
	errorx.PanicIfError(err)
	defer func() {
		logger.Println("[DB] closing the connection with datastore. err: ", db.Close())
	}()
	errorx.PanicIfError(db.Ping())

	dir, err := os.Getwd()
	errorx.PanicIfError(err)

	filePath := filepath.Join(dir, "/datastore/postgres_store/migration")

	// migrate
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	errorx.PanicIfError(err)

	migrateInstance, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", filePath),
		"postgres", driver)
	errorx.PanicIfError(err)
	_ = migrateInstance.Up()

	inMemoryStorageFeed, err := redis_store.NewRedisConnection(redis_store.RedisConnectionParams{
		Host:          conf.RedisHost,
		Port:          conf.RedisPort,
		LogicDatabase: conf.RedisFeedLogicDatabase,
	})
	errorx.PanicIfError(err)
	defer inMemoryStorageFeed.Close()

	inMemoryStorageAds, err := redis_store.NewRedisConnection(redis_store.RedisConnectionParams{
		Host:          conf.RedisHost,
		Port:          conf.RedisPort,
		LogicDatabase: conf.RedisAdLogicDatabase,
	})
	errorx.PanicIfError(err)
	defer inMemoryStorageFeed.Close()

	// repo
	var repoPosts poststore.Querier = poststore.New(db)

	// in-memory datastore
	feedInMemoryStorer := feedstore.NewStorage(logger, inMemoryStorageFeed)
	adInMemoryStorer := adstore.NewStorage(logger, inMemoryStorageAds)

	// services
	servicePosts := postfeed.NewPostService(logger, repoPosts, feedInMemoryStorer, adInMemoryStorer)

	// sync job runner
	runner := synchron.NewSyncRunner(logger, repoPosts, feedInMemoryStorer)

	// REST API handler
	restHandler := handler.NewPostServer(servicePosts)
	restHandler.SetRouter()

	go func() {
		logger.Printf("[REST] stated REST service at %v ... \n", fmt.Sprintf(":%v", conf.ListenAddr))
		err := restHandler.Router.Run(fmt.Sprintf(":%v", conf.ListenAddr))
		errorx.PanicIfError(err)
	}()

	// this eliminates the case when there are posts in SQL datastore,
	// but cache is empty
	runner.UpdateFeed()

	var quit = make(chan struct{}, 1)
	go runner.RunSchedule(context.Background(), *conf, quit)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// send command to runner to wrap up its job
	quit <- struct{}{}
	// wait runner to finish its jobs
	<-quit
}
