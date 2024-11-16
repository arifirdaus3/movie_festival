package model

type UserMovieVote struct {
	UserID  uint
	User    User `gorm:"foreignKey:UserID;references:ID;primaryKey"`
	MovieID uint
	Movie   Movie `gorm:"foreignKey:MovieID;references:ID;primaryKey"`
	Type    bool  // true upvote, false downvote
}

type UserVote struct {
	MovieID uint `json:"movieID"`
	Type    bool `json:"type"` // true upvote, false downvote
}

type RequestMovieVote struct {
	UserID  uint `json:"userID"`
	MovieID uint `json:"movieID"`
	Type    bool `json:"type"` // true upvote, false downvote
}

type VotedMovieCount struct {
	MovieID uint  `json:"movieID"`
	Count   int64 `json:"count"`
}
