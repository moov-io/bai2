# bai2
BAI2 file format is a standardized set of codes which comes in text format

BAI2- a widely accepted and used Bank Statement Format for Bank Reconciliation. To have uniformity in the Bank formats, BAI (Bank Administration Institute) developed a generic format and it is widely accepted by most of the Banks in USA. That format is called BAI2 format. BAI2 format is getting used by Corporates to import the Bank statements in the system and will be used for reconciliation. BAI2 file format is a standardized set of codes which comes in text format. It contains the data related to AP Payments (negotiable and void status) and Receipts (remittance status).

Reference 

https://www.bai.org/docs/default-source/libraries/site-general-downloads/cash_management_2005.pdf

https://www.tdcommercialbanking.com/document/PDF/bai.pdf


## BAI Format - Overview
Some customers may download their files in Bank Administration Institute (BAI) format which was developed as the basis for uniform formats and terminology by the Institute. Customers can use their own software application to generate spreadsheets or reports from the data in the BAI file. In BAI format, each entry is separated by a comma and the last character is followed by a slash (/). The comma is not counted as part of the entry length

### Record Types
BAI files via Web Business Banking Balance Reporting include the following types of records: 

| Record Code | Record Name |
| ----------- | ----------- | 
| 01 | File Header |
| 02 | Group header | 
| 03 | Account Identifier and Summary/Status |
| 16 | Account Transaction Detail |
| 88 | Continuation of Account Summary Record |
| 49 | Account Trailer |
| 98 | Group Trailer |
| 99 | File Trailer |


### Table 1: File Header
The File Header is the first record in a BAI format file. It always has a record code of 01 

example:```01,0004,12345,060321,0829,001,80,1,2/```

| ENTRY NO. | CHARACTER POSITION | ENTRY LENGTH | DESCRIPTION | EXAMPLE |
| ---- | ---- | ---- | ------------- | ---- |
| 1 | 1 | 2 | Record Code | 01,* | 
| 2 | 4 | 4 | Sender | 0004,* (TD Bank)| 
| 3 | 9 | 5 | Receiver (customer number) | 12345, |
| 4 | 15 | 6 | Date File Created <br>(YYMMDD) | 060321, |
| 5 | 22 | 4 | Time File Created (HHMM) | 2400, |
| 6 | 27 | 3 | File ID Number | 001, (daily) |
| 7 | 31 | 2 | Physical Record Length | 80,* |
| 8 | 34 | 1 | Block Size | 1,* |
| 9 | 36 | 1 | Version Number | 2/* ** | 
<div style="text-align: right">
*Value never changes <br>
** / indicates end of logical record
</div>

### Table 2: Group Header
The Group Header is the second record in a BAI format file. It always has a record code of 02.

example:```02,12345,0004,1,060317,,CAD,/``` 

| ENTRY NO. | CHARACTER POSITION | ENTRY LENGTH | DESCRIPTION | EXAMPLE |
| --- | --- | --- | -----| ----- |
| 1 | 1 | 2 | Record Code | 02,* |
| 2 | 4 | 5 | Receiver (customer number) | 12345, |
| 3 | 10 | 4 | Sender | 0004,* (TD Bank) |
| 4 | 15 | 1 | Group Status | 1,* |
| 5 | 17 | 6 | As-of-Date (YYMMDD) | 060317, |
| 6 | 24 | 0 | As-of-Time (n/a) | , |
| 7 | 25 | 3 | Currency Code | CAD, or USD, |
| 8 | 29 | 0 | As-of-Date Modifier (n/a) | /** | 
<div style="text-align: right">
*Value never changes <br>
** / indicates end of logical record
</div>


### Table 3: Account Identifier and Summary/Status for Current (CDA), Personal (PDA), and Loan accounts 
This record contains information on opening and closing balances for CDA and PDA accounts. It always
has a record code of 03. 

example:```03,10200123456,CAD,040,+000000000000,,,045,+000000000000,,/```

| ENTRY NO. | CHARACTER POSITION | ENTRY LENGTH | DESCRIPTION | EXAMPLE |
| --- | --- | --- | ---------- | ----------| 
| 1 | 1 | 2 | Record Code | 03,* |
| 2 | 4 | 11 | Account Number | 10200123456, | 
| 3 | 16 | 3 | Currency Code | CAD, or USD, |
| 4 | 20 | 3 | Type Code | 040,* |
| 5 | 24 | 13 | Opening Balance | +123456789012, |
| 6 | 38 | 0 | Item Count (n/a) | , |
| 7 | 39 | 0 | Funds Type (n/a) | , |
| 8 | 40 | 3 | Type Code | 045,* |
| 9 | 44 | 13 | Closing Balance | +123456789012, |
| 10 | 58 | 0 | Item Count (n/a) | , |
| 11 | 59 | 0 | Funds Type (n/a) | /** | 
<div style="text-align: right">
*Value never changes <br>
** / indicates end of logical record
</div>

### Table 4: Continuation of Account Summary Record - Current (CDA) and Personal (PDA) Accounts
his record is a continuation of the account identifier and summary record (see Table 3) and contains value date and summary information (i.e. total credits and debits as well as total credit and debit dollar amounts). It always has a record code of 88. 

example:```88,100,000000000208500,00003,V,060316,,400,000000000208500,00008,
V,060316,/```

| ENTRY NO. | CHARACTER POSITION | ENTRY LENGTH | DESCRIPTION | EXAMPLE |
| --- | --- | --- | ---------- | ----------| 
| 1 | 1 | 2 | Record Code | 88,* |
| 2 | 4 | 3 | Type Code (total credits) | 100,* |
| 3 | 8 | 15 | Total Credit $ Amount | 123456789012345, |
| 4 | 24 | 5 | Total # of Credits (item count) | 12345, |
| 5 | 30 | 1 | Funds Type | V,= value dated |
| 6 | 32 | 6 | Value Date (YYMMDD) | 060316, |
| 7 | 39 | 0 | Value Time (HHMM) | , |
| 8 | 40 | 3 | Type Code | 400,* |
| 9 | 44 | 15 | Total Debit $ Amount | 123456789012345, | 
| 10 | 60 | 5 | Total # of Debits (item count) | 12345, | 
| 11 | 66 | 1 | Funds Type | V,= value dated |
| 12 | 68 | 6 | Value Date (YYMMDD) | 060316, |
| 13 | 75 | 0 | Value Time (HHMM) | /** |
<div style="text-align: right">
*Value never changes <br>
** / indicates end of logical record
</div>

_The 03 record for loan accounts contains value date information but NO summary information. There is no
88 record for loan accounts._

