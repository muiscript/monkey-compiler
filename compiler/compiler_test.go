package compiler

import (
	"monkey-compiler/ast"
	"monkey-compiler/code"
	"monkey-compiler/lexer"
	"monkey-compiler/object"
	"monkey-compiler/parser"
	"testing"
)

type compilerTestCase struct {
	desc                 string
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	testCases := []compilerTestCase{
		{
			desc:              "1+2",
			input:             "1 + 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "1;2",
			input:             "1; 2",
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "-1",
			input:             "-1;",
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, testCases)
}

func TestBooleanExpression(t *testing.T) {
	testCases := []compilerTestCase{
		{
			desc:              "true",
			input:             "true;",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "false",
			input:             "false;",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "5>3",
			input:             "5 > 3;",
			expectedConstants: []interface{}{5, 3},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "5<3",
			input:             "5 < 3;",
			expectedConstants: []interface{}{3, 5},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "5==3",
			input:             "5 == 3;",
			expectedConstants: []interface{}{5, 3},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "5!=3",
			input:             "5 != 3;",
			expectedConstants: []interface{}{5, 3},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			desc:              "!false",
			input:             "!false;",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpBang),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, testCases)
}

func TestConditional(t *testing.T) {
	testCases := []compilerTestCase{
		{
			desc:              "if-statement-with-true-condition",
			input:             "if (true) { 10 }; 33;",
			expectedConstants: []interface{}{10, 33},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),              // 00
				code.Make(code.OpJumpNotTruthy, 10), // 01
				code.Make(code.OpConstant, 0),       // 04
				code.Make(code.OpJump, 11),          // 07
				code.Make(code.OpNull),              // 10
				code.Make(code.OpPop),               // 11
				code.Make(code.OpConstant, 1),       // 12
				code.Make(code.OpPop),               // 15
			},
		},
		{
			desc:              "if-else-statement-with-true-condition",
			input:             "if (true) { 10 } else { 20 }; 33;",
			expectedConstants: []interface{}{10, 20, 33},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpJumpNotTruthy, 10),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpJump, 13),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, testCases)
}

func TestGlobalLetStatement(t *testing.T) {
	testCases := []compilerTestCase{
		{
			desc: "assignments",
			input: `
			let one = 1;
			let two = 2;
			`,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 1),
			},
		},
		{
			desc: "assignment-and-retrieval",
			input: `
			let one = 1;
			one;
			`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			desc: "assignment-assignment-and-retrieval",
			input: `
			let one = 1;
			let two = one;
			two;
			`,
			expectedConstants: []interface{}{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpSetGlobal, 1),
				code.Make(code.OpGetGlobal, 1),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, testCases)
}

func runCompilerTests(t *testing.T, testCases []compilerTestCase) {
	t.Helper()

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			program := parse(tc.input)

			compiler := New()
			err := compiler.Compile(program)
			if err != nil {
				t.Fatalf("compile error: %s", err.Error())
			}

			byteCode := compiler.ByteCode()

			expectedInstructions := concatInstructions(tc.expectedInstructions)
			if len(byteCode.Instructions) != len(expectedInstructions) {
				t.Fatalf("instruction wrong.\nwant=%s\ngot=%s", expectedInstructions, byteCode.Instructions)
			}
			for i, b := range byteCode.Instructions {
				if b != expectedInstructions[i] {
					t.Fatalf("instruction wrong\nwant=%s\ngot=%s", expectedInstructions, byteCode.Instructions)
					break
				}
			}

			if len(byteCode.Constants) != len(tc.expectedConstants) {
				t.Fatalf("constants wrong. want=%+v, got=%+v", tc.expectedConstants, byteCode.Constants)
			}
			for i, c := range tc.expectedConstants {
				switch c := c.(type) {
				case int:
					testIntegerObject(t, int64(c), byteCode.Constants[i])
				}
			}
		})
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func concatInstructions(instructions []code.Instructions) code.Instructions {
	out := code.Instructions{}
	for _, ins := range instructions {
		out = append(out, ins...)
	}
	return out
}

func testIntegerObject(t *testing.T, expected int64, actual object.Object) {
	t.Helper()

	actualInteger, ok := actual.(*object.Integer)
	if !ok {
		t.Fatalf("could not convert to Integer: %+v", actual)
	}

	if actualInteger.Value != expected {
		t.Fatalf("integer valud wrong. want=%d, got=%d", expected, actualInteger.Value)
	}
}
