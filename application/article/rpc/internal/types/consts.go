package types

const (
	SortPublishTime = iota //0按时间发布时间排序
	SortLikeCount          //1按点赞数排序
)

const (
	DefaultPageSize = 20
	DefaultLimit    = 200

	DefaultSortLikeCursor = 1 << 30
)

const (
	// ArticleStatusPending 待审核
	ArticleStatusPending = iota
	// ArticleStatusNotPass 审核不通过
	ArticleStatusNotPass
	// ArticleStatusVisible 可见
	ArticleStatusVisible
	// ArticleStatusUserDelete 用户删除
	ArticleStatusUserDelete
)
