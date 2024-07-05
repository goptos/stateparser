package nodes

import (
	"strings"

	"github.com/goptos/stateparser/lexer/tokens"
	"github.com/goptos/utils"
)

var verbose = (*utils.Verbose).New(nil)

type NodeType string

const (
	StartElement      NodeType = "StartElement"
	Component         NodeType = "Component"
	EndElement        NodeType = "EndElement"
	Comment           NodeType = "Comment"
	Text              NodeType = "Text"
	DynText           NodeType = "DynText"
	Attribute         NodeType = "Attribute"
	ArgumentAttribute NodeType = "ArgumentAttribute"
	DynAttribute      NodeType = "DynAttribute"
	EventAttribute    NodeType = "EventAttribute"
	KeywordAttribute  NodeType = "KeywordAttribute"
)

type StartElementNode struct {
	_type         NodeType
	name          string
	children      []Node
	isSelfClosing bool
}

type ComponentNode struct {
	_type         NodeType
	name          string
	children      []Node
	isSelfClosing bool
}

type EndElementNode struct {
	_type         NodeType
	name          string
	startElemNode Node
}

type CommentNode struct {
	_type NodeType
	data  string
}

type TextNode struct {
	_type NodeType
	data  string
}

type DynTextNode struct {
	_type  NodeType
	effect string
}

type AttributeNode struct {
	_type NodeType
	name  string
	value string
}

type ArgumentAttributeNode struct {
	_type NodeType
	name  string
}

type DynAttributeNode struct {
	_type  NodeType
	name   string
	value  string
	effect string
}

type EventAttributeNode struct {
	_type  NodeType
	name   string
	event  string
	effect string
}

type KeywordAttributeNode struct {
	_type  NodeType
	name   string
	effect string
}

func NewAmbiguousRootNode(token tokens.Token) Node {
	if token.GetIsComponent() {
		return NewComponentNode(token)
	}
	return NewStartElementNode(token)
}

func NewStartElementNode(token tokens.Token) *StartElementNode {
	var children = []Node{}
	for _, attribute := range token.GetAttributes() {
		switch attribute.Type {
		case tokens.ArgumentAttribute:
			children = append(children, NewArgumentAttributeNode(&attribute))
		case tokens.EventAttribute:
			children = append(children, NewEventAttributeNode(&attribute))
		case tokens.DynamicAttribute:
			children = append(children, NewDynAttributeNode(&attribute))
		case tokens.KeywordAttribute:
			children = append(children, NewKeywordAttributeNode(&attribute))
		default:
			children = append(children, NewAttributeNode(&attribute))
		}
	}
	return &StartElementNode{
		_type:         StartElement,
		name:          token.GetName(),
		children:      children,
		isSelfClosing: token.GetIsSelfClosing()}
}

func NewComponentNode(token tokens.Token) *ComponentNode {
	var children = []Node{}
	for _, attribute := range token.GetAttributes() {
		switch attribute.Type {
		case tokens.EventAttribute:
			children = append(children, NewEventAttributeNode(&attribute))
		case tokens.DynamicAttribute:
			children = append(children, NewDynAttributeNode(&attribute))
		case tokens.KeywordAttribute:
			children = append(children, NewKeywordAttributeNode(&attribute))
		default:
			children = append(children, NewAttributeNode(&attribute))
		}
	}
	return &ComponentNode{
		_type:         Component,
		name:          token.GetName(),
		children:      children,
		isSelfClosing: token.GetIsSelfClosing()}
}

func NewEndElementNode(token tokens.Token, node Node) *EndElementNode {
	return &EndElementNode{
		_type:         EndElement,
		name:          token.GetName(),
		startElemNode: node}
}

func NewCommentNode(token tokens.Token) *CommentNode {
	return &CommentNode{
		_type: Comment,
		data:  token.GetData()}
}

func NewTextNode(token tokens.Token) *TextNode {
	return &TextNode{
		_type: Text,
		data:  token.GetData()}
}

func NewDynTextNode(token tokens.Token) *DynTextNode {
	return &DynTextNode{
		_type:  DynText,
		effect: token.GetData()}
}

func NewAttributeNode(attribute *tokens.Attribute) *AttributeNode {
	return &AttributeNode{
		_type: Attribute,
		name:  attribute.Name,
		value: attribute.Value,
	}
}

func NewArgumentAttributeNode(attribute *tokens.Attribute) *ArgumentAttributeNode {
	return &ArgumentAttributeNode{
		_type: Attribute,
		name:  attribute.Name,
	}
}

func NewDynAttributeNode(attribute *tokens.Attribute) *DynAttributeNode {
	return &DynAttributeNode{
		_type:  DynAttribute,
		name:   strings.Split(attribute.Name, ":")[0],
		value:  strings.Split(attribute.Name, ":")[1],
		effect: attribute.Value,
	}
}

func NewEventAttributeNode(attribute *tokens.Attribute) *EventAttributeNode {
	return &EventAttributeNode{
		_type:  EventAttribute,
		name:   strings.Split(attribute.Name, ":")[0],
		event:  strings.Split(attribute.Name, ":")[1],
		effect: attribute.Value,
	}
}

func NewKeywordAttributeNode(attribute *tokens.Attribute) *KeywordAttributeNode {
	return &KeywordAttributeNode{
		_type:  KeywordAttribute,
		name:   attribute.Name,
		effect: attribute.Value,
	}
}

type Node interface {
	GetType() NodeType
	GetName() string
	GetChildren() []Node
	GetStartElementNode() Node
	GetData() string
	GetEffect() string
	GetValue() string
	GetEvent() string
	GetIsSelfClosing() bool
	AppendToChildren(Node)
	Print(*int) error
}

// GetType()

func (_self *StartElementNode) GetType() NodeType {
	return _self._type
}

func (_self *ComponentNode) GetType() NodeType {
	return _self._type
}

func (_self *EndElementNode) GetType() NodeType {
	return _self._type
}

func (_self *CommentNode) GetType() NodeType {
	return _self._type
}

func (_self *TextNode) GetType() NodeType {
	return _self._type
}

func (_self *DynTextNode) GetType() NodeType {
	return _self._type
}

func (_self *AttributeNode) GetType() NodeType {
	return _self._type
}

func (_self *ArgumentAttributeNode) GetType() NodeType {
	return _self._type
}

func (_self *DynAttributeNode) GetType() NodeType {
	return _self._type
}

func (_self *EventAttributeNode) GetType() NodeType {
	return _self._type
}

func (_self *KeywordAttributeNode) GetType() NodeType {
	return _self._type
}

// GetName()

func (_self *StartElementNode) GetName() string {
	return _self.name
}

func (_self *ComponentNode) GetName() string {
	return _self.name
}

func (_self *EndElementNode) GetName() string {
	return _self.name
}

func (_self *CommentNode) GetName() string {
	utils.Assert(false, "token has name property", 2)
	return ""
}

func (_self *TextNode) GetName() string {
	utils.Assert(false, "token has name property", 2)
	return ""
}

func (_self *DynTextNode) GetName() string {
	utils.Assert(false, "token has name property", 2)
	return ""
}

func (_self *AttributeNode) GetName() string {
	return _self.name
}

func (_self *ArgumentAttributeNode) GetName() string {
	return _self.name
}

func (_self *DynAttributeNode) GetName() string {
	return _self.name
}

func (_self *EventAttributeNode) GetName() string {
	return _self.name
}

func (_self *KeywordAttributeNode) GetName() string {
	return _self.name
}

// GetChildren()

func (_self *StartElementNode) GetChildren() []Node {
	return _self.children
}

func (_self *ComponentNode) GetChildren() []Node {
	return _self.children
}

func (_self *EndElementNode) GetChildren() []Node {
	utils.Assert(false, "token has children property", 2)
	return []Node{}
}

func (_self *CommentNode) GetChildren() []Node {
	utils.Assert(false, "token has children property", 2)
	return []Node{}
}

func (_self *TextNode) GetChildren() []Node {
	utils.Assert(false, "token has children property", 2)
	return []Node{}
}

func (_self *DynTextNode) GetChildren() []Node {
	utils.Assert(false, "token has children property", 2)
	return []Node{}
}

func (_self *AttributeNode) GetChildren() []Node {
	utils.Assert(false, "token has children property", 2)
	return []Node{}
}

func (_self *ArgumentAttributeNode) GetChildren() []Node {
	utils.Assert(false, "token has children property", 2)
	return []Node{}
}

func (_self *DynAttributeNode) GetChildren() []Node {
	utils.Assert(false, "token has children property", 2)
	return []Node{}
}

func (_self *EventAttributeNode) GetChildren() []Node {
	utils.Assert(false, "token has children property", 2)
	return []Node{}
}

func (_self *KeywordAttributeNode) GetChildren() []Node {
	utils.Assert(false, "token has children property", 2)
	return []Node{}
}

// GetStartElementNode()

func (_self *StartElementNode) GetStartElementNode() Node {
	utils.Assert(false, "token has startElementNode property", 2)
	return &StartElementNode{}
}

func (_self *ComponentNode) GetStartElementNode() Node {
	utils.Assert(false, "token has startElementNode property", 2)
	return &StartElementNode{}
}

func (_self *EndElementNode) GetStartElementNode() Node {
	return _self.startElemNode
}

func (_self *CommentNode) GetStartElementNode() Node {
	utils.Assert(false, "token has startElementNode property", 2)
	return &StartElementNode{}
}

func (_self *TextNode) GetStartElementNode() Node {
	utils.Assert(false, "token has startElementNode property", 2)
	return &StartElementNode{}
}

func (_self *DynTextNode) GetStartElementNode() Node {
	utils.Assert(false, "token has startElementNode property", 2)
	return &StartElementNode{}
}

func (_self *AttributeNode) GetStartElementNode() Node {
	utils.Assert(false, "token has startElementNode property", 2)
	return &StartElementNode{}
}

func (_self *ArgumentAttributeNode) GetStartElementNode() Node {
	utils.Assert(false, "token has startElementNode property", 2)
	return &StartElementNode{}
}

func (_self *DynAttributeNode) GetStartElementNode() Node {
	utils.Assert(false, "token has startElementNode property", 2)
	return &StartElementNode{}
}

func (_self *EventAttributeNode) GetStartElementNode() Node {
	utils.Assert(false, "token has startElementNode property", 2)
	return &StartElementNode{}
}

func (_self *KeywordAttributeNode) GetStartElementNode() Node {
	utils.Assert(false, "token has startElementNode property", 2)
	return &StartElementNode{}
}

// GetData()

func (_self *StartElementNode) GetData() string {
	utils.Assert(false, "token has data property", 2)
	return ""
}

func (_self *ComponentNode) GetData() string {
	utils.Assert(false, "token has data property", 2)
	return ""
}

func (_self *EndElementNode) GetData() string {
	utils.Assert(false, "token has data property", 2)
	return ""
}

func (_self *CommentNode) GetData() string {
	return _self.data
}

func (_self *TextNode) GetData() string {
	return _self.data
}

func (_self *DynTextNode) GetData() string {
	utils.Assert(false, "token has data property", 2)
	return ""
}

func (_self *AttributeNode) GetData() string {
	utils.Assert(false, "token has data property", 2)
	return ""
}

func (_self *ArgumentAttributeNode) GetData() string {
	utils.Assert(false, "token has data property", 2)
	return ""
}

func (_self *DynAttributeNode) GetData() string {
	utils.Assert(false, "token has data property", 2)
	return ""
}

func (_self *EventAttributeNode) GetData() string {
	utils.Assert(false, "token has data property", 2)
	return ""
}

func (_self *KeywordAttributeNode) GetData() string {
	utils.Assert(false, "token has data property", 2)
	return ""
}

// GetEffect()

func (_self *StartElementNode) GetEffect() string {
	utils.Assert(false, "token has effect property", 2)
	return ""
}

func (_self *ComponentNode) GetEffect() string {
	utils.Assert(false, "token has effect property", 2)
	return ""
}

func (_self *EndElementNode) GetEffect() string {
	utils.Assert(false, "token has effect property", 2)
	return ""
}

func (_self *CommentNode) GetEffect() string {
	utils.Assert(false, "token has effect property", 2)
	return ""
}

func (_self *TextNode) GetEffect() string {
	utils.Assert(false, "token has effect property", 2)
	return ""
}

func (_self *DynTextNode) GetEffect() string {
	return _self.effect
}

func (_self *AttributeNode) GetEffect() string {
	utils.Assert(false, "token has effect property", 2)
	return ""
}

func (_self *ArgumentAttributeNode) GetEffect() string {
	utils.Assert(false, "token has effect property", 2)
	return ""
}

func (_self *DynAttributeNode) GetEffect() string {
	return _self.effect
}

func (_self *EventAttributeNode) GetEffect() string {
	return _self.effect
}

func (_self *KeywordAttributeNode) GetEffect() string {
	return _self.effect
}

// GetValue()

func (_self *StartElementNode) GetValue() string {
	utils.Assert(false, "token has value property", 2)
	return ""
}

func (_self *ComponentNode) GetValue() string {
	utils.Assert(false, "token has value property", 2)
	return ""
}

func (_self *EndElementNode) GetValue() string {
	utils.Assert(false, "token has value property", 2)
	return ""
}

func (_self *CommentNode) GetValue() string {
	utils.Assert(false, "token has value property", 2)
	return ""
}

func (_self *TextNode) GetValue() string {
	utils.Assert(false, "token has value property", 2)
	return ""
}

func (_self *DynTextNode) GetValue() string {
	utils.Assert(false, "token has value property", 2)
	return ""
}

func (_self *AttributeNode) GetValue() string {
	return _self.value
}

func (_self *ArgumentAttributeNode) GetValue() string {
	utils.Assert(false, "token has value property", 2)
	return ""
}

func (_self *DynAttributeNode) GetValue() string {
	return _self.value
}

func (_self *EventAttributeNode) GetValue() string {
	utils.Assert(false, "token has value property", 2)
	return ""
}

func (_self *KeywordAttributeNode) GetValue() string {
	utils.Assert(false, "token has value property", 2)
	return ""
}

// GetEvent()

func (_self *StartElementNode) GetEvent() string {
	utils.Assert(false, "token has event property", 2)
	return ""
}

func (_self *ComponentNode) GetEvent() string {
	utils.Assert(false, "token has event property", 2)
	return ""
}

func (_self *EndElementNode) GetEvent() string {
	utils.Assert(false, "token has event property", 2)
	return ""
}

func (_self *CommentNode) GetEvent() string {
	utils.Assert(false, "token has event property", 2)
	return ""
}

func (_self *TextNode) GetEvent() string {
	utils.Assert(false, "token has event property", 2)
	return ""
}

func (_self *DynTextNode) GetEvent() string {
	utils.Assert(false, "token has event property", 2)
	return ""
}

func (_self *AttributeNode) GetEvent() string {
	utils.Assert(false, "token has event property", 2)
	return ""
}

func (_self *ArgumentAttributeNode) GetEvent() string {
	utils.Assert(false, "token has event property", 2)
	return ""
}

func (_self *DynAttributeNode) GetEvent() string {
	utils.Assert(false, "token has event property", 2)
	return ""
}

func (_self *EventAttributeNode) GetEvent() string {
	return _self.event
}

func (_self *KeywordAttributeNode) GetEvent() string {
	utils.Assert(false, "token has event property", 2)
	return ""
}

// GetIsSelfClosing

func (_self *StartElementNode) GetIsSelfClosing() bool {
	return _self.isSelfClosing
}

func (_self *ComponentNode) GetIsSelfClosing() bool {
	return _self.isSelfClosing
}

func (_self *EndElementNode) GetIsSelfClosing() bool {
	utils.Assert(false, "token has isSelfClosing property", 2)
	return false
}

func (_self *CommentNode) GetIsSelfClosing() bool {
	utils.Assert(false, "token has isSelfClosing property", 2)
	return false
}

func (_self *TextNode) GetIsSelfClosing() bool {
	utils.Assert(false, "token has isSelfClosing property", 2)
	return false
}

func (_self *DynTextNode) GetIsSelfClosing() bool {
	utils.Assert(false, "token has isSelfClosing property", 2)
	return false
}

func (_self *AttributeNode) GetIsSelfClosing() bool {
	utils.Assert(false, "token has isSelfClosing property", 2)
	return false
}

func (_self *ArgumentAttributeNode) GetIsSelfClosing() bool {
	utils.Assert(false, "token has isSelfClosing property", 2)
	return false
}

func (_self *DynAttributeNode) GetIsSelfClosing() bool {
	utils.Assert(false, "token has isSelfClosing property", 2)
	return false
}

func (_self *EventAttributeNode) GetIsSelfClosing() bool {
	utils.Assert(false, "token has isSelfClosing property", 2)
	return false
}

func (_self *KeywordAttributeNode) GetIsSelfClosing() bool {
	utils.Assert(false, "token has isSelfClosing property", 2)
	return false
}

// AppendToChildren()

func (_self *StartElementNode) AppendToChildren(n Node) {
	_self.children = append(_self.children, n)
}

func (_self *ComponentNode) AppendToChildren(n Node) {
	_self.children = append(_self.children, n)
}

func (_self *EndElementNode) AppendToChildren(n Node) {
	utils.Assert(false, "token has children property", 2)
}

func (_self *CommentNode) AppendToChildren(n Node) {
	utils.Assert(false, "token has children property", 2)
}

func (_self *TextNode) AppendToChildren(n Node) {
	utils.Assert(false, "token has children property", 2)
}

func (_self *DynTextNode) AppendToChildren(n Node) {
	utils.Assert(false, "token has children property", 2)
}

func (_self *AttributeNode) AppendToChildren(n Node) {
	utils.Assert(false, "token has children property", 2)
}

func (_self *ArgumentAttributeNode) AppendToChildren(n Node) {
	utils.Assert(false, "token has children property", 2)
}

func (_self *DynAttributeNode) AppendToChildren(n Node) {
	utils.Assert(false, "token has children property", 2)
}

func (_self *EventAttributeNode) AppendToChildren(n Node) {
	utils.Assert(false, "token has children property", 2)
}

func (_self *KeywordAttributeNode) AppendToChildren(n Node) {
	utils.Assert(false, "token has children property", 2)
}

// Print()

func (_self *StartElementNode) Print(depth *int) error {
	var indent = ""
	for i := 0; i < *depth; i++ {
		indent = indent + " "
	}
	var selfClosing string
	if _self.isSelfClosing {
		selfClosing = "(SelfClosing)"
	}
	verbose.Printf(2, indent+"%s    %s    %d Children    %s\n",
		_self._type,
		_self.name,
		len(_self.children),
		selfClosing)
	return nil
}

func (_self *ComponentNode) Print(depth *int) error {
	var indent = ""
	for i := 0; i < *depth; i++ {
		indent = indent + " "
	}
	var selfClosing string
	if _self.isSelfClosing {
		selfClosing = "(SelfClosing)"
	}
	verbose.Printf(2, indent+"%s    %s    %d Children    %s    %s\n",
		_self._type,
		_self.name,
		len(_self.children),
		"(Component)",
		selfClosing)
	return nil
}

func (_self *EndElementNode) Print(depth *int) error {
	var indent = ""
	for i := 0; i <= *depth; i++ {
		indent = indent + " "
	}
	verbose.Printf(2, indent+"%s    %s    %s\n",
		_self._type,
		_self.name,
		_self.startElemNode.GetName())
	return nil
}

func (_self *CommentNode) Print(depth *int) error {
	var indent = ""
	for i := 0; i <= *depth; i++ {
		indent = indent + " "
	}
	verbose.Printf(2, indent+"%s    %s\n",
		_self._type,
		_self.data)
	return nil
}

func (_self *TextNode) Print(depth *int) error {
	var indent = ""
	for i := 0; i <= *depth; i++ {
		indent = indent + " "
	}
	verbose.Printf(2, indent+"%s    %s\n",
		_self._type,
		_self.data)
	return nil
}

func (_self *DynTextNode) Print(depth *int) error {
	var indent = ""
	for i := 0; i <= *depth; i++ {
		indent = indent + " "
	}
	verbose.Printf(2, indent+"%s    {%s}\n",
		_self._type,
		_self.effect)
	return nil
}

func (_self *AttributeNode) Print(depth *int) error {
	var indent = ""
	for i := 0; i <= *depth; i++ {
		indent = indent + " "
	}
	verbose.Printf(2, indent+"%s    %s    %s\n",
		_self._type,
		_self.name,
		_self.value)
	return nil
}

func (_self *ArgumentAttributeNode) Print(depth *int) error {
	var indent = ""
	for i := 0; i <= *depth; i++ {
		indent = indent + " "
	}
	verbose.Printf(2, indent+"%s    %s\n",
		_self._type,
		_self.name)
	return nil
}

func (_self *DynAttributeNode) Print(depth *int) error {
	var indent = ""
	for i := 0; i <= *depth; i++ {
		indent = indent + " "
	}
	verbose.Printf(2, indent+"%s    %s    %s    {%s}\n",
		_self._type,
		_self.name,
		_self.value,
		_self.effect)
	return nil
}

func (_self *EventAttributeNode) Print(depth *int) error {
	var indent = ""
	for i := 0; i <= *depth; i++ {
		indent = indent + " "
	}
	verbose.Printf(2, indent+"%s    %s    %s    {%s}\n",
		_self._type,
		_self.name,
		_self.event,
		_self.effect)
	return nil
}

func (_self *KeywordAttributeNode) Print(depth *int) error {
	var indent = ""
	for i := 0; i <= *depth; i++ {
		indent = indent + " "
	}
	verbose.Printf(2, indent+"%s    %s    {%s}\n",
		_self._type,
		_self.name,
		_self.effect)
	return nil
}
