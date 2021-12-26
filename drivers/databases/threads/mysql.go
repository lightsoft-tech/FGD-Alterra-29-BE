package threads

import (
	"context"
	"fgd-alterra-29/business/threads"
	"fgd-alterra-29/drivers/databases/comments"
	"fmt"

	"gorm.io/gorm"
)

type MysqlThreadRepository struct {
	Conn *gorm.DB
}

func NewMysqlThreadRepository(conn *gorm.DB) threads.Repository {
	return &MysqlThreadRepository{
		Conn: conn,
	}
}

func (DB *MysqlThreadRepository) GetProfileThreads(ctx context.Context, id int) ([]threads.Domain, error) {
	var Thread []Threads

	SQReplier := DB.Conn.Table("comments").Select("name").Joins("join users on comments.user_id = users.id").Where("users.id = comments.user_id").Limit(1)
	Replier := DB.Conn.Table("comments").Select("(?) as name", SQReplier).Where("comments.thread_id = threads.id").Limit(1)

	Comment := DB.Conn.Table("comments").Select("comment").Where("comments.thread_id = threads.id").Limit(1)

	result := DB.Conn.Table("threads").Where("threads.user_id = 1").Select("title, (?) as Comment, (?) as RecentReplier", Comment, Replier).
		Find(&Thread)

	if result.Error != nil {
		return []threads.Domain{}, result.Error
	}

	return ToListDomain(Thread), nil
}

func (DB *MysqlThreadRepository) GetRecommendationThreads(ctx context.Context, id int) ([]threads.Domain, error) {
	var Thread []Threads
	var Comment []comments.Comments

	Q_Like := DB.Conn.Table("thread_likes").Where("thread_id = threads.id").Select("count(thread_likes.user_id)").Group("thread_id")
	Q_Comment := DB.Conn.Table("comments").Where("thread_id = threads.id").Select("count(comment)").Group("thread_id")

	result := DB.Conn.Table("threads").Select("*, title, content, (?) as Q_Like, (?) as Q_Comment", Q_Like, Q_Comment).Preload("Comments", func(db *gorm.DB) *gorm.DB {
		return db.Table("comments").Select("*").Preload("Replies").Find(&Comment)
	}).
		Find(&Thread)

	if result.Error != nil {
		return []threads.Domain{}, result.Error
	}

	fmt.Println("ThreadPreload", Thread)

	return ToListDomain(Thread), nil
}
