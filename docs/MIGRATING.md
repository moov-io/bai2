# Migrating from `pkg/lib`

`pkg/lib` was renamed to `pkg/bai2`. BAI3 / BTR3 (X9.121 version 3) lives in `pkg/bai3`.
A version-detecting reader sits on the module root.

`pkg/lib` remains as a compatibility shim. Existing `import "github.com/moov-io/bai2/pkg/lib"`
code continues to compile. New code should import `pkg/bai2` or the root `Read` helper.

## Import paths

```
// before
import "github.com/moov-io/bai2/pkg/lib"
f := lib.NewBai2()

// after (BAI2 only)
import "github.com/moov-io/bai2/pkg/bai2"
f := bai2.NewBai2()

// after (auto-detect BAI2 vs BAI3)
import moovbai "github.com/moov-io/bai2"
f, err := moovbai.Read(r)
```

## Auto-detect

`moovbai.Read` inspects the 01 Version Number:

| Version | Parser |
|---|---|
| 2 | `pkg/bai2` |
| 3 | `pkg/bai3` |

`ReadOptions.IgnoreVersion` keeps the historical behavior: parse as BAI2 even when
the header is not `2` (banks that stamp BAI2 files with `3`). To force a parser:

```
moovbai.ReadWithOptions(r, moovbai.ReadOptions{ForceVersion: 3})
```

## BAI3 differences that matter

- Record 02 is a **bank header**. Currency is positional-null (not defaulted to USD).
- Record 03 **requires** a currency code.
- Group status on 02 is always Update (`1`); other values are retired.
- Funds type `D` is retired in BTR3 (still parsed).
- Type 89/90 batch/invoice records from BTRS v1 are not implemented.

## CLI / HTTP

`bai2 parse|print|format` auto-detects version. `--ignoreVersion` still means
“parse as BAI2”. HTTP `/parse`, `/print`, `/format` do the same.

`/format` JSON for BAI3 uses `Banks` (record 02 bank headers) instead of
`Groups`. The generated `pkg/client` models remain BAI2-shaped.
