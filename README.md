# go_reloaded - Text Processor in Go

This project is a **text processing tool** written in Go.  
It reads a text file, applies transformations based on special markers, and outputs the processed text.

---

## Features

- Convert numbers from **hexadecimal `(hex)`** and **binary `(bin)`** to decimal.
- Apply text transformations:
  - Uppercase `(up)` or `(up, n)`
  - Lowercase `(low)` or `(low, n)`
  - Capitalize `(cap)` or `(cap, n)`
- Automatically change `a` to `an` when followed by a vowel.
- Handle punctuation spacing properly, including quotes `'`.
- Flexible and extensible for additional markers.

---
