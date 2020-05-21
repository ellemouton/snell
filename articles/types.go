package articles

type Info struct {
	ID          int64
	Name        string
	Description string
	Price       int64
	ContentID   int64
}

type Content struct {
	ID   int64
	Text string
}
