package stateparser

import (
	"fmt"
	"strings"

	"github.com/goptos/stateparser/ast"
	"github.com/goptos/stateparser/ast/nodes"
	"github.com/goptos/stateparser/stacks"
)

type nodeInfo struct {
	isEach          bool
	hasIf           bool
	isComponent     bool
	isSelfClosing   bool
	ifFunction      string
	collectFunction string
	keyFunction     string
	viewComponent   string
}

type Parser struct {
	Ast           *ast.Ast
	Result        string
	statements    []string
	index         int
	nodeInfoStack stacks.Stack[nodeInfo]
}

func New() *Parser {
	return &Parser{
		Ast:           nil,
		Result:        "",
		statements:    []string{},
		index:         0,
		nodeInfoStack: stacks.New[nodeInfo](),
	}
}

func (_self *Parser) reset() {
	_self.Ast = nil
	_self.Result = ""
	_self.statements = []string{}
	_self.index = -1
	_self.nodeInfoStack = stacks.New[nodeInfo]()
}

func (_self *Parser) updateNodeInfo(node *nodes.StartElementNode) {
	var nodeInfo = nodeInfo{
		isEach:          false,
		hasIf:           false,
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
				nodeInfo.ifFunction = childNode.GetEffect()
			}
			if childNode.GetName() == "each" {
				nodeInfo.collectFunction = childNode.GetEffect()
			}
			if childNode.GetName() == "key" {
				nodeInfo.keyFunction = childNode.GetEffect()
			}
		case nodes.StartElement:
			if childNode.GetIsComponent() && childNode.GetIsSelfClosing() {
				nodeInfo.viewComponent = childNode.GetName()
			}
		}
	}
	if nodeInfo.collectFunction != "" &&
		nodeInfo.keyFunction != "" &&
		nodeInfo.viewComponent != "" {
		nodeInfo.isEach = true
	}
	if nodeInfo.ifFunction != "" {
		nodeInfo.hasIf = true
	}
	_self.nodeInfoStack.Push(nodeInfo)
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
	var nodeInfo = _self.nodeInfoStack.Pop()
	if nodeInfo.hasIf {
		_self.appendToPrevStatement(".DynChild(cx, %s, %s)",
			nodeInfo.ifFunction,
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
		(*nodes.StartElementNode).Print(node, depth)
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
			if _self.nodeInfoStack.Peak().isEach {
				_self.prependToStatement("system.Each(")
				_self.appendToStatement(", cx, %s, %s, %s.View)",
					_self.nodeInfoStack.Peak().collectFunction,
					_self.nodeInfoStack.Peak().keyFunction,
					_self.nodeInfoStack.Peak().viewComponent)
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
		if _self.index == 0 && _self.nodeInfoStack.Peak().hasIf {
			_self.prependToStatement("system.DynElem(")
			_self.appendToStatement(", cx, %s)", _self.nodeInfoStack.Peak().ifFunction)
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
		(*nodes.TextNode).Print(node, depth)
		_self.appendToStatement(".Text(\"%s\")",
			node.GetData())
		return nil
	}
	/*
		`{count.Get()}` => `.DynText(cx, func() string { return fmt.Sprintf("%v", count.Get()) })`
	*/
	_self.Ast.DynTextNodeProcessor = func(node *nodes.DynTextNode, depth *int) error {
		(*nodes.DynTextNode).Print(node, depth)
		_self.appendToStatement(".DynText(cx, func() string { return fmt.Sprintf(\"%%v\", %s) })",
			node.GetEffect())
		return nil
	}
	/*
		`id="sub-button"` => `.Attr("id", "sub-button")`
	*/
	_self.Ast.AttributeNodeProcessor = func(node *nodes.AttributeNode, depth *int) error {
		(*nodes.AttributeNode).Print(node, depth)
		if _self.nodeInfoStack.Peak().isComponent {
			return nil
		}
		_self.appendToStatement(".Attr(\"%s\", \"%s\")",
			node.GetName(),
			node.GetValue())
		return nil
	}
	/*
		`on:click={ func(Event) {} }` => `.On("click", func(Event))`
	*/
	_self.Ast.EventAttributeNodeProcessor = func(node *nodes.EventAttributeNode, depth *int) error {
		(*nodes.EventAttributeNode).Print(node, depth)
		_self.appendToStatement(".On(\"%s\", %s)",
			node.GetEvent(),
			node.GetEffect())
		return nil
	}
	/*
		`class:dark={ func() bool {} }` => `.DynAttr("class", "dark", func() bool)`
	*/
	_self.Ast.DynAttributeNodeProcessor = func(node *nodes.DynAttributeNode, depth *int) error {
		(*nodes.DynAttributeNode).Print(node, depth)
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
		(*nodes.EndElementNode).Print(node, depth)
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
	_self.appendToStatement("\r")
	_self.Result = _self.statements[0]
	return nil
}
