# guernica

---

Guernica is a password generator that provides ways of setting the alphabet and the password length, generating batches of passwords, passphrases from the EFF wordlist, and more.

## Dependencies

- go (>= v.1.23.0)

## Build

Simply use go command to build guernica:

```sh
go build
```

## Usage

```
guernica [flags]

Flags:
-a, --alphabet string            password alphabet <digits,lower,upper,alnum,full,complete> (default "full")
    --ambiguous-exclude string   custom character set to strip when excluding ambiguous characters
-c, --clipboard                  copy the generated password to the clipboard
-n, --count int                  number of passwords to generate (default 1)
    --custom-alphabet string     use an arbitrary character set as alphabet
-h, --help                       help for guernica
-l, --length int                 password length (default 20)
    --loop                       generate passwords forever until interrupted
-x, --no-ambiguous               exclude ambiguous characters (l, 1, I, O, 0, ...)
-e, --no-entropy                 do not print the entropy of generated passwords
-r, --require-all                ensure at least one character from every alphabet class
-s, --separator string           word separator for passphrases (default " ")
    --wordlist-file string       path to a custom wordlist for passphrases
-w, --words int                  generate a passphrase of N words instead of a password
```

## Examples

Generate a 20-character password with the default alphabet:

```sh
guernica
```

Generate a 32-character password using only lowercase letters:

```sh
guernica -l 32 -a lower
```

Use an arbitrary character set as the alphabet:

```sh
guernica --custom-alphabet "0123456789abcdefghijklmnopqrstuvwxyz"
```

Copy a password to the clipboard:

```sh
guernica --clipboard
```

Generate 10 passwords at once:

```sh
guernica --count 10
```

Generate a six-word passphrase (EFF large wordlist, space-separated):

```sh
guernica --words 6
```

Generate a passphrase with a custom separator and word count:

```sh
guernica --words 4 --separator -
```

Use your own wordlist file:

```sh
guernica --words 5 --wordlist-file /path/to/words.txt
```

Exclude ambiguous characters and require at least one character from each alphabet class:

```sh
guernica --no-ambiguous --require-all -l 16
```

## Notes

- Entropy is shown as the theoretical maximum (`log2(alphabet_size^length)`, or `log2(wordlist_size^words)` for passphrases) and is suppressed with `--no-entropy`.
- `--clipboard` requires a clipboard tool: `wl-copy`, `xclip`, or `xsel` on Linux, `pbcopy` on macOS, and `clip` on Windows. It errors out when none is available.
- The embedded passphrase wordlist is the EFF large wordlist (7,776 words), from https://www.eff.org/dice.
