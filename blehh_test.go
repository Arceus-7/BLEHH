package main

import (
	"strings"
	"testing"
)

// ─── wrap() tests ───────────────────────────────────────────────────────────

func TestWrap(t *testing.T) {
	tests := []struct {
		name     string
		acc      int
		delta    int
		expected int
	}{
		// Basic forward steps
		{"1+1", 1, 1, 2},
		{"3+2", 3, 2, 5},
		{"5+1", 5, 1, 6},

		// Wrapping forward (past 6 → back to low values)
		{"6+1 wraps to 1", 6, 1, 1},
		{"5+2 wraps to 1", 5, 2, 1},
		{"4+3 wraps to 1", 4, 3, 1},
		{"6+3 wraps to 3", 6, 3, 3},

		// Basic backward steps
		{"3-1", 3, -1, 2},
		{"6-2", 6, -2, 4},

		// Wrapping backward (past 1 → back to high values)
		{"1-1 wraps to 6", 1, -1, 6},
		{"1-2 wraps to 5", 1, -2, 5},
		{"2-3 wraps to 5", 2, -3, 5},

		// Zero delta
		{"no change", 4, 0, 4},

		// Large deltas (multiple wraps)
		{"1+6 full circle", 1, 6, 1},
		{"1+12 two full circles", 1, 12, 1},
		{"1+7 one full + 1", 1, 7, 2},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := wrap(tc.acc, tc.delta)
			if got != tc.expected {
				t.Errorf("wrap(%d, %d) = %d, want %d", tc.acc, tc.delta, got, tc.expected)
			}
		})
	}
}

// ─── isOdd() tests ─────────────────────────────────────────────────────────

func TestIsOdd(t *testing.T) {
	for _, v := range []int{1, 3, 5} {
		if !isOdd(v) {
			t.Errorf("isOdd(%d) should be true", v)
		}
	}
	for _, v := range []int{2, 4, 6} {
		if isOdd(v) {
			t.Errorf("isOdd(%d) should be false", v)
		}
	}
}

// ─── Command tests ─────────────────────────────────────────────────────────

func TestBCommand(t *testing.T) {
	// acc starts at 1 (odd) → B adds 1 → acc = 2
	out, err := interpret("B", 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if out != "" {
		t.Fatalf("B alone should produce no output, got %q", out)
	}

	// acc=1 (odd) → B → acc=2 (even) → B → acc+2=4
	// Verify via O: acc=4 (even) → print letter 'D' (4th letter)
	out, err = interpret("BBO", 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if out != "D" {
		t.Errorf("BBO: expected \"D\", got %q", out)
	}
}

func TestLCommand(t *testing.T) {
	// acc=1 (odd) → L subtracts 1 → wrap(1,-1)=6
	// acc=6 (even) → O → print letter 'F' (6th letter)
	out, err := interpret("LO", 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if out != "F" {
		t.Errorf("LO: expected \"F\", got %q", out)
	}

	// acc=1 → L → 6 (even) → L subtracts 2 → wrap(6,-2)=4
	// acc=4 (even) → O → letter 'D' (4th letter)
	out, err = interpret("LLO", 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if out != "D" {
		t.Errorf("LLO: expected \"D\", got %q", out)
	}
}

func TestOCommand(t *testing.T) {
	// acc=1 (odd) → O prints "1" (decimal)
	out, err := interpret("O", 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if out != "1" {
		t.Errorf("O with acc=1: expected %q, got %q", "1", out)
	}

	// acc=1 → B → acc=2 (even) → O prints letter 'B' (2nd letter)
	out, err = interpret("BO", 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if out != "B" {
		t.Errorf("BO: expected \"B\", got %q", out)
	}
}

func TestPCommand(t *testing.T) {
	// acc=1 (odd) → P → acc+1 = 2 (crosses to even)
	// Then O → prints letter 'B' (2nd)
	out, err := interpret("PO", 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if out != "B" {
		t.Errorf("PO: expected \"B\", got %q", out)
	}

	// acc=1 → B → acc=2 (even) → P → acc-1 = 1 (crosses back to odd)
	// O → prints "1"
	out, err = interpret("BPO", 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if out != "1" {
		t.Errorf("BPO: expected \"1\", got %q", out)
	}

	// acc=1 → B → 2 → B → 4 → P → 3 (odd) → O → prints "3"
	out, err = interpret("BBPO", 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if out != "3" {
		t.Errorf("BBPO: expected \"3\", got %q", out)
	}
}

// ─── Loop tests ─────────────────────────────────────────────────────────────

func TestLoopOddEntry(t *testing.T) {
	// acc=1 (odd). Enter loop, enteredOdd=true.
	// Inside loop: B → acc=2 → B → acc=4 → B → acc=6 → B → acc=2 ...
	// Wait — we need acc to reach 1 for exit.
	// acc=1 → (  [enteredOdd=true, exit when acc==1]
	//   B → 2, B → 4, B → 6, B → 2 — this cycles.
	// Actually let's use L to go down.
	// acc=1 → B → 2 → ... let me think of a terminating case.
	//
	// Simpler: acc=1 → B → 2 (odd→+1).
	//          Loop: acc=2 → (  [enteredEven=true, exit when acc==6]
	//            B → 4, B → 6 → exit!
	// So "B(BB)O" should leave acc=6 and print 'F'.
	out, err := interpret("B(BB)O", 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if out != "F" {
		t.Errorf("B(BB)O: expected \"F\", got %q", out)
	}
}

func TestLoopEvenEntry(t *testing.T) {
	// acc=1 → B → 2 (even). Enter loop [enteredEven, exit when acc==6].
	// BB inside: 2→4→6. At ')' acc==6, so exit.
	// Output O: acc=6 (even) → letter 'F'.
	out, err := interpret("B(BB)O", 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if out != "F" {
		t.Errorf("expected \"F\", got %q", out)
	}
}

func TestLoopRepeats(t *testing.T) {
	// We want a loop that actually repeats at least once.
	// acc=1 → B → 2. Enter loop [even, exit when acc==6].
	// Iteration 1: B → 4. ')' → acc!=6, repeat.
	// Iteration 2: B → 6. ')' → acc==6, exit.
	// O → letter 'F'
	out, err := interpret("B(B)O", 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if out != "F" {
		t.Errorf("B(B)O: expected \"F\", got %q", out)
	}
}

// ─── Non-command characters are ignored ─────────────────────────────────────

func TestIgnoredCharacters(t *testing.T) {
	// "O" with extra noise around it should just print "1".
	out, err := interpret("  O \t\n xyz! ", 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if out != "1" {
		t.Errorf("expected %q, got %q", "1", out)
	}
}

// ─── Step limit ─────────────────────────────────────────────────────────────

func TestStepLimit(t *testing.T) {
	// "(B)" starting at acc=1:
	//   ( [enteredOdd, exit when acc==1]
	//   B → 2,  ')' → acc!=1 → repeat
	//   B → 4,  ')' → acc!=1 → repeat
	//   B → 6,  ')' → acc!=1 → repeat
	//   B → 2,  ... infinite cycle: 2→4→6→2→4→6 never hits 1
	// The step limit should kick in.
	_, err := interpret("(B)", 100, false)
	if err == nil {
		t.Fatal("expected step limit error for infinite loop (B)")
	}
	if !strings.Contains(err.Error(), "step limit") {
		t.Errorf("expected step limit error, got: %v", err)
	}
}

// ─── Multiple O outputs ────────────────────────────────────────────────────

func TestMultipleOutputs(t *testing.T) {
	// acc=1(odd) O→"1", B→2(even) O→"B", B→4(even) O→"D"
	out, err := interpret("OBOBO", 0, false)
	if err != nil {
		t.Fatal(err)
	}
	expected := "1BD"
	if out != expected {
		t.Errorf("OBOBO: expected %q, got %q", expected, out)
	}
}

// ─── Empty program ─────────────────────────────────────────────────────────

func TestEmptyProgram(t *testing.T) {
	out, err := interpret("", 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if out != "" {
		t.Errorf("empty program should produce no output, got %q", out)
	}
}
