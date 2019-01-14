package article

import (
	"testing"
)

func TestFetchArticles(t *testing.T) {
	reader := &ArticleReader{
		API: "https://jenkins-zh.github.io/index.json",
	}
	articles, err := reader.FetchArticles()
	if err != nil {
		t.Errorf("fetch error %v", err)
	} else if len(articles) == 0 {
		t.Errorf("fetch zero article")
	} else {
		for i, article := range articles {
			if article.Title == "" || article.Description == "" ||
				article.URI == "" {
				t.Errorf("article [%d] title, description or uri is empty", i)
			}
		}
	}

	ar, err := reader.FindByTitle("行为")
	if err != nil {
		t.Errorf("%v", err)
	}

	for _, a := range ar {
		t.Errorf("%v", a)
	}
}
