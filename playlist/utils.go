package playlist

import (
	"fmt"
	"io"
	"strings"
)

func isIntegerFloat64(v float64) bool {
	return float64(int64(v)) == v
}

func checkAttributeName(name string) (err error) {
	// RFC 8216, 4.2
	for _, r := range name {
		switch {
		case r == '-':
		case 'A' <= r && r <= 'Z':
		case '0' <= r && r <= '9':
		default:
			return errInvalidAttributeName
		}
	}
	return
}

func splitAttributes(s string, maxnum int) []string {
	if s == "" {
		return []string{""}
	}

	if maxnum <= 0 {
		maxnum = strings.Count(s, ",") + 1
	}

	ss := make([]string, 0, maxnum)

	var readstr bool
	var start int

LOOP:
	for i, c := range s {
		// (xgf): We cannot directly use strings.Split(s, ",")，
		// becuase the attribute value may contain the comma character ",".
		switch c {
		case '"':
			readstr = !readstr
			if readstr {
				continue LOOP
			}

		case ',':
			if readstr {
				continue LOOP
			}

			if len(ss) == maxnum-1 {
				ss = append(ss, s[start:])
				start = len(s)
				break LOOP
			} else {
				ss = append(ss, s[start:i])
				start = i + 1
			}
		}
	}

	if start < len(s) {
		ss = append(ss, s[start:])
	}

	return ss
}

func parseAttribute(s string, name, value *string) (err error) {
	if index := strings.IndexByte(s, '='); index < 0 {
		return errInvalidAttribute
	} else {
		switch *name, *value = s[:index], s[index+1:]; {
		case *name == "":
			return errInvalidAttributeName

		case *value == "":
			return errInvalidAttributeValue
		}
	}

	var _name _UnquotedString
	if err = _name.decode(*name); err != nil {
		return fmt.Errorf("%w: %w", errInvalidAttributeName, err)
	} else if err = checkAttributeName(_name.get()); err != nil {
		return err
	}

	*name = _name.get()
	return
}

func iterAttributes(s string, maxnum int, fn func(name, value string) (err error)) (err error) {
	items := splitAttributes(s, maxnum)
	for _, item := range items {
		var name, value string
		if err = parseAttribute(item, &name, &value); err != nil {
			return
		}

		if err = fn(name, value); err != nil {
			err = fmt.Errorf("%s: %w", name, err)
			return
		}
	}
	return
}

func tryWrite[T _Value](w io.Writer, err error, value T) error {
	if err != nil || value.IsZero() {
		return err
	}
	return value.encode(w)
}

func tryWriteAny(w io.Writer, err error, values ...any) error {
	if err != nil || len(values) == 0 {
		return err
	}
	for _, v := range values {
		switch _v := v.(type) {
		case string:
			err = tryWriteString(w, err, _v)

		case _Value:
			err = tryWrite(w, err, _v)

		default:
			panic(fmt.Errorf("unsupported type %T", v))
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func tryWriteString(w io.Writer, err error, s string) error {
	if err != nil || s == "" {
		return err
	}

	_, err = io.WriteString(w, s)
	return err
}

func tryWriteAttr[T _Value](w io.Writer, err error, name string, value T) error {
	if err != nil || value.IsZero() {
		return err
	}

	err = checkAttributeName(name)

	err = tryWrite(w, err, _UnquotedString(name))
	err = tryWriteString(w, err, "=")
	err = tryWrite(w, err, value)

	if err != nil {
		err = fmt.Errorf("%s: %w", name, err)
	}

	return err
}

func tryWriteAttrs(w io.Writer, err error, first bool, attrs ..._Attr) error {
	if err != nil {
		return err
	}

	var i int
	for _, attr := range attrs {
		if !attr.IsZero() {
			if !first || i > 0 {
				err = tryWriteString(w, err, ",")
			}

			err = tryWriteAttr(w, err, attr.Name.get(), attr.Value)
			if err != nil {
				break
			}

			i++
		}
	}
	return err
}

func tryWriteTag[T _Value](w io.Writer, err error, tag Tag, attr T) error {
	if err != nil || attr.IsZero() {
		return err
	}

	err = tryWriteString(w, err, string(tag))
	if !_isbool(attr) {
		err = tryWriteString(w, err, ":")
		err = tryWrite(w, err, attr)
	}
	err = tryWriteString(w, err, "\n")

	if err != nil {
		err = fmt.Errorf("%s: %w", tag, err)
	}

	return err
}

func _isbool(v _Value) bool {
	_, ok := v.(_Bool)
	return ok
}
