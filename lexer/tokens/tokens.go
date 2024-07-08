package tokens

import (
	"github.com/goptos/utils"
)

var verbose = (*utils.Verbose).New(nil)

func runeCount(s string) int {
	var count = 0
	for range s {
		count++
	}
	return count
}

type TokenType string

const (
	StartTag  TokenType = "StartTag"
	EndTag    TokenType = "EndTag"
	Comment   TokenType = "Comment"
	Text      TokenType = "Text"
	Code      TokenType = "Code"
	EndOfFile TokenType = "EndOfFile"
)

type Position struct {
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
}

type AttributeType string

const (
	NormalAttribute   AttributeType = "NormalAttribute"
	EventAttribute    AttributeType = "EventAttribute"
	DynamicAttribute  AttributeType = "DynamicAttribute"
	KeywordAttribute  AttributeType = "KeywordAttribute"
	ArgumentAttribute AttributeType = "ArgumentAttribute"
)

type Attribute struct {
	Type          AttributeType
	NamePosition  Position
	ValuePosition Position
	Name          string
	Value         string
}

type StartTagToken struct {
	_type         TokenType
	position      Position
	name          string
	attributes    []Attribute
	isComponent   bool
	isSelfClosing bool
}

type EndTagToken struct {
	_type    TokenType
	position Position
	name     string
}

type CommentToken struct {
	_type    TokenType
	position Position
	data     string
}

type TextToken struct {
	_type    TokenType
	position Position
	data     string
}

type CodeToken struct {
	_type    TokenType
	position Position
	data     string
}

type EndOfFileToken struct {
	_type    TokenType
	position Position
}

func NewStartTagToken(ln int, ch int) *StartTagToken {
	return &StartTagToken{
		_type: StartTag,
		position: Position{
			StartLine:   ln,
			StartColumn: ch,
			EndLine:     ln,
			EndColumn:   ch},
		name:          "",
		attributes:    []Attribute{},
		isComponent:   false,
		isSelfClosing: false}
}

func NewEndTagToken(ln int, ch int) *EndTagToken {
	return &EndTagToken{
		_type: EndTag,
		position: Position{
			StartLine:   ln,
			StartColumn: ch,
			EndLine:     ln,
			EndColumn:   ch},
		name: ""}
}

func NewCommentToken(ln int, ch int) *CommentToken {
	return &CommentToken{
		_type: Comment,
		position: Position{
			StartLine:   ln,
			StartColumn: ch,
			EndLine:     ln,
			EndColumn:   ch},
		data: ""}
}

func NewTextToken(ln int, ch int) *TextToken {
	return &TextToken{
		_type: Text,
		position: Position{
			StartLine:   ln,
			StartColumn: ch,
			EndLine:     ln,
			EndColumn:   ch},
		data: ""}
}

func NewCodeToken(ln int, ch int) *CodeToken {
	return &CodeToken{
		_type: Code,
		position: Position{
			StartLine:   ln,
			StartColumn: ch,
			EndLine:     ln,
			EndColumn:   ch},
		data: ""}
}

func NewEndOfFileToken(ln int, ch int) *EndOfFileToken {
	return &EndOfFileToken{
		_type: EndOfFile,
		position: Position{
			StartLine:   ln,
			StartColumn: ch,
			EndLine:     ln,
			EndColumn:   ch}}
}

type Token interface {
	GetType() TokenType
	GetPosition() Position
	GetName() string
	GetData() string
	GetAttributes() []Attribute
	GetAttributeType() AttributeType
	GetIsSelfClosing() bool
	GetIsComponent() bool
	SetPosition(Position)
	NewAttribute(int, int)
	AppendToName(string)
	AppendToData(string)
	AppendToAttributeName(string)
	AppendToAttributeValue(string)
	SetAttributeType(AttributeType)
	SetAttributeNamePosition(Position)
	SetAttributeValuePosition(Position)
	SetIsComponent(bool)
	SetIsSelfClosing(bool)
	Print()
}

// GetType()

func (_self *StartTagToken) GetType() TokenType {
	return _self._type
}

func (_self *EndTagToken) GetType() TokenType {
	return _self._type
}

func (_self *CommentToken) GetType() TokenType {
	return _self._type
}

func (_self *TextToken) GetType() TokenType {
	return _self._type
}

func (_self *CodeToken) GetType() TokenType {
	return _self._type
}

func (_self *EndOfFileToken) GetType() TokenType {
	return _self._type
}

// GetPosition()

func (_self *StartTagToken) GetPosition() Position {
	return _self.position
}

func (_self *EndTagToken) GetPosition() Position {
	return _self.position
}

func (_self *CommentToken) GetPosition() Position {
	return _self.position
}

func (_self *TextToken) GetPosition() Position {
	return _self.position
}

func (_self *CodeToken) GetPosition() Position {
	return _self.position
}

func (_self *EndOfFileToken) GetPosition() Position {
	return _self.position
}

// GetName()

func (_self *StartTagToken) GetName() string {
	return _self.name
}

func (_self *EndTagToken) GetName() string {
	return _self.name
}

func (_self *CommentToken) GetName() string {
	utils.Assert(false, "token has name property", 2)
	return ""
}

func (_self *TextToken) GetName() string {
	utils.Assert(false, "token has name property", 2)
	return ""
}

func (_self *CodeToken) GetName() string {
	utils.Assert(false, "token has name property", 2)
	return ""
}

func (_self *EndOfFileToken) GetName() string {
	utils.Assert(false, "token has name property", 2)
	return ""
}

// GetData()

func (_self *StartTagToken) GetData() string {
	utils.Assert(false, "token has data property", 2)
	return ""
}

func (_self *EndTagToken) GetData() string {
	utils.Assert(false, "token has data property", 2)
	return ""
}

func (_self *CommentToken) GetData() string {
	return _self.data
}

func (_self *TextToken) GetData() string {
	return _self.data
}

func (_self *CodeToken) GetData() string {
	return _self.data
}

func (_self *EndOfFileToken) GetData() string {
	utils.Assert(false, "token has data property", 2)
	return ""
}

// GetAttributes()

func (_self *StartTagToken) GetAttributes() []Attribute {
	return _self.attributes
}

func (_self *EndTagToken) GetAttributes() []Attribute {
	utils.Assert(false, "token has attributes property", 2)
	return []Attribute{}
}

func (_self *CommentToken) GetAttributes() []Attribute {
	utils.Assert(false, "token has attributes property", 2)
	return []Attribute{}
}

func (_self *TextToken) GetAttributes() []Attribute {
	utils.Assert(false, "token has attributes property", 2)
	return []Attribute{}
}

func (_self *CodeToken) GetAttributes() []Attribute {
	utils.Assert(false, "token has attributes property", 2)
	return []Attribute{}
}

func (_self *EndOfFileToken) GetAttributes() []Attribute {
	utils.Assert(false, "token has attributes property", 2)
	return []Attribute{}
}

// GetAttributeType()

func (_self *StartTagToken) GetAttributeType() AttributeType {
	return _self.attributes[len(_self.attributes)-1].Type
}

func (_self *EndTagToken) GetAttributeType() AttributeType {
	utils.Assert(false, "token has attributes property", 2)
	return NormalAttribute
}

func (_self *CommentToken) GetAttributeType() AttributeType {
	utils.Assert(false, "token has attributes property", 2)
	return NormalAttribute
}

func (_self *TextToken) GetAttributeType() AttributeType {
	utils.Assert(false, "token has attributes property", 2)
	return NormalAttribute
}

func (_self *CodeToken) GetAttributeType() AttributeType {
	utils.Assert(false, "token has attributes property", 2)
	return NormalAttribute
}

func (_self *EndOfFileToken) GetAttributeType() AttributeType {
	utils.Assert(false, "token has attributes property", 2)
	return NormalAttribute
}

// GetIsSelfClosing()

func (_self *StartTagToken) GetIsSelfClosing() bool {
	return _self.isSelfClosing
}

func (_self *EndTagToken) GetIsSelfClosing() bool {
	utils.Assert(false, "token has isSelfClosing property", 2)
	return false
}

func (_self *CommentToken) GetIsSelfClosing() bool {
	utils.Assert(false, "token has isSelfClosing property", 2)
	return false
}

func (_self *TextToken) GetIsSelfClosing() bool {
	utils.Assert(false, "token has isSelfClosing property", 2)
	return false
}

func (_self *CodeToken) GetIsSelfClosing() bool {
	utils.Assert(false, "token has isSelfClosing property", 2)
	return false
}

func (_self *EndOfFileToken) GetIsSelfClosing() bool {
	utils.Assert(false, "token has isSelfClosing property", 2)
	return false
}

// GetIsComponent()

func (_self *StartTagToken) GetIsComponent() bool {
	return _self.isComponent
}

func (_self *EndTagToken) GetIsComponent() bool {
	utils.Assert(false, "token has isComponent property", 2)
	return false
}

func (_self *CommentToken) GetIsComponent() bool {
	utils.Assert(false, "token has isComponent property", 2)
	return false
}

func (_self *TextToken) GetIsComponent() bool {
	utils.Assert(false, "token has isComponent property", 2)
	return false
}

func (_self *CodeToken) GetIsComponent() bool {
	utils.Assert(false, "token has isComponent property", 2)
	return false
}

func (_self *EndOfFileToken) GetIsComponent() bool {
	utils.Assert(false, "token has isComponent property", 2)
	return false
}

// SetPosition()

func (_self *StartTagToken) SetPosition(position Position) {
	_self.position = position
}

func (_self *EndTagToken) SetPosition(position Position) {
	_self.position = position
}

func (_self *CommentToken) SetPosition(position Position) {
	_self.position = position
}

func (_self *TextToken) SetPosition(position Position) {
	_self.position = position
}

func (_self *CodeToken) SetPosition(position Position) {
	_self.position = position
}

func (_self *EndOfFileToken) SetPosition(position Position) {
	_self.position = position
}

// NewAttribute()

func (_self *StartTagToken) NewAttribute(ln int, ch int) {
	_self.attributes = append(_self.attributes, Attribute{
		Type: NormalAttribute,
		NamePosition: Position{
			StartLine:   ln,
			StartColumn: ch,
			EndLine:     ln,
			EndColumn:   ch - 1,
		},
		ValuePosition: Position{
			StartLine:   ln,
			StartColumn: ch + 2,
			EndLine:     ln,
			EndColumn:   ch + 1,
		},
		Name:  "",
		Value: ""})
}

func (_self *EndTagToken) NewAttribute(ln int, ch int) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *CommentToken) NewAttribute(ln int, ch int) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *TextToken) NewAttribute(ln int, ch int) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *CodeToken) NewAttribute(ln int, ch int) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *EndOfFileToken) NewAttribute(ln int, ch int) {
	utils.Assert(false, "token has data property", 2)
}

// AppendToName()

func (_self *StartTagToken) AppendToName(s string) {
	_self.name = _self.name + s
}

func (_self *EndTagToken) AppendToName(s string) {
	_self.name = _self.name + s
}

func (_self *CommentToken) AppendToName(s string) {
	utils.Assert(false, "token has name property", 2)
}

func (_self *TextToken) AppendToName(s string) {
	utils.Assert(false, "token has name property", 2)
}

func (_self *CodeToken) AppendToName(s string) {
	utils.Assert(false, "token has name property", 2)
}

func (_self *EndOfFileToken) AppendToName(s string) {
	utils.Assert(false, "token has name property", 2)
}

// AppendToData()

func (_self *StartTagToken) AppendToData(s string) {
	utils.Assert(false, "token has data property", 2)
}

func (_self *EndTagToken) AppendToData(s string) {
	utils.Assert(false, "token has data property", 2)
}

func (_self *CommentToken) AppendToData(s string) {
	_self.data = _self.data + s
}

func (_self *TextToken) AppendToData(s string) {
	_self.data = _self.data + s
}

func (_self *CodeToken) AppendToData(s string) {
	_self.data = _self.data + s
}

func (_self *EndOfFileToken) AppendToData(s string) {
	utils.Assert(false, "token has data property", 2)
}

// AppendToAttributeName()

func (_self *StartTagToken) AppendToAttributeName(s string) {
	utils.Assert(runeCount(s) == 1, "AppendToAttributeName(s string) where s is only 1 character", 2)
	_self.attributes[len(_self.attributes)-1].Name = _self.attributes[len(_self.attributes)-1].Name + s
	_self.attributes[len(_self.attributes)-1].NamePosition.EndColumn++
	_self.attributes[len(_self.attributes)-1].ValuePosition.StartColumn++
	_self.attributes[len(_self.attributes)-1].ValuePosition.EndColumn++
}

func (_self *EndTagToken) AppendToAttributeName(s string) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *CommentToken) AppendToAttributeName(s string) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *TextToken) AppendToAttributeName(s string) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *CodeToken) AppendToAttributeName(s string) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *EndOfFileToken) AppendToAttributeName(s string) {
	utils.Assert(false, "token has attributes property", 2)
}

// AppendToAttributeValue()

func (_self *StartTagToken) AppendToAttributeValue(s string) {
	utils.Assert(runeCount(s) == 1, "AppendToAttributeName(s string) where s is only 1 character", 2)
	_self.attributes[len(_self.attributes)-1].Value = _self.attributes[len(_self.attributes)-1].Value + s
	_self.attributes[len(_self.attributes)-1].ValuePosition.EndColumn++
}

func (_self *EndTagToken) AppendToAttributeValue(s string) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *CommentToken) AppendToAttributeValue(s string) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *TextToken) AppendToAttributeValue(s string) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *CodeToken) AppendToAttributeValue(s string) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *EndOfFileToken) AppendToAttributeValue(s string) {
	utils.Assert(false, "token has attributes property", 2)
}

// SetAttributeType()

func (_self *StartTagToken) SetAttributeType(t AttributeType) {
	_self.attributes[len(_self.attributes)-1].Type = t
	switch t {
	case ArgumentAttribute:
		_self.attributes[len(_self.attributes)-1].ValuePosition.StartLine = 0
		_self.attributes[len(_self.attributes)-1].ValuePosition.StartColumn = 0
		_self.attributes[len(_self.attributes)-1].ValuePosition.EndLine = 0
		_self.attributes[len(_self.attributes)-1].ValuePosition.EndColumn = 0
	default:
		_self.attributes[len(_self.attributes)-1].ValuePosition.StartColumn--
		_self.attributes[len(_self.attributes)-1].ValuePosition.EndColumn++
	}
}

func (_self *EndTagToken) SetAttributeType(t AttributeType) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *CommentToken) SetAttributeType(t AttributeType) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *TextToken) SetAttributeType(t AttributeType) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *CodeToken) SetAttributeType(t AttributeType) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *EndOfFileToken) SetAttributeType(t AttributeType) {
	utils.Assert(false, "token has attributes property", 2)
}

// SetAttributeNamePosition()

func (_self *StartTagToken) SetAttributeNamePosition(position Position) {
	_self.attributes[len(_self.attributes)-1].NamePosition = position
}

func (_self *EndTagToken) SetAttributeNamePosition(position Position) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *CommentToken) SetAttributeNamePosition(position Position) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *TextToken) SetAttributeNamePosition(position Position) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *CodeToken) SetAttributeNamePosition(position Position) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *EndOfFileToken) SetAttributeNamePosition(position Position) {
	utils.Assert(false, "token has data property", 2)
}

// SetAttributeValuePosition()

func (_self *StartTagToken) SetAttributeValuePosition(position Position) {
	_self.attributes[len(_self.attributes)-1].ValuePosition = position
}

func (_self *EndTagToken) SetAttributeValuePosition(position Position) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *CommentToken) SetAttributeValuePosition(position Position) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *TextToken) SetAttributeValuePosition(position Position) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *CodeToken) SetAttributeValuePosition(position Position) {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *EndOfFileToken) SetAttributeValuePosition(position Position) {
	utils.Assert(false, "token has data property", 2)
}

// SetIsComponent()

func (_self *StartTagToken) SetIsComponent(b bool) {
	_self.isComponent = b
}

func (_self *EndTagToken) SetIsComponent(b bool) {
	utils.Assert(false, "token has isComponent property", 2)
}

func (_self *CommentToken) SetIsComponent(b bool) {
	utils.Assert(false, "token has isComponent property", 2)
}

func (_self *TextToken) SetIsComponent(b bool) {
	utils.Assert(false, "token has isComponent property", 2)
}

func (_self *CodeToken) SetIsComponent(b bool) {
	utils.Assert(false, "token has isComponent property", 2)
}

func (_self *EndOfFileToken) SetIsComponent(b bool) {
	utils.Assert(false, "token has isComponent property", 2)
}

// SetIsSelfClosing()

func (_self *StartTagToken) SetIsSelfClosing(b bool) {
	_self.isSelfClosing = b
}

func (_self *EndTagToken) SetIsSelfClosing(b bool) {
	utils.Assert(false, "token has isSelfClosing property", 2)
}

func (_self *CommentToken) SetIsSelfClosing(b bool) {
	utils.Assert(false, "token has isSelfClosing property", 2)
}

func (_self *TextToken) SetIsSelfClosing(b bool) {
	utils.Assert(false, "token has isSelfClosing property", 2)
}

func (_self *CodeToken) SetIsSelfClosing(b bool) {
	utils.Assert(false, "token has isSelfClosing property", 2)
}

func (_self *EndOfFileToken) SetIsSelfClosing(b bool) {
	utils.Assert(false, "token has isSelfClosing property", 2)
}

// Print()

func (_self *StartTagToken) Print() {
	var component string
	if _self.isComponent {
		component = "(Component)"
	}
	var selfClosing string
	if _self.isSelfClosing {
		selfClosing = "(SelfClosing)"
	}
	verbose.Printf(0, "%v", _self.position)
	verbose.Printf(0, "%s\t%s\t%d Attributes\t%s\t%s\n",
		_self._type,
		_self.name,
		len(_self.attributes),
		component,
		selfClosing)
	for i, attribute := range _self.attributes {
		switch attribute.Type {
		case ArgumentAttribute:
			verbose.Printf(0, " attribute%d\t%s\t%v%s\n", i, attribute.Type, attribute.NamePosition, attribute.Name)
		case NormalAttribute:
			verbose.Printf(0, " attribute%d\t%s \t%v%s\t%v%s\n", i, attribute.Type, attribute.NamePosition, attribute.Name, attribute.ValuePosition, attribute.Value)
		case EventAttribute:
			verbose.Printf(0, " attribute%d\t%s  \t%v%s\t%v{%s}\n", i, attribute.Type, attribute.NamePosition, attribute.Name, attribute.ValuePosition, attribute.Value)
		default:
			verbose.Printf(0, " attribute%d\t%s\t%v%s\t%v{%s}\n", i, attribute.Type, attribute.NamePosition, attribute.Name, attribute.ValuePosition, attribute.Value)
		}
	}
}

func (_self *EndTagToken) Print() {
	verbose.Printf(0, "%v", _self.position)
	verbose.Printf(0, "%s\t%s\n", _self._type, _self.name)
}

func (_self *CommentToken) Print() {
	verbose.Printf(0, "%v", _self.position)
	verbose.Printf(0, "%s\t%s\n", _self._type, _self.data)
}

func (_self *TextToken) Print() {
	verbose.Printf(0, "%v", _self.position)
	verbose.Printf(0, "%s   \t%s\n", _self._type, _self.data)
}

func (_self *CodeToken) Print() {
	verbose.Printf(0, "%v", _self.position)
	verbose.Printf(0, "%s\t%s\n", _self._type, _self.data)
}

func (_self *EndOfFileToken) Print() {
	verbose.Printf(0, "%v", _self.position)
	verbose.Printf(0, "%s\n", _self._type)
}
