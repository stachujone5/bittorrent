package bencode

import (
	"errors"
	"strconv"
)

type parser struct {
	data []byte
	pos  int
}

func Unmarshal(data []byte) (any, error) {
	p := &parser{data: data, pos: 0}
	return p.parse()
}

func (p *parser) parse() (any, error) {
	if p.pos >= len(p.data) {
		return nil, errors.New("unexpected end of data")
	}

	switch p.data[p.pos] {
	case 'i':
		return p.parseInt()
	case 'l':
		return p.parseList()
	case 'd':
		return p.parseDict()
	default:
		if p.data[p.pos] >= '0' && p.data[p.pos] <= '9' {
			return p.parseString()
		}
		return nil, errors.New("invalid bencode format")
	}
}

func (p *parser) parseInt() (int, error) {
	// integer: i<value>e
	if p.data[p.pos] != 'i' {
		return 0, errors.New("expected 'i'")
	}
	p.pos++ // skip 'i'
	start := p.pos
	for p.pos < len(p.data) && p.data[p.pos] != 'e' {
		p.pos++
	}
	if p.pos >= len(p.data) {
		return 0, errors.New("unterminated integer")
	}
	valStr := string(p.data[start:p.pos])
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, err
	}
	p.pos++ // skip 'e'
	return val, nil
}

func (p *parser) parseString() (string, error) {
	// string: <length>:<data>
	start := p.pos
	for p.pos < len(p.data) && p.data[p.pos] != ':' {
		p.pos++
	}
	if p.pos >= len(p.data) {
		return "", errors.New("unterminated string length")
	}
	lenStr := string(p.data[start:p.pos])
	length, err := strconv.Atoi(lenStr)
	if err != nil {
		return "", err
	}
	p.pos++ // skip ':'
	if p.pos+length > len(p.data) {
		return "", errors.New("string data too short")
	}
	str := string(p.data[p.pos : p.pos+length])
	p.pos += length
	return str, nil
}

func (p *parser) parseList() ([]any, error) {
	if p.data[p.pos] != 'l' {
		return nil, errors.New("expected 'l'")
	}
	p.pos++ // skip 'l'
	var list []any
	for p.pos < len(p.data) && p.data[p.pos] != 'e' {
		elem, err := p.parse()
		if err != nil {
			return nil, err
		}
		list = append(list, elem)
	}
	if p.pos >= len(p.data) {
		return nil, errors.New("unterminated list")
	}
	p.pos++ // skip 'e'
	return list, nil
}

func (p *parser) parseDict() (map[string]any, error) {
	if p.data[p.pos] != 'd' {
		return nil, errors.New("expected 'd'")
	}
	p.pos++ // skip 'd'
	dict := make(map[string]any)
	for p.pos < len(p.data) && p.data[p.pos] != 'e' {
		key, err := p.parseString()
		if err != nil {
			return nil, err
		}
		val, err := p.parse()
		if err != nil {
			return nil, err
		}
		dict[key] = val
	}
	if p.pos >= len(p.data) {
		return nil, errors.New("unterminated dict")
	}
	p.pos++ // skip 'e'
	return dict, nil
}
