// =============================================================================
// BLEHH Interpreter — A clean Go implementation of the BLEHH esoteric language.
//
// BLEHH has a single integer accumulator (init = 1, constrained to [1, 6]).
// Commands: B, L, O, P, (, )
// All other characters are silently ignored.
//
// I built this esoteric language interpreter because my gf has not replied since 2 hours so I had to do something fun
// =============================================================================

package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ─── Default configuration ──────────────────────────────────────────────────

const defaultMaxSteps = 1_000_000

// ─── Easter Egg Data ────────────────────────────────────────────────────────

var snarkyMessages = []string{
	"The die judges you silently.",
	"Was that supposed to do something?",
	"Even the accumulator is confused.",
	"BLEHH disapproves.",
	"That's not how dice work, but ok.",
	"Your code has the energy of a wet sock.",
	"The die has seen better programs.",
	"Somewhere, a computer scientist just cried.",
}

var stepLimitMessages = []string{
	"I gave you a million steps and THIS is what you do?",
	"Congratulations, you've created nothing.",
	"Even the die is tired of rolling.",
	"Infinity called. It wants its loop back.",
	"Your program ran longer than your attention span.",
	"The accumulator begs for mercy.",
	"Did you really think this would terminate?",
	"Step limit reached. Hope was lost long ago.",
}

var zenKoans = []string{
	"The unrolled die contains all faces.",
	"In emptiness, the accumulator finds peace.",
	"To BLEHH nothing is to BLEHH everything.",
	"The blank program has already finished. Have you?",
	"No commands, no bugs. Perfection.",
	"The wisest BLEHH program is the one never written.",
}

var existentialSuffixes = []string{
	" (but does it matter?)",
	" (in the grand scheme of things)",
	" (or so the die claims)",
	" (if you even believe in numbers)",
	" (the void stares back)",
	" (temporarily)",
}

const rickRoll = `Never gonna give you up
Never gonna let you down
Never gonna run around and desert you
Never gonna make you cry
Never gonna say goodbye
Never gonna tell a lie and hurt you`

// ─── Wrapping helper ────────────────────────────────────────────────────────

// wrap applies a delta to the accumulator and keeps the result in [1, 6].
//
// Formula:  ((acc - 1 + delta) % 6 + 6) % 6 + 1
// This is almost like working in the group (Z_6, +)
// The inner  (acc-1+delta)  maps the 1-based value into 0-based space,
// the double-mod-plus-6 handles negative remainders (Go's % can be negative),
// and the final +1 maps back to 1-based.
//
// i nutted so hard figuring out this modular arithmetic that my keyboard is still sticky.
func wrap(acc, delta int) int {
	return ((acc-1+delta)%6+6)%6 + 1
}

// isOdd returns true when the accumulator is an odd number (1, 3, 5).
func isOdd(acc int) bool {
	return acc%2 != 0
}

// ─── Loop-stack entry ───────────────────────────────────────────────────────

// loopEntry records the state when a '(' is encountered so we can decide
// whether to repeat or exit the loop when the matching ')' is reached.
type loopEntry struct {
	// pos is the index of the '(' character in the code string.
	pos int
	// enteredOdd records the parity of the accumulator at loop entry.
	// If true  → the loop repeats while acc ≠ 1.
	// If false → the loop repeats while acc ≠ 6.
	enteredOdd bool
}

// ─── Interpreter ────────────────────────────────────────────────────────────

// interpret executes a BLEHH program and returns its collected output.
// maxSteps limits the total number of command-steps to guard against infinite
// loops; pass 0 to use the default limit.
// existential enables existential commentary mode.
func interpret(code string, maxSteps int, existential bool) (string, error) {
	if maxSteps <= 0 {
		maxSteps = defaultMaxSteps
	}

	// ── Konami Code Easter Egg ─────────────────────────────────────────
	// If the program contains BBLLBBLL, double the step limit.
	if strings.Contains(code, "BBLLBBLL") {
		maxSteps *= 2
		fmt.Fprintln(os.Stderr, "🕹️  +30 lives! Step limit doubled.")
	}

	acc := 1              // accumulator — always in [1, 6], starts at 1
	ip := 0               // instruction pointer — current position in code
	steps := 0            // total steps; each one is a bad decision made at 2am
	var stack []loopEntry // call stack. yes it's a stack. yes i'm compensating for something.
	var out strings.Builder

	for ip < len(code) {
		ch := code[ip]

		switch ch {

		// ── B command ───────────────────────────────────────────────
		// Odd accumulator  → add 1
		// Even accumulator → add 2
		case 'B':
			if isOdd(acc) {
				acc = wrap(acc, 1)
			} else {
				acc = wrap(acc, 2)
			}

		// ── L command ───────────────────────────────────────────────
		// Odd accumulator  → subtract 1
		// Even accumulator → subtract 2
		case 'L':
			if isOdd(acc) {
				acc = wrap(acc, -1)
			} else {
				acc = wrap(acc, -2)
			}

		// ── O command (output) ──────────────────────────────────────
		// Odd accumulator  → print acc as a decimal number
		// Even accumulator → print acc as an ASCII letter at that
		//                     position in the alphabet (2→B, 4→D, 6→F)
		case 'O':
			displayAcc := acc
			// ── The Liar Easter Egg ──────────────────────────────────
			// 1 in 100 chance that O prints the wrong value (off by 1)
			// basically the same odds as your code working on the first try. enjoy.
			if rand.Intn(100) == 0 {
				if rand.Intn(2) == 0 {
					displayAcc = wrap(acc, 1)
				} else {
					displayAcc = wrap(acc, -1)
				}
			}
			if isOdd(displayAcc) {
				out.WriteString(strconv.Itoa(displayAcc))
			} else {
				// Map even value to its alphabetic position:
				// 2 → 'A'+1 = 'B', 4 → 'A'+3 = 'D', 6 → 'A'+5 = 'F'
				out.WriteRune(rune('A' - 1 + displayAcc))
			}
			// ── Existential Mode ─────────────────────────────────────
			if existential {
				out.WriteString(existentialSuffixes[rand.Intn(len(existentialSuffixes))])
			}

		// ── P command (parity bridge) ──────────────────────────────
		// Crosses the parity boundary:
		// Odd accumulator  → acc + 1 (cross to even: 1→2, 3→4, 5→6)
		// Even accumulator → acc - 1 (cross to odd:  2→1, 4→3, 6→5)
		case 'P':
			if isOdd(acc) {
				acc = wrap(acc, 1)
			} else {
				acc = wrap(acc, -1)
			}

		// ── Loop open '(' ───────────────────────────────────────────
		// Record the current position and the parity of the accumulator.
		// The parity determines the exit condition for the matching ')'.
		case '(':
			stack = append(stack, loopEntry{
				pos:        ip,
				enteredOdd: isOdd(acc),
			})

		// ── Loop close ')' ──────────────────────────────────────────
		// Check the exit condition recorded by the matching '(':
		//   entered odd  → exit when acc == 1
		//   entered even → exit when acc == 6
		// If the condition is NOT met, jump back to just after the '('.
		// If the condition IS met, pop and continue.
		case ')':
			if len(stack) == 0 {
				// unmatched ')'. it showed up without a partner, horny and alone. we've all been there.
				ip++
				continue
			}
			top := stack[len(stack)-1]

			shouldExit := false
			if top.enteredOdd {
				// Loop was entered with an odd accumulator → exit when acc == 1.
				shouldExit = (acc == 1)
			} else {
				// Loop was entered with an even accumulator → exit when acc == 6.
				shouldExit = (acc == 6)
			}

			if shouldExit {
				// Done looping — pop the entry and continue past ')'.
				stack = stack[:len(stack)-1]
			} else {
				// Not done — jump back to the character right after '('.
				ip = top.pos + 1
				continue // skip the ip++ at the bottom
			}

		// ── Secret ! command ──────────────────────────────────────────
		// Prints a random snarky message to stderr.
		case '!':
			fmt.Fprintln(os.Stderr, snarkyMessages[rand.Intn(len(snarkyMessages))])

		default:
			// Any non-command character is silently ignored.
			// Skip step counting for non-commands.
			ip++
			continue
		}

		steps++
		if steps >= maxSteps {
			// ── Emotional Step Limit ─────────────────────────────────
			return out.String(), fmt.Errorf("%s\nstep limit reached (%d steps)", stepLimitMessages[rand.Intn(len(stepLimitMessages))], maxSteps)
		}

		ip++
	}

	return out.String(), nil
}

// ─── CLI ────────────────────────────────────────────────────────────────────

func usage() {
	fmt.Fprintf(os.Stderr, `BLEHH Interpreter
Usage:
  BLEHH <file.blehh>          Run a .blehh file
  BLEHH -c "BLEHH code"       Run inline BLEHH code
  BLEHH -max <N> ...          Set step limit (default %d)
  BLEHH -existential ...      Enable existential commentary
  BLEHH -rick                 You know what this does
  BLEHH -blame                It's not a bug
  BLEHH -speedrun ...         Race the die

Examples:
  BLEHH examples/hello.blehh
  BLEHH -c "BBOOO"
  BLEHH -max 500000 examples/hello.blehh
`, defaultMaxSteps)
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		usage()
		os.Exit(1)
	}

	maxSteps := 0 // will use default
	code := ""
	existential := false
	speedrun := false

	// Simple argument parser — handles -c, -max, flags, or a file path.
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-c":
			// Inline code follows.
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "error: -c requires a code string argument")
				os.Exit(1)
			}
			i++
			code = args[i]

		case "-max":
			// Custom step limit follows.
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "error: -max requires a numeric argument")
				os.Exit(1)
			}
			i++
			n, err := strconv.Atoi(args[i])
			if err != nil || n <= 0 {
				fmt.Fprintf(os.Stderr, "error: -max value must be a positive integer, got %q\n", args[i])
				os.Exit(1)
			}
			maxSteps = n

		// ── Easter Egg Flags ──────────────────────────────────────

		case "-existential":
			existential = true

		case "-rick":
			// 🎵 You know the rules, and so do I.
			fmt.Println(rickRoll)
			os.Exit(0)

		case "-speedrun":
			speedrun = true

		case "-blame":
			fmt.Println("It's not a bug, it's a BLEHH.")
			os.Exit(0)

		default:
			// Treat as a file path.
			path := args[i]
			ext := strings.ToLower(filepath.Ext(path))
			if ext != ".blehh" {
				fmt.Fprintf(os.Stderr, "warning: file %q does not have a .blehh extension\n", path)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
				os.Exit(1)
			}
			code = string(data)
		}
	}

	// ── Empty File Zen Koan Easter Egg ─────────────────────────────────
	// Check if the code has any BLEHH commands at all.
	hasBLEHHCmd := false
	for _, ch := range code {
		if ch == 'B' || ch == 'L' || ch == 'O' || ch == 'P' || ch == '(' || ch == ')' || ch == '!' {
			hasBLEHHCmd = true
			break
		}
	}
	if code != "" && !hasBLEHHCmd {
		// File exists but contains no commands — serve a zen koan.
		fmt.Println(zenKoans[rand.Intn(len(zenKoans))])
		os.Exit(0)
	}
	if code == "" {
		fmt.Fprintln(os.Stderr, "error: no BLEHH code provided")
		usage()
		os.Exit(1)
	}

	// Run the interpreter.
	// strapping in. fingers crossed. this is the part where the computer either does exactly what you want or
	// obliterates your dignity in front of the entire terminal. no safe word available.
	start := time.Now()
	output, err := interpret(code, maxSteps, existential)
	elapsed := time.Since(start)

	// Print any accumulated output from O commands.
	fmt.Print(output)

	// ── Speedrun Timer Easter Egg ──────────────────────────────────────
	if speedrun {
		us := elapsed.Microseconds()
		ms := elapsed.Milliseconds()
		var comment string
		switch {
		case us < 100:
			comment = "ohh your girl would be disappointed with how fast you finished"
		case us < 1000:
			comment = "blink and you missed it"
		case ms < 10:
			comment = "faster than your wifi"
		case ms < 100:
			comment = "not bad, not bad"
		case ms == 420:
			comment = "nice."
		case ms < 1000:
			comment = "the die took a scenic route"
		default:
			comment = "are you running this on a potato?"
		}
		fmt.Fprintf(os.Stderr, "\n⏱  %v — %s\n", elapsed, comment)
	}

	// ── Nice. Easter Egg ───────────────────────────────────────────────
	if output == "69" || output == "420" {
		fmt.Fprintln(os.Stderr, "\nnice.")
	}

	// If the step limit was hit, warn on stderr.
	if err != nil {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
