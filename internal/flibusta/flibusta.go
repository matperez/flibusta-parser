package flibusta

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Book struct {
	ID         int
	Title      string
	ReadCount  int
	Authors    []Author
	Annotation string
	Genres     []Genre
}

type Author struct {
	ID   int
	Name string
}

type Genre struct {
	ID   int
	Name string
}

type Client interface {
	GetBook(int) (*Book, error)
	Auth(username, password string) error
}

type Flibusta struct {
	client *http.Client
}

//NewClient creates new http client
func NewClient() (*http.Client, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, errors.Wrap(err, "error creating the http client")
	}
	client := &http.Client{
		// Prevent redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}
	return client, nil
}

//NewFlibusta creates new flibusta client
func NewFlibusta() (Client, error) {
	client, err := NewClient()
	if err != nil {
		return nil, errors.Wrap(err, "error creating the flibusta client")
	}
	return &Flibusta{client: client}, nil
}

//Auth authorizes a client
func (f *Flibusta) Auth(username, password string) error {
	if len(username) == 0 || len(password) == 0 {
		return errors.New("the username and the password must be set")
	}
	res, err := f.client.Get("https://flibusta.is")
	if err != nil {
		return errors.Wrap(err, "error getting unauthorized page")
	}
	defer res.Body.Close()
	params, err := parseAuthPage(res.Body)
	if err != nil {
		return errors.Wrap(err, "error getting the auth params")
	}
	params.data.Add("name", username)
	params.data.Add("pass", password)
	req, err := http.NewRequest(http.MethodPost, "https://flibusta.is"+params.loginUrl, strings.NewReader(params.data.Encode()))
	if err != nil {
		return errors.Wrap(err, "error making preparing an auth request")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Origin", "https://flibusta.is")
	req.Header.Add("Referer", "https://flibusta.is/")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36")
	res, err = f.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "error making an auth request")
	}
	defer res.Body.Close()
	if res.StatusCode != 302 {
		return errors.Errorf("error doing an auth request: http status code is %d", res.StatusCode)
	}
	res, err = f.client.Get(res.Header.Get("Location"))
	if err != nil {
		return errors.Wrap(err, "error checking for authorization status")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.New("error reading authorized page response body")
	}
	if !bytes.Contains(body, []byte(params.data.Get("user"))) {
		return errors.New("error checking if auth was success")
	}
	return nil
}

func (f *Flibusta) GetBook(id int) (*Book, error) {
	resp, err := f.client.Get("https://flibusta.is/b/" + strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("error getting the book content: the request was redirected")
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parsePageContent(string(content))
}

//parsePageContent fetches the book info from a page content
func parsePageContent(content string) (*Book, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return nil, errors.Wrap(err, "error parsing the page content")
	}
	var page Book

	// получаем название книги
	match := regexp.MustCompile(`/>(.*?)<span style=size`).FindStringSubmatch(content)
	if match == nil {
		match = regexp.MustCompile(`>(.*?)<span style=size`).FindStringSubmatch(content)
		if match == nil {
			return nil, errors.New("error getting the book title")
		}
	}
	page.Title = strings.TrimSpace(match[1])

	// получаем ID книги
	watchLinkPattern := regexp.MustCompile(`polka/watch/add/(\d+)`)
	match = watchLinkPattern.FindStringSubmatch(content)
	if match == nil {
		return nil, errors.New("error getting the book ID")
	}
	page.ID, err = strconv.Atoi(match[1])
	if err != nil {
		return nil, errors.Wrap(err, "error converting the book ID to an int")
	}

	spaceAndLineEndPattern := regexp.MustCompile(`\s{2,}|\n`)

	match = regexp.MustCompile(`книга прочитана (\d+)`).FindStringSubmatch(content)
	if match != nil {
		page.ReadCount, _ = strconv.Atoi(match[1])
	}

	// получаем список авторов. ищем все ссылки на авторов после тега скрипт вначале страницы
	page.Authors = []Author{}
	doc.Find("script~a[href*='/a/']").Each(func(i int, selection *goquery.Selection) {
		var author Author
		author.Name = spaceAndLineEndPattern.ReplaceAllString(selection.Text(), " ")
		match = regexp.MustCompile(`/a/(\d+)`).FindStringSubmatch(selection.AttrOr("href", ""))
		if match == nil {
			return
		}
		author.ID, _ = strconv.Atoi(match[1])
		page.Authors = append(page.Authors, author)
	})

	// получаем аннотацию. ищем от заголовка до ближайшей ссылки, либо линии-разделителя
	annotationPattern := regexp.MustCompile(`(?s)<h2>Аннотация</h2>(.*?)(?:<hr/>|<a href)`)
	match = annotationPattern.FindStringSubmatch(content)
	if match == nil {
		return nil, errors.New("error getting the book annotation")
	}
	page.Annotation = strings.TrimSpace(match[1])
	// удаляем ведущие переводы строк, если они есть
	match = regexp.MustCompile(`(?s)(.*?)(?i:<br>)`).FindStringSubmatch(page.Annotation)
	if match != nil {
		page.Annotation = strings.TrimSpace(match[1])
	}

	// получаем список жанров
	doc.Find("a.genre").Each(func(i int, selection *goquery.Selection) {
		var genre Genre
		genre.Name = strings.TrimSpace(spaceAndLineEndPattern.ReplaceAllString(selection.Text(), " "))
		match = regexp.MustCompile(`/g/(\d+)`).FindStringSubmatch(selection.AttrOr("href", ""))
		if match == nil {
			return
		}
		genre.ID, _ = strconv.Atoi(match[1])
		page.Genres = append(page.Genres, genre)
	})

	// возвращаем результат
	return &page, nil
}

//authParams authorization page data
type authParams struct {
	loginUrl string
	data     url.Values
}

//parseAuthPage fetches the auth page params
func parseAuthPage(content io.Reader) (*authParams, error) {
	params := authParams{
		data: url.Values{},
	}
	doc, err := goquery.NewDocumentFromReader(content)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing auth page content")
	}
	form := doc.Find("#user-login-form")
	params.loginUrl = form.AttrOr("action", "")
	if params.loginUrl == "" {
		return nil, errors.New("error getting the form action")
	}
	params.data.Add("op", "Вход в систему")
	formId, ok := form.Find("input[name=form_id]").Attr("value")
	if !ok {
		return nil, errors.New("form_id not found")
	}
	params.data.Add("form_id", formId)
	formBuildID, ok := form.Find("input[name=form_build_id]").Attr("value")
	if !ok {
		return nil, errors.New("form_build_id not found")
	}
	params.data.Add("form_build_id", formBuildID)
	openIDReturnTo, ok := form.Find("input[name='openid.return_to']").Attr("value")
	if !ok {
		return nil, errors.New("openid.return_to not found")
	}
	params.data.Add("openid.return_to", openIDReturnTo)
	return &params, nil
}
