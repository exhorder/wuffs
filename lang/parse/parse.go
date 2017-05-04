// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package parse

// TODO: write a formal grammar for the language.

import (
	"fmt"

	a "github.com/google/puffs/lang/ast"
	t "github.com/google/puffs/lang/token"
)

func ParseFile(src []t.Token, m *t.IDMap, filename string) (*a.Node, error) {
	p := &parser{
		src:      src,
		m:        m,
		filename: filename,
	}
	if len(src) > 0 {
		p.lastLine = src[len(src)-1].Line
	}
	return p.parseFile()
}

type parser struct {
	src      []t.Token
	m        *t.IDMap
	filename string
	lastLine uint32
}

func (p *parser) line() uint32 {
	if len(p.src) != 0 {
		return p.src[0].Line
	}
	return p.lastLine
}

func (p *parser) peekID() t.ID {
	if len(p.src) > 0 {
		return p.src[0].ID
	}
	return 0
}

func (p *parser) parseFile() (*a.Node, error) {
	decls := []*a.Node(nil)
	for len(p.src) > 0 {
		d, err := p.parseTopLevelDecl()
		if err != nil {
			return nil, err
		}
		decls = append(decls, d)
	}
	return &a.Node{
		Kind:  a.KFile,
		List0: decls,
	}, nil
}

func (p *parser) parseTopLevelDecl() (*a.Node, error) {
	switch p.src[0].ID {
	case t.IDFunc:
		flags := a.Flags(0)
		p.src = p.src[1:]
		id0, id1, err := p.parseQualifiedIdent()
		if err != nil {
			return nil, err
		}
		inParams, err := p.parseParamList()
		if err != nil {
			return nil, err
		}
		if p.peekID() == t.IDQuestion {
			flags |= a.FlagsSuspendible
			p.src = p.src[1:]
		}
		outParams, err := p.parseParamList()
		if err != nil {
			return nil, err
		}
		block, err := p.parseBlock()
		if err != nil {
			return nil, err
		}
		if x := p.peekID(); x != t.IDSemicolon {
			got := p.m.ByKey(x.Key())
			return nil, fmt.Errorf("parse: expected (implicit) \";\", got %q at %s:%d", got, p.filename, p.line())
		}
		p.src = p.src[1:]
		return &a.Node{
			Kind:  a.KFunc,
			Flags: flags,
			ID0:   id0,
			ID1:   id1,
			List0: inParams,
			List1: outParams,
			List2: block,
		}, nil
	}
	return nil, fmt.Errorf("parse: unrecognized top level declaration at %s:%d", p.filename, p.src[0].Line)
}

// parseQualifiedIdent parses "foo.bar" or "bar".
func (p *parser) parseQualifiedIdent() (t.ID, t.ID, error) {
	x, err := p.parseIdent()
	if err != nil {
		return 0, 0, err
	}

	if p.peekID() != t.IDDot {
		return 0, x, nil
	}
	p.src = p.src[1:]

	y, err := p.parseIdent()
	if err != nil {
		return 0, 0, err
	}
	return x, y, nil
}

func (p *parser) parseIdent() (t.ID, error) {
	if len(p.src) == 0 {
		return 0, fmt.Errorf("parse: expected identifier at %s:%d", p.filename, p.line())
	}
	x := p.src[0]
	if !x.IsIdent() {
		got := p.m.ByKey(x.Key())
		return 0, fmt.Errorf("parse: expected identifier, got %q at %s:%d", got, p.filename, p.line())
	}
	p.src = p.src[1:]
	return x.ID, nil
}

func (p *parser) parseParamList() ([]*a.Node, error) {
	if x := p.peekID(); x != t.IDOpenParen {
		got := p.m.ByKey(x.Key())
		return nil, fmt.Errorf("parse: expected \"(\", got %q for parameter list at %s:%d",
			got, p.filename, p.line())
	}
	p.src = p.src[1:]

	params := []*a.Node(nil)
	for len(p.src) > 0 {
		if p.src[0].ID == t.IDCloseParen {
			p.src = p.src[1:]
			return params, nil
		}

		param, err := p.parseParam()
		if err != nil {
			return nil, err
		}
		params = append(params, param)

		switch x := p.peekID(); x {
		case t.IDCloseParen:
			p.src = p.src[1:]
			return params, nil
		case t.IDComma:
			p.src = p.src[1:]
		default:
			got := p.m.ByKey(x.Key())
			return nil, fmt.Errorf("parse: expected \")\", got %q for parameter list at %s:%d",
				got, p.filename, p.line())
		}
	}
	return nil, fmt.Errorf("parse: expected \")\" for parameter list at %s:%d", p.filename, p.line())
}

func (p *parser) parseParam() (*a.Node, error) {
	id, err := p.parseIdent()
	if err != nil {
		return nil, err
	}
	rhs, err := p.parseType()
	if err != nil {
		return nil, err
	}
	return &a.Node{
		Kind: a.KParam,
		ID0:  id,
		RHS:  rhs,
	}, nil
}

func (p *parser) parseType() (*a.Node, error) {
	switch p.peekID() {
	case t.IDPtr:
		p.src = p.src[1:]
		rhs, err := p.parseType()
		if err != nil {
			return nil, err
		}
		return &a.Node{
			Kind: a.KType,
			ID0:  t.IDPtr,
			RHS:  rhs,
		}, nil
	}

	id0, id1, err := p.parseQualifiedIdent()
	if err != nil {
		return nil, err
	}
	return &a.Node{
		Kind: a.KType,
		ID0:  id0,
		ID1:  id1,
	}, nil
}

func (p *parser) parseBlock() ([]*a.Node, error) {
	if x := p.peekID(); x != t.IDOpenCurly {
		got := p.m.ByKey(x.Key())
		return nil, fmt.Errorf("parse: expected \"{\", got %q for block at %s:%d", got, p.filename, p.line())
	}
	p.src = p.src[1:]

	block := []*a.Node(nil)
	for len(p.src) > 0 {
		if p.src[0].ID == t.IDCloseCurly {
			p.src = p.src[1:]
			return block, nil
		}

		s, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		block = append(block, s)

		if x := p.peekID(); x != t.IDSemicolon {
			got := p.m.ByKey(x.Key())
			return nil, fmt.Errorf("parse: expected (implicit) \";\", got %q at %s:%d", got, p.filename, p.line())
		}
		p.src = p.src[1:]
	}
	return nil, fmt.Errorf("parse: expected \"}\" for block at %s:%d", p.filename, p.line())
}

func (p *parser) parseStatement() (*a.Node, error) {
	// TODO: parse statements other than x = y.

	lhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if x := p.peekID(); x != t.IDEq {
		got := p.m.ByKey(x.Key())
		return nil, fmt.Errorf("parse: expected \"=\", got %q for statement at %s:%d", got, p.filename, p.line())
	}
	p.src = p.src[1:]

	rhs, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &a.Node{
		Kind: a.KAssign,
		LHS:  lhs,
		RHS:  rhs,
	}, nil
}

func (p *parser) parseExpr() (*a.Node, error) {
	lhs, err := p.parseOperand()
	if err != nil {
		return nil, err
	}
	if x := p.peekID(); x.IsBinaryOp() {
		p.src = p.src[1:]
		parseFunc := (*parser).parseOperand
		if x == t.IDAs {
			parseFunc = (*parser).parseType
		}
		rhs, err := parseFunc(p)
		if err != nil {
			return nil, err
		}
		return &a.Node{
			Kind: a.KExpr,
			ID0:  x,
			LHS:  lhs,
			RHS:  rhs,
		}, nil
	}
	return lhs, nil
}

func (p *parser) parseOperand() (*a.Node, error) {
	switch x := p.peekID(); {
	case x == t.IDOpenParen:
		// TODO: parse "(x + y)".

	case x.IsUnaryOp():
		p.src = p.src[1:]
		rhs, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &a.Node{
			Kind: a.KExpr,
			ID0:  x,
			RHS:  rhs,
		}, nil
	}

	id, err := p.parseIdent()
	if err != nil {
		return nil, err
	}

	switch p.peekID() {
	case t.IDOpenParen:
		// TODO: parse "f(i)".
	case t.IDOpenBracket:
		// TODO: parse "a[i]".
	case t.IDDot:
		// TODO: parse "x.y".
	}

	return &a.Node{
		Kind: a.KExpr,
		ID1:  id,
	}, nil
}