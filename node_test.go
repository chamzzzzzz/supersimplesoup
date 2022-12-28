package supersimplesoup

import (
	"strings"
	"testing"
)

const testHTML = `
<html>
	<head>
		<title>supersimplesoup</title>
	</head>
	<body>
		<div>
			<ul>
				<li>
					<a id="1st id" href="1st href" title="1st title" class="link">1st text</a>
				</li>
				<li>
					<a id="2nd id" href="2nd href" title="2nd title" class="link">2nd text</a>
				</li>
				<li>
					<a id="3rd id" href="3rd href" title="3rd title" class="link">3rd text</a>
				</li>
				<li>
					<a id="4th id" href="4th href" title="4th title" class="link">4th text</a>
				</li>
				<li>
					<a id="5th id" href="5th href" title="5th title" class="link">5th text</a>
				</li>
			</ul>
		</div>
	</body>
</html>
`

func TestWalk(t *testing.T) {
	n, err := Parse(strings.NewReader(testHTML))
	if err != nil {
		t.Fatal(err)
	}
	err = Walk(n, func(node *Node) error {
		t.Logf("%d - %q", node.Type, node.Data)
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func TestFind(t *testing.T) {
	type expected struct {
		id       string
		class    string
		title    string
		href     string
		text     string
		fulltext string
	}
	type testcase struct {
		tag      string
		attrkv   []string
		expected expected
	}

	n, err := Parse(strings.NewReader(testHTML))
	if err != nil {
		t.Fatal(err)
	}

	testcases := []testcase{
		{"title", nil, expected{"", "", "", "", "supersimplesoup", "supersimplesoup"}},
		{"a", nil, expected{"1st id", "link", "1st title", "1st href", "1st text", "1st text"}},
		{"a", []string{"id"}, expected{"1st id", "link", "1st title", "1st href", "1st text", "1st text"}},
		{"a", []string{"id", "1st id"}, expected{"1st id", "link", "1st title", "1st href", "1st text", "1st text"}},
		{"a", []string{"id", "2nd id"}, expected{"2nd id", "link", "2nd title", "2nd href", "2nd text", "2nd text"}},
		{"a", []string{"id", "3rd id"}, expected{"3rd id", "link", "3rd title", "3rd href", "3rd text", "3rd text"}},
		{"a", []string{"id", "4th id"}, expected{"4th id", "link", "4th title", "4th href", "4th text", "4th text"}},
		{"a", []string{"id", "5th id"}, expected{"5th id", "link", "5th title", "5th href", "5th text", "5th text"}},
	}

	for _, testcase := range testcases {
		found, err := n.Find(testcase.tag, testcase.attrkv...)
		if err != nil {
			t.Errorf("%s => unexpected error %q", testcase.tag, err)
		} else {
			actual := found.ID()
			if actual != testcase.expected.id {
				t.Errorf("%s => expected %q, got %q", testcase.tag, testcase.expected.id, actual)
			}
			actual = found.Class()
			if actual != testcase.expected.class {
				t.Errorf("%s => expected %q, got %q", testcase.tag, testcase.expected.class, actual)
			}
			actual = found.Title()
			if actual != testcase.expected.title {
				t.Errorf("%s => expected %q, got %q", testcase.tag, testcase.expected.title, actual)
			}
			actual = found.Href()
			if actual != testcase.expected.href {
				t.Errorf("%s => expected %q, got %q", testcase.tag, testcase.expected.href, actual)
			}
			actual = found.Text()
			if actual != testcase.expected.text {
				t.Errorf("%s => expected %q, got %q", testcase.tag, testcase.expected.text, actual)
			}
			actual = found.FullText()
			if actual != testcase.expected.fulltext {
				t.Errorf("%s => expected %q, got %q", testcase.tag, testcase.expected.fulltext, actual)
			}
		}
	}
}
