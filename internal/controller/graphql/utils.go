package graphql

import (
	"github.com/oustrix/ozon_journal/internal/controller/graphql/model"
	"github.com/oustrix/ozon_journal/internal/entity"
)

func postToGraphQL(post *entity.Post) *model.Post {
	comments := make([]*model.Comment, 0, len(post.Comments))
	for _, c := range post.Comments {
		comments = append(comments, commentToGraphQL(&c))
	}

	return &model.Post{
		ID:          post.ID,
		Title:       post.Title,
		Content:     post.Content,
		PublishedAt: post.PublishedAt,
		AuthorID:    post.AuthorID,
		Commentable: post.Commentable,
		Comments:    comments,
	}
}

func commentToGraphQL(comment *entity.Comment) *model.Comment {
	return &model.Comment{
		ID:              comment.ID,
		Content:         comment.Content,
		AuthorID:        comment.AuthorID,
		PostID:          comment.PostID,
		PublishedAt:     comment.PublishedAt,
		ParentCommentID: comment.ParentCommentID,
	}
}
