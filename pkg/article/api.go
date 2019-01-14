package article

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type ArticleReader struct {
	API string
}

func (a *ArticleReader) FetchArticles() (articles []Article, err error) {
	var apiURL *url.URL

	apiURL, err = url.Parse(a.API)
	if err != nil {
		return
	}

	var resp *http.Response
	resp, err = http.Get(apiURL.String())
	if err != nil {
		return
	}

	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var allArticles []Article
	err = json.Unmarshal(data, &allArticles)

	for _, article := range allArticles {
		if article.Title == "" || article.Description == "" ||
			article.URI == "" {
			continue
		}
		articles = append(articles, article)
	}
	return
}

func (a *ArticleReader) FindByTitle(title string) (articles []Article, err error) {
	var allArticles []Article

	allArticles, err = a.FetchArticles()
	if err != nil {
		return
	}

	for _, article := range allArticles {
		if strings.Contains(article.Title, title) {
			articles = append(articles, article)
		}
	}
	return
}
