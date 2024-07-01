package main

import (
	"fmt"
	"runtime"
	"strings"
	"unicode"
)

const EOF = "EOF"

const (
	BeforeCodeState                               string = "BeforeCodeState"
	CodeState                                     string = "CodeState"
	AfterCodeState                                string = "AfterCodeState"
	EndOfFileState                                string = "EndOfFileState"
	DataState                                     string = "DataState"
	RCDATAState                                   string = "RCDATAState"
	RAWTEXTState                                  string = "RAWTEXTState"
	ScriptDataState                               string = "ScriptDataState"
	PLAINTEXTState                                string = "PLAINTEXTState"
	TagOpenState                                  string = "TagOpenState"
	EndTagOpenState                               string = "EndTagOpenState"
	TagNameState                                  string = "TagNameState"
	RCDATALessThanSignState                       string = "RCDATALessThanSignState"
	RCDATAEndTagOpenState                         string = "RCDATAEndTagOpenState"
	RCDATAEndTagNameState                         string = "RCDATAEndTagNameState"
	RAWTEXTLessThanSignState                      string = "RAWTEXTLessThanSignState"
	RAWTEXTEndTagOpenState                        string = "RAWTEXTEndTagOpenState"
	RAWTEXTEndTagNameState                        string = "RAWTEXTEndTagNameState"
	ScriptDataLessThanSignState                   string = "ScriptDataLessThanSignState"
	ScriptDataEndTagOpenState                     string = "ScriptDataEndTagOpenState"
	ScriptDataEndTagNameState                     string = "ScriptDataEndTagNameState"
	ScriptDataEscapeStartState                    string = "ScriptDataEscapeStartState"
	ScriptDataEscapeStartDashState                string = "ScriptDataEscapeStartDashState"
	ScriptDataEscapedState                        string = "ScriptDataEscapedState"
	ScriptDataEscapedDashState                    string = "ScriptDataEscapedDashState"
	ScriptDataEscapedDashDashState                string = "ScriptDataEscapedDashDashState"
	ScriptDataEscapedLessThanSignState            string = "ScriptDataEscapedLessThanSignState"
	ScriptDataEscapedEndTagOpenState              string = "ScriptDataEscapedEndTagOpenState"
	ScriptDataEscapedEndTagNameState              string = "ScriptDataEscapedEndTagNameState"
	ScriptDataDoubleEscapeStartState              string = "ScriptDataDoubleEscapeStartState"
	ScriptDataDoubleEscapedState                  string = "ScriptDataDoubleEscapedState"
	ScriptDataDoubleEscapedDashState              string = "ScriptDataDoubleEscapedDashState"
	ScriptDataDoubleEscapedDashDashState          string = "ScriptDataDoubleEscapedDashDashState"
	ScriptDataDoubleEscapedLessThanSignState      string = "ScriptDataDoubleEscapedLessThanSignState"
	ScriptDataDoubleEscapeEndState                string = "ScriptDataDoubleEscapeEndState"
	BeforeAttributeNameState                      string = "BeforeAttributeNameState"
	AttributeNameState                            string = "AttributeNameState"
	AfterAttributeNameState                       string = "AfterAttributeNameState"
	BeforeAttributeValueState                     string = "BeforeAttributeValueState"
	AttributeValueDoubleQuotedState               string = "AttributeValueDoubleQuotedState"
	AttributeValueSingleQuotedState               string = "AttributeValueSingleQuotedState"
	AttributeValueUnquotedState                   string = "AttributeValueUnquotedState"
	AfterAttributeValueQuotedState                string = "AfterAttributeValueQuotedState"
	SelfClosingStartTagState                      string = "SelfClosingStartTagState"
	BogusCommentState                             string = "BogusCommentState"
	MarkupDeclarationOpenState                    string = "MarkupDeclarationOpenState"
	CommentStartState                             string = "CommentStartState"
	CommentStartDashState                         string = "CommentStartDashState"
	CommentState                                  string = "CommentState"
	CommentLessThanSignState                      string = "CommentLessThanSignState"
	CommentLessThanSignBangState                  string = "CommentLessThanSignBangState"
	CommentLessThanSignBangDashState              string = "CommentLessThanSignBangDashState"
	CommentLessThanSignBangDashDashState          string = "CommentLessThanSignBangDashDashState"
	CommentEndDashState                           string = "CommentEndDashState"
	CommentEndState                               string = "CommentEndState"
	CommentEndBangState                           string = "CommentEndBangState"
	DOCTYPEState                                  string = "DOCTYPEState"
	BeforeDOCTYPENameState                        string = "BeforeDOCTYPENameState"
	DOCTYPENameState                              string = "DOCTYPENameState"
	AfterDOCTYPENameState                         string = "AfterDOCTYPENameState"
	AfterDOCTYPEPublicKeywordState                string = "AfterDOCTYPEPublicKeywordState"
	BeforeDOCTYPEPublicIdentifierState            string = "BeforeDOCTYPEPublicIdentifierState"
	DOCTYPEPublicIdentifierDoubleQuotedState      string = "DOCTYPEPublicIdentifierDoubleQuotedState"
	DOCTYPEPublicIdentifierSingleQuotedState      string = "DOCTYPEPublicIdentifierSingleQuotedState"
	AfterDOCTYPEPublicIdentifierState             string = "AfterDOCTYPEPublicIdentifierState"
	BetweenDOCTYPEPublicAndSystemIdentifiersState string = "BetweenDOCTYPEPublicAndSystemIdentifiersState"
	AfterDOCTYPESystemKeywordState                string = "AfterDOCTYPESystemKeywordState"
	BeforeDOCTYPESystemIdentifierState            string = "BeforeDOCTYPESystemIdentifierState"
	DOCTYPESystemIdentifierDoubleQuotedState      string = "DOCTYPESystemIdentifierDoubleQuotedState"
	DOCTYPESystemIdentifierSingleQuotedState      string = "DOCTYPESystemIdentifierSingleQuotedState"
	AfterDOCTYPESystemIdentifierState             string = "AfterDOCTYPESystemIdentifierState"
	BogusDOCTYPEState                             string = "BogusDOCTYPEState"
	CDATASectionState                             string = "CDATASectionState"
	CDATASectionBracketState                      string = "CDATASectionBracketState"
	CDATASectionEndState                          string = "CDATASectionEndState"
	CharacterReferenceState                       string = "CharacterReferenceState"
	NamedCharacterReferenceState                  string = "NamedCharacterReferenceState"
	AmbiguousAmpersandState                       string = "AmbiguousAmpersandState"
	NumericCharacterReferenceState                string = "NumericCharacterReferenceState"
	HexadecimalCharacterReferenceStartState       string = "HexadecimalCharacterReferenceStartState"
	DecimalCharacterReferenceStartState           string = "DecimalCharacterReferenceStartState"
	HexadecimalCharacterReferenceState            string = "HexadecimalCharacterReferenceState"
	DecimalCharacterReferenceState                string = "DecimalCharacterReferenceState"
	NumericCharacterReferenceEndState             string = "NumericCharacterReferenceEndState"
) // https://html.spec.whatwg.org/#tokenization

const (
	StartTag  string = "StartTag"
	EndTag    string = "EndTag"
	Comment   string = "Comment"
	Character string = "Character"
	Code      string = "Code"
	EndOfFile string = "EndOfFile"
)

type Attribute struct {
	Name  string
	Value string
}

type StartTagType struct {
	TagName     string
	SelfClosing bool
	Attributes  []Attribute
}

func (*StartTagType) New() StartTagType {
	return StartTagType{
		TagName:     "",
		SelfClosing: false,
		Attributes:  []Attribute{}}
}

type EndTagType struct {
	TagName string
}

func (*EndTagType) New() EndTagType {
	return EndTagType{
		TagName: ""}
}

type CommentType struct {
	Data string
}

func (*CommentType) New() CommentType {
	return CommentType{
		Data: ""}
}

type CharacterType struct {
	Data string
}

func (*CharacterType) New() CharacterType {
	return CharacterType{
		Data: ""}
}

type CodeType struct {
	Data string
}

func (*CodeType) New() CodeType {
	return CodeType{
		Data: ""}
}

type EndOfFileType struct {
}

func (*EndOfFileType) New() EndOfFileType {
	return EndOfFileType{}
}

type TokenType interface {
	StartTagType |
		EndTagType |
		CommentType |
		CharacterType |
		CodeType |
		EndOfFileType
}

type Token struct {
	Type  string
	Value interface{} //T TokenType
}

type DocReader struct {
	CodeBuffer      string
	CodeIndentCount int
	Char            string
	Chars           []string
	Curser          int
	Length          int
	LineNumber      int
	PeakBuffer      string
	ReturnState     string
	Rune            rune
	Runes           []rune
	State           string
	StringBuffer    string
	TemporaryBuffer string
	Token           *Token
	Tokens          []Token
}

func (_self *DocReader) New(source string) DocReader {
	var chars = intoStringArray(source)
	var runes = intoRuneArray(source)
	return DocReader{
		CodeBuffer:      "",
		CodeIndentCount: 0,
		Char:            chars[0],
		Chars:           chars,
		Curser:          0,
		Length:          len(chars),
		LineNumber:      0,
		PeakBuffer:      "",
		ReturnState:     "",
		Rune:            runes[0],
		Runes:           runes,
		State:           DataState,
		StringBuffer:    "",
		TemporaryBuffer: "",
		Token:           nil,
		Tokens:          make([]Token, 0)}
}

func (_self *DocReader) Consume() {
	if _self.Curser+1 >= _self.Length {
		_self.Char = EOF
	} else {
		_self.Char = _self.Chars[_self.Curser]
		_self.Rune = _self.Runes[_self.Curser]
		_self.Curser++
	}
}

func (_self *DocReader) ConsumeN(n int) {
	for n > 1 {
		_self.Curser++
		n--
	}
	_self.Consume()
}

func (_self *DocReader) ReConsume() {
	_self.Curser--
}

func (_self *DocReader) Ignore() {
	// Do nothing!
}

func (_self *DocReader) Peak(n int) {
	if _self.Curser+n >= _self.Length {
		return
	}
	_self.PeakBuffer = ""
	for i := _self.Curser; i < _self.Curser+n; i++ {
		_self.PeakBuffer = _self.PeakBuffer + _self.Chars[i]
	}
}

func (_self *DocReader) NewStartTagToken() {
	_self.Token = &Token{Type: StartTag, Value: (*StartTagType).New(nil)}
}

func (_self *DocReader) NewEndTagToken() {
	_self.Token = &Token{Type: EndTag, Value: (*EndTagType).New(nil)}
}

func (_self *DocReader) NewCommentToken() {
	_self.Token = &Token{Type: Comment, Value: (*CommentType).New(nil)}
}

func (_self *DocReader) NewCharacterToken() {
	_self.Token = &Token{Type: Character, Value: (*CharacterType).New(nil)}
}

func (_self *DocReader) NewCodeToken() {
	_self.Token = &Token{Type: Code, Value: (*CodeType).New(nil)}
}

func (_self *DocReader) NewEndOfFileToken() {
	_self.Token = &Token{Type: EndOfFile, Value: (*EndOfFileType).New(nil)}
}

func (_self *DocReader) ClearTemporaryBuffer() {
	_self.TemporaryBuffer = ""
}

func (_self *DocReader) ClearCodeBuffer() {
	_self.CodeIndentCount = 0
	_self.CodeBuffer = ""
}

func (_self *DocReader) ClearStringBuffer() {
	_self.Buffer = ""
}

func (_self *DocReader) AppendToTemporaryBuffer(s string) {
	_self.TemporaryBuffer = _self.TemporaryBuffer + s
}

func (_self *DocReader) AppendToCodeBuffer(s string) {
	switch s {
	case "{":
		_self.CodeIndentCount++
	case "}":
		_self.CodeIndentCount--
	}
	_self.CodeBuffer = _self.CodeBuffer + s
}

func (_self *DocReader) AppendToStringBuffer(s string) {
	_self.Buffer = _self.Buffer + s
}

// When a state says to flush code points consumed as a character reference, it means that for each code point in the temporary buffer (in the order they were added to the buffer) user agent must append the code point from the buffer to the current attribute's value if the character reference was consumed as part of an attribute, or emit the code point as a character token otherwise.
func (_self *DocReader) FlushCodePointsConsumedAsACharReference() {
	// if consumed as part of an attribute,
	if _self.ReturnState == AttributeValueDoubleQuotedState ||
		_self.ReturnState == AttributeValueSingleQuotedState ||
		_self.ReturnState == AttributeValueUnquotedState {
		// append the code point from the buffer to the current attribute's value
		_self.AppendToTokenAttributeValue(_self.TemporaryBuffer)
		return
	}
	// else, emit the code point as a character token
	_self.NewCharacterToken()
	_self.AppendToTokenData(_self.TemporaryBuffer)
	_self.EmitToken()
}

func (_self *DocReader) FlushCodePointsConsumedAsCode() {
	// if consumed as part of an attribute,
	if _self.ReturnState == AfterAttributeValueQuotedState {
		// append the code point from the buffer to the current attribute's value
		_self.AppendToTokenAttributeValue(_self.CodeBuffer)
		return
	}
	// else, emit the code point as a character token
	_self.NewCodeToken()
	_self.AppendToTokenData(_self.CodeBuffer)
	_self.EmitToken()
}

func (_self *DocReader) FlushCodePointsConsumedAsAChar() {
	if len(_self.Buffer) <= 0 {
		return
	}
	_self.NewCharacterToken()
	_self.AppendToTokenData(_self.Buffer)
	_self.EmitToken()
	_self.ClearBuffer()
}

func (_self *DocReader) AppendToTokenData(s string) {
	commentToken, ok := _self.Token.Value.(CommentType)
	if ok {
		commentToken.Data = commentToken.Data + s
		_self.Token.Value = commentToken
	}
	characterToken, ok := _self.Token.Value.(CharacterType)
	if ok {
		characterToken.Data = characterToken.Data + s
		_self.Token.Value = characterToken
	}
	codeToken, ok := _self.Token.Value.(CodeType)
	if ok {
		codeToken.Data = codeToken.Data + s
		_self.Token.Value = codeToken
	}
}

func (_self *DocReader) AppendToTokenName(s string) {
	startTagToken, ok := _self.Token.Value.(StartTagType)
	if ok {
		startTagToken.TagName = startTagToken.TagName + s
		_self.Token.Value = startTagToken
	}
	endTagToken, ok := _self.Token.Value.(EndTagType)
	if ok {
		endTagToken.TagName = endTagToken.TagName + s
		_self.Token.Value = endTagToken
	}
}

func (_self *DocReader) SetTokenSelfClosing(b bool) {
	startTagToken, ok := _self.Token.Value.(StartTagType)
	Assert(ok, "current token is of a StartTagType", 2)
	startTagToken.SelfClosing = b
	_self.Token.Value = startTagToken
}

func (_self *DocReader) NewTokenAttribute() {
	startTagToken, ok := _self.Token.Value.(StartTagType)
	Assert(ok, "current token is of a StartTagType", 2)
	startTagToken.Attributes = append(startTagToken.Attributes,
		Attribute{Name: "", Value: ""})
	_self.Token.Value = startTagToken
}

func (_self *DocReader) AppendToTokenAttributeName(s string) {
	startTagToken, ok := _self.Token.Value.(StartTagType)
	Assert(ok, "current token is of a StartTagType", 2)
	startTagToken.Attributes[len(startTagToken.Attributes)-1].Name =
		startTagToken.Attributes[len(startTagToken.Attributes)-1].Name + s
	_self.Token.Value = startTagToken
}

func (_self *DocReader) AppendToTokenAttributeValue(s string) {
	startTagToken, ok := _self.Token.Value.(StartTagType)
	Assert(ok, "current token is of a StartTagType", 2)
	startTagToken.Attributes[len(startTagToken.Attributes)-1].Value =
		startTagToken.Attributes[len(startTagToken.Attributes)-1].Value + s
	_self.Token.Value = startTagToken
}

func (_self *DocReader) EmitToken() {
	fmt.Printf("%v\n", *_self.Token)
	_self.Tokens = append(_self.Tokens, *_self.Token)
	_self.Token = nil
}

func isAsciiAlpha(r rune) bool {
	// https://infra.spec.whatwg.org/#ascii-alpha
	return unicode.IsLetter(r)
}

func isAsciiDigit(r rune) bool {
	// https://infra.spec.whatwg.org/#ascii-digit
	return unicode.IsDigit(r)
}

func isAsciiAlphanumeric(r rune) bool {
	// https://infra.spec.whatwg.org/#ascii-alphanumeric
	return isAsciiAlpha(r) && isAsciiDigit(r)
}

func isAsciiUpperAlpha(r rune) bool {
	// https://infra.spec.whatwg.org/#ascii-upper-alpha
	return unicode.IsUpper(r)
}

func isWhiteSpace(r rune) bool {
	// U+0009 CHARACTER TABULATION (tab)
	// U+000A LINE FEED (LF)
	// U+000C FORM FEED (FF)
	// U+0020 SPACE
	if string(r) == "\t" {
		return true
	}
	if string(r) == "\n" {
		return true
	}
	return unicode.IsSpace(r)
}

func intoRuneArray(str string) []rune {
	var arr = []rune{}
	for _, r := range str {
		arr = append(arr, r)
	}
	arr = append(arr, 0)
	return arr
}

func intoStringArray(str string) []string {
	var arr = []string{}
	for _, r := range str {
		arr = append(arr, string(r))
	}
	arr = append(arr, EOF)
	return arr
}

func Lexer(source string) error {
	var reader = (*DocReader).New(nil, source)

	var watchDogCounter = reader.Length * 10
	for reader.State != EndOfFileState {
		watchDogCounter--
		if watchDogCounter <= 0 {
			break
		}

		// fmt.Printf("state: %s - ", reader.State)
		switch reader.State {

		case DataState: // https://html.spec.whatwg.org/#data-state
			reader.Consume()
			switch reader.Char {
			case "{":
				reader.FlushCodePointsConsumedAsAChar()
				reader.ClearCodeBuffer()
				reader.ReturnState = DataState
				reader.State = BeforeCodeState
			// case "&":
			// 	reader.ReturnState = DataState
			// 	reader.State = CharacterReferenceState
			case "<":
				reader.FlushCodePointsConsumedAsAChar()
				reader.State = TagOpenState
			case EOF:
				reader.FlushCodePointsConsumedAsAChar()
				reader.NewEndOfFileToken()
				reader.EmitToken()
				reader.State = EndOfFileState
			default:
				// reader.NewCharacterToken()
				// reader.AppendToTokenData(reader.Char)
				// reader.EmitToken()
				reader.AppendToBuffer(reader.Char)
			}

		case RCDATAState: // https://html.spec.whatwg.org/#rcdata-state

		case RAWTEXTState: // https://html.spec.whatwg.org/#rawtext-state

		case ScriptDataState: // https://html.spec.whatwg.org/#script-data-state

		case PLAINTEXTState: // https://html.spec.whatwg.org/#plaintext-state

		case TagOpenState: // https://html.spec.whatwg.org/#tag-open-state
			reader.Consume()
			if isAsciiAlpha(reader.Rune) {
				reader.NewStartTagToken()
				reader.ReConsume()
				reader.State = TagNameState
				continue
			}
			switch reader.Char {
			case "!":
				reader.State = MarkupDeclarationOpenState
			case "/":
				reader.State = EndTagOpenState
			case "?":
				fmt.Printf("%s: unexpected-question-mark-instead-of-tag-name\n", TagOpenState)
				reader.NewCommentToken()
				reader.ReConsume()
				reader.State = BogusCommentState
			case EOF:
				fmt.Printf("%s: eof-before-tag-name\n", TagOpenState)
				reader.NewCharacterToken()
				reader.AppendToTokenData("<")
				reader.EmitToken()
				reader.NewEndOfFileToken()
				reader.EmitToken()
			default:
				fmt.Printf("%s: invalid-first-character-of-tag-name\n", TagOpenState)
				reader.NewCharacterToken()
				reader.AppendToTokenData("<")
				reader.EmitToken()
				reader.ReConsume()
				reader.State = DataState
			}

		case EndTagOpenState: // https://html.spec.whatwg.org/#tag-open-state
			reader.Consume()
			if isAsciiAlpha(reader.Rune) {
				reader.NewEndTagToken()
				reader.ReConsume()
				reader.State = TagNameState
				continue
			}
			switch reader.Char {
			case ">":
				fmt.Printf("%s: missing-end-tag-name\n", EndTagOpenState)
				reader.State = DataState

			case EOF:
				fmt.Printf("%s: eof-before-tag-name\n", EndTagOpenState)
				reader.NewCharacterToken()
				reader.AppendToTokenData("<")
				reader.EmitToken()
				reader.NewCharacterToken()
				reader.AppendToTokenData("/")
				reader.EmitToken()
				reader.NewEndOfFileToken()
				reader.EmitToken()

			default:
				fmt.Printf("%s: invalid-first-character-of-tag-name\n", EndTagOpenState)
				reader.NewCommentToken()
				reader.ReConsume()
				reader.State = BogusCommentState
			}

		case TagNameState: // https://html.spec.whatwg.org/#tag-name-state
			reader.Consume()
			if isWhiteSpace(reader.Rune) {
				reader.State = BeforeAttributeNameState
				continue
			}
			if isAsciiUpperAlpha(reader.Rune) {
				reader.AppendToTokenName(strings.ToLower(reader.Char))
				continue
			}
			switch reader.Char {
			case "/":
				reader.State = SelfClosingStartTagState
			case ">":
				reader.EmitToken()
				reader.State = DataState
			case EOF:
				fmt.Printf("%s: eof-in-tag\n", TagNameState)
				reader.NewEndOfFileToken()
				reader.EmitToken()
			default:
				reader.AppendToTokenName(reader.Char)
			}

		case RCDATALessThanSignState: // https://html.spec.whatwg.org/#rcdata-less-than-sign-state

		case RCDATAEndTagOpenState: // https://html.spec.whatwg.org/#rcdata-end-tag-open-state

		case RCDATAEndTagNameState: // https://html.spec.whatwg.org/#rcdata-end-tag-name-state

		case RAWTEXTLessThanSignState: // https://html.spec.whatwg.org/#rawtext-less-than-sign-state

		case RAWTEXTEndTagOpenState: // https://html.spec.whatwg.org/#rawtext-end-tag-open-state

		case RAWTEXTEndTagNameState: // https://html.spec.whatwg.org/#rawtext-end-tag-name-state

		case ScriptDataLessThanSignState: //

		case ScriptDataEndTagOpenState: //

		case ScriptDataEndTagNameState: //

		case ScriptDataEscapeStartState: //

		case ScriptDataEscapeStartDashState: //

		case ScriptDataEscapedState: //

		case ScriptDataEscapedDashState: //

		case ScriptDataEscapedDashDashState: //

		case ScriptDataEscapedLessThanSignState: //

		case ScriptDataEscapedEndTagOpenState: //

		case ScriptDataEscapedEndTagNameState: //

		case ScriptDataDoubleEscapeStartState: //

		case ScriptDataDoubleEscapedState: //

		case ScriptDataDoubleEscapedDashState: //

		case ScriptDataDoubleEscapedDashDashState: //

		case ScriptDataDoubleEscapedLessThanSignState: //

		case ScriptDataDoubleEscapeEndState: //

		case BeforeAttributeNameState: // https://html.spec.whatwg.org/#before-attribute-name-state
			reader.Consume()
			if isWhiteSpace(reader.Rune) {
				reader.Ignore()
				continue
			}
			switch reader.Char {
			case "/":
				reader.ReConsume()
				reader.State = AfterAttributeNameState
			case ">":
				reader.ReConsume()
				reader.State = AfterAttributeNameState
			case EOF:
				reader.ReConsume()
				reader.State = AfterAttributeNameState
			case "=":
				fmt.Printf("%s: unexpected-equals-sign-before-attribute-name\n", BeforeAttributeNameState)
				reader.NewTokenAttribute()
				reader.AppendToTokenAttributeName(reader.Char)
				reader.State = AttributeNameState
			default:
				reader.NewTokenAttribute()
				reader.ReConsume()
				reader.State = AttributeNameState
			}

		case AttributeNameState: // https://html.spec.whatwg.org/#attribute-name-state
			reader.Consume()
			if isWhiteSpace(reader.Rune) {
				reader.ReConsume()
				reader.State = AfterAttributeNameState
				continue
			}
			if isAsciiUpperAlpha(reader.Rune) {
				reader.AppendToTokenAttributeName(strings.ToLower(reader.Char))
				continue
			}
			switch reader.Char {
			case "/":
				reader.ReConsume()
				reader.State = AfterAttributeNameState
			case ">":
				reader.ReConsume()
				reader.State = AfterAttributeNameState
			case EOF:
				reader.ReConsume()
				reader.State = AfterAttributeNameState
			case "=":
				reader.State = BeforeAttributeValueState
			case `"`:
				fmt.Printf("%s: unexpected-character-in-attribute-name %s\n", AttributeNameState, reader.Char)
				reader.AppendToTokenAttributeName(reader.Char)
			case "'":
				fmt.Printf("%s: unexpected-character-in-attribute-name %s\n", AttributeNameState, reader.Char)
				reader.AppendToTokenAttributeName(reader.Char)
			case "<":
				fmt.Printf("%s: unexpected-character-in-attribute-name %s\n", AttributeNameState, reader.Char)
				reader.AppendToTokenAttributeName(reader.Char)
			default:
				reader.AppendToTokenAttributeName(reader.Char)
			}

		case AfterAttributeNameState: // https://html.spec.whatwg.org/#after-attribute-name-state
			reader.Consume()
			if isWhiteSpace(reader.Rune) {
				reader.Ignore()
				continue
			}
			switch reader.Char {
			case "/":
				reader.State = SelfClosingStartTagState
			case "=":
				reader.State = BeforeAttributeValueState
			case ">":
				reader.EmitToken()
				reader.State = DataState
			case EOF:
				fmt.Printf("%s: eof-in-tag\n", AfterAttributeNameState)
				reader.NewEndOfFileToken()
				reader.EmitToken()
			default:
				reader.NewTokenAttribute()
				reader.ReConsume()
				reader.State = AttributeNameState
			}

		case BeforeAttributeValueState: // https://html.spec.whatwg.org/#before-attribute-value-state
			reader.Consume()
			if isWhiteSpace(reader.Rune) {
				reader.Ignore()
				continue
			}
			switch reader.Char {
			case "{":
				reader.ClearCodeBuffer()
				reader.ReturnState = AfterAttributeValueQuotedState
				reader.State = BeforeCodeState
			case `"`:
				reader.State = AttributeValueDoubleQuotedState
			case "'":
				reader.State = AttributeValueSingleQuotedState
			case ">":
				fmt.Printf("%s: missing-attribute-value\n", BeforeAttributeValueState)
				reader.EmitToken()
				reader.State = DataState
			default:
				reader.ReConsume()
				reader.State = AttributeValueUnquotedState
			}

		case BeforeCodeState: // NOT IN SPEC
			reader.Consume()
			switch reader.Char {
			case "{":
				reader.AppendToCodeBuffer(reader.Char)
				reader.State = CodeState
			default:
				reader.ReConsume()
				reader.State = CodeState
			}

		case CodeState: // NOT IN SPEC
			reader.Consume()
			switch reader.Char {
			case "{":
				reader.ReConsume()
				reader.State = BeforeCodeState
			case EOF:
				fmt.Printf("%s: eof-in-code\n", CodeState)
				reader.FlushCodePointsConsumedAsCode()
				reader.NewEndOfFileToken()
				reader.EmitToken()
			case "}":
				reader.ReConsume()
				reader.State = AfterCodeState
			default:
				reader.AppendToCodeBuffer(reader.Char)
			}

		case AfterCodeState: // NOT IN SPEC
			reader.Consume()
			switch reader.Char {
			case "}":
				if reader.CodeIndentCount <= 0 {
					reader.FlushCodePointsConsumedAsCode()
					reader.State = reader.ReturnState
					continue
				}
				reader.AppendToCodeBuffer(reader.Char)
				reader.State = CodeState
			}

		case AttributeValueDoubleQuotedState: // https://html.spec.whatwg.org/#attribute-value-(double-quoted)-state
			reader.Consume()
			switch reader.Char {
			case `"`:
				reader.State = AfterAttributeValueQuotedState
			// case "&":
			// 	reader.ReturnState = AttributeValueDoubleQuotedState
			// 	reader.State = CharacterReferenceState
			case EOF:
				fmt.Printf("%s: eof-in-tag\n", AttributeValueDoubleQuotedState)
				reader.NewEndOfFileToken()
				reader.EmitToken()
			default:
				reader.AppendToTokenAttributeValue(reader.Char)
			}

		case AttributeValueSingleQuotedState: // https://html.spec.whatwg.org/#attribute-value-(single-quoted)-state
			reader.Consume()
			switch reader.Char {
			case "'":
				reader.State = AfterAttributeValueQuotedState
			// case "&":
			// 	reader.ReturnState = AttributeValueSingleQuotedState
			// 	reader.State = CharacterReferenceState
			case EOF:
				fmt.Printf("%s: eof-in-tag\n", AttributeValueSingleQuotedState)
				reader.NewEndOfFileToken()
				reader.EmitToken()
			default:
				reader.AppendToTokenAttributeValue(reader.Char)
			}

		case AttributeValueUnquotedState: // https://html.spec.whatwg.org/#attribute-value-(unquoted)-state
			reader.Consume()
			if isWhiteSpace(reader.Rune) {
				reader.State = BeforeAttributeNameState
				continue
			}
			switch reader.Char {
			// case "&":
			// 	reader.ReturnState = AttributeValueUnquotedState
			// 	reader.State = CharacterReferenceState
			case ">":
				reader.EmitToken()
				reader.State = DataState
			case `"`:
				fmt.Printf("%s: unexpected-character-in-unquoted-attribute-value %s\n", AttributeValueUnquotedState, reader.Char)
				reader.AppendToTokenAttributeValue(reader.Char)
			case "'":
				fmt.Printf("%s: unexpected-character-in-unquoted-attribute-value %s\n", AttributeValueUnquotedState, reader.Char)
				reader.AppendToTokenAttributeValue(reader.Char)
			case "<":
				fmt.Printf("%s: unexpected-character-in-unquoted-attribute-value %s\n", AttributeValueUnquotedState, reader.Char)
				reader.AppendToTokenAttributeValue(reader.Char)
			case "=":
				fmt.Printf("%s: unexpected-character-in-unquoted-attribute-value %s\n", AttributeValueUnquotedState, reader.Char)
				reader.AppendToTokenAttributeValue(reader.Char)
			case "`":
				fmt.Printf("%s: unexpected-character-in-unquoted-attribute-value %s\n", AttributeValueUnquotedState, reader.Char)
				reader.AppendToTokenAttributeValue(reader.Char)
			case EOF:
				fmt.Printf("%s: eof-in-tag\n", AttributeValueUnquotedState)
				reader.NewEndOfFileToken()
				reader.EmitToken()
			default:
				reader.AppendToTokenAttributeValue(reader.Char)
			}

		case AfterAttributeValueQuotedState: // https://html.spec.whatwg.org/#after-attribute-value-(quoted)-state
			reader.Consume()
			if isWhiteSpace(reader.Rune) {
				reader.State = BeforeAttributeNameState
				continue
			}
			switch reader.Char {
			case "/":
				reader.State = SelfClosingStartTagState
			case ">":
				reader.EmitToken()
				reader.State = DataState
			case EOF:
				fmt.Printf("%s: eof-in-tag\n", AfterAttributeValueQuotedState)
				reader.NewEndOfFileToken()
				reader.EmitToken()
			default:
				fmt.Printf("%s: missing-whitespace-between-attributes\n", AfterAttributeValueQuotedState)
				reader.ReConsume()
				reader.State = BeforeAttributeNameState
			}

		case SelfClosingStartTagState: // https://html.spec.whatwg.org/#self-closing-start-tag-state
			reader.Consume()
			switch reader.Char {
			case ">":
				reader.SetTokenSelfClosing(true)
				reader.EmitToken()
				reader.State = DataState
			case EOF:
				fmt.Printf("%s: eof-in-tag\n", SelfClosingStartTagState)
				reader.NewEndOfFileToken()
				reader.EmitToken()
			default:
				fmt.Printf("%s: unexpected-solidus-in-tag\n", SelfClosingStartTagState)
				reader.ReConsume()
				reader.State = BeforeAttributeNameState
			}

		case BogusCommentState: // https://html.spec.whatwg.org/#bogus-comment-state
			reader.Consume()
			switch reader.Char {
			case ">":
				reader.EmitToken()
				reader.State = DataState
			case EOF:
				reader.EmitToken()
				reader.NewEndOfFileToken()
				reader.EmitToken()
			default:
				reader.AppendToTokenData(reader.Char)
			}

		case MarkupDeclarationOpenState: // https://html.spec.whatwg.org/#markup-declaration-open-state
			var n = 2
			reader.Peak(n)
			switch reader.PeakBuffer {
			case "--":
				reader.NewCommentToken()
				reader.ConsumeN(n)
				reader.State = CommentStartState
			default:
				fmt.Printf("%s: incorrectly-opened-comment\n", MarkupDeclarationOpenState)
				reader.NewCommentToken()
				reader.State = BogusCommentState
			}

		case CommentStartState: // https://html.spec.whatwg.org/#comment-start-state
			reader.Consume()
			switch reader.Char {
			case "-":
				reader.State = CommentStartDashState
			case ">":
				fmt.Printf("%s: abrupt-closing-of-empty-comment\n", CommentStartState)
				reader.EmitToken()
				reader.State = DataState
			default:
				reader.ReConsume()
				reader.State = CommentState
			}

		case CommentStartDashState: // https://html.spec.whatwg.org/#comment-start-dash-state
			reader.Consume()
			switch reader.Char {
			case "-":
				reader.State = CommentEndState
			case ">":
				fmt.Printf("%s: abrupt-closing-of-empty-comment\n", CommentStartDashState)
				reader.EmitToken()
				reader.State = DataState
			case EOF:
				fmt.Printf("%s: eof-in-comment\n", CommentStartDashState)
				reader.EmitToken()
				reader.NewEndOfFileToken()
				reader.EmitToken()
			default:
				reader.AppendToTokenData("-")
				reader.ReConsume()
				reader.State = CommentState
			}

		case CommentState: // https://html.spec.whatwg.org/#comment-state
			reader.Consume()
			switch reader.Char {
			case "<":
				reader.AppendToTokenData(reader.Char)
				reader.State = CommentLessThanSignState
			case "-":
				reader.State = CommentEndDashState
			case EOF:
				fmt.Printf("%s: eof-in-comment\n", CommentState)
				reader.EmitToken()
				reader.NewEndOfFileToken()
				reader.EmitToken()
			default:
				reader.AppendToTokenData(reader.Char)
			}

		case CommentLessThanSignState: // https://html.spec.whatwg.org/#comment-less-than-sign-state
			reader.Consume()
			switch reader.Char {
			case "!":
				reader.AppendToTokenData(reader.Char)
				reader.State = CommentLessThanSignBangState
			case "<":
				reader.AppendToTokenData(reader.Char)
			default:
				reader.ReConsume()
				reader.State = CommentState
			}

		case CommentLessThanSignBangState: // https://html.spec.whatwg.org/#comment-less-than-sign-bang-state
			reader.Consume()
			switch reader.Char {
			case "-":
				reader.State = CommentLessThanSignBangDashState
			default:
				reader.ReConsume()
				reader.State = CommentState
			}

		case CommentLessThanSignBangDashState: // https://html.spec.whatwg.org/#comment-less-than-sign-bang-dash-state
			reader.Consume()
			switch reader.Char {
			case "-":
				reader.State = CommentLessThanSignBangDashDashState
			default:
				reader.ReConsume()
				reader.State = CommentState
			}

		case CommentLessThanSignBangDashDashState: // https://html.spec.whatwg.org/#comment-less-than-sign-bang-dash-dash-state
			reader.Consume()
			switch reader.Char {
			case ">":
				reader.ReConsume()
				reader.State = CommentState
			case EOF:
				reader.ReConsume()
				reader.State = CommentState
			default:
				fmt.Printf("%s: nested-comment\n", CommentLessThanSignBangDashDashState)
				reader.ReConsume()
				reader.State = CommentEndState
			}

		case CommentEndDashState: // https://html.spec.whatwg.org/#comment-end-dash-state
			reader.Consume()
			switch reader.Char {
			case "-":
				reader.State = CommentEndState
			case EOF:
				fmt.Printf("%s: eof-in-comment\n", CommentEndDashState)
				reader.EmitToken()
				reader.NewEndOfFileToken()
				reader.EmitToken()
			default:
				reader.AppendToTokenData("-")
				reader.ReConsume()
				reader.State = CommentState
			}

		case CommentEndState: // https://html.spec.whatwg.org/#comment-end-state
			reader.Consume()
			switch reader.Char {
			case ">":
				reader.EmitToken()
				reader.State = DataState
			case "!":
				reader.State = CommentEndBangState
			case "-":
				reader.AppendToTokenData("-")
			case EOF:
				fmt.Printf("%s: eof-in-comment\n", CommentEndState)
				reader.EmitToken()
				reader.NewEndOfFileToken()
				reader.EmitToken()
			default:
				reader.AppendToTokenData("-")
				reader.ReConsume()
				reader.State = CommentState
			}

		case CommentEndBangState: // https://html.spec.whatwg.org/#comment-end-bang-state
			reader.Consume()
			switch reader.Char {
			case "-":
				reader.AppendToTokenData("--!")
				reader.State = CommentEndDashState
			case ">":
				fmt.Printf("%s: incorrectly-closed-comment\n", CommentEndBangState)
				reader.EmitToken()
				reader.State = DataState
			case EOF:
				fmt.Printf("%s: eof-in-comment\n", CommentEndBangState)
				reader.EmitToken()
				reader.NewEndOfFileToken()
				reader.EmitToken()
			default:
				reader.AppendToTokenData("--!")
				reader.ReConsume()
				reader.State = CommentState
			}

		case DOCTYPEState:

		case BeforeDOCTYPENameState:

		case DOCTYPENameState:

		case AfterDOCTYPENameState:

		case AfterDOCTYPEPublicKeywordState:

		case BeforeDOCTYPEPublicIdentifierState:

		case DOCTYPEPublicIdentifierDoubleQuotedState:

		case DOCTYPEPublicIdentifierSingleQuotedState:

		case AfterDOCTYPEPublicIdentifierState:

		case BetweenDOCTYPEPublicAndSystemIdentifiersState:

		case AfterDOCTYPESystemKeywordState:

		case BeforeDOCTYPESystemIdentifierState:

		case DOCTYPESystemIdentifierDoubleQuotedState:

		case DOCTYPESystemIdentifierSingleQuotedState:

		case AfterDOCTYPESystemIdentifierState:

		case BogusDOCTYPEState:

		case CDATASectionState: // https://html.spec.whatwg.org/#cdata-section-state

		case CDATASectionBracketState: // https://html.spec.whatwg.org/#cdata-section-bracket-state

		case CDATASectionEndState: // https://html.spec.whatwg.org/#cdata-section-end-state

		case CharacterReferenceState: // https://html.spec.whatwg.org/#character-reference-state
			reader.ClearTemporaryBuffer()
			reader.AppendToTemporaryBuffer("&")
			reader.Consume()
			if isAsciiAlphanumeric(reader.Rune) {
				reader.ReConsume()
				reader.State = NamedCharacterReferenceState
				continue
			}
			switch reader.Char {
			case "#":
				reader.AppendToTemporaryBuffer(reader.Char)
				reader.State = NumericCharacterReferenceState
			default:
				reader.FlushCodePointsConsumedAsACharReference()
				reader.ReConsume()
				reader.State = reader.ReturnState
			}

		case NamedCharacterReferenceState: // https://html.spec.whatwg.org/#named-character-reference-state
			reader.FlushCodePointsConsumedAsACharReference()
			reader.ReConsume()
			reader.State = reader.ReturnState

		case AmbiguousAmpersandState: // https://html.spec.whatwg.org/#ambiguous-ampersand-state
			reader.FlushCodePointsConsumedAsACharReference()
			reader.ReConsume()
			reader.State = reader.ReturnState

		case NumericCharacterReferenceState: // https://html.spec.whatwg.org/#numeric-character-reference-state
			reader.FlushCodePointsConsumedAsACharReference()
			reader.ReConsume()
			reader.State = reader.ReturnState

		case HexadecimalCharacterReferenceStartState: // https://html.spec.whatwg.org/#hexadecimal-character-reference-start-state

		case DecimalCharacterReferenceStartState: // https://html.spec.whatwg.org/#decimal-character-reference-start-state

		case HexadecimalCharacterReferenceState: // https://html.spec.whatwg.org/#hexadecimal-character-reference-state

		case DecimalCharacterReferenceState: // https://html.spec.whatwg.org/#decimal-character-reference-state

		case NumericCharacterReferenceEndState: // https://html.spec.whatwg.org/#numeric-character-reference-end-state

		}
	}
	if watchDogCounter <= 0 {
		return fmt.Errorf("watchDogCounter ran out")
	}
	return nil
}

func main() {
	// err := Lexer(`<ul if={func() []string { return []string{"a","b","c"}}} />`)
	err := Lexer(`<div>hello</div><div>world</div><div>Value: {count.Get()}</div>`)
	// err := Lexer(tricky)
	if err != nil {
		fmt.Println(err)
	}
}

const tricky = `<div>
	<!--this is a comment {} "" -->
	<p>This is some text !##^%)!&#)!&$)(*&!&^~)^87686912634</p>
	<a href="https://www.theinterwebs.com">link text</a>
	<br />
	<input/>
	<div class="me andme andmeandme"></div>
	<div class=me>
		<ul if={func() []string { return []string{"a","b","c"}}}>
			<Li />
			<Li/>
			<Button count />
			<Button count/>
		</ul>
		<p>Value: {count.Get()}</p>
		<p class:dark={isDark}></p>
	</div>
</div>`

func Assert(condition bool, msg string, skip int) {
	if !condition {
		_, file, line, _ := runtime.Caller(skip)
		panic(fmt.Sprintf("ASSERT [the following condition was not met in '%s' at line %d]: %s", file, line, msg))
	}
}

// Return items of array from index i to index j (inclusive of i and j).
func Pick[T any](a *[]T, i *int, j *int) []T {
	Assert(*i >= 0, fmt.Sprintf("i >= 0 when calling 'func pick(a, i, j)', i:%d j:%d", *i, *j), 2)
	Assert(*i <= *j, fmt.Sprintf("i <= j when calling 'func pick(a, i, j)', i:%d j:%d", *i, *j), 2)
	Assert(*j < len(*a), fmt.Sprintf("j < len(a) when calling 'func pick(a, i, j)', j:%d len(a):%d", *j, len(*a)), 2)
	return (*a)[*i : *j+1]
}

func Collect(arr *[]string, i *int, j int) string {
	var str = ""
	var picked = Pick(arr, i, &j)
	for k := 0; k < len(picked); k++ {
		str = str + picked[k]
	}
	return str
}
