package fetch

import "fmt"

type Post struct {
	URL  string
	Name string
}

func GetPosts() []Post {

	fmt.Println("getting posts")
	return []Post{}
}
