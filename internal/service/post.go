package service

import (
	"log"

	"github.com/canghai118/posts-list/internal/model"
	"gorm.io/gorm"
)

func NewPostService(db *gorm.DB) *Post {
	return &Post{db: db}
}

type Post struct {
	db *gorm.DB
}

func (p *Post) Get(id int) (*model.Post, error) {
	var post model.Post
	err := p.db.First(&post, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &post, nil
}

func (p *Post) AddLike(postId int, addCount int) error {
	// update posts set like_count = like_count + ? where id = ?
	return p.db.Model(&model.Post{}).Where("id = ?", postId).Update("like_count", gorm.Expr("like_count + ?", addCount)).Error
}

func (p *Post) ScanAll(from int, batchSize int, ch chan *model.Post) error {
	var postId int = 0
	var posts []*model.Post
	for {
		err := p.db.Where("id > ?", postId).Limit(batchSize).Order("id ASC").Find(&posts).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil
			}

			return err
		}
		if len(posts) == 0 {
			return nil
		}
		lastPost := posts[len(posts)-1]
		postId = lastPost.Id

		log.Printf("Scanned posts up to %d\n", postId)

		for _, post := range posts {
			ch <- post
		}

		if len(posts) < batchSize {
			return nil
		}
	}
}

func (p *Post) GetByIds(ids []int) ([]*model.Post, error) {
	var posts []*model.Post
	err := p.db.Where("id IN ?", ids).Order("like_count desc").Find(&posts).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return posts, nil
}

func (p *Post) AddBatch(posts []*model.Post) error {
	return p.db.CreateInBatches(posts, len(posts)).Error
}
