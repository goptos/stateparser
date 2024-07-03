package tokens

import (
	"github.com/goptos/utils"
)

var verbose = (*utils.Verbose).New(nil)

type TokenType string

const (
	StartTag  TokenType = "StartTag"
	EndTag    TokenType = "EndTag"
	Comment   TokenType = "Comment"
	Text      TokenType = "Text"
	Code      TokenType = "Code"
	EndOfFile TokenType = "EndOfFile"
)

type Attribute struct {
	Name  string
	Value string
}

type StartTagToken struct {
	_type         TokenType
	name          string
	attributes    []Attribute
	isComponent   bool
	isSelfClosing bool
}

type EndTagToken struct {
	_type TokenType
	name  string
}

type CommentToken struct {
	_type TokenType
	data  string
}

type TextToken struct {
	_type TokenType
	data  string
}

type CodeToken struct {
	_type TokenType
	data  string
}

type EndOfFileToken struct {
	_type TokenType
}

func NewStartTagToken() *StartTagToken {
	return &StartTagToken{
		_type:         StartTag,
		name:          "",
		attributes:    []Attribute{},
		isComponent:   false,
		isSelfClosing: false}
}

func NewEndTagToken() *EndTagToken {
	return &EndTagToken{
		_type: EndTag,
		name:  ""}
}

func NewCommentToken() *CommentToken {
	return &CommentToken{
		_type: Comment,
		data:  ""}
}

func NewTextToken() *TextToken {
	return &TextToken{
		_type: Text,
		data:  ""}
}

func NewCodeToken() *CodeToken {
	return &CodeToken{
		_type: Code,
		data:  ""}
}

func NewEndOfFileToken() *EndOfFileToken {
	return &EndOfFileToken{
		_type: EndOfFile}
}

type Token interface {
	GetType() TokenType
	GetName() string
	GetData() string
	GetAttributes() []Attribute
	GetIsSelfClosing() bool
	GetIsComponent() bool
	NewAttribute()
	AppendToName(string)
	AppendToData(string)
	AppendToAttributeName(string)
	AppendToAttributeValue(string)
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

// NewAttribute()

func (_self *StartTagToken) NewAttribute() {
	_self.attributes = append(_self.attributes, Attribute{
		Name:  "",
		Value: ""})
}

func (_self *EndTagToken) NewAttribute() {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *CommentToken) NewAttribute() {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *TextToken) NewAttribute() {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *CodeToken) NewAttribute() {
	utils.Assert(false, "token has attributes property", 2)
}

func (_self *EndOfFileToken) NewAttribute() {
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
	_self.attributes[len(_self.attributes)-1].Name = _self.attributes[len(_self.attributes)-1].Name + s
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
	_self.attributes[len(_self.attributes)-1].Value = _self.attributes[len(_self.attributes)-1].Value + s
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
	if _self.isSelfClosing {
		component = "(Component)"
	}
	var selfClosing string
	if _self.isSelfClosing {
		selfClosing = "(SelfClosing)"
	}
	verbose.Printf(0, "%s\t%s\t%d Attributes\t%s\t%s\n",
		_self._type,
		_self.name,
		len(_self.attributes),
		component,
		selfClosing)
	for i, attribute := range _self.attributes {
		verbose.Printf(0, " attribute%d\t%s\t%s\n", i, attribute.Name, attribute.Value)
	}
}

func (_self *EndTagToken) Print() {
	verbose.Printf(0, "%s\t%s\n", _self._type, _self.name)
}

func (_self *CommentToken) Print() {
	verbose.Printf(0, "%s\t%s\n", _self._type, _self.data)
}

func (_self *TextToken) Print() {
	verbose.Printf(0, "%s   \t%s\n", _self._type, _self.data)
}

func (_self *CodeToken) Print() {
	verbose.Printf(0, "%s\t%s\n", _self._type, _self.data)
}

func (_self *EndOfFileToken) Print() {
	verbose.Printf(0, "%s\n", _self._type)
}
