package flibusta

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/url"
	"reflect"
	"testing"
)

func Test_parseAuthPage(t *testing.T) {
	t.Parallel()
	content, err := ioutil.ReadFile("test-pages/guest-index.html")
	if err != nil {
		log.Fatal(err)
	}
	want := authParams{
		loginUrl: "/node?destination=node",
		data:     url.Values{},
	}
	want.data.Add("op", "Вход в систему")
	want.data.Add("form_id", "user_login_block")
	want.data.Add("form_build_id", "form-HI6XJ4kDqyAfitAPMwsbVb9k-_GsEfYe0HW7P3p5EZA")
	want.data.Add("openid.return_to", "http://flibusta.is/openid/authenticate?destination=node")
	t.Run("Getting auth params", func(t *testing.T) {
		got, err := parseAuthPage(bytes.NewReader(content))
		if err != nil {
			t.Errorf("parseAuthPage() error = %v", err)
			return
		}
		if !reflect.DeepEqual(got.loginUrl, want.loginUrl) {
			t.Errorf("parseAuthPage() got login url = %v, want %v", got, want)
		}
		if !reflect.DeepEqual(got.data, want.data) {
			t.Errorf("parseAuthPage() got data = %v, want %v", got, want)
		}
	})

}

func Test_parsePageContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filename string
		want     *Book
	}{
		{
			name:     "Parsing: Внетелесный опыт",
			filename: "test-pages/book-9.html",
			want: &Book{
				ID:        9,
				ReadCount: 704,
				Title:     "Внетелесный опыт",
				Authors: []Author{
					{
						ID:   24445,
						Name: "Аарон",
					},
					{
						ID:   32193,
						Name: "Сентхил Кумар",
					},
				},
				Annotation: `отсутствует`,
				Genres: []Genre{
					{
						ID:   97,
						Name: "Религия, религиозная литература",
					},
				},
			},
		},
		{
			name:     "Parsing: Игровой движок [Программирование и внутреннее устройство]",
			filename: "test-pages/book-611196.html",
			want: &Book{
				ID:        611196,
				ReadCount: 288,
				Title:     "Игровой движок [Программирование и внутреннее устройство]",
				Authors: []Author{
					{
						ID:   237578,
						Name: "Джейсон Грегори",
					},
				},
				Annotation: `<p>Книга Джейсона Грегори не случайно является бестселлером. Двадцать лет работы автора над первоклассными играми в Midway, Electronic Arts и Naughty Dog позволяют поделиться знаниями о теории и практике разработки ПО для игрового движка. Игровое программирование — сложная и огромная тема, охватывающая множество вопросов.<br />
                    Граница между игровым движком и игрой размыта. В этой книге основное внимание уделено движку, основным низкоуровневым системам, системам разрешения коллизий, симуляции физики, анимации персонажей, аудио, а также базовому слою геймплея, включающему объектную модель игры, редактор мира, системы событий и скриптинга</p>`,
				Genres: []Genre{
					{
						ID:   83,
						Name: "Зарубежная компьютерная, околокомпьютерная литература",
					},
					{
						ID:   81,
						Name: "Программирование, программы, базы данных",
					},
				},
			},
		},
		{
			name:     "Parsing: Design Driven Testing: Test Smarter, Not Harder",
			filename: "test-pages/book-235391.html",
			want: &Book{
				ID:        235391,
				ReadCount: 83,
				Title:     "Design Driven Testing: Test Smarter, Not Harder",
				Authors: []Author{
					{
						ID:   77447,
						Name: "Matt Stephens",
					},
					{
						ID:   77448,
						Name: "Doug Rosenberg",
					},
				},
				Annotation: `<p>Apress, 2010, 344 pp.<br />
                    ISBN-10:  	1430229438<br />
                    ISBN-13:  	9781430229438<br />
                    The groundbreaking book Design Driven Testing brings sanity back to the software development process by flipping around the concept of Test Driven Development (TDD)—restoring the concept of using testing to verify a design instead of pretending that unit tests are a replacement for design. Anyone who feels that TDD is “Too Damn Difficult” will appreciate this book.</p>
                <p>Design Driven Testing shows that, by combining a forward-thinking development process with cutting-edge automation, testing can be a finely targeted, business-driven, rewarding effort. In other words, you’ll learn how to test smarter, not harder.</p>
                <p>    Applies a feedback-driven approach to each stage of the project lifecycle.<br />
                    Illustrates a lightweight and effective approach using a core subset of UML.<br />
                    Follows a real-life example project using Java and Flex/ActionScript.<br />
                    Presents bonus chapters for advanced DDTers covering unit-test antipatterns (and their opposite, “test-conscious” design patterns), and showing how to create your own test transformation templates in Enterprise Architect.</p>
                <p>What you’ll learn</p>
                <p>    Create unit and behavioral tests using JUnit, NUnit, FlexUnit.<br />
                    Generate acceptance tests for all usage paths through use case thread expansion.<br />
                    Generate requirement tests for functional requirements.<br />
                    Run complex acceptance tests across the enterprise.<br />
                    Isolate individual control points for self-contained unit/behavioral tests.<br />
                    Apply behavior-driven development frameworks like JBehave and NBehave</p>
                <p>Who this book is for</p>
                <p>Design Driven Testing should appeal to developers, project managers, testers, business analysts, architects...in fact anyone who builds software that needs to be tested. While equally applicable on both large and small projects, Design Driven Testing is especially helpful to those developers who need to verify their software against formal requirements. Such developers will benefit greatly from the rational and disciplined approach espoused by the authors.</p>`,
				Genres: []Genre{
					{
						ID:   81,
						Name: "Программирование, программы, базы данных",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		content, err := ioutil.ReadFile(tt.filename)
		if err != nil {
			log.Fatal(err)
		}
		t.Run(tt.name, func(t *testing.T) {
			want := tt.want
			got, err := parsePageContent(string(content))
			if err != nil {
				t.Errorf("Parse() error = %v, wantErr %v", err, false)
				return
			}
			if !reflect.DeepEqual(got.Title, want.Title) {
				t.Errorf("Parse() got title = %v, want title %v", got.Title, want.Title)
			}
			if !reflect.DeepEqual(got.ID, want.ID) {
				t.Errorf("Parse() got ID = %v, want ID %v", got.ID, want.ID)
			}
			if !reflect.DeepEqual(got.ReadCount, want.ReadCount) {
				t.Errorf("Parse() got read count = %v, want read count %v", got.ReadCount, want.ReadCount)
			}
			if !reflect.DeepEqual(got.Authors, want.Authors) {
				t.Errorf("Parse() got authors = %v, want authors %v", got.Authors, want.Authors)
			}
			if !reflect.DeepEqual(got.Genres, want.Genres) {
				t.Errorf("Parse() got genres = %v, want genres %v", got.Genres, want.Genres)
			}
			if !reflect.DeepEqual(got.Annotation, want.Annotation) {
				t.Errorf("Parse() got annotation = %v, want annotation %v", got.Annotation, want.Annotation)
			}
		})
	}
}
