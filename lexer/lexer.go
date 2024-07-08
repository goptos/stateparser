package lexer

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/goptos/stateparser/lexer/tokens"
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
	beforeTextCodeState                  string = "beforeTextCodeState"
	textCodeState                        string = "textCodeState"
	afterTextCodeState                   string = "afterTextCodeState"
	beforeAttributeValueCodeState        string = "beforeAttributeValueCodeState"
	attributeValueCodeState              string = "attributeValueCodeState"
	afterAttributeValueCodeState         string = "afterAttributeValueCodeState"
	endOfFileState                       string = "endOfFileState"
	dataState                            string = "dataState"
	beforeTextState                      string = "beforeTextState"
	textState                            string = "textState"
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
	codeIndentCount       int
	char                  string
	chars                 []string
	curser                int
	KeywordAttributeNames map[string]interface{}
	length                int
	lineNumber            int
	lineNumberMap         map[int]int
	peakBuffer            string
	_rune                 rune
	rubbishBuffer         string
	runes                 []rune
	Source                string
	state                 string
	token                 tokens.Token
	Tokens                []tokens.Token
}

func New(source string) *Lexer {
	var runes, chars = intoArray(source)
	return &Lexer{
		codeIndentCount:       0,
		char:                  chars[0],
		chars:                 chars,
		curser:                0,
		KeywordAttributeNames: make(map[string]interface{}),
		length:                len(chars),
		lineNumber:            1,
		lineNumberMap:         make(map[int]int),
		peakBuffer:            "",
		rubbishBuffer:         "",
		_rune:                 runes[0],
		runes:                 runes,
		state:                 dataState,
		Source:                source,
		token:                 nil,
		Tokens:                []tokens.Token{}}
}

func (_self *Lexer) clearRubbishBuffer() {
	_self.rubbishBuffer = ""
}

func (_self *Lexer) appendToRubbishBuffer(s string) {
	_self.rubbishBuffer = _self.rubbishBuffer + s
}

func (_self *Lexer) flushRubbishBufferToToken() {
	_self.token.AppendToData(_self.rubbishBuffer)
	_self.rubbishBuffer = ""
}

func (_self *Lexer) consume() {
	if _self.curser+1 >= _self.length {
		_self.char = EOF
		_self.lineNumber++
	} else {
		_self.char = _self.chars[_self.curser]
		_self._rune = _self.runes[_self.curser]
		_self.curser++
	}
	_self.lineNumberMap[_self.lineNumber]++
	if _self.char == "\n" {
		_self.lineNumber++
	}
}

func (_self *Lexer) reConsume() {
	if _self.char == "\n" {
		_self.lineNumber--
	}
	_self.lineNumberMap[_self.lineNumber]--
	_self.curser--
}

func (_self *Lexer) consumeN(n int) {
	for n > 0 {
		_self.consume()
		n--
	}
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

func (_self *Lexer) emitToken() {
	var position = _self.token.GetPosition()
	position.EndLine = _self.lineNumber
	position.EndColumn = _self.lineNumberMap[_self.lineNumber]
	switch _self.token.GetType() {
	case tokens.StartTag:
		position.StartColumn--
	case tokens.EndTag:
		position.StartColumn -= 2
	case tokens.Comment:
		position.StartColumn--
	case tokens.Text:
		position.EndLine = position.StartLine
		position.EndColumn = position.StartColumn
		for _, r := range _self.token.GetData() {
			position.EndColumn++
			if string(r) == "\n" {
				position.EndLine++
				position.EndColumn = 1
			}
		}
		position.EndColumn--
	case tokens.Code:
	case tokens.EndOfFile:
		position.StartColumn--
		position.EndColumn--
	}
	_self.token.SetPosition(position)
	if verbose.Level >= 3 {
		_self.token.Print()
	}
	_self.Tokens = append(_self.Tokens, _self.token)
	_self.token = nil
}

func (_self *Lexer) Tokenise() error {
	verbose.Printf(4, "::: Lexer.Tokenise() :::\n")
	verbose.Printf(4, "KeywordAttributeNames:\n%v\n", _self.KeywordAttributeNames)
	for _self.state != endOfFileState {
		switch _self.state {

		case dataState: // https://html.spec.whatwg.org/#data-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			if isAsciiWhiteSpace(_self._rune) {
				_self.ignore()
				continue
			}
			switch _self.char {
			case "{":
				_self.token = tokens.NewCodeToken(_self.lineNumber, _self.lineNumberMap[_self.lineNumber])
				_self.codeIndentCount = 0
				_self.state = beforeTextCodeState
			case "<":
				_self.state = tagOpenState
			case EOF:
				_self.token = tokens.NewEndOfFileToken(_self.lineNumber, _self.lineNumberMap[_self.lineNumber])
				_self.emitToken()
				_self.state = endOfFileState
			default:
				_self.token = tokens.NewTextToken(_self.lineNumber, _self.lineNumberMap[_self.lineNumber])
				_self.reConsume()
				_self.state = textState
			}

		case textState: // NOT IN SPEC
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			if isAsciiWhiteSpace(_self._rune) {
				_self.appendToRubbishBuffer(_self.char)
				continue
			}
			switch _self.char {
			case "{":
				_self.clearRubbishBuffer()
				_self.emitToken()
				_self.reConsume()
				_self.state = dataState
			case "<":
				_self.clearRubbishBuffer()
				_self.emitToken()
				_self.reConsume()
				_self.state = dataState
			case EOF:
				_self.clearRubbishBuffer()
				_self.emitToken()
				_self.token = tokens.NewEndOfFileToken(_self.lineNumber, _self.lineNumberMap[_self.lineNumber])
				_self.emitToken()
				_self.state = endOfFileState
			default:
				_self.flushRubbishBufferToToken()
				_self.token.AppendToData(_self.char)
			}

		case tagOpenState: // https://html.spec.whatwg.org/#tag-open-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			if isAsciiAlpha(_self._rune) {
				_self.token = tokens.NewStartTagToken(_self.lineNumber, _self.lineNumberMap[_self.lineNumber])
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
				verbose.Printf(0, "error in %s: unexpected-question-mark-instead-of-tag-name\n", _self.state)
				return fmt.Errorf("error in %s: unexpected-question-mark-instead-of-tag-name", _self.state)
			case EOF:
				verbose.Printf(0, "error in %s: eof-before-tag-name\n", _self.state)
				return fmt.Errorf("error in %s: eof-before-tag-name", _self.state)
			default:
				verbose.Printf(0, "error in %s: invalid-first-character-of-tag-name\n", _self.state)
				return fmt.Errorf("error in %s: invalid-first-character-of-tag-name", _self.state)
			}

		case endTagOpenState: // https://html.spec.whatwg.org/#tag-open-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			if isAsciiAlpha(_self._rune) {
				_self.token = tokens.NewEndTagToken(_self.lineNumber, _self.lineNumberMap[_self.lineNumber])
				_self.reConsume()
				_self.state = tagNameState
				continue
			}
			switch _self.char {
			case ">":
				verbose.Printf(0, "error in %s: missing-end-tag-name\n", _self.state)
				return fmt.Errorf("error in %s: missing-end-tag-name", _self.state)

			case EOF:
				verbose.Printf(0, "error in %s: eof-before-tag-name\n", _self.state)
				return fmt.Errorf("error in %s: eof-before-tag-name", _self.state)

			default:
				verbose.Printf(0, "error in %s: invalid-first-character-of-tag-name\n", _self.state)
				return fmt.Errorf("error in %s: invalid-first-character-of-tag-name", _self.state)
			}

		case tagNameState: // https://html.spec.whatwg.org/#tag-name-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			if isAsciiWhiteSpace(_self._rune) {
				_self.state = beforeAttributeNameState
				continue
			}
			if isAsciiUpperAlpha(_self._rune) {
				if runeCount(_self.token.GetName()) == 0 {
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
				verbose.Printf(0, "error in %s: eof-in-tag\n", _self.state)
				return fmt.Errorf("error in %s: eof-in-tag", _self.state)
			default:
				_self.token.AppendToName(_self.char)
			}

		case beforeAttributeNameState: // https://html.spec.whatwg.org/#before-attribute-name-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
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
				verbose.Printf(0, "error in %s: unexpected-equals-sign-before-attribute-name\n", _self.state)
				return fmt.Errorf("error in %s: unexpected-equals-sign-before-attribute-name", _self.state)
			default:
				_self.token.NewAttribute(_self.lineNumber, _self.lineNumberMap[_self.lineNumber])
				_self.reConsume()
				_self.state = attributeNameState
			}

		case attributeNameState: // https://html.spec.whatwg.org/#attribute-name-state
			for mapK := range _self.KeywordAttributeNames {
				_self.peak(runeCount(mapK))
				if _self.peakBuffer == mapK {
					_self.token.SetAttributeType(tokens.KeywordAttribute)
				}
			}
			_self.peak(runeCount("on"))
			if _self.peakBuffer == "on" {
				_self.token.SetAttributeType(tokens.EventAttribute)
			}
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			if isAsciiWhiteSpace(_self._rune) {
				_self.token.SetAttributeType(tokens.ArgumentAttribute)
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
				verbose.Printf(0, "error in %s: unexpected-character-in-attribute-name %s\n", _self.state, _self.char)
				return fmt.Errorf("error in %s: unexpected-character-in-attribute-name %s", _self.state, _self.char)
			case "'":
				verbose.Printf(0, "error in %s: unexpected-character-in-attribute-name %s\n", _self.state, _self.char)
				return fmt.Errorf("error in %s: unexpected-character-in-attribute-name %s", _self.state, _self.char)
			case "<":
				verbose.Printf(0, "error in %s: unexpected-character-in-attribute-name %s\n", _self.state, _self.char)
				return fmt.Errorf("error in %s: unexpected-character-in-attribute-name %s", _self.state, _self.char)
			case ":":
				if _self.token.GetAttributeType() != tokens.EventAttribute {
					_self.token.SetAttributeType(tokens.DynamicAttribute)
				}
				_self.token.AppendToAttributeName(_self.char)
			default:
				_self.token.AppendToAttributeName(_self.char)
			}

		case afterAttributeNameState: // https://html.spec.whatwg.org/#after-attribute-name-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
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
				verbose.Printf(0, "error in %s: eof-in-tag\n", _self.state)
				return fmt.Errorf("error in %s: eof-in-tag", _self.state)
			default:
				_self.token.NewAttribute(_self.lineNumber, _self.lineNumberMap[_self.lineNumber])
				_self.reConsume()
				_self.state = attributeNameState
			}

		case beforeAttributeValueState: // https://html.spec.whatwg.org/#before-attribute-value-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			if isAsciiWhiteSpace(_self._rune) {
				_self.ignore()
				continue
			}
			switch _self.char {
			case "{":
				_self.codeIndentCount = 0
				_self.state = beforeAttributeValueCodeState
			case `"`:
				_self.state = attributeValueDoubleQuotedState
			case "'":
				_self.state = attributeValueSingleQuotedState
			case ">":
				verbose.Printf(0, "error in %s: missing-attribute-value\n", _self.state)
				return fmt.Errorf("error in %s: missing-attribute-value", _self.state)
			default:
				_self.reConsume()
				_self.state = attributeValueUnquotedState
			}

		case beforeTextCodeState: // NOT IN SPEC
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			switch _self.char {
			case "{":
				_self.codeIndentCount++
				_self.token.AppendToData(_self.char)
				_self.state = textCodeState
			default:
				_self.reConsume()
				_self.state = textCodeState
			}

		case beforeAttributeValueCodeState: // NOT IN SPEC
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			switch _self.char {
			case "{":
				_self.codeIndentCount++
				_self.token.AppendToAttributeValue(_self.char)
				_self.state = attributeValueCodeState
			default:
				_self.reConsume()
				_self.state = attributeValueCodeState
			}

		case textCodeState: // NOT IN SPEC
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			switch _self.char {
			case "{":
				_self.reConsume()
				_self.state = beforeTextCodeState
			case EOF:
				verbose.Printf(0, "error in %s: eof-in-code\n", _self.state)
				return fmt.Errorf("error in %s: eof-in-code", _self.state)
			case "}":
				_self.reConsume()
				_self.state = afterTextCodeState
			default:
				_self.token.AppendToData(_self.char)
			}

		case attributeValueCodeState: // NOT IN SPEC
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			switch _self.char {
			case "{":
				_self.reConsume()
				_self.state = beforeAttributeValueCodeState
			case EOF:
				verbose.Printf(0, "error in %s: eof-in-code\n", _self.state)
				return fmt.Errorf("error in %s: eof-in-code", _self.state)
			case "}":
				_self.reConsume()
				_self.state = afterAttributeValueCodeState
			default:
				_self.token.AppendToAttributeValue(_self.char)
			}

		case afterTextCodeState: // NOT IN SPEC
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			switch _self.char {
			case "}":
				if _self.codeIndentCount <= 0 {
					_self.emitToken()
					_self.state = dataState
					continue
				}
				_self.codeIndentCount--
				_self.token.AppendToData(_self.char)
				_self.state = textCodeState
			}

		case afterAttributeValueCodeState: // NOT IN SPEC
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			switch _self.char {
			case "}":
				if _self.codeIndentCount <= 0 {
					_self.state = afterAttributeValueQuotedState
					continue
				}
				_self.codeIndentCount--
				_self.token.AppendToAttributeValue(_self.char)
				_self.state = attributeValueCodeState
			}

		case attributeValueDoubleQuotedState: // https://html.spec.whatwg.org/#attribute-value-(double-quoted)-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			switch _self.char {
			case `"`:
				_self.state = afterAttributeValueQuotedState
			case EOF:
				verbose.Printf(0, "error in %s: eof-in-tag\n", _self.state)
				return fmt.Errorf("error in %s: eof-in-tag", _self.state)
			default:
				_self.token.AppendToAttributeValue(_self.char)
			}

		case attributeValueSingleQuotedState: // https://html.spec.whatwg.org/#attribute-value-(single-quoted)-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			switch _self.char {
			case "'":
				_self.state = afterAttributeValueQuotedState
			case EOF:
				verbose.Printf(0, "error in %s: eof-in-tag\n", _self.state)
				return fmt.Errorf("error in %s: eof-in-tag", _self.state)
			default:
				_self.token.AppendToAttributeValue(_self.char)
			}

		case attributeValueUnquotedState: // https://html.spec.whatwg.org/#attribute-value-(unquoted)-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			if isAsciiWhiteSpace(_self._rune) {
				_self.state = beforeAttributeNameState
				continue
			}
			switch _self.char {
			case ">":
				_self.emitToken()
				_self.state = dataState
			case `"`:
				verbose.Printf(0, "error in %s: unexpected-character-in-unquoted-attribute-value %s\n", _self.state, _self.char)
				return fmt.Errorf("error in %s: unexpected-character-in-unquoted-attribute-value %s", _self.state, _self.char)
			case "'":
				verbose.Printf(0, "error in %s: unexpected-character-in-unquoted-attribute-value %s\n", _self.state, _self.char)
				return fmt.Errorf("error in %s: unexpected-character-in-unquoted-attribute-value %s", _self.state, _self.char)
			case "<":
				verbose.Printf(0, "error in %s: unexpected-character-in-unquoted-attribute-value %s\n", _self.state, _self.char)
				return fmt.Errorf("error in %s: unexpected-character-in-unquoted-attribute-value %s", _self.state, _self.char)
			case "=":
				verbose.Printf(0, "error in %s: unexpected-character-in-unquoted-attribute-value %s\n", _self.state, _self.char)
				return fmt.Errorf("error in %s: unexpected-character-in-unquoted-attribute-value %s", _self.state, _self.char)
			case "`":
				verbose.Printf(0, "error in %s: unexpected-character-in-unquoted-attribute-value %s\n", _self.state, _self.char)
				return fmt.Errorf("error in %s: unexpected-character-in-unquoted-attribute-value %s", _self.state, _self.char)
			case EOF:
				verbose.Printf(0, "error in %s: eof-in-tag\n", _self.state)
				return fmt.Errorf("error in %s: eof-in-tag", _self.state)
			default:
				_self.token.AppendToAttributeValue(_self.char)
			}

		case afterAttributeValueQuotedState: // https://html.spec.whatwg.org/#after-attribute-value-(quoted)-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
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
				verbose.Printf(0, "error in %s: eof-in-tag\n", _self.state)
				return fmt.Errorf("error in %s: eof-in-tag", _self.state)
			default:
				verbose.Printf(0, "error in %s: missing-whitespace-between-attributes\n", _self.state)
				return fmt.Errorf("error in %s: missing-whitespace-between-attributes", _self.state)
			}

		case selfClosingStartTagState: // https://html.spec.whatwg.org/#self-closing-start-tag-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			switch _self.char {
			case ">":
				_self.token.SetIsSelfClosing(true)
				_self.emitToken()
				_self.state = dataState
			case EOF:
				verbose.Printf(0, "error in %s: eof-in-tag\n", _self.state)
				return fmt.Errorf("error in %s: eof-in-tag", _self.state)
			default:
				verbose.Printf(0, "error in %s: unexpected-solidus-in-tag\n", _self.state)
				return fmt.Errorf("error in %s: unexpected-solidus-in-tag", _self.state)
			}

		case bogusCommentState: // https://html.spec.whatwg.org/#bogus-comment-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			switch _self.char {
			case ">":
				_self.emitToken()
				_self.state = dataState
			case EOF:
				_self.emitToken()
				_self.token = tokens.NewEndOfFileToken(_self.lineNumber, _self.lineNumberMap[_self.lineNumber])
				_self.emitToken()
			default:
				_self.token.AppendToData(_self.char)
			}

		case markupDeclarationOpenState: // https://html.spec.whatwg.org/#markup-declaration-open-state
			var n = 2
			_self.peak(n)
			switch _self.peakBuffer {
			case "--":
				verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, "--")
				_self.token = tokens.NewCommentToken(_self.lineNumber, _self.lineNumberMap[_self.lineNumber])
				_self.consumeN(n)
				_self.state = commentStartState
			default:
				verbose.Printf(0, "error in %s: incorrectly-opened-comment\n", _self.state)
				return fmt.Errorf("error in %s: incorrectly-opened-comment", _self.state)
			}

		case commentStartState: // https://html.spec.whatwg.org/#comment-start-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			switch _self.char {
			case "-":
				_self.state = commentStartDashState
			case ">":
				verbose.Printf(0, "error in %s: abrupt-closing-of-empty-comment\n", _self.state)
				return fmt.Errorf("error in %s: abrupt-closing-of-empty-comment", _self.state)
			default:
				_self.reConsume()
				_self.state = commentState
			}

		case commentStartDashState: // https://html.spec.whatwg.org/#comment-start-dash-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			switch _self.char {
			case "-":
				_self.state = commentEndState
			case ">":
				verbose.Printf(0, "error in %s: abrupt-closing-of-empty-comment\n", _self.state)
				return fmt.Errorf("error in %s: abrupt-closing-of-empty-comment", _self.state)
			case EOF:
				verbose.Printf(0, "error in %s: eof-in-comment\n", _self.state)
				return fmt.Errorf("error in %s: eof-in-comment", _self.state)
			default:
				_self.token.AppendToData("-")
				_self.reConsume()
				_self.state = commentState
			}

		case commentState: // https://html.spec.whatwg.org/#comment-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			switch _self.char {
			case "<":
				_self.token.AppendToData(_self.char)
				_self.state = commentLessThanSignState
			case "-":
				_self.state = commentEndDashState
			case EOF:
				verbose.Printf(0, "error in %s: eof-in-comment\n", _self.state)
				return fmt.Errorf("error in %s: eof-in-comment", _self.state)
			default:
				_self.token.AppendToData(_self.char)
			}

		case commentLessThanSignState: // https://html.spec.whatwg.org/#comment-less-than-sign-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
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
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			switch _self.char {
			case "-":
				_self.state = commentLessThanSignBangDashState
			default:
				_self.reConsume()
				_self.state = commentState
			}

		case commentLessThanSignBangDashState: // https://html.spec.whatwg.org/#comment-less-than-sign-bang-dash-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			switch _self.char {
			case "-":
				_self.state = commentLessThanSignBangDashDashState
			default:
				_self.reConsume()
				_self.state = commentState
			}

		case commentLessThanSignBangDashDashState: // https://html.spec.whatwg.org/#comment-less-than-sign-bang-dash-dash-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			switch _self.char {
			case ">":
				_self.reConsume()
				_self.state = commentState
			case EOF:
				_self.reConsume()
				_self.state = commentState
			default:
				verbose.Printf(0, "error in %s: nested-comment\n", _self.state)
				return fmt.Errorf("error in %s: nested-comment", _self.state)
			}

		case commentEndDashState: // https://html.spec.whatwg.org/#comment-end-dash-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			switch _self.char {
			case "-":
				_self.state = commentEndState
			case EOF:
				verbose.Printf(0, "error in %s: eof-in-comment\n", _self.state)
				return fmt.Errorf("error in %s: eof-in-comment", _self.state)
			default:
				_self.token.AppendToData("-")
				_self.reConsume()
				_self.state = commentState
			}

		case commentEndState: // https://html.spec.whatwg.org/#comment-end-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			switch _self.char {
			case ">":
				_self.emitToken()
				_self.state = dataState
			case "!":
				_self.state = commentEndBangState
			case "-":
				_self.token.AppendToData("-")
			case EOF:
				verbose.Printf(0, "error in %s: eof-in-comment\n", _self.state)
				return fmt.Errorf("error in %s: eof-in-comment", _self.state)
			default:
				_self.token.AppendToData("-")
				_self.reConsume()
				_self.state = commentState
			}

		case commentEndBangState: // https://html.spec.whatwg.org/#comment-end-bang-state
			_self.consume()
			verbose.Printf(6, "~ in %s consuming: %q\n", _self.state, _self.char)
			switch _self.char {
			case "-":
				_self.token.AppendToData("--!")
				_self.state = commentEndDashState
			case ">":
				verbose.Printf(0, "error in %s: incorrectly-closed-comment\n", _self.state)
				return fmt.Errorf("error in %s: incorrectly-closed-comment", _self.state)
			case EOF:
				verbose.Printf(0, "error in %s: eof-in-comment\n", _self.state)
				return fmt.Errorf("error in %s: eof-in-comment", _self.state)
			default:
				_self.token.AppendToData("--!")
				_self.reConsume()
				_self.state = commentState
			}
		}
	}
	return nil
}
