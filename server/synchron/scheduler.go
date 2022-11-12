package synchron

import (
	"context"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/yerassyldanay/makala/pkg/configx"
)

func (r SyncRun) RunSchedule(ctx context.Context, conf configx.Configuration, quit chan struct{}) {
	scheduler := gocron.NewScheduler(time.UTC)
	if scheduler == nil {
		log.Fatalln("scheduler is nil")
		return
	}
	defer func(sc *gocron.Scheduler) {
		r.Log.Println("[CRON] stopping scheduler...")
		sc.Stop()
	}(scheduler)

	// create a new feed and update old one
	jobPrepareFeed, err := scheduler.Cron(conf.CronSyncFeed).Do(func() {
		r.Log.Printf("[CRON] [%s] stated creating feed in cache... \n", time.Now().Format(conf.TimeFormat))
		r.UpdateFeed()
		r.Log.Printf("[CRON] [%s] finished creating feed in cache... \n", time.Now().Format(conf.TimeFormat))
	})

	if err != nil {
		log.Fatalln("failed to run job. err: ", err)
	}

	scheduler.StartImmediately()
	scheduler.StartAsync()

	r.Log.Printf("[CRON] next run is scheduled to %s \n", jobPrepareFeed.NextRun().Format(conf.TimeFormat))

	<-quit
	quit <- struct{}{}
}
