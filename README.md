arithmetic—quiz on simple arithmetic.

A "functional" port of [the classic Unix game](http://man.cat-v.org/unix_8th/6/arithmetic), for no particular reason, to Golang. Largely inspired by Eamon McManus' [version that ships with bsdgames](https://github.com/jsm28/bsd-games/blob/master/arithmetic/arithmetic.c).

**Key differences:**
1. Added support for Unicode symbols (÷ and ×). Use the -u flag to replace the standard ASCII symbols.
2. Both the problem and its answer will be within the provided range.
3. Anywhere between 15 and 30 questions are asked in a go - this helps create some variation.
4. Limited the maximum range such that overflows are avoided. This can be seen if -r is set to 0.  The original program had a limit of 100, while McManus' version had no limit to the maximum range.
5. Penalisation is a little more lax.
