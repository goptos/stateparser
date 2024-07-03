package stateparser

import (
	"fmt"
	"strings"

	"github.com/goptos/stateparser/ast"
	"github.com/goptos/stateparser/ast/nodes"
)

type nodeInfo struct {
	isEach          bool
	isDynElem       bool
	isComponent     bool
	isSelfClosing   bool
	ifFunction      string
	collectFunction string
	keyFunction     string
	viewComponent   string
}

type Parser struct {
	Ast        *ast.Ast
	Result     string
	statements []string
	index      int
	nodeInfo   *nodeInfo
}

func New() *Parser {
	return &Parser{
		Ast:        nil,
		Result:     "",
		statements: []string{},
		index:      0,
		nodeInfo:   nil,
	}
}

func (_self *Parser) reset() {
	_self.Ast = nil
	_self.Result = ""
	_self.statements = []string{}
	_self.index = -1
	_self.nodeInfo = nil
}

func (_self *Parser) updateNodeInfo(node *nodes.StartElementNode) {
	_self.nodeInfo = &nodeInfo{
		isEach:          false,
		isDynElem:       false,
		isComponent:     node.GetIsComponent(),
		isSelfClosing:   node.GetIsSelfClosing(),
		ifFunction:      "",
		collectFunction: "",
		keyFunction:     "",
		viewComponent:   "",
	}
	for _, childNode := range node.GetChildren() {
		switch childNode.GetType() {
		case nodes.KeywordAttribute:
			if childNode.GetName() == "if" {
				_self.nodeInfo.ifFunction = childNode.GetEffect()
			}
			if childNode.GetName() == "each" {
				_self.nodeInfo.collectFunction = childNode.GetEffect()
			}
			if childNode.GetName() == "key" {
				_self.nodeInfo.keyFunction = childNode.GetEffect()
			}
		case nodes.StartElement:
			if childNode.GetIsComponent() && childNode.GetIsSelfClosing() {
				_self.nodeInfo.viewComponent = childNode.GetName()
			}
		}
	}
	if _self.nodeInfo.collectFunction != "" &&
		_self.nodeInfo.keyFunction != "" &&
		_self.nodeInfo.viewComponent != "" {
		_self.nodeInfo.isEach = true
	}
	if _self.nodeInfo.ifFunction != "" {
		_self.nodeInfo.isDynElem = true
	}
}

func (_self *Parser) newStatement(s string, args ...interface{}) {
	_self.statements = append(_self.statements, fmt.Sprintf(s, args...))
	_self.index = len(_self.statements) - 1
}

func (_self *Parser) appendToStatement(s string, args ...interface{}) {
	_self.statements[_self.index] = _self.statements[_self.index] +
		fmt.Sprintf(s, args...)
}

func (_self *Parser) prependToStatement(s string, args ...interface{}) {
	_self.statements[_self.index] = fmt.Sprintf(s, args...) +
		_self.statements[_self.index]
}

func (_self *Parser) appendToPrevStatement(s string, args ...interface{}) {
	if _self.index-1 < 0 {
		return
	}
	_self.statements[_self.index-1] = _self.statements[_self.index-1] +
		fmt.Sprintf(s, args...)
}

func (_self *Parser) squashStatement() {
	if _self.index-1 < 0 {
		return
	}
	if _self.nodeInfo.isDynElem {
		_self.appendToPrevStatement(".DynChild(cx, %s, %s)",
			_self.nodeInfo.ifFunction,
			_self.statements[_self.index])
	} else {
		_self.appendToPrevStatement(".Child(%s)",
			_self.statements[_self.index])
	}
	// drop the last statement
	_self.statements = _self.statements[0:_self.index]
	// re-index the array
	_self.index = len(_self.statements) - 1
}

func (_self *Parser) statementContains(s string, args ...interface{}) bool {
	if _self.index-1 < 0 {
		return false
	}
	return strings.Contains(_self.statements[_self.index], fmt.Sprintf(s, args...))
}

func (_self *Parser) ParseView(source string) error {
	_self.reset()
	_self.Ast = ast.New(source)
	_self.Ast.KeywordAttributeNames["if"] = nil
	_self.Ast.KeywordAttributeNames["each"] = nil
	_self.Ast.KeywordAttributeNames["key"] = nil
	/*
		`<...`
	*/
	_self.Ast.StartElementNodeProcessor = func(node *nodes.StartElementNode, depth *int) error {
		_self.updateNodeInfo(node)
		/*
			`<div>` => `(*Elem).New(nil, "div")`
		*/
		if !node.GetIsComponent() {
			_self.newStatement("(*Elem).New(nil, \"%s\")", node.GetName())
			/*
				`<ul each={cF} key={kF}><Li /></ul>` =>
				`system.Each((*Elem).New(nil, "ul"), cx, cF, kF, Li.View)`
			*/
			if _self.nodeInfo.isEach {
				_self.prependToStatement("system.Each(")
				_self.appendToStatement(", cx, %s, %s, %s.View)",
					_self.nodeInfo.collectFunction,
					_self.nodeInfo.keyFunction,
					_self.nodeInfo.viewComponent)
			}
		}
		/*
			`<Button arg1 arg2 />` => `Button.View(cx, arg1, arg2)`
		*/
		if node.GetIsComponent() {
			if _self.statementContains(", %s.View)", node.GetName()) {
				return nil
			}
			_self.newStatement("%s.View(cx", node.GetName())
			for _, childNode := range node.GetChildren() {
				if childNode.GetType() == nodes.Attribute {
					_self.appendToStatement(", %s", childNode.GetName())
				}
			}
			_self.appendToStatement(")")
		}
		if _self.index == 0 && _self.nodeInfo.isDynElem {
			_self.prependToStatement("system.DynElem(")
			_self.appendToStatement(", cx, %s)", _self.nodeInfo.ifFunction)
		}
		if node.GetIsSelfClosing() {
			_self.squashStatement()
		}
		return nil
	}
	/*
		`Hello` => `.Text("Hello")`
	*/
	_self.Ast.TextNodeProcessor = func(node *nodes.TextNode, depth *int) error {
		_self.appendToStatement(".Text(\"%s\")",
			node.GetData())
		return nil
	}
	/*
		`{count.Get()}` => `.DynText(cx, func() string { return fmt.Sprintf("%v", count.Get()) })`
	*/
	_self.Ast.DynTextNodeProcessor = func(node *nodes.DynTextNode, depth *int) error {
		_self.appendToStatement(".DynText(cx, func() string { return fmt.Sprintf(\"%%v\", %s) })",
			node.GetEffect())
		return nil
	}
	/*
		`id="sub-button"` => `.Attr("id", "sub-button")`
	*/
	_self.Ast.AttributeNodeProcessor = func(node *nodes.AttributeNode, depth *int) error {
		if _self.nodeInfo.isComponent {
			return nil
		}
		_self.appendToStatement(".Attr(\"%s\", %s)",
			node.GetName(),
			node.GetValue())
		return nil
	}
	/*
		`on:click={ func(Event) {} }` => `.On("click", func(Event))`
	*/
	_self.Ast.EventAttributeNodeProcessor = func(node *nodes.EventAttributeNode, depth *int) error {
		_self.appendToStatement(".On(\"%s\", %s)",
			node.GetEvent(),
			node.GetEffect())
		return nil
	}
	/*
		`class:dark={ func() bool {} }` => `.DynAttr("class", "dark", func() bool)`
	*/
	_self.Ast.DynAttributeNodeProcessor = func(node *nodes.DynAttributeNode, depth *int) error {
		_self.appendToStatement(".DynAttr(\"%s\", \"%s\", %s)",
			node.GetName(),
			node.GetValue(),
			node.GetEffect())
		return nil
	}
	/*
		</...>
	*/
	_self.Ast.EndElementNodeProcessor = func(node *nodes.EndElementNode, depth *int) error {
		_self.squashStatement()
		return nil
	}
	err := _self.Ast.Create()
	if err != nil {
		return err
	}
	err = _self.Ast.Process()
	if err != nil {
		return err
	}
	_self.appendToStatement("\r\n")
	_self.Result = _self.statements[0]
	return nil
}
