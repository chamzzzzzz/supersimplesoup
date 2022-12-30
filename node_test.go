package supersimplesoup

import (
	"reflect"
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
			<ul id="ul-id-1" title="ul-title-1" class="ul-class-1">
				<li id="li-id-1" title="li-title-1" class="li-class-1">
					<a id="a-id-1" href="a-href-1" title="a-title-1" class="a-class-1">a-text-1</a>
					<a id="a-id-2" href="a-href-2" title="a-title-2" class="a-class-1">a-text-2</a>
				</li>
				<li id="li-id-2" title="li-title-2" class="li-class-1">
					<a id="a-id-3" href="a-href-3" title="a-title-3" class="a-class-1">a-text-3</a>
					<a id="a-id-4" href="a-href-4" title="a-title-4" class="a-class-1">a-text-4</a>
				</li>
			</ul>
			<ul id="ul-id-2" title="ul-title-2" class="ul-class-2">
				<li id="li-id-3" title="li-title-3" class="li-class-2">
					<a id="a-id-5" href="a-href-5" title="a-title-5" class="a-class-2">a-text-5</a>
					<a id="a-id-6" href="a-href-6" title="a-title-6" class="a-class-2">a-text-6</a>
				</li>
				<li id="li-id-4" title="li-title-4" class="li-class-2">
					<a id="a-id-7" href="a-href-7" title="a-title-7" class="a-class-2">a-text-7</a>
					<a id="a-id-8" href="a-href-8" title="a-title-8" class="a-class-2">a-text-8</a>
				</li>
			</ul>
		</div>
	</body>
</html>
`

func getTestNode() *Node {
	n, err := Parse(strings.NewReader(testHTML))
	if err != nil {
		panic(err)
	}
	return n
}

var root *Node = getTestNode()

func TestWalk(t *testing.T) {
	want := 19
	got := 0
	Walk(root, func(node *Node) error {
		if node.IsElementNode() {
			got++
		}
		return nil
	})
	if want != got {
		t.Errorf("element node count, want %d, got %d", want, got)
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		tag    string
		attrkv []string
	}{
		{"title", nil},
		{"div", nil},
		{"ul", nil},
		{"ul", []string{"id"}},
		{"ul", []string{"id", "ul-id-1"}},
		{"ul", []string{"id", "ul-id-2"}},
		{"li", nil},
		{"li", []string{"id"}},
		{"li", []string{"id", "li-id-1"}},
		{"li", []string{"id", "li-id-2"}},
		{"li", []string{"id", "li-id-3"}},
		{"li", []string{"id", "li-id-4"}},
		{"a", nil},
		{"a", []string{"id"}},
		{"a", []string{"id", "a-id-1"}},
		{"a", []string{"id", "a-id-2"}},
		{"a", []string{"id", "a-id-3"}},
		{"a", []string{"id", "a-id-4"}},
		{"a", []string{"id", "a-id-5"}},
		{"a", []string{"id", "a-id-6"}},
		{"a", []string{"id", "a-id-7"}},
		{"a", []string{"id", "a-id-8"}},
		{"a", []string{"title", "a-title-1"}},
		{"a", []string{"title", "a-title-2"}},
		{"a", []string{"title", "a-title-3"}},
		{"a", []string{"title", "a-title-4"}},
		{"a", []string{"title", "a-title-5"}},
		{"a", []string{"title", "a-title-6"}},
		{"a", []string{"title", "a-title-7"}},
		{"a", []string{"title", "a-title-8"}},
		{"a", []string{"class", "a-class-1"}},
		{"a", []string{"class", "a-class-2"}},
	}

	want := 1
	for _, test := range tests {
		got := 0
		node, err := root.Find(test.tag, test.attrkv...)
		if err == nil && node != nil {
			got = 1
		}
		if want != got {
			t.Errorf("`%s` element node count, want %d, got %d", prettyTagAttr(test.tag, test.attrkv), want, got)
		}
	}
}

func TestQuery(t *testing.T) {
	tests := []struct {
		tag    string
		attrkv []string
	}{
		{"title", nil},
		{"div", nil},
		{"ul", nil},
		{"ul", []string{"id"}},
		{"ul", []string{"id", "ul-id-1"}},
		{"ul", []string{"id", "ul-id-2"}},
		{"li", nil},
		{"li", []string{"id"}},
		{"li", []string{"id", "li-id-1"}},
		{"li", []string{"id", "li-id-2"}},
		{"li", []string{"id", "li-id-3"}},
		{"li", []string{"id", "li-id-4"}},
		{"a", nil},
		{"a", []string{"id"}},
		{"a", []string{"id", "a-id-1"}},
		{"a", []string{"id", "a-id-2"}},
		{"a", []string{"id", "a-id-3"}},
		{"a", []string{"id", "a-id-4"}},
		{"a", []string{"id", "a-id-5"}},
		{"a", []string{"id", "a-id-6"}},
		{"a", []string{"id", "a-id-7"}},
		{"a", []string{"id", "a-id-8"}},
		{"a", []string{"title", "a-title-1"}},
		{"a", []string{"title", "a-title-2"}},
		{"a", []string{"title", "a-title-3"}},
		{"a", []string{"title", "a-title-4"}},
		{"a", []string{"title", "a-title-5"}},
		{"a", []string{"title", "a-title-6"}},
		{"a", []string{"title", "a-title-7"}},
		{"a", []string{"title", "a-title-8"}},
		{"a", []string{"class", "a-class-1"}},
		{"a", []string{"class", "a-class-2"}},
	}

	want := 1
	for _, test := range tests {
		got := 0
		node := root.Query(test.tag, test.attrkv...)
		if node != nil {
			got = 1
		}
		if want != got {
			t.Errorf("`%s` element node count, want %d, got %d", prettyTagAttr(test.tag, test.attrkv), want, got)
		}
	}
}

func TestQueryAll(t *testing.T) {
	tests := []struct {
		tag    string
		attrkv []string
		want   int
	}{
		{"title", nil, 1},
		{"div", nil, 1},
		{"ul", nil, 2},
		{"ul", []string{"id"}, 2},
		{"ul", []string{"id", "ul-id-1"}, 1},
		{"ul", []string{"id", "ul-id-2"}, 1},
		{"li", nil, 4},
		{"li", []string{"id"}, 4},
		{"li", []string{"id", "li-id-1"}, 1},
		{"li", []string{"id", "li-id-2"}, 1},
		{"li", []string{"id", "li-id-3"}, 1},
		{"li", []string{"id", "li-id-4"}, 1},
		{"a", nil, 8},
		{"a", []string{"id"}, 8},
		{"a", []string{"id", "a-id-1"}, 1},
		{"a", []string{"id", "a-id-2"}, 1},
		{"a", []string{"id", "a-id-3"}, 1},
		{"a", []string{"id", "a-id-4"}, 1},
		{"a", []string{"id", "a-id-5"}, 1},
		{"a", []string{"id", "a-id-6"}, 1},
		{"a", []string{"id", "a-id-7"}, 1},
		{"a", []string{"id", "a-id-8"}, 1},
		{"a", []string{"title", "a-title-1"}, 1},
		{"a", []string{"title", "a-title-2"}, 1},
		{"a", []string{"title", "a-title-3"}, 1},
		{"a", []string{"title", "a-title-4"}, 1},
		{"a", []string{"title", "a-title-5"}, 1},
		{"a", []string{"title", "a-title-6"}, 1},
		{"a", []string{"title", "a-title-7"}, 1},
		{"a", []string{"title", "a-title-8"}, 1},
		{"a", []string{"class", "a-class-1"}, 4},
		{"a", []string{"class", "a-class-2"}, 4},
	}

	for _, test := range tests {
		nodes := root.QueryAll(test.tag, test.attrkv...)
		got := len(nodes)
		if test.want != got {
			t.Errorf("`%s` element node count, want %d, got %d", prettyTagAttr(test.tag, test.attrkv), test.want, got)
		}
	}
}

func TestQueryChain(t *testing.T) {
	want := []string{"a-id-1"}
	var got []string
	if node := root.Query("ul").Query("li").Query("a"); node != nil {
		got = append(got, node.ID())
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("`%s` element attribute id, want %v, got %v", "ul.li.a", want, got)
	}

	want = []string{"a-id-1", "a-id-2"}
	got = nil
	for _, node := range root.Query("ul").Query("li").QueryAll("a") {
		got = append(got, node.ID())
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("`%s` element attribute id, want %v, got %v", "ul.li.$$a", want, got)
	}

	want = []string{"a-id-1", "a-id-3"}
	got = nil
	for _, node := range root.Query("ul").QueryAll("li").Query("a") {
		got = append(got, node.ID())
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("`%s` element attribute id, want %v, got %v", "ul.$$li.a", want, got)
	}

	want = []string{"a-id-1", "a-id-2", "a-id-3", "a-id-4"}
	got = nil
	for _, node := range root.Query("ul").QueryAll("li").QueryAll("a") {
		got = append(got, node.ID())
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("`%s` element attribute id, want %v, got %v", "ul.$$li.$$a", want, got)
	}

	want = []string{"a-id-1", "a-id-5"}
	got = nil
	for _, node := range root.QueryAll("ul").Query("li").Query("a") {
		got = append(got, node.ID())
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("`%s` element attribute id, want %v, got %v", "$$ul.li.a", want, got)
	}

	want = []string{"a-id-1", "a-id-2", "a-id-5", "a-id-6"}
	got = nil
	for _, node := range root.QueryAll("ul").Query("li").QueryAll("a") {
		got = append(got, node.ID())
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("`%s` element attribute id, want %v, got %v", "$$ul.li.$$a", want, got)
	}

	want = []string{"a-id-1", "a-id-3", "a-id-5", "a-id-7"}
	got = nil
	for _, node := range root.QueryAll("ul").QueryAll("li").Query("a") {
		got = append(got, node.ID())
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("`%s` element attribute id, want %v, got %v", "$$ul.$$li.a", want, got)
	}

	want = []string{"a-id-1", "a-id-2", "a-id-3", "a-id-4", "a-id-5", "a-id-6", "a-id-7", "a-id-8"}
	got = nil
	for _, node := range root.QueryAll("ul").QueryAll("li").QueryAll("a") {
		got = append(got, node.ID())
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("`%s` element attribute id, want %v, got %v", "$$ul.$$li.$$a", want, got)
	}
}

func TestText(t *testing.T) {
	tests := []struct {
		tag    string
		attrkv []string
		want   string
	}{
		{"title", nil, "supersimplesoup"},
		{"div", nil, ""},
		{"ul", nil, ""},
		{"ul", []string{"id"}, ""},
		{"ul", []string{"id", "ul-id-1"}, ""},
		{"ul", []string{"id", "ul-id-2"}, ""},
		{"li", nil, ""},
		{"li", []string{"id"}, ""},
		{"li", []string{"id", "li-id-1"}, ""},
		{"li", []string{"id", "li-id-2"}, ""},
		{"li", []string{"id", "li-id-3"}, ""},
		{"li", []string{"id", "li-id-4"}, ""},
		{"a", nil, "a-text-1"},
		{"a", []string{"id"}, "a-text-1"},
		{"a", []string{"id", "a-id-1"}, "a-text-1"},
		{"a", []string{"id", "a-id-2"}, "a-text-2"},
		{"a", []string{"id", "a-id-3"}, "a-text-3"},
		{"a", []string{"id", "a-id-4"}, "a-text-4"},
		{"a", []string{"id", "a-id-5"}, "a-text-5"},
		{"a", []string{"id", "a-id-6"}, "a-text-6"},
		{"a", []string{"id", "a-id-7"}, "a-text-7"},
		{"a", []string{"id", "a-id-8"}, "a-text-8"},
		{"a", []string{"title", "a-title-1"}, "a-text-1"},
		{"a", []string{"title", "a-title-2"}, "a-text-2"},
		{"a", []string{"title", "a-title-3"}, "a-text-3"},
		{"a", []string{"title", "a-title-4"}, "a-text-4"},
		{"a", []string{"title", "a-title-5"}, "a-text-5"},
		{"a", []string{"title", "a-title-6"}, "a-text-6"},
		{"a", []string{"title", "a-title-7"}, "a-text-7"},
		{"a", []string{"title", "a-title-8"}, "a-text-8"},
		{"a", []string{"class", "a-class-1"}, "a-text-1"},
		{"a", []string{"class", "a-class-2"}, "a-text-5"},
	}

	for _, test := range tests {
		node := root.Query(test.tag, test.attrkv...)
		if node == nil {
			t.Errorf("`%s` element node count, want %d, got %d", prettyTagAttr(test.tag, test.attrkv), 1, 0)
			continue
		}
		got := node.Text()
		if test.want != got {
			t.Errorf("`%s` element node text, want %q, got %q", prettyTagAttr(test.tag, test.attrkv), test.want, got)
		}
	}
}

func TestFullText(t *testing.T) {
	tests := []struct {
		tag    string
		attrkv []string
		want   string
	}{
		{"title", nil, "supersimplesoup"},
		{"li", nil, "\n\t\t\t\t\ta-text-1\n\t\t\t\t\ta-text-2\n\t\t\t\t"},
		{"li", []string{"id"}, "\n\t\t\t\t\ta-text-1\n\t\t\t\t\ta-text-2\n\t\t\t\t"},
		{"li", []string{"id", "li-id-1"}, "\n\t\t\t\t\ta-text-1\n\t\t\t\t\ta-text-2\n\t\t\t\t"},
		{"li", []string{"id", "li-id-2"}, "\n\t\t\t\t\ta-text-3\n\t\t\t\t\ta-text-4\n\t\t\t\t"},
		{"li", []string{"id", "li-id-3"}, "\n\t\t\t\t\ta-text-5\n\t\t\t\t\ta-text-6\n\t\t\t\t"},
		{"li", []string{"id", "li-id-4"}, "\n\t\t\t\t\ta-text-7\n\t\t\t\t\ta-text-8\n\t\t\t\t"},
		{"a", nil, "a-text-1"},
		{"a", []string{"id"}, "a-text-1"},
		{"a", []string{"id", "a-id-1"}, "a-text-1"},
		{"a", []string{"id", "a-id-2"}, "a-text-2"},
		{"a", []string{"id", "a-id-3"}, "a-text-3"},
		{"a", []string{"id", "a-id-4"}, "a-text-4"},
		{"a", []string{"id", "a-id-5"}, "a-text-5"},
		{"a", []string{"id", "a-id-6"}, "a-text-6"},
		{"a", []string{"id", "a-id-7"}, "a-text-7"},
		{"a", []string{"id", "a-id-8"}, "a-text-8"},
		{"a", []string{"title", "a-title-1"}, "a-text-1"},
		{"a", []string{"title", "a-title-2"}, "a-text-2"},
		{"a", []string{"title", "a-title-3"}, "a-text-3"},
		{"a", []string{"title", "a-title-4"}, "a-text-4"},
		{"a", []string{"title", "a-title-5"}, "a-text-5"},
		{"a", []string{"title", "a-title-6"}, "a-text-6"},
		{"a", []string{"title", "a-title-7"}, "a-text-7"},
		{"a", []string{"title", "a-title-8"}, "a-text-8"},
		{"a", []string{"class", "a-class-1"}, "a-text-1"},
		{"a", []string{"class", "a-class-2"}, "a-text-5"},
	}

	for _, test := range tests {
		node := root.Query(test.tag, test.attrkv...)
		if node == nil {
			t.Errorf("`%s` element node count, want %d, got %d", prettyTagAttr(test.tag, test.attrkv), 1, 0)
			continue
		}
		got := node.FullText()
		if test.want != got {
			t.Errorf("`%s` element node full text, want %q, got %q", prettyTagAttr(test.tag, test.attrkv), test.want, got)
		}
	}
}

func TestAttribute(t *testing.T) {
	tests := []struct {
		tag    string
		attrkv []string
		want   []string
	}{
		{"title", nil, []string{"", "", "", ""}},
		{"div", nil, []string{"", "", "", ""}},
		{"ul", nil, []string{"ul-id-1", "ul-class-1", "ul-title-1", ""}},
		{"ul", []string{"id"}, []string{"ul-id-1", "ul-class-1", "ul-title-1", ""}},
		{"ul", []string{"id", "ul-id-1"}, []string{"ul-id-1", "ul-class-1", "ul-title-1", ""}},
		{"ul", []string{"id", "ul-id-2"}, []string{"ul-id-2", "ul-class-2", "ul-title-2", ""}},
		{"li", nil, []string{"li-id-1", "li-class-1", "li-title-1", ""}},
		{"li", []string{"id"}, []string{"li-id-1", "li-class-1", "li-title-1", ""}},
		{"li", []string{"id", "li-id-1"}, []string{"li-id-1", "li-class-1", "li-title-1", ""}},
		{"li", []string{"id", "li-id-2"}, []string{"li-id-2", "li-class-1", "li-title-2", ""}},
		{"li", []string{"id", "li-id-3"}, []string{"li-id-3", "li-class-2", "li-title-3", ""}},
		{"li", []string{"id", "li-id-4"}, []string{"li-id-4", "li-class-2", "li-title-4", ""}},
		{"a", nil, []string{"a-id-1", "a-class-1", "a-title-1", "a-href-1"}},
		{"a", []string{"id"}, []string{"a-id-1", "a-class-1", "a-title-1", "a-href-1"}},
		{"a", []string{"id", "a-id-1"}, []string{"a-id-1", "a-class-1", "a-title-1", "a-href-1"}},
		{"a", []string{"id", "a-id-2"}, []string{"a-id-2", "a-class-1", "a-title-2", "a-href-2"}},
		{"a", []string{"id", "a-id-3"}, []string{"a-id-3", "a-class-1", "a-title-3", "a-href-3"}},
		{"a", []string{"id", "a-id-4"}, []string{"a-id-4", "a-class-1", "a-title-4", "a-href-4"}},
		{"a", []string{"id", "a-id-5"}, []string{"a-id-5", "a-class-2", "a-title-5", "a-href-5"}},
		{"a", []string{"id", "a-id-6"}, []string{"a-id-6", "a-class-2", "a-title-6", "a-href-6"}},
		{"a", []string{"id", "a-id-7"}, []string{"a-id-7", "a-class-2", "a-title-7", "a-href-7"}},
		{"a", []string{"id", "a-id-8"}, []string{"a-id-8", "a-class-2", "a-title-8", "a-href-8"}},
		{"a", []string{"title", "a-title-1"}, []string{"a-id-1", "a-class-1", "a-title-1", "a-href-1"}},
		{"a", []string{"title", "a-title-2"}, []string{"a-id-2", "a-class-1", "a-title-2", "a-href-2"}},
		{"a", []string{"title", "a-title-3"}, []string{"a-id-3", "a-class-1", "a-title-3", "a-href-3"}},
		{"a", []string{"title", "a-title-4"}, []string{"a-id-4", "a-class-1", "a-title-4", "a-href-4"}},
		{"a", []string{"title", "a-title-5"}, []string{"a-id-5", "a-class-2", "a-title-5", "a-href-5"}},
		{"a", []string{"title", "a-title-6"}, []string{"a-id-6", "a-class-2", "a-title-6", "a-href-6"}},
		{"a", []string{"title", "a-title-7"}, []string{"a-id-7", "a-class-2", "a-title-7", "a-href-7"}},
		{"a", []string{"title", "a-title-8"}, []string{"a-id-8", "a-class-2", "a-title-8", "a-href-8"}},
		{"a", []string{"class", "a-class-1"}, []string{"a-id-1", "a-class-1", "a-title-1", "a-href-1"}},
		{"a", []string{"class", "a-class-2"}, []string{"a-id-5", "a-class-2", "a-title-5", "a-href-5"}},
	}

	for _, test := range tests {
		node := root.Query(test.tag, test.attrkv...)
		if node == nil {
			t.Errorf("`%s` element node count, want %d, got %d", prettyTagAttr(test.tag, test.attrkv), 1, 0)
			continue
		}
		got := node.ID()
		if test.want[0] != got {
			t.Errorf("`%s` element node attribute id, want %q, got %q", prettyTagAttr(test.tag, test.attrkv), test.want[0], got)
		}
		got = node.Class()
		if test.want[1] != got {
			t.Errorf("`%s` element node attribute class, want %q, got %q", prettyTagAttr(test.tag, test.attrkv), test.want[1], got)
		}
		got = node.Title()
		if test.want[2] != got {
			t.Errorf("`%s` element node attribute title, want %q, got %q", prettyTagAttr(test.tag, test.attrkv), test.want[2], got)
		}
		got = node.Href()
		if test.want[3] != got {
			t.Errorf("`%s` element node attribute href, want %q, got %q", prettyTagAttr(test.tag, test.attrkv), test.want[3], got)
		}
	}
}
