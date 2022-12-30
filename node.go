/*
Package supersimplesoup implements a super simple soup like DOM API.
*/
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

// Parse returns the parse tree for the HTML from the given Reader.
//
// The input is assumed to be UTF-8 encoded.
func Parse(r io.Reader) (*Node, error) {
	if n, err := html.Parse(r); err != nil {
		return nil, err
	} else {
		return (*Node)(n), nil
	}
}

// WalkFunc is the type of the function called by Walk to visit each node.
//
// The error result returned by the function controls how Walk continues.
// 	- If the function returns the special value SkipNode, Walk skips the current node.
// 	- If the function returns the special value SkipAll, Walk stops entirely and returns nil.
//	- If the function returns a non-nil error, Walk stops entirely and returns that error.
type WalkFunc func(node *Node) error

// Walk walks the node tree rooted at root, calling fn for each node in the tree, including root.
//
// The nodes are walked in depth first order, which makes the output deterministic.
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

// ParentNode returns the parent node of this node.
func (n *Node) ParentNode() *Node {
	return (*Node)(n.Parent)
}

// FirstChildNode returns the first direct child node of this node.
func (n *Node) FirstChildNode() *Node {
	return (*Node)(n.FirstChild)
}

// LastChildNode returns the last direct child node of this node.
func (n *Node) LastChildNode() *Node {
	return (*Node)(n.LastChild)
}

// PrevSiblingNode returns the previous sibling node of this node.
func (n *Node) PrevSiblingNode() *Node {
	return (*Node)(n.PrevSibling)
}

// NextSiblingNode returns the next sibling node of this node.
func (n *Node) NextSiblingNode() *Node {
	return (*Node)(n.NextSibling)
}

// Children returns all the direct child nodes of this node.
func (n *Node) ChildrenNodes() (children []*Node) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, (*Node)(c))
	}
	return
}

// IsElementNode returns whether is an element node.
func (n *Node) IsElementNode() bool {
	return n.Type == html.ElementNode
}

// IsTextNode returns whether is a text node.
func (n *Node) IsTextNode() bool {
	return n.Type == html.TextNode
}

// Attributes returns all the attributes key-value map of this node.
func (n *Node) Attributes() map[string]string {
	if !n.IsElementNode() || len(n.Attr) == 0 {
		return nil
	}
	attrs := make(map[string]string, len(n.Attr))
	for _, attr := range n.Attr {
		attrs[attr.Key] = attr.Val
	}
	return attrs
}

// Attribute return the key specified attribute of this node.
func (n *Node) Attribute(key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// ID returns the id attribute of this node.
func (n *Node) ID() string {
	return n.Attribute("id")
}

// Class returns the class attribute of this node.
func (n *Node) Class() string {
	return n.Attribute("class")
}

// Href returns the href attribute of this node.
func (n *Node) Href() string {
	return n.Attribute("href")
}

// Title returns the title attribute of this node.
func (n *Node) Title() string {
	return n.Attribute("title")
}

// HTML returns the HTML source code of this node.
func (n *Node) HTML() string {
	var buf bytes.Buffer
	html.Render(&buf, (*html.Node)(n))
	return buf.String()
}

// Text returns the text joined by all the direct child text nodes of this node.
func (n *Node) Text() (text string) {
	for _, c := range n.ChildrenNodes() {
		if !c.IsTextNode() {
			continue
		}
		if blankRegexp.MatchString(c.Data) {
			continue
		}
		text = text + c.Data
	}
	return
}

// Text returns the text joined by all the child text nodes in depth order of this node.
func (n *Node) FullText() string {
	var buf bytes.Buffer
	n.Walk(func(node *Node) error {
		if node.IsTextNode() {
			buf.WriteString(node.Data)
		}
		return nil
	})
	return buf.String()
}

// Walk walks the node tree rooted at this node, calling fn for each node in the tree, including this node.
func (n *Node) Walk(fn WalkFunc) error {
	return Walk(n, fn)
}

// Find returns the first child element node matched by the specified tag and optional attribute key and value of this node.
//
// It returns an error if no child element node is matched.
func (n *Node) Find(tag string, attrkv ...string) (*Node, error) {
	if n == nil {
		return nil, fmt.Errorf("not allow to find on a blank node")
	}
	if ns := query(n, tag, attrkv, 1); len(ns) > 0 {
		return ns[0], nil
	} else {
		return nil, fmt.Errorf("not found element `%s`", prettyTagAttr(tag, attrkv))
	}
}

// Query returns the first child element node matched by the specified tag and optional attribute key and value of this node.
//
// It returns nil if no child element node is matched.
//
// Allow chaining call.
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

// QueryAll returns the child element nodes matched by the specified tag and optional attribute key and value of this node.
//
// It returns nil if no child element node is matched.
//
// Allow chaining call.
func (n *Node) QueryAll(tag string, attrkv ...string) Nodes {
	if n == nil {
		return nil
	}
	return query(n, tag, attrkv, 0)
}

type Nodes []*Node

// Query returns all the first child element node matched by the specified tag and optional attribute key and value on each node of this nodes.
//
// It returns nil if no child element node is matched.
//
// Allow chaining call.
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

// QueryAll returns all the child element nodes matched by the specified tag and optional attribute key and value on each node of this nodes.
//
// It returns nil if no child element node is matched.
//
// Allow chaining call.
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
		return tag
	} else {
		return fmt.Sprintf("%s[%s=%s]", tag, key, val)
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
	if !n.IsElementNode() {
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
