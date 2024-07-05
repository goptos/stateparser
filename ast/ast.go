package ast

import (
	"fmt"

	"github.com/goptos/stateparser/ast/nodes"
	"github.com/goptos/stateparser/lexer"
	"github.com/goptos/stateparser/lexer/tokens"
	"github.com/goptos/utils"
)

var verbose = (*utils.Verbose).New(nil)

type Ast struct {
	keywordAttributeNames          map[string]interface{}
	Lexer                          *lexer.Lexer
	Root                           nodes.Node
	StartElementNodeProcessor      func(node *nodes.StartElementNode, depth *int) error
	ComponentNodeProcessor         func(node *nodes.ComponentNode, depth *int) error
	EndElementNodeProcessor        func(node *nodes.EndElementNode, depth *int) error
	CommentNodeProcessor           func(node *nodes.CommentNode, depth *int) error
	TextNodeProcessor              func(node *nodes.TextNode, depth *int) error
	DynTextNodeProcessor           func(node *nodes.DynTextNode, depth *int) error
	AttributeNodeProcessor         func(node *nodes.AttributeNode, depth *int) error
	ArgumentAttributeNodeProcessor func(node *nodes.ArgumentAttributeNode, depth *int) error
	DynAttributeNodeProcessor      func(node *nodes.DynAttributeNode, depth *int) error
	EventAttributeNodeProcessor    func(node *nodes.EventAttributeNode, depth *int) error
	KeywordAttributeNodeProcessor  func(node *nodes.KeywordAttributeNode, depth *int) error
}

func New(source string) *Ast {
	return &Ast{
		keywordAttributeNames:          make(map[string]interface{}),
		Lexer:                          lexer.New(source),
		Root:                           nil,
		StartElementNodeProcessor:      (*nodes.StartElementNode).Print,
		ComponentNodeProcessor:         (*nodes.ComponentNode).Print,
		EndElementNodeProcessor:        (*nodes.EndElementNode).Print,
		CommentNodeProcessor:           (*nodes.CommentNode).Print,
		TextNodeProcessor:              (*nodes.TextNode).Print,
		DynTextNodeProcessor:           (*nodes.DynTextNode).Print,
		AttributeNodeProcessor:         (*nodes.AttributeNode).Print,
		ArgumentAttributeNodeProcessor: (*nodes.ArgumentAttributeNode).Print,
		DynAttributeNodeProcessor:      (*nodes.DynAttributeNode).Print,
		EventAttributeNodeProcessor:    (*nodes.EventAttributeNode).Print,
		KeywordAttributeNodeProcessor:  (*nodes.KeywordAttributeNode).Print}
}

func (_self *Ast) AddKeywordAttributeName(s string) {
	_self.keywordAttributeNames[s] = nil
	_self.Lexer.KeywordAttributeNames[s] = nil
}

func (_self *Ast) Create() error {
	err := _self.Lexer.Tokenise()
	if err != nil {
		return err
	}
	verbose.Printf(3, "::: Ast.Create() :::\n")
	for i := 0; i < len(_self.Lexer.Tokens); i++ {
		if _self.Lexer.Tokens[i].GetType() != tokens.StartTag {
			continue
		}
		var index = i
		root, err := _self.createR(&index)
		if err != nil {
			return err
		}
		_self.Root = root
		break
	}
	if _self.Root == nil {
		return fmt.Errorf("must be a HTML element or a Component")
	}
	return nil
}

func (_self *Ast) createR(index *int) (nodes.Node, error) {
	verbose.Printf(3, "%d\t%s (%s)\n", *index, _self.Lexer.Tokens[*index].GetName(), _self.Lexer.Tokens[*index].GetType())
	utils.Assert(
		_self.Lexer.Tokens[*index].GetType() == tokens.StartTag,
		"token is a tokens.StartTag type",
		2)
	var ambiguousRootNode = nodes.NewAmbiguousRootNode(_self.Lexer.Tokens[*index])
	*index++
	if ambiguousRootNode.GetIsSelfClosing() {
		return ambiguousRootNode, nil
	}
	for *index < len(_self.Lexer.Tokens) {
		var token = _self.Lexer.Tokens[*index]
		switch token.GetType() {
		case tokens.StartTag:
			var child, err = _self.createR(index)
			if err != nil {
				return nil, err
			}
			ambiguousRootNode.AppendToChildren(child)
		case tokens.EndTag:
			ambiguousRootNode.AppendToChildren(nodes.NewEndElementNode(token, ambiguousRootNode))
			*index++
			return ambiguousRootNode, nil
		case tokens.Comment:
			ambiguousRootNode.AppendToChildren(nodes.NewCommentNode(token))
			*index++
		case tokens.Text:
			ambiguousRootNode.AppendToChildren(nodes.NewTextNode(token))
			*index++
		case tokens.Code:
			ambiguousRootNode.AppendToChildren(nodes.NewDynTextNode(token))
			*index++
		case tokens.EndOfFile:
			*index++
		default:
			return nil, fmt.Errorf("unknown TokenType %q", token.GetType())
		}
	}
	return ambiguousRootNode, nil
}

func (_self *Ast) Process() error {
	verbose.Printf(2, "::: Ast.Process() :::\n")
	var depth = 0
	err := _self.processR(_self.Root, &depth)
	if err != nil {
		return err
	}
	return nil
}

func (_self *Ast) processR(ambiguousNode nodes.Node, depth *int) error {
	utils.Assert(
		ambiguousNode.GetType() == nodes.StartElement ||
			ambiguousNode.GetType() == nodes.Component,
		"current node is either a nodes.StartElement or a nodes.Component type",
		2)
	switch ambiguousNode.GetType() {
	case nodes.StartElement:
		err := _self.StartElementNodeProcessor(ambiguousNode.(*nodes.StartElementNode), depth)
		if err != nil {
			return err
		}
	case nodes.Component:
		err := _self.ComponentNodeProcessor(ambiguousNode.(*nodes.ComponentNode), depth)
		if err != nil {
			return err
		}
	}
	for i := 0; i < len(ambiguousNode.GetChildren()); i++ {
		var node = ambiguousNode.GetChildren()[i]
		switch node.GetType() {
		case nodes.StartElement:
			*depth++
			err := _self.processR(node.(*nodes.StartElementNode), depth)
			if err != nil {
				return err
			}
			if node.GetIsSelfClosing() {
				*depth--
			}
		case nodes.Component:
			*depth++
			err := _self.processR(node.(*nodes.ComponentNode), depth)
			if err != nil {
				return err
			}
			if node.GetIsSelfClosing() {
				*depth--
			}
		case nodes.EndElement:
			err := _self.EndElementNodeProcessor(node.(*nodes.EndElementNode), depth)
			if err != nil {
				return err
			}
			*depth--
		case nodes.Comment:
			err := _self.CommentNodeProcessor(node.(*nodes.CommentNode), depth)
			if err != nil {
				return err
			}
		case nodes.Text:
			err := _self.TextNodeProcessor(node.(*nodes.TextNode), depth)
			if err != nil {
				return err
			}
		case nodes.DynText:
			err := _self.DynTextNodeProcessor(node.(*nodes.DynTextNode), depth)
			if err != nil {
				return err
			}
		case nodes.Attribute:
			err := _self.AttributeNodeProcessor(node.(*nodes.AttributeNode), depth)
			if err != nil {
				return err
			}
		case nodes.ArgumentAttribute:
			err := _self.ArgumentAttributeNodeProcessor(node.(*nodes.ArgumentAttributeNode), depth)
			if err != nil {
				return err
			}
		case nodes.DynAttribute:
			err := _self.DynAttributeNodeProcessor(node.(*nodes.DynAttributeNode), depth)
			if err != nil {
				return err
			}
		case nodes.EventAttribute:
			err := _self.EventAttributeNodeProcessor(node.(*nodes.EventAttributeNode), depth)
			if err != nil {
				return err
			}
		case nodes.KeywordAttribute:
			err := _self.KeywordAttributeNodeProcessor(node.(*nodes.KeywordAttributeNode), depth)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
