package buildindex

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/canghai118/posts-list/internal/model"
	"github.com/canghai118/posts-list/internal/service"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbConnection string = "root:demopassword@tcp(127.0.0.1:3306)/posts_list_demo?parseTime=true&loc=Local"
var redisConnection string = "localhost:6379"
var redisPageKey = "posts_page"
var concurrency = 1

var BuildIndexCmd = &cobra.Command{
	Use: "build-index",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := gorm.Open(mysql.Open(dbConnection))
		if err != nil {
			panic(err)
		}

		postSerive := service.NewPostService(db)

		rdb := redis.NewClient(&redis.Options{
			Addr:     redisConnection,
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		pageService := service.NewPageService(rdb, redisPageKey)

		startTime := time.Now()
		indexCount := uint64(0)

		ch := make(chan *model.Post, 100)
		go func() {
			err := postSerive.ScanAll(0, 100_000, ch)
			close(ch)
			if err != nil {
				log.Fatalf("Error scanning posts: %v\n", err)
			}
		}()

		ctx := context.Background()

		wg := sync.WaitGroup{}
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for post := range ch {
					err := pageService.UpdateIndex(ctx, post.Id, float64(post.LikeCount))
					if err != nil {
						log.Fatalf("Error updating page index: %v\n", err)
					}
					atomic.AddUint64(&indexCount, 1)
				}
			}()
		}

		wg.Wait()

		cost := time.Since(startTime)
		log.Printf("Index built, total: %d, cost: %v\n", indexCount, cost)
	},
}

func init() {
	BuildIndexCmd.Flags().StringVarP(&dbConnection, "db-connection", "d", dbConnection, "Database connection string")
	BuildIndexCmd.Flags().StringVarP(&redisConnection, "redis-connection", "r", redisConnection, "Redis connection string")
	BuildIndexCmd.Flags().StringVarP(&redisPageKey, "redis-page-key", "k", redisPageKey, "Redis key to store page index")
	BuildIndexCmd.Flags().IntVarP(&concurrency, "concurrency", "c", concurrency, "Number of concurrent workers")
}
