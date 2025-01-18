package mysql

import (
	"web_app/models"
)

func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(post_id, author_id, community_id, title, content) values(?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlStr, p.ID, p.AuthorID, p.CommunityID, p.Title, p.Content)
	return
}

func GetPostByID(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select post_id, author_id, community_id, title, content, create_time from post where post_id = ?`
	err = db.Get(post, sqlStr, pid)
	return
}

func GetPostList(offset, limit int64) (post *[]models.Post, err error) {
	post = new([]models.Post)
	sqlStr := `select post_id, author_id, community_id, title, content, create_time from post limit ?,?`
	err = db.Select(post, sqlStr, (offset-1)*limit, limit)
	return
}
