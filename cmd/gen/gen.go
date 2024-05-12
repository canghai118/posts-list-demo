package gen

import (
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/canghai118/posts-list/internal/model"
	"github.com/canghai118/posts-list/internal/service"
	"github.com/canghai118/posts-list/internal/utils"
	"github.com/spf13/cobra"
)

var dbConnection string = "root:demopassword@tcp(127.0.0.1:3306)/posts_list_demo?parseTime=true&loc=Local"
var batchSize uint64 = 1000
var count = 1000
var cocurrent = 1

var GenCmd = &cobra.Command{
	Use: "gen",
	Run: func(cmd *cobra.Command, args []string) {
		db := utils.MustGetGorm(dbConnection)
		postService := service.NewPostService(db)

		remain := uint64(count)

		inserted := uint64(0)

		wg := sync.WaitGroup{}
		for i := 0; i < cocurrent; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for {

					currentRemain := atomic.LoadUint64(&remain)
					if currentRemain <= 0 {
						return
					}
					var count uint64 = batchSize
					if count > currentRemain {
						count = currentRemain
					}

					if !atomic.CompareAndSwapUint64(&remain, currentRemain, currentRemain-count) {
						continue
					}
					var posts []*model.Post
					for i := uint64(0); i < count; i++ {
						posts = append(posts, generateRandomPost())
					}

					err := postService.AddBatch(posts)
					if err != nil {
						panic(err)
					}
					atomic.AddUint64(&inserted, count)
					log.Printf("Total inserted: %d\n", atomic.LoadUint64(&inserted))
				}
			}()
		}
		wg.Wait()
	},
}

var baseTime = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

func generateRandomPost() *model.Post {
	randomPublishTime := baseTime.Add(time.Duration(rand.Int63n(365*24*3600)) * time.Second)

	postRandom := generateRandomString(8)

	return &model.Post{
		Title:       "Random Post " + postRandom,
		Content:     "Random Post Content " + postRandom,
		LikeCount:   normalDistRandom(10_000, 5_000),
		PublishTime: randomPublishTime,
	}
}

func generateRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func normalDistRandom(mean float64, stdDev float64) int {
	randomFloat := rand.NormFloat64()*stdDev + mean
	randomInt := int(randomFloat)
	if randomInt < 0 {
		randomInt = -randomInt // 取绝对值确保是正数
	}
	return randomInt
}

func init() {
	GenCmd.PersistentFlags().StringVar(&dbConnection, "db-connection", dbConnection, "Database connection string")
	GenCmd.PersistentFlags().Uint64Var(&batchSize, "batch-size", batchSize, "Batch size for writen operation")
	GenCmd.PersistentFlags().IntVar(&count, "count", count, "Number of posts to generate")
	GenCmd.PersistentFlags().IntVar(&cocurrent, "cocurrent", cocurrent, "Number of concurrent workers")
}
