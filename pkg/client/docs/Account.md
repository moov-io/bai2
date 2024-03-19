# Account

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccountNumber** | Pointer to **string** |  | [optional] 
**CurrencyCode** | Pointer to **string** |  | [optional] 
**Summaries** | Pointer to [**[]AccountSummary**](AccountSummary.md) |  | [optional] 
**AccountControlTotal** | Pointer to **string** |  | [optional] 
**NumberOfRecords** | Pointer to **int32** |  | [optional] 
**Details** | Pointer to [**[]Detail**](Detail.md) |  | [optional] 

## Methods

### NewAccount

`func NewAccount() *Account`

NewAccount instantiates a new Account object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAccountWithDefaults

`func NewAccountWithDefaults() *Account`

NewAccountWithDefaults instantiates a new Account object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccountNumber

`func (o *Account) GetAccountNumber() string`

GetAccountNumber returns the AccountNumber field if non-nil, zero value otherwise.

### GetAccountNumberOk

`func (o *Account) GetAccountNumberOk() (*string, bool)`

GetAccountNumberOk returns a tuple with the AccountNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountNumber

`func (o *Account) SetAccountNumber(v string)`

SetAccountNumber sets AccountNumber field to given value.

### HasAccountNumber

`func (o *Account) HasAccountNumber() bool`

HasAccountNumber returns a boolean if a field has been set.

### GetCurrencyCode

`func (o *Account) GetCurrencyCode() string`

GetCurrencyCode returns the CurrencyCode field if non-nil, zero value otherwise.

### GetCurrencyCodeOk

`func (o *Account) GetCurrencyCodeOk() (*string, bool)`

GetCurrencyCodeOk returns a tuple with the CurrencyCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrencyCode

`func (o *Account) SetCurrencyCode(v string)`

SetCurrencyCode sets CurrencyCode field to given value.

### HasCurrencyCode

`func (o *Account) HasCurrencyCode() bool`

HasCurrencyCode returns a boolean if a field has been set.

### GetSummaries

`func (o *Account) GetSummaries() []AccountSummary`

GetSummaries returns the Summaries field if non-nil, zero value otherwise.

### GetSummariesOk

`func (o *Account) GetSummariesOk() (*[]AccountSummary, bool)`

GetSummariesOk returns a tuple with the Summaries field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSummaries

`func (o *Account) SetSummaries(v []AccountSummary)`

SetSummaries sets Summaries field to given value.

### HasSummaries

`func (o *Account) HasSummaries() bool`

HasSummaries returns a boolean if a field has been set.

### GetAccountControlTotal

`func (o *Account) GetAccountControlTotal() string`

GetAccountControlTotal returns the AccountControlTotal field if non-nil, zero value otherwise.

### GetAccountControlTotalOk

`func (o *Account) GetAccountControlTotalOk() (*string, bool)`

GetAccountControlTotalOk returns a tuple with the AccountControlTotal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountControlTotal

`func (o *Account) SetAccountControlTotal(v string)`

SetAccountControlTotal sets AccountControlTotal field to given value.

### HasAccountControlTotal

`func (o *Account) HasAccountControlTotal() bool`

HasAccountControlTotal returns a boolean if a field has been set.

### GetNumberOfRecords

`func (o *Account) GetNumberOfRecords() int32`

GetNumberOfRecords returns the NumberOfRecords field if non-nil, zero value otherwise.

### GetNumberOfRecordsOk

`func (o *Account) GetNumberOfRecordsOk() (*int32, bool)`

GetNumberOfRecordsOk returns a tuple with the NumberOfRecords field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumberOfRecords

`func (o *Account) SetNumberOfRecords(v int32)`

SetNumberOfRecords sets NumberOfRecords field to given value.

### HasNumberOfRecords

`func (o *Account) HasNumberOfRecords() bool`

HasNumberOfRecords returns a boolean if a field has been set.

### GetDetails

`func (o *Account) GetDetails() []Detail`

GetDetails returns the Details field if non-nil, zero value otherwise.

### GetDetailsOk

`func (o *Account) GetDetailsOk() (*[]Detail, bool)`

GetDetailsOk returns a tuple with the Details field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDetails

`func (o *Account) SetDetails(v []Detail)`

SetDetails sets Details field to given value.

### HasDetails

`func (o *Account) HasDetails() bool`

HasDetails returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


