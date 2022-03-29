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

ex: ```01,0004,12345,060321,0829,001,80,1,2/```

| ENTRY NO. | CHARACTER POSITION | ENTRY LENGTH | DESCRIPTION | EXAMPLE |
| ---- | ---- | ---- | ------------- | ---- |
| 1 | 1 | 2 | Record Code | 01,* | 

2 4 4 Sender 0004,* (TD Bank)
3 9 5 Receiver (customer number) 12345,
4 15 6 Date File Created
(YYMMDD)
060321,
5 22 4 Time File Created (HHMM) 2400,
6 27 3 File ID Number 001, (daily)
7 31 2 Physical Record Length 80,*
8 34 1 Block Size 1,*
9 36 1 Version Number 2/* ** 



