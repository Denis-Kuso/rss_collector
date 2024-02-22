package fetch 


type Post struct {
	URL string
	Name string
}

func GetPosts() ([]Post);
