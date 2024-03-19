# AccountSummary

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TypeCode** | Pointer to **string** |  | [optional] 
**Amount** | Pointer to **string** |  | [optional] 
**ItemCount** | Pointer to **int32** |  | [optional] 
**FundsType** | Pointer to [**FundsType**](FundsType.md) |  | [optional] 

## Methods

### NewAccountSummary

`func NewAccountSummary() *AccountSummary`

NewAccountSummary instantiates a new AccountSummary object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAccountSummaryWithDefaults

`func NewAccountSummaryWithDefaults() *AccountSummary`

NewAccountSummaryWithDefaults instantiates a new AccountSummary object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTypeCode

`func (o *AccountSummary) GetTypeCode() string`

GetTypeCode returns the TypeCode field if non-nil, zero value otherwise.

### GetTypeCodeOk

`func (o *AccountSummary) GetTypeCodeOk() (*string, bool)`

GetTypeCodeOk returns a tuple with the TypeCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTypeCode

`func (o *AccountSummary) SetTypeCode(v string)`

SetTypeCode sets TypeCode field to given value.

### HasTypeCode

`func (o *AccountSummary) HasTypeCode() bool`

HasTypeCode returns a boolean if a field has been set.

### GetAmount

`func (o *AccountSummary) GetAmount() string`

GetAmount returns the Amount field if non-nil, zero value otherwise.

### GetAmountOk

`func (o *AccountSummary) GetAmountOk() (*string, bool)`

GetAmountOk returns a tuple with the Amount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmount

`func (o *AccountSummary) SetAmount(v string)`

SetAmount sets Amount field to given value.

### HasAmount

`func (o *AccountSummary) HasAmount() bool`

HasAmount returns a boolean if a field has been set.

### GetItemCount

`func (o *AccountSummary) GetItemCount() int32`

GetItemCount returns the ItemCount field if non-nil, zero value otherwise.

### GetItemCountOk

`func (o *AccountSummary) GetItemCountOk() (*int32, bool)`

GetItemCountOk returns a tuple with the ItemCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItemCount

`func (o *AccountSummary) SetItemCount(v int32)`

SetItemCount sets ItemCount field to given value.

### HasItemCount

`func (o *AccountSummary) HasItemCount() bool`

HasItemCount returns a boolean if a field has been set.

### GetFundsType

`func (o *AccountSummary) GetFundsType() FundsType`

GetFundsType returns the FundsType field if non-nil, zero value otherwise.

### GetFundsTypeOk

`func (o *AccountSummary) GetFundsTypeOk() (*FundsType, bool)`

GetFundsTypeOk returns a tuple with the FundsType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFundsType

`func (o *AccountSummary) SetFundsType(v FundsType)`

SetFundsType sets FundsType field to given value.

### HasFundsType

`func (o *AccountSummary) HasFundsType() bool`

HasFundsType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


