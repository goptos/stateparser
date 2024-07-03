package lexer

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/goptos/stateparser/lexer/tokens"
	"github.com/goptos/utils"
)

var verbose = (*utils.Verbose).New(nil)

const EOF = "EOF"

// https://infra.spec.whatwg.org/#ascii-alpha
func isAsciiAlpha(r rune) bool {
	return unicode.IsLetter(r)
}

// https://infra.spec.whatwg.org/#ascii-upper-alpha
func isAsciiUpperAlpha(r rune) bool {
	return unicode.IsUpper(r)
}

// https://infra.spec.whatwg.org/#ascii-whitespace
func isAsciiWhiteSpace(r rune) bool {
	if string(r) == "\t" {
		return true
	}
	if string(r) == "\n" {
		return true
	}
	return unicode.IsSpace(r)
}

func intoArray(s string) ([]rune, []string) {
	var aR = []rune{}
	var aC = []string{}
	for _, r := range s {
		aR = append(aR, r)
		aC = append(aC, string(r))
	}
	aR = append(aR, 0)
	aC = append(aC, EOF)
	return aR, aC
}

// https://html.spec.whatwg.org/#tokenization
const (
	beforeCodeState                      string = "beforeCodeState"
	codeState                            string = "codeState"
	afterCodeState                       string = "afterCodeState"
	endOfFileState                       string = "endOfFileState"
	dataState                            string = "dataState"
	tagOpenState                         string = "tagOpenState"
	endTagOpenState                      string = "endTagOpenState"
	tagNameState                         string = "tagNameState"
	beforeAttributeNameState             string = "beforeAttributeNameState"
	attributeNameState                   string = "attributeNameState"
	afterAttributeNameState              string = "afterAttributeNameState"
	beforeAttributeValueState            string = "beforeAttributeValueState"
	attributeValueDoubleQuotedState      string = "attributeValueDoubleQuotedState"
	attributeValueSingleQuotedState      string = "attributeValueSingleQuotedState"
	attributeValueUnquotedState          string = "attributeValueUnquotedState"
	afterAttributeValueQuotedState       string = "afterAttributeValueQuotedState"
	selfClosingStartTagState             string = "selfClosingStartTagState"
	bogusCommentState                    string = "bogusCommentState"
	markupDeclarationOpenState           string = "markupDeclarationOpenState"
	commentStartState                    string = "commentStartState"
	commentStartDashState                string = "commentStartDashState"
	commentState                         string = "commentState"
	commentLessThanSignState             string = "commentLessThanSignState"
	commentLessThanSignBangState         string = "commentLessThanSignBangState"
	commentLessThanSignBangDashState     string = "commentLessThanSignBangDashState"
	commentLessThanSignBangDashDashState string = "commentLessThanSignBangDashDashState"
	commentEndDashState                  string = "commentEndDashState"
	commentEndState                      string = "commentEndState"
	commentEndBangState                  string = "commentEndBangState"
)

type Lexer struct {
	codeBuffer      string
	codeIndentCount int
	char            string
	chars           []string
	curser          int
	length          int
	peakBuffer      string
	returnState     string
	_rune           rune
	runes           []rune
	Source          string
	state           string
	textBuffer      string
	token           tokens.Token
	Tokens          []tokens.Token
}

func New(source string) *Lexer {
	var runes, chars = intoArray(source)
	return &Lexer{
		codeBuffer:      "",
		codeIndentCount: 0,
		char:            chars[0],
		chars:           chars,
		curser:          0,
		length:          len(chars),
		peakBuffer:      "",
		returnState:     "",
		_rune:           runes[0],
		runes:           runes,
		state:           dataState,
		Source:          source,
		textBuffer:      "",
		token:           nil,
		Tokens:          []tokens.Token{}}
}

func (_self *Lexer) consume() {
	if _self.curser+1 >= _self.length {
		_self.char = EOF
	} else {
		_self.char = _self.chars[_self.curser]
		_self._rune = _self.runes[_self.curser]
		_self.curser++
	}
}

func (_self *Lexer) consumeN(n int) {
	for n > 1 {
		_self.curser++
		n--
	}
	_self.consume()
}

func (_self *Lexer) reConsume() {
	_self.curser--
}

func (_self *Lexer) ignore() {
	// Do nothing!
}

func (_self *Lexer) peak(n int) {
	if _self.curser+n >= _self.length {
		return
	}
	_self.peakBuffer = ""
	for i := _self.curser; i < _self.curser+n; i++ {
		_self.peakBuffer = _self.peakBuffer + _self.chars[i]
	}
}

func (_self *Lexer) clearTextBuffer() {
	_self.textBuffer = ""
}

func (_self *Lexer) clearCodeBuffer() {
	_self.codeIndentCount = 0
	_self.codeBuffer = ""
}

func (_self *Lexer) appendToTextBuffer(s string) {
	switch s {
	case " ":
		if len(_self.textBuffer) <= 0 {
			return
		}
	case "\n":
		if len(_self.textBuffer) <= 0 {
			return
		}
		s = "\\n"
	case "\t":
		if len(_self.textBuffer) <= 0 {
			return
		}
		s = "\\t"
	}
	_self.textBuffer = _self.textBuffer + s
}

func (_self *Lexer) appendToCodeBuffer(s string) {
	switch s {
	case "{":
		_self.codeIndentCount++
	case "}":
		_self.codeIndentCount--
	}
	_self.codeBuffer = _self.codeBuffer + s
}

func (_self *Lexer) flushCodeBufferToToken() {
	if _self.returnState == afterAttributeValueQuotedState {
		_self.token.AppendToAttributeValue(_self.codeBuffer)
		return
	}
	_self.token = tokens.NewCodeToken()
	_self.token.AppendToData(_self.codeBuffer)
	_self.emitToken()
}

func (_self *Lexer) flushTextBufferToToken() {
	if len(_self.textBuffer) <= 0 {
		return
	}
	_self.token = tokens.NewTextToken()
	_self.token.AppendToData(_self.textBuffer)
	_self.emitToken()
	_self.clearTextBuffer()
}

func (_self *Lexer) emitToken() {
	if verbose.Level >= 4 {
		_self.token.Print()
	}
	_self.Tokens = append(_self.Tokens, _self.token)
	_self.token = nil
}

func (_self *Lexer) Tokenise() error {
	verbose.Printf(4, "::: Lexer.Tokenise() :::\n")
	for _self.state != endOfFileState {
		switch _self.state {

		case dataState: // https://html.spec.whatwg.org/#data-state
			_self.consume()
			switch _self.char {
			case "{":
				_self.flushTextBufferToToken()
				_self.clearCodeBuffer()
				_self.returnState = dataState
				_self.state = beforeCodeState
			case "<":
				_self.flushTextBufferToToken()
				_self.state = tagOpenState
			case EOF:
				_self.flushTextBufferToToken()
				_self.token = tokens.NewEndOfFileToken()
				_self.emitToken()
				_self.state = endOfFileState
			default:
				_self.appendToTextBuffer(_self.char)
			}

		case tagOpenState: // https://html.spec.whatwg.org/#tag-open-state
			_self.consume()
			if isAsciiAlpha(_self._rune) {
				_self.token = tokens.NewStartTagToken()
				_self.reConsume()
				_self.state = tagNameState
				continue
			}
			switch _self.char {
			case "!":
				_self.state = markupDeclarationOpenState
			case "/":
				_self.state = endTagOpenState
			case "?":
				verbose.Printf(0, "%s: unexpected-question-mark-instead-of-tag-name\n", tagOpenState)
				return fmt.Errorf("%s: unexpected-question-mark-instead-of-tag-name", tagOpenState)
			case EOF:
				verbose.Printf(0, "%s: eof-before-tag-name\n", tagOpenState)
				return fmt.Errorf("%s: eof-before-tag-name", tagOpenState)
			default:
				verbose.Printf(0, "%s: invalid-first-character-of-tag-name\n", tagOpenState)
				return fmt.Errorf("%s: invalid-first-character-of-tag-name", tagOpenState)
			}

		case endTagOpenState: // https://html.spec.whatwg.org/#tag-open-state
			_self.consume()
			if isAsciiAlpha(_self._rune) {
				_self.token = tokens.NewEndTagToken()
				_self.reConsume()
				_self.state = tagNameState
				continue
			}
			switch _self.char {
			case ">":
				verbose.Printf(0, "%s: missing-end-tag-name\n", endTagOpenState)
				return fmt.Errorf("%s: missing-end-tag-name", endTagOpenState)

			case EOF:
				verbose.Printf(0, "%s: eof-before-tag-name\n", endTagOpenState)
				return fmt.Errorf("%s: eof-before-tag-name", endTagOpenState)

			default:
				verbose.Printf(0, "%s: invalid-first-character-of-tag-name\n", endTagOpenState)
				return fmt.Errorf("%s: invalid-first-character-of-tag-name", endTagOpenState)
			}

		case tagNameState: // https://html.spec.whatwg.org/#tag-name-state
			_self.consume()
			if isAsciiWhiteSpace(_self._rune) {
				_self.state = beforeAttributeNameState
				continue
			}
			if isAsciiUpperAlpha(_self._rune) {
				if len(_self.token.GetName()) == 0 {
					_self.token.SetIsComponent(true)
					_self.token.AppendToName(_self.char)
					continue
				}
				_self.token.AppendToName(strings.ToLower(_self.char))
				continue
			}
			switch _self.char {
			case "/":
				_self.state = selfClosingStartTagState
			case ">":
				_self.emitToken()
				_self.state = dataState
			case EOF:
				verbose.Printf(0, "%s: eof-in-tag\n", tagNameState)
				return fmt.Errorf("%s: eof-in-tag", tagNameState)
			default:
				_self.token.AppendToName(_self.char)
			}

		case beforeAttributeNameState: // https://html.spec.whatwg.org/#before-attribute-name-state
			_self.consume()
			if isAsciiWhiteSpace(_self._rune) {
				_self.ignore()
				continue
			}
			switch _self.char {
			case "/":
				_self.reConsume()
				_self.state = afterAttributeNameState
			case ">":
				_self.reConsume()
				_self.state = afterAttributeNameState
			case EOF:
				_self.reConsume()
				_self.state = afterAttributeNameState
			case "=":
				verbose.Printf(0, "%s: unexpected-equals-sign-before-attribute-name\n", beforeAttributeNameState)
				return fmt.Errorf("%s: unexpected-equals-sign-before-attribute-name", beforeAttributeNameState)
			default:
				_self.token.NewAttribute()
				_self.reConsume()
				_self.state = attributeNameState
			}

		case attributeNameState: // https://html.spec.whatwg.org/#attribute-name-state
			_self.consume()
			if isAsciiWhiteSpace(_self._rune) {
				_self.reConsume()
				_self.state = afterAttributeNameState
				continue
			}
			if isAsciiUpperAlpha(_self._rune) {
				_self.token.AppendToAttributeName(strings.ToLower(_self.char))
				continue
			}
			switch _self.char {
			case "/":
				_self.reConsume()
				_self.state = afterAttributeNameState
			case ">":
				_self.reConsume()
				_self.state = afterAttributeNameState
			case EOF:
				_self.reConsume()
				_self.state = afterAttributeNameState
			case "=":
				_self.state = beforeAttributeValueState
			case `"`:
				verbose.Printf(0, "%s: unexpected-character-in-attribute-name %s\n", attributeNameState, _self.char)
				return fmt.Errorf("%s: unexpected-character-in-attribute-name %s", attributeNameState, _self.char)
			case "'":
				verbose.Printf(0, "%s: unexpected-character-in-attribute-name %s\n", attributeNameState, _self.char)
				return fmt.Errorf("%s: unexpected-character-in-attribute-name %s", attributeNameState, _self.char)
			case "<":
				verbose.Printf(0, "%s: unexpected-character-in-attribute-name %s\n", attributeNameState, _self.char)
				return fmt.Errorf("%s: unexpected-character-in-attribute-name %s", attributeNameState, _self.char)
			default:
				_self.token.AppendToAttributeName(_self.char)
			}

		case afterAttributeNameState: // https://html.spec.whatwg.org/#after-attribute-name-state
			_self.consume()
			if isAsciiWhiteSpace(_self._rune) {
				_self.ignore()
				continue
			}
			switch _self.char {
			case "/":
				_self.state = selfClosingStartTagState
			case "=":
				_self.state = beforeAttributeValueState
			case ">":
				_self.emitToken()
				_self.state = dataState
			case EOF:
				verbose.Printf(0, "%s: eof-in-tag\n", afterAttributeNameState)
				return fmt.Errorf("%s: eof-in-tag", afterAttributeNameState)
			default:
				_self.token.NewAttribute()
				_self.reConsume()
				_self.state = attributeNameState
			}

		case beforeAttributeValueState: // https://html.spec.whatwg.org/#before-attribute-value-state
			_self.consume()
			if isAsciiWhiteSpace(_self._rune) {
				_self.ignore()
				continue
			}
			switch _self.char {
			case "{":
				_self.clearCodeBuffer()
				_self.returnState = afterAttributeValueQuotedState
				_self.state = beforeCodeState
			case `"`:
				_self.state = attributeValueDoubleQuotedState
			case "'":
				_self.state = attributeValueSingleQuotedState
			case ">":
				verbose.Printf(0, "%s: missing-attribute-value\n", beforeAttributeValueState)
				return fmt.Errorf("%s: missing-attribute-value", beforeAttributeValueState)
			default:
				_self.reConsume()
				_self.state = attributeValueUnquotedState
			}

		case beforeCodeState: // NOT IN SPEC
			_self.consume()
			switch _self.char {
			case "{":
				_self.appendToCodeBuffer(_self.char)
				_self.state = codeState
			default:
				_self.reConsume()
				_self.state = codeState
			}

		case codeState: // NOT IN SPEC
			_self.consume()
			switch _self.char {
			case "{":
				_self.reConsume()
				_self.state = beforeCodeState
			case EOF:
				verbose.Printf(0, "%s: eof-in-code\n", codeState)
				return fmt.Errorf("%s: eof-in-code", codeState)
			case "}":
				_self.reConsume()
				_self.state = afterCodeState
			default:
				_self.appendToCodeBuffer(_self.char)
			}

		case afterCodeState: // NOT IN SPEC
			_self.consume()
			switch _self.char {
			case "}":
				if _self.codeIndentCount <= 0 {
					_self.flushCodeBufferToToken()
					_self.state = _self.returnState
					continue
				}
				_self.appendToCodeBuffer(_self.char)
				_self.state = codeState
			}

		case attributeValueDoubleQuotedState: // https://html.spec.whatwg.org/#attribute-value-(double-quoted)-state
			_self.consume()
			switch _self.char {
			case `"`:
				_self.state = afterAttributeValueQuotedState
			case EOF:
				verbose.Printf(0, "%s: eof-in-tag\n", attributeValueDoubleQuotedState)
				return fmt.Errorf("%s: eof-in-tag", attributeValueDoubleQuotedState)
			default:
				_self.token.AppendToAttributeValue(_self.char)
			}

		case attributeValueSingleQuotedState: // https://html.spec.whatwg.org/#attribute-value-(single-quoted)-state
			_self.consume()
			switch _self.char {
			case "'":
				_self.state = afterAttributeValueQuotedState
			case EOF:
				verbose.Printf(0, "%s: eof-in-tag\n", attributeValueSingleQuotedState)
				return fmt.Errorf("%s: eof-in-tag", attributeValueSingleQuotedState)
			default:
				_self.token.AppendToAttributeValue(_self.char)
			}

		case attributeValueUnquotedState: // https://html.spec.whatwg.org/#attribute-value-(unquoted)-state
			_self.consume()
			if isAsciiWhiteSpace(_self._rune) {
				_self.state = beforeAttributeNameState
				continue
			}
			switch _self.char {
			case ">":
				_self.emitToken()
				_self.state = dataState
			case `"`:
				verbose.Printf(0, "%s: unexpected-character-in-unquoted-attribute-value %s\n", attributeValueUnquotedState, _self.char)
				return fmt.Errorf("%s: unexpected-character-in-unquoted-attribute-value %s", attributeValueUnquotedState, _self.char)
			case "'":
				verbose.Printf(0, "%s: unexpected-character-in-unquoted-attribute-value %s\n", attributeValueUnquotedState, _self.char)
				return fmt.Errorf("%s: unexpected-character-in-unquoted-attribute-value %s", attributeValueUnquotedState, _self.char)
			case "<":
				verbose.Printf(0, "%s: unexpected-character-in-unquoted-attribute-value %s\n", attributeValueUnquotedState, _self.char)
				return fmt.Errorf("%s: unexpected-character-in-unquoted-attribute-value %s", attributeValueUnquotedState, _self.char)
			case "=":
				verbose.Printf(0, "%s: unexpected-character-in-unquoted-attribute-value %s\n", attributeValueUnquotedState, _self.char)
				return fmt.Errorf("%s: unexpected-character-in-unquoted-attribute-value %s", attributeValueUnquotedState, _self.char)
			case "`":
				verbose.Printf(0, "%s: unexpected-character-in-unquoted-attribute-value %s\n", attributeValueUnquotedState, _self.char)
				return fmt.Errorf("%s: unexpected-character-in-unquoted-attribute-value %s", attributeValueUnquotedState, _self.char)
			case EOF:
				verbose.Printf(0, "%s: eof-in-tag\n", attributeValueUnquotedState)
				return fmt.Errorf("%s: eof-in-tag", attributeValueUnquotedState)
			default:
				_self.token.AppendToAttributeValue(_self.char)
			}

		case afterAttributeValueQuotedState: // https://html.spec.whatwg.org/#after-attribute-value-(quoted)-state
			_self.consume()
			if isAsciiWhiteSpace(_self._rune) {
				_self.state = beforeAttributeNameState
				continue
			}
			switch _self.char {
			case "/":
				_self.state = selfClosingStartTagState
			case ">":
				_self.emitToken()
				_self.state = dataState
			case EOF:
				verbose.Printf(0, "%s: eof-in-tag\n", afterAttributeValueQuotedState)
				return fmt.Errorf("%s: eof-in-tag", afterAttributeValueQuotedState)
			default:
				verbose.Printf(0, "%s: missing-whitespace-between-attributes\n", afterAttributeValueQuotedState)
				return fmt.Errorf("%s: missing-whitespace-between-attributes", afterAttributeValueQuotedState)
			}

		case selfClosingStartTagState: // https://html.spec.whatwg.org/#self-closing-start-tag-state
			_self.consume()
			switch _self.char {
			case ">":
				_self.token.SetIsSelfClosing(true)
				_self.emitToken()
				_self.state = dataState
			case EOF:
				verbose.Printf(0, "%s: eof-in-tag\n", selfClosingStartTagState)
				return fmt.Errorf("%s: eof-in-tag", selfClosingStartTagState)
			default:
				verbose.Printf(0, "%s: unexpected-solidus-in-tag\n", selfClosingStartTagState)
				return fmt.Errorf("%s: unexpected-solidus-in-tag", selfClosingStartTagState)
			}

		case bogusCommentState: // https://html.spec.whatwg.org/#bogus-comment-state
			_self.consume()
			switch _self.char {
			case ">":
				_self.emitToken()
				_self.state = dataState
			case EOF:
				_self.emitToken()
				_self.token = tokens.NewEndOfFileToken()
				_self.emitToken()
			default:
				_self.token.AppendToData(_self.char)
			}

		case markupDeclarationOpenState: // https://html.spec.whatwg.org/#markup-declaration-open-state
			var n = 2
			_self.peak(n)
			switch _self.peakBuffer {
			case "--":
				_self.token = tokens.NewCommentToken()
				_self.consumeN(n)
				_self.state = commentStartState
			default:
				verbose.Printf(0, "%s: incorrectly-opened-comment\n", markupDeclarationOpenState)
				return fmt.Errorf("%s: incorrectly-opened-comment", markupDeclarationOpenState)
			}

		case commentStartState: // https://html.spec.whatwg.org/#comment-start-state
			_self.consume()
			switch _self.char {
			case "-":
				_self.state = commentStartDashState
			case ">":
				verbose.Printf(0, "%s: abrupt-closing-of-empty-comment\n", commentStartState)
				return fmt.Errorf("%s: abrupt-closing-of-empty-comment", commentStartState)
			default:
				_self.reConsume()
				_self.state = commentState
			}

		case commentStartDashState: // https://html.spec.whatwg.org/#comment-start-dash-state
			_self.consume()
			switch _self.char {
			case "-":
				_self.state = commentEndState
			case ">":
				verbose.Printf(0, "%s: abrupt-closing-of-empty-comment\n", commentStartDashState)
				return fmt.Errorf("%s: abrupt-closing-of-empty-comment", commentStartDashState)
			case EOF:
				verbose.Printf(0, "%s: eof-in-comment\n", commentStartDashState)
				return fmt.Errorf("%s: eof-in-comment", commentStartDashState)
			default:
				_self.token.AppendToData("-")
				_self.reConsume()
				_self.state = commentState
			}

		case commentState: // https://html.spec.whatwg.org/#comment-state
			_self.consume()
			switch _self.char {
			case "<":
				_self.token.AppendToData(_self.char)
				_self.state = commentLessThanSignState
			case "-":
				_self.state = commentEndDashState
			case EOF:
				verbose.Printf(0, "%s: eof-in-comment\n", commentState)
				return fmt.Errorf("%s: eof-in-comment", commentState)
			default:
				_self.token.AppendToData(_self.char)
			}

		case commentLessThanSignState: // https://html.spec.whatwg.org/#comment-less-than-sign-state
			_self.consume()
			switch _self.char {
			case "!":
				_self.token.AppendToData(_self.char)
				_self.state = commentLessThanSignBangState
			case "<":
				_self.token.AppendToData(_self.char)
			default:
				_self.reConsume()
				_self.state = commentState
			}

		case commentLessThanSignBangState: // https://html.spec.whatwg.org/#comment-less-than-sign-bang-state
			_self.consume()
			switch _self.char {
			case "-":
				_self.state = commentLessThanSignBangDashState
			default:
				_self.reConsume()
				_self.state = commentState
			}

		case commentLessThanSignBangDashState: // https://html.spec.whatwg.org/#comment-less-than-sign-bang-dash-state
			_self.consume()
			switch _self.char {
			case "-":
				_self.state = commentLessThanSignBangDashDashState
			default:
				_self.reConsume()
				_self.state = commentState
			}

		case commentLessThanSignBangDashDashState: // https://html.spec.whatwg.org/#comment-less-than-sign-bang-dash-dash-state
			_self.consume()
			switch _self.char {
			case ">":
				_self.reConsume()
				_self.state = commentState
			case EOF:
				_self.reConsume()
				_self.state = commentState
			default:
				verbose.Printf(0, "%s: nested-comment\n", commentLessThanSignBangDashDashState)
				return fmt.Errorf("%s: nested-comment", commentLessThanSignBangDashDashState)
			}

		case commentEndDashState: // https://html.spec.whatwg.org/#comment-end-dash-state
			_self.consume()
			switch _self.char {
			case "-":
				_self.state = commentEndState
			case EOF:
				verbose.Printf(0, "%s: eof-in-comment\n", commentEndDashState)
				return fmt.Errorf("%s: eof-in-comment", commentEndDashState)
			default:
				_self.token.AppendToData("-")
				_self.reConsume()
				_self.state = commentState
			}

		case commentEndState: // https://html.spec.whatwg.org/#comment-end-state
			_self.consume()
			switch _self.char {
			case ">":
				_self.emitToken()
				_self.state = dataState
			case "!":
				_self.state = commentEndBangState
			case "-":
				_self.token.AppendToData("-")
			case EOF:
				verbose.Printf(0, "%s: eof-in-comment\n", commentEndState)
				return fmt.Errorf("%s: eof-in-comment", commentEndState)
			default:
				_self.token.AppendToData("-")
				_self.reConsume()
				_self.state = commentState
			}

		case commentEndBangState: // https://html.spec.whatwg.org/#comment-end-bang-state
			_self.consume()
			switch _self.char {
			case "-":
				_self.token.AppendToData("--!")
				_self.state = commentEndDashState
			case ">":
				verbose.Printf(0, "%s: incorrectly-closed-comment\n", commentEndBangState)
				return fmt.Errorf("%s: incorrectly-closed-comment", commentEndBangState)
			case EOF:
				verbose.Printf(0, "%s: eof-in-comment\n", commentEndBangState)
				return fmt.Errorf("%s: eof-in-comment", commentEndBangState)
			default:
				_self.token.AppendToData("--!")
				_self.reConsume()
				_self.state = commentState
			}
		}
	}
	return nil
}
