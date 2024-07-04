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
	Ast        *ast.Ast
	Result     string
	statements stacks.Stack[string]
	nodeInfo   stacks.Stack[nodeInfo]
}

func New() *Parser {
	return &Parser{
		Ast:        nil,
		Result:     "",
		statements: stacks.New[string](),
		nodeInfo:   stacks.New[nodeInfo](),
	}
}

func (_self *Parser) reset() {
	_self.Ast = nil
	_self.Result = ""
	_self.statements = stacks.New[string]()
	_self.nodeInfo = stacks.New[nodeInfo]()
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
	_self.nodeInfo.Push(nodeInfo)
}

func (_self *Parser) newStatement(s string, args ...interface{}) {
	_self.statements.Push(fmt.Sprintf(s, args...))
}

func (_self *Parser) appendToStatement(s string, args ...interface{}) {
	_self.statements.Push(_self.statements.Pop() + fmt.Sprintf(s, args...))
}

func (_self *Parser) prependToStatement(s string, args ...interface{}) {
	_self.statements.Push(fmt.Sprintf(s, args...) + _self.statements.Pop())
}

func (_self *Parser) statementContains(s string, args ...interface{}) bool {
	if _self.statements.Depth() < 0 {
		return false
	}
	return strings.Contains(_self.statements.Peak(), fmt.Sprintf(s, args...))
}

func (_self *Parser) squashStatement() {
	if _self.statements.Depth() == 0 && _self.nodeInfo.Peak().hasIf {
		_self.prependToStatement("`<invalid view: if statements on root elements are not supported>` //")
	}
	if _self.statements.Depth() == 0 {
		return
	}
	var statement = _self.statements.Pop()
	if _self.nodeInfo.Peak().hasIf {
		_self.appendToStatement(".DynChild(cx, %s, %s)",
			_self.nodeInfo.Peak().ifFunction,
			statement)
		return
	}
	_self.appendToStatement(".Child(%s)",
		statement)
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
			if _self.nodeInfo.Peak().isEach {
				_self.prependToStatement("system.Each(")
				_self.appendToStatement(", cx, %s, %s, %s.View)",
					_self.nodeInfo.Peak().collectFunction,
					_self.nodeInfo.Peak().keyFunction,
					_self.nodeInfo.Peak().viewComponent)
			}
		}
		/*
			`<Button arg1 arg2 />` => `Button.View(cx, arg1, arg2)`
		*/
		if node.GetIsComponent() {
			if _self.statementContains(", %s.View)", node.GetName()) {
				_self.nodeInfo.Pop()
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
		if node.GetIsSelfClosing() {
			_self.squashStatement()
			_self.nodeInfo.Pop()
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
		if _self.nodeInfo.Peak().isComponent {
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
		_self.appendToStatement(".DynAttr(cx, %s, \"%s\", \"%s\")",
			node.GetEffect(),
			node.GetName(),
			node.GetValue())
		return nil
	}
	/*
		</...>
	*/
	_self.Ast.EndElementNodeProcessor = func(node *nodes.EndElementNode, depth *int) error {
		(*nodes.EndElementNode).Print(node, depth)
		_self.squashStatement()
		_self.nodeInfo.Pop()
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
	_self.Result = _self.statements.Pop() + "\r"
	return nil
}
