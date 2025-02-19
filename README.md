# guernica

---

Guernica is a password generator that provides ways of setting the alphabet and the password length.

## Dependencies

- go (>= v.1.23.0)

## Build

Simply use go command to build guernica:

```sh
go build
```

## Usage

```
guernica <flags>

Flags:
-a, --alphabet string password alphabet <digits,lower,upper,alnum,full,complete> (default "full")
-h, --help help for guernica
-l, --length int password length (default 20)
```
