package models

// 定义请求的参数结构体

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 投票数据
type ParamVoteData struct {
	PostID    int64 `json:"post_id,string" binding:"required"`                 // 帖子ID
	Direction int   `json:"direction,string" binding:"required,oneof=1 0 -1" ` // 赞成票(1)还是反对票(-1)取消投票(0)
}
