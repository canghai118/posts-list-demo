package serve

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/canghai118/posts-list/internal/service"
	"github.com/canghai118/posts-list/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var dbConnection string = "root:demopassword@tcp(127.0.0.1:3306)/posts_list_demo?parseTime=true&loc=Local"
var redisConnection string = "localhost:6379"
var redisPageKey = "posts_page"
var listen string = ":8080"

var ServeCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		gin.SetMode(gin.ReleaseMode)

		db := utils.MustGetGorm(dbConnection)

		rdb := utils.MustGetRedis(redisConnection)

		pageService := service.NewPageService(rdb, redisPageKey)
		postService := service.NewPostService(db)

		router := gin.Default()

		router.GET("/api/posts", func(g *gin.Context) {
			pageStr := g.Query("page")
			sizeStr := g.Query("size")

			page, _ := strconv.ParseInt(pageStr, 10, 64)
			size, _ := strconv.ParseInt(sizeStr, 10, 64)

			if page == 0 {
				page = 1
			}
			if size == 0 {
				size = 100
			}

			var startTime time.Time

			startTime = time.Now()
			ids, err := pageService.GetPage(g.Request.Context(), int(page), int(size))
			if err != nil {
				g.JSON(500, gin.H{"error": err.Error()})
				return
			}
			fmt.Printf("page: %d, size:%d, ids: %v\n", page, size, ids)
			total, err := pageService.GetTotoal(g.Request.Context())
			if err != nil {
				g.JSON(500, gin.H{"error": err.Error()})
				return
			}

			pageQueryCost := time.Since(startTime)

			startTime = time.Now()
			posts, err := postService.GetByIds(ids)
			if err != nil {
				g.JSON(500, gin.H{"error": err.Error()})
				return
			}
			postsQueryCost := time.Since(startTime)

			totalPage := total / size

			g.JSON(200, gin.H{
				"total":     total,
				"totalPage": totalPage,
				"page":      page,
				"size":      size,
				"lists":     posts,
				"peformance": gin.H{
					"pageQueryCost":  pageQueryCost.String(),
					"postsQueryCost": postsQueryCost.String(),
				},
			})
		})

		router.POST("/api/posts/addLike", func(g *gin.Context) {
			postIdStr := g.Query("id")
			postId, err := strconv.Atoi(postIdStr)
			if err != nil {
				g.JSON(400, gin.H{"error": "Invalid post id"})
				return
			}
			likeCountStr := g.Query("like")
			likeCount, _ := strconv.Atoi(likeCountStr)
			if likeCount == 0 {
				likeCount = 1
			}
			post, err := postService.Get(postId)
			if err != nil {
				g.JSON(500, gin.H{"error": err.Error()})
				return
			}
			if post == nil {
				g.JSON(404, gin.H{"error": "Post not found"})
				return
			}

			err = postService.AddLike(postId, likeCount)
			if err != nil {
				g.JSON(500, gin.H{"error": err.Error()})
				return
			}
			post, err = postService.Get(postId)
			if err != nil {
				g.JSON(500, gin.H{"error": err.Error()})
				return
			}
			if post == nil {
				g.JSON(404, gin.H{"error": "Post not found"})
				return
			}
			err = pageService.UpdateIndex(g.Request.Context(), post.Id, float64(post.LikeCount))
			if err != nil {
				g.JSON(500, gin.H{"error": fmt.Errorf("failed to update page index: %w", err).Error()})
				return
			}

			g.JSON(200, post)

		})

		log.Printf("Listening on %s\n", listen)
		if err := router.Run(listen); err != nil {
			panic(err)
		}
	},
}

func init() {
	ServeCmd.PersistentFlags().StringVar(&dbConnection, "db-connection", dbConnection, "Database connection string")
	ServeCmd.PersistentFlags().StringVar(&redisConnection, "redis-connection", redisConnection, "Redis connection string")
	ServeCmd.PersistentFlags().StringVar(&redisPageKey, "redis-page-key", redisPageKey, "Redis key for the page index")
	ServeCmd.PersistentFlags().StringVar(&listen, "listen", listen, "Address to listen on")
}
