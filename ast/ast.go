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
	KeywordAttributeNames         map[string]interface{}
	Lexer                         *lexer.Lexer
	Root                          *nodes.StartElementNode
	StartElementNodeProcessor     func(node *nodes.StartElementNode, depth *int) error
	EndElementNodeProcessor       func(node *nodes.EndElementNode, depth *int) error
	CommentNodeProcessor          func(node *nodes.CommentNode, depth *int) error
	TextNodeProcessor             func(node *nodes.TextNode, depth *int) error
	DynTextNodeProcessor          func(node *nodes.DynTextNode, depth *int) error
	AttributeNodeProcessor        func(node *nodes.AttributeNode, depth *int) error
	DynAttributeNodeProcessor     func(node *nodes.DynAttributeNode, depth *int) error
	EventAttributeNodeProcessor   func(node *nodes.EventAttributeNode, depth *int) error
	KeywordAttributeNodeProcessor func(node *nodes.KeywordAttributeNode, depth *int) error
}

func New(source string) *Ast {
	return &Ast{
		KeywordAttributeNames:         make(map[string]interface{}),
		Lexer:                         lexer.New(source),
		Root:                          nil,
		StartElementNodeProcessor:     (*nodes.StartElementNode).Print,
		EndElementNodeProcessor:       (*nodes.EndElementNode).Print,
		CommentNodeProcessor:          (*nodes.CommentNode).Print,
		TextNodeProcessor:             (*nodes.TextNode).Print,
		DynTextNodeProcessor:          (*nodes.DynTextNode).Print,
		AttributeNodeProcessor:        (*nodes.AttributeNode).Print,
		DynAttributeNodeProcessor:     (*nodes.DynAttributeNode).Print,
		EventAttributeNodeProcessor:   (*nodes.EventAttributeNode).Print,
		KeywordAttributeNodeProcessor: (*nodes.KeywordAttributeNode).Print}
}

func (_self *Ast) Create() error {
	err := _self.Lexer.Tokenise()
	if err != nil {
		return err
	}
	verbose.Printf(4, "::: Ast.Create() :::\n")
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

func (_self *Ast) createR(index *int) (*nodes.StartElementNode, error) {
	verbose.Printf(4, "%d\t%s\n", *index, _self.Lexer.Tokens[*index].GetName())
	var startTagNode = nodes.NewStartElementNode(_self.Lexer.Tokens[*index], _self.KeywordAttributeNames)
	*index++
	if startTagNode.GetIsSelfClosing() {
		return startTagNode, nil
	}
	for *index < len(_self.Lexer.Tokens) {
		var token = _self.Lexer.Tokens[*index]
		switch token.GetType() {
		case tokens.StartTag:
			var child, err = _self.createR(index)
			if err != nil {
				return nil, err
			}
			startTagNode.AppendToChildren(child)
		case tokens.EndTag:
			startTagNode.AppendToChildren(nodes.NewEndElementNode(token, *startTagNode))
			*index++
			return startTagNode, nil
		case tokens.Comment:
			startTagNode.AppendToChildren(nodes.NewCommentNode(token))
			*index++
		case tokens.Text:
			startTagNode.AppendToChildren(nodes.NewTextNode(token))
			*index++
		case tokens.Code:
			startTagNode.AppendToChildren(nodes.NewDynTextNode(token))
			*index++
		case tokens.EndOfFile:
			*index++
		default:
			return nil, fmt.Errorf("unknown TokenType %q", token.GetType())
		}
	}
	return startTagNode, nil
}

func (_self *Ast) Process() error {
	verbose.Printf(4, "::: Ast.Process() :::\n")
	var depth = 0
	err := _self.processR(_self.Root, &depth)
	if err != nil {
		return err
	}
	return nil
}

func (_self *Ast) processR(startElementNode *nodes.StartElementNode, depth *int) error {
	err := _self.StartElementNodeProcessor(startElementNode, depth)
	if err != nil {
		return err
	}
	for i := 0; i < len(startElementNode.GetChildren()); i++ {
		var node = startElementNode.GetChildren()[i]
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
		case nodes.EndElement:
			err = _self.EndElementNodeProcessor(node.(*nodes.EndElementNode), depth)
			if err != nil {
				return err
			}
			*depth--
		case nodes.Comment:
			err = _self.CommentNodeProcessor(node.(*nodes.CommentNode), depth)
			if err != nil {
				return err
			}
		case nodes.Text:
			err = _self.TextNodeProcessor(node.(*nodes.TextNode), depth)
			if err != nil {
				return err
			}
		case nodes.DynText:
			err = _self.DynTextNodeProcessor(node.(*nodes.DynTextNode), depth)
			if err != nil {
				return err
			}
		case nodes.Attribute:
			err = _self.AttributeNodeProcessor(node.(*nodes.AttributeNode), depth)
			if err != nil {
				return err
			}
		case nodes.DynAttribute:
			err = _self.DynAttributeNodeProcessor(node.(*nodes.DynAttributeNode), depth)
			if err != nil {
				return err
			}
		case nodes.EventAttribute:
			err = _self.EventAttributeNodeProcessor(node.(*nodes.EventAttributeNode), depth)
			if err != nil {
				return err
			}
		case nodes.KeywordAttribute:
			err = _self.KeywordAttributeNodeProcessor(node.(*nodes.KeywordAttributeNode), depth)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
