# Group

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Receiver** | Pointer to **string** |  | [optional] 
**Originator** | Pointer to **string** |  | [optional] 
**GroupStatus** | Pointer to **int32** |  | [optional] 
**AsOfDate** | Pointer to **string** |  | [optional] 
**CurrencyCode** | Pointer to **string** |  | [optional] 
**GroupControlTotal** | Pointer to **string** |  | [optional] 
**NumberOfAccounts** | Pointer to **int32** |  | [optional] 
**NumberOfRecords** | Pointer to **int32** |  | [optional] 
**Accounts** | Pointer to [**[]Account**](Account.md) |  | [optional] 

## Methods

### NewGroup

`func NewGroup() *Group`

NewGroup instantiates a new Group object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGroupWithDefaults

`func NewGroupWithDefaults() *Group`

NewGroupWithDefaults instantiates a new Group object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetReceiver

`func (o *Group) GetReceiver() string`

GetReceiver returns the Receiver field if non-nil, zero value otherwise.

### GetReceiverOk

`func (o *Group) GetReceiverOk() (*string, bool)`

GetReceiverOk returns a tuple with the Receiver field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReceiver

`func (o *Group) SetReceiver(v string)`

SetReceiver sets Receiver field to given value.

### HasReceiver

`func (o *Group) HasReceiver() bool`

HasReceiver returns a boolean if a field has been set.

### GetOriginator

`func (o *Group) GetOriginator() string`

GetOriginator returns the Originator field if non-nil, zero value otherwise.

### GetOriginatorOk

`func (o *Group) GetOriginatorOk() (*string, bool)`

GetOriginatorOk returns a tuple with the Originator field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOriginator

`func (o *Group) SetOriginator(v string)`

SetOriginator sets Originator field to given value.

### HasOriginator

`func (o *Group) HasOriginator() bool`

HasOriginator returns a boolean if a field has been set.

### GetGroupStatus

`func (o *Group) GetGroupStatus() int32`

GetGroupStatus returns the GroupStatus field if non-nil, zero value otherwise.

### GetGroupStatusOk

`func (o *Group) GetGroupStatusOk() (*int32, bool)`

GetGroupStatusOk returns a tuple with the GroupStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGroupStatus

`func (o *Group) SetGroupStatus(v int32)`

SetGroupStatus sets GroupStatus field to given value.

### HasGroupStatus

`func (o *Group) HasGroupStatus() bool`

HasGroupStatus returns a boolean if a field has been set.

### GetAsOfDate

`func (o *Group) GetAsOfDate() string`

GetAsOfDate returns the AsOfDate field if non-nil, zero value otherwise.

### GetAsOfDateOk

`func (o *Group) GetAsOfDateOk() (*string, bool)`

GetAsOfDateOk returns a tuple with the AsOfDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAsOfDate

`func (o *Group) SetAsOfDate(v string)`

SetAsOfDate sets AsOfDate field to given value.

### HasAsOfDate

`func (o *Group) HasAsOfDate() bool`

HasAsOfDate returns a boolean if a field has been set.

### GetCurrencyCode

`func (o *Group) GetCurrencyCode() string`

GetCurrencyCode returns the CurrencyCode field if non-nil, zero value otherwise.

### GetCurrencyCodeOk

`func (o *Group) GetCurrencyCodeOk() (*string, bool)`

GetCurrencyCodeOk returns a tuple with the CurrencyCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrencyCode

`func (o *Group) SetCurrencyCode(v string)`

SetCurrencyCode sets CurrencyCode field to given value.

### HasCurrencyCode

`func (o *Group) HasCurrencyCode() bool`

HasCurrencyCode returns a boolean if a field has been set.

### GetGroupControlTotal

`func (o *Group) GetGroupControlTotal() string`

GetGroupControlTotal returns the GroupControlTotal field if non-nil, zero value otherwise.

### GetGroupControlTotalOk

`func (o *Group) GetGroupControlTotalOk() (*string, bool)`

GetGroupControlTotalOk returns a tuple with the GroupControlTotal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGroupControlTotal

`func (o *Group) SetGroupControlTotal(v string)`

SetGroupControlTotal sets GroupControlTotal field to given value.

### HasGroupControlTotal

`func (o *Group) HasGroupControlTotal() bool`

HasGroupControlTotal returns a boolean if a field has been set.

### GetNumberOfAccounts

`func (o *Group) GetNumberOfAccounts() int32`

GetNumberOfAccounts returns the NumberOfAccounts field if non-nil, zero value otherwise.

### GetNumberOfAccountsOk

`func (o *Group) GetNumberOfAccountsOk() (*int32, bool)`

GetNumberOfAccountsOk returns a tuple with the NumberOfAccounts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumberOfAccounts

`func (o *Group) SetNumberOfAccounts(v int32)`

SetNumberOfAccounts sets NumberOfAccounts field to given value.

### HasNumberOfAccounts

`func (o *Group) HasNumberOfAccounts() bool`

HasNumberOfAccounts returns a boolean if a field has been set.

### GetNumberOfRecords

`func (o *Group) GetNumberOfRecords() int32`

GetNumberOfRecords returns the NumberOfRecords field if non-nil, zero value otherwise.

### GetNumberOfRecordsOk

`func (o *Group) GetNumberOfRecordsOk() (*int32, bool)`

GetNumberOfRecordsOk returns a tuple with the NumberOfRecords field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumberOfRecords

`func (o *Group) SetNumberOfRecords(v int32)`

SetNumberOfRecords sets NumberOfRecords field to given value.

### HasNumberOfRecords

`func (o *Group) HasNumberOfRecords() bool`

HasNumberOfRecords returns a boolean if a field has been set.

### GetAccounts

`func (o *Group) GetAccounts() []Account`

GetAccounts returns the Accounts field if non-nil, zero value otherwise.

### GetAccountsOk

`func (o *Group) GetAccountsOk() (*[]Account, bool)`

GetAccountsOk returns a tuple with the Accounts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccounts

`func (o *Group) SetAccounts(v []Account)`

SetAccounts sets Accounts field to given value.

### HasAccounts

`func (o *Group) HasAccounts() bool`

HasAccounts returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


