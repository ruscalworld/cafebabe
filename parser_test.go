package cafebabe

import (
	"reflect"
	"testing"
)

type simple struct {
	A U1
	B U2
	C U4

	Ignored    U2 `cb:"-"`
	unexported U1
}

type nested struct {
	S1 simple
	S2 simple

	unexported U1
}

type withArr struct {
	A simple
	L U2
	S []simple

	unexported U1
}

type var1 struct {
	A U2
	B U1
}

type var2 struct {
	A U4
}

type variadic struct {
	E U1
	V any `cb:"variadic"`
}

func (v *variadic) PrepareV() reflect.Type {
	switch v.E {
	case 1:
		return reflect.TypeFor[var1]()
	case 2:
		return reflect.TypeFor[var2]()
	}

	return nil
}

type testCase[T any] struct {
	input    []byte
	expected T
}

type structTest[T any] struct {
	builder func() T
	cases   []testCase[T]
}

func (st *structTest[T]) run(i int, t *testing.T) {
	for j, c := range st.cases {
		target := st.builder()
		err := Unmarshal(c.input, &target)
		if err != nil {
			t.Errorf("test %d, case %d: %s", i, j, err)
			continue
		}

		if !reflect.DeepEqual(target, c.expected) {
			t.Errorf("test %d, case %d: expected %v, got %v", i, j, c.expected, target)
		}
	}
}

type test interface {
	run(i int, t *testing.T)
}

func TestUnmarshal(t *testing.T) {
	var tests = []test{
		&structTest[simple]{
			builder: func() simple {
				return simple{}
			},
			cases: []testCase[simple]{
				{
					input:    []byte{1, 0, 2, 0, 0, 0, 4},
					expected: simple{A: 1, B: 2, C: 4},
				},
			},
		},
		&structTest[nested]{
			builder: func() nested {
				return nested{}
			},
			cases: []testCase[nested]{
				{
					input:    []byte{1, 0, 2, 0, 0, 0, 4, 2, 0, 4, 0, 0, 0, 8},
					expected: nested{S1: simple{A: 1, B: 2, C: 4}, S2: simple{A: 2, B: 4, C: 8}},
				},
			},
		},
		&structTest[withArr]{
			builder: func() withArr {
				return withArr{}
			},
			cases: []testCase[withArr]{
				{
					input: []byte{1, 0, 2, 0, 0, 0, 4, 0, 2, 1, 0, 2, 0, 0, 0, 4, 2, 0, 4, 0, 0, 0, 8},
					expected: withArr{
						A: simple{A: 1, B: 2, C: 4},
						L: 2,
						S: []simple{
							{A: 1, B: 2, C: 4},
							{A: 2, B: 4, C: 8},
						},
					},
				},
			},
		},
		&structTest[variadic]{
			builder: func() variadic {
				return variadic{}
			},
			cases: []testCase[variadic]{
				{
					input: []byte{1, 0, 1, 2},
					expected: variadic{
						E: 1,
						V: var1{A: 1, B: 2},
					},
				},
				{
					input: []byte{2, 0, 0, 0, 4},
					expected: variadic{
						E: 2,
						V: var2{A: 4},
					},
				},
			},
		},
	}

	for i, ts := range tests {
		ts.run(i, t)
	}
}
