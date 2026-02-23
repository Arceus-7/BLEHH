# BLEHH

<p align="center">
  <img src="logo.jpeg" alt="BLEHH Logo" width="200">
</p>

## An Esoteric Programming Language

### v1.0 — Official Language Specification

## 1. What is BLEHH?

BLEHH is a minimalist esoteric programming language. It has exactly 5 commands. It operates on a single value called the accumulator, which is always a dice face — a number between 1 and 6. That's it. There is no memory, no variables, no functions, no types.

The twist that makes BLEHH fun and weird: every command behaves differently depending on whether the current accumulator value is ODD (1, 3, 5) or EVEN (2, 4, 6). The same program can produce wildly different results depending on where you start. This is called the Parity Twist.

## 2. Design Philosophy

- Exactly 5 commands. No more, ever.
- Single accumulator. No tape, no stack, no registers.
- Values live between 1 and 6 only (wrapping enforced).
- Parity changes everything. Odd and even values behave differently.
- File extension: `.blehh`
- Any character that is not a command is silently ignored (use this for comments!).

## 3. The Accumulator

The accumulator is the one and only piece of state in a BLEHH program. It starts at 1 at the beginning of every program. It can only ever hold values from 1 to 6 (inclusive). If an operation would take it below 1 it wraps to 6, and if it would go above 6 it wraps to 1.

Examples of wrapping:

- At 6, B command (even) → +2 → wraps to 2  (6 → skip 1 → land on 2)
- At 1, L command (odd)  → -1 → wraps to 6
- At 5, B command (odd)  → +1 → wraps to 6
- At 6, L command (even) → -2 → wraps to 4

## 4. Command Reference

BLEHH has exactly 5 commands:

| Command | Value is ODD (1, 3, 5) | Value is EVEN (2, 4, 6) |
|---------|------------------------|--------------------------|
| B       | Add 1 to accumulator   | Add 2 to accumulator     |
| L       | Subtract 1 from accumulator | Subtract 2 from accumulator |
| O       | Print accumulator as a number | Print letter at acc's alphabet position (2→B, 4→D, 6→F) |
| P       | acc + 1 (cross to even: 1→2, 3→4, 5→6) | acc - 1 (cross to odd: 2→1, 4→3, 6→5) |
| ( )     | Loop while value ≠ 1   | Loop while value ≠ 6     |

Note: Wrapping always stays within 1–6. Values never go outside this range.

## 5. Command Details

### B — Bump

Increases the accumulator. The amount depends on parity at the time of execution.

- Odd  → acc = acc + 1  (wraps 6 → 1)
- Even → acc = acc + 2  (wraps: 5→1, 6→2)

### L — Lower

Decreases the accumulator. The amount depends on parity at the time of execution.

- Odd  → acc = acc - 1  (wraps 1 → 6)
- Even → acc = acc - 2  (wraps: 2→6, 1→5)

### O — Output

Prints the current accumulator value. HOW it prints depends on parity.

- Odd  → prints the number itself      e.g. acc=1 prints: 1
- Even → prints the letter at that position in the alphabet: 2→B, 4→D, 6→F

This means BLEHH programs output a mix of digits (1, 3, 5) and letters (B, D, F) depending on accumulator parity. The alphabet mapping keeps output readable while preserving the parity twist!

### P — Parity Bridge

Crosses the parity boundary. This is the only way to move between the odd and even worlds.

- Odd  → acc = acc + 1  (1→2, 3→4, 5→6)
- Even → acc = acc - 1  (2→1, 4→3, 6→5)

Without P, once you leave odd you can never return. P is the bridge that connects both sides of the die.

### ( ) — Loop

Loops the code between ( and ) repeatedly. The exit condition depends on the parity of the accumulator when the loop was ENTERED.

- Entered while Odd  → loops while acc ≠ 1  (exits when acc reaches 1)
- Entered while Even → loops while acc ≠ 6  (exits when acc reaches 6)

Warning: The exit condition is determined at loop entry. If you enter at an odd value, you loop until you hit 1, even if parity changes mid-loop. Nested loops each evaluate their own condition independently.

## 6. Wrapping Rules

The accumulator always stays in range [1, 6]. Wrapping is sequential — values step one at a time through the range, looping around.

| Operation      | Starting Value | Result |
|----------------|---------------|--------|
| B (odd, +1)    | 6             | 1      |
| B (even, +2)   | 5             | 1      |
| B (even, +2)   | 6             | 2      |
| L (odd, -1)    | 1             | 6      |
| L (even, -2)   | 2             | 6      |
| L (even, -2)   | 1             | 5      |

## 7. Comments & Whitespace

Any character that is not B, L, O, P, (, or ) is completely ignored. This means you can write comments freely inline with your code.

```
BBB this bumps three times OOO then outputs three times
-- start here -- B B B -- now output -- O
```

> **Caution:** Be careful not to use the letters B, L, O, or P (uppercase) or parentheses in your comments — they will be executed as commands!

Whitespace (spaces, newlines, tabs) is also ignored. Format your code however you like.

## 8. Example Programs

### Print '1'

The simplest program. Accumulator starts at 1 (odd), O prints numerically.

```
O
```

### Print '3'

Start at 1 (odd), bump twice (+1 each time since odd), arrive at 3, print.

```
B B O
```

Trace: start=1(odd) → B → 2(even) → B → 4(even) → O → prints ASCII for 4

Hmm! Parity shifted mid-way. This is the trap. Let's try again:

```
B L B O   -- experiment to find the right path to 3
```

This is the puzzle of BLEHH — parity shifts as you move, so planning ahead is essential.

### Infinite loop (careful!)

```
B (B)   -- enter loop at even, loops until 6, but B keeps changing things
```

### Reset and print

```
B B P O   -- bump twice, then reset, then print
```

## 9. Implementation Notes

For anyone building a BLEHH interpreter, here's the minimum you need:

- State: a single integer initialized to 1, constrained to [1, 6].
- Parser: read characters one at a time, skip anything not in {B, L, O, P, (, )}.
- Loop handling: track loop start positions with a stack. On (, push current position. On ), check exit condition based on parity at loop entry (store this when pushing).
- Parity check: (acc % 2 !== 0) = odd, (acc % 2 === 0) = even.
- Wrapping: ((acc - 1 + delta) % 6) + 1 for positive delta, with similar modular arithmetic for negative.
- File extension: `.blehh`
- Encoding: UTF-8 plain text.

No standard library. No imports. No nothing. Just BLEHH.

## 10. Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) 1.21 or later

### Build

```bash
git clone <repo-url>
cd BLEHH
go build -o BLEHH .
```

On Windows this produces `BLEHH.exe`.

### Run a .blehh file

```bash
./BLEHH examples/hello.blehh
```

### Run inline code

```bash
./BLEHH -c "OBOBO"
```

### Set a step limit

Protect against infinite loops by capping the number of steps:

```bash
./BLEHH -max 500 examples/the_trap.blehh
```

The default limit is 1,000,000 steps. If exceeded, the interpreter prints a warning to stderr and exits.

### Example programs

The `examples/` folder contains several programs to explore:

| File | What it does |
|------|-------------|
| `hello.blehh` | Exercises all 5 commands |
| `count_all_faces.blehh` | Tries to print every dice face 1–6 |
| `parity_trap.blehh` | Guess the output — you'll be wrong |
| `roller_coaster.blehh` | Rides the accumulator up and down |
| `yo_yo.blehh` | A loop that exits faster than you expect |
| `staircase.blehh` | Discovers the one-way parity door |
| `odd_forever.blehh` | Proof that odd is a one-move lifespan |
| `nested_loops.blehh` | Nested loops with independent exits |
| `the_trap.blehh` | Deliberately infinite — use `-max`! |
| `dice_roll.blehh` | Even-face cycle: 2 → 4 → 6 → 2 → ... |

### Writing your own

1. Create a file with the `.blehh` extension
2. Write using only the 5 commands: `B`, `L`, `O`, `P`, `(`, `)`
3. Everything else is ignored — use it for comments, but avoid uppercase B, L, O, P and parentheses in comment text
4. Run it: `./BLEHH myprogram.blehh`

---

> *The die has 6 faces. This interpreter has 10 secrets.*
> *Start looking where there's nothing to find.*

BLEHH v1.0  •  An esoteric language spec  •  5 commands, infinite frustration
