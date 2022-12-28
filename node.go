package supersimplesoup

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"regexp"
	"strings"
)

var blankRegexp = regexp.MustCompile(`^\s+$`)

var (
	SkipNode = errors.New("skip this node")
	SkipAll  = errors.New("skip everything and stop the walk")
)

type Node html.Node

func Parse(r io.Reader) (*Node, error) {
	if n, err := html.Parse(r); err != nil {
		return nil, err
	} else {
		return (*Node)(n), nil
	}
}

type WalkFunc func(node *Node) error

func Walk(root *Node, fn WalkFunc) error {
	err := walk(root, fn)
	if err == SkipNode || err == SkipAll {
		return nil
	}
	return err
}

func walk(node *Node, fn WalkFunc) error {
	if node == nil {
		return nil
	}
	err := fn(node)
	if err != nil {
		return err
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if err := walk((*Node)(c), fn); err != nil {
			if err != SkipNode {
				return err
			}
		}
	}
	return nil
}

func (n *Node) ParentNode() *Node {
	return (*Node)(n.Parent)
}

func (n *Node) FirstChildNode() *Node {
	return (*Node)(n.FirstChild)
}

func (n *Node) LastChildNode() *Node {
	return (*Node)(n.LastChild)
}

func (n *Node) PrevSiblingNode() *Node {
	return (*Node)(n.PrevSibling)
}

func (n *Node) NextSiblingNode() *Node {
	return (*Node)(n.NextSibling)
}

func (n *Node) ChildrenNodes() (children []*Node) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, (*Node)(c))
	}
	return
}

func (n *Node) Attributes() map[string]string {
	if n.Type != html.ElementNode || len(n.Attr) == 0 {
		return nil
	}
	attrs := make(map[string]string, len(n.Attr))
	for _, attr := range n.Attr {
		attrs[attr.Key] = attr.Val
	}
	return attrs
}

func (n *Node) Attribute(key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func (n *Node) ID() string {
	return n.Attribute("id")
}

func (n *Node) Class() string {
	return n.Attribute("class")
}

func (n *Node) Href() string {
	return n.Attribute("href")
}

func (n *Node) Title() string {
	return n.Attribute("title")
}

func (n *Node) HTML() string {
	var buf bytes.Buffer
	html.Render(&buf, (*html.Node)(n))
	return buf.String()
}

func (n *Node) Text() (text string) {
	for _, c := range n.ChildrenNodes() {
		if c.Type != html.TextNode {
			continue
		}
		if blankRegexp.MatchString(c.Data) {
			continue
		}
		text = text + c.Data
	}
	return
}

func (n *Node) FullText() string {
	var buf bytes.Buffer
	n.Walk(func(node *Node) error {
		if node.Type == html.TextNode {
			buf.WriteString(node.Data)
		}
		return nil
	})
	return buf.String()
}

func (n *Node) Walk(fn WalkFunc) error {
	return Walk(n, fn)
}

func (n *Node) Find(tag string, attrkv ...string) (*Node, error) {
	if n == nil {
		return nil, fmt.Errorf("not allow to find on a blank node")
	}
	if ns := query(n, tag, attrkv, 1); len(ns) > 0 {
		return ns[0], nil
	} else {
		return nil, fmt.Errorf("not found element %s", prettyTagAttr(tag, attrkv))
	}
}

func (n *Node) Query(tag string, attrkv ...string) *Node {
	if n == nil {
		return nil
	}
	if ns := query(n, tag, attrkv, 1); len(ns) > 0 {
		return ns[0]
	} else {
		return nil
	}
}

func (n *Node) QueryAll(tag string, attrkv ...string) Nodes {
	if n == nil {
		return nil
	}
	return query(n, tag, attrkv, 0)
}

type Nodes []*Node

func (ns Nodes) Query(tag string, attrkv ...string) (found Nodes) {
	if ns == nil {
		return
	}
	for _, n := range ns {
		for _, node := range query(n, tag, attrkv, 1) {
			found = append(found, node)
		}
	}
	return
}

func (ns Nodes) QueryAll(tag string, attrkv ...string) (found Nodes) {
	if ns == nil {
		return
	}
	for _, n := range ns {
		for _, node := range query(n, tag, attrkv, 0) {
			found = append(found, node)
		}
	}
	return
}

func plainAttr(attrkv []string) (string, string) {
	if n := len(attrkv); n == 0 {
		return "", ""
	} else if n >= 2 {
		return attrkv[0], attrkv[1]
	} else {
		return attrkv[0], ""
	}
}

func prettyTagAttr(tag string, attrkv []string) string {
	key, val := plainAttr(attrkv)
	if key == "" && val == "" {
		return fmt.Sprintf("`%s`", tag)
	} else {
		return fmt.Sprintf("`%s` with attribute `%s=%s`", tag, key, val)
	}
}

func query(n *Node, tag string, attrkv []string, m int) (found Nodes) {
	Walk(n, func(node *Node) error {
		if node == n {
			return nil
		}
		if match(node, tag, attrkv) {
			found = append(found, node)
			if m > 0 && len(found) >= m {
				return SkipAll
			}
		}
		return nil
	})
	return
}

func match(n *Node, tag string, attrkv []string) bool {
	if n.Type != html.ElementNode {
		return false
	}
	if tag != "" && tag != n.Data {
		return false
	}
	key, val := plainAttr(attrkv)
	if key == "" {
		return true
	}
	for i := 0; i < len(n.Attr); i++ {
		if key != n.Attr[i].Key {
			continue
		}
		if val == "" || val == n.Attr[i].Val {
			return true
		}
		for _, f := range strings.Fields(val) {
			if !strings.Contains(n.Attr[i].Val, f) {
				return false
			}
		}
		return true
	}
	return false
}
