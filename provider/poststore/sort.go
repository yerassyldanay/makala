package poststore

type PostSortContainer []FeedPost

func (a PostSortContainer) Len() int           { return len(a) }
func (a PostSortContainer) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PostSortContainer) Less(i, j int) bool { return a[i].Score > a[j].Score }
