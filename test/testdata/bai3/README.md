# BAI3 / BTR3 golden files

These files are assembled from the **public** record samples in the
[X9.121 BTRS Version 3 Format Guide](https://x9.org/wp-content/uploads/2018/07/X9.121-2016-BTRS-Version-3.0.pdf)
(also linked from [x9.org/standards/btrs](https://x9.org/standards/btrs/download-btrs/)).

They are not a copy of the copyrighted standard. Each file stitches published
01/02/03/16/88/49/98/99 examples into a parseable envelope so the parser can be
regression-tested against those field layouts.

| File | Source samples |
|---|---|
| `x9-mandatory.txt` | 01 mandatory fields; 02 mandatory fields; 03 account+currency only |
| `x9-status-summary-detail.txt` | 03 with status 010 + summary 190; 16 funds type 0 |
| `x9-bank-header-all-fields.txt` | 02 with ultimate receiver and as-of-time; 16 funds type S |
| `x9-continuation-16.txt` | 16 funds type S continued by 88 bank/customer ref and text |
| `x9-03-split-summary.txt` | 03 status/summary codes split across 88 records |
