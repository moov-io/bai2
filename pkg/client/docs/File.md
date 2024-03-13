# File

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Sender** | Pointer to **string** |  | [optional] 
**Receiver** | Pointer to **string** |  | [optional] 
**FileCreatedDate** | Pointer to **string** |  | [optional] 
**FileCreatedTime** | Pointer to **string** |  | [optional] 
**FileIdNumber** | Pointer to **string** |  | [optional] 
**PhysicalRecordLength** | Pointer to **int32** |  | [optional] 
**BlockSize** | Pointer to **int32** |  | [optional] 
**VersionNumber** | Pointer to **int32** |  | [optional] 
**FileControlTotal** | Pointer to **string** |  | [optional] 
**NumberOfGroups** | Pointer to **int32** |  | [optional] 
**NumberOfRecords** | Pointer to **int32** |  | [optional] 
**Groups** | Pointer to [**[]Group**](Group.md) |  | [optional] 

## Methods

### NewFile

`func NewFile() *File`

NewFile instantiates a new File object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewFileWithDefaults

`func NewFileWithDefaults() *File`

NewFileWithDefaults instantiates a new File object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSender

`func (o *File) GetSender() string`

GetSender returns the Sender field if non-nil, zero value otherwise.

### GetSenderOk

`func (o *File) GetSenderOk() (*string, bool)`

GetSenderOk returns a tuple with the Sender field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSender

`func (o *File) SetSender(v string)`

SetSender sets Sender field to given value.

### HasSender

`func (o *File) HasSender() bool`

HasSender returns a boolean if a field has been set.

### GetReceiver

`func (o *File) GetReceiver() string`

GetReceiver returns the Receiver field if non-nil, zero value otherwise.

### GetReceiverOk

`func (o *File) GetReceiverOk() (*string, bool)`

GetReceiverOk returns a tuple with the Receiver field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReceiver

`func (o *File) SetReceiver(v string)`

SetReceiver sets Receiver field to given value.

### HasReceiver

`func (o *File) HasReceiver() bool`

HasReceiver returns a boolean if a field has been set.

### GetFileCreatedDate

`func (o *File) GetFileCreatedDate() string`

GetFileCreatedDate returns the FileCreatedDate field if non-nil, zero value otherwise.

### GetFileCreatedDateOk

`func (o *File) GetFileCreatedDateOk() (*string, bool)`

GetFileCreatedDateOk returns a tuple with the FileCreatedDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFileCreatedDate

`func (o *File) SetFileCreatedDate(v string)`

SetFileCreatedDate sets FileCreatedDate field to given value.

### HasFileCreatedDate

`func (o *File) HasFileCreatedDate() bool`

HasFileCreatedDate returns a boolean if a field has been set.

### GetFileCreatedTime

`func (o *File) GetFileCreatedTime() string`

GetFileCreatedTime returns the FileCreatedTime field if non-nil, zero value otherwise.

### GetFileCreatedTimeOk

`func (o *File) GetFileCreatedTimeOk() (*string, bool)`

GetFileCreatedTimeOk returns a tuple with the FileCreatedTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFileCreatedTime

`func (o *File) SetFileCreatedTime(v string)`

SetFileCreatedTime sets FileCreatedTime field to given value.

### HasFileCreatedTime

`func (o *File) HasFileCreatedTime() bool`

HasFileCreatedTime returns a boolean if a field has been set.

### GetFileIdNumber

`func (o *File) GetFileIdNumber() string`

GetFileIdNumber returns the FileIdNumber field if non-nil, zero value otherwise.

### GetFileIdNumberOk

`func (o *File) GetFileIdNumberOk() (*string, bool)`

GetFileIdNumberOk returns a tuple with the FileIdNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFileIdNumber

`func (o *File) SetFileIdNumber(v string)`

SetFileIdNumber sets FileIdNumber field to given value.

### HasFileIdNumber

`func (o *File) HasFileIdNumber() bool`

HasFileIdNumber returns a boolean if a field has been set.

### GetPhysicalRecordLength

`func (o *File) GetPhysicalRecordLength() int32`

GetPhysicalRecordLength returns the PhysicalRecordLength field if non-nil, zero value otherwise.

### GetPhysicalRecordLengthOk

`func (o *File) GetPhysicalRecordLengthOk() (*int32, bool)`

GetPhysicalRecordLengthOk returns a tuple with the PhysicalRecordLength field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhysicalRecordLength

`func (o *File) SetPhysicalRecordLength(v int32)`

SetPhysicalRecordLength sets PhysicalRecordLength field to given value.

### HasPhysicalRecordLength

`func (o *File) HasPhysicalRecordLength() bool`

HasPhysicalRecordLength returns a boolean if a field has been set.

### GetBlockSize

`func (o *File) GetBlockSize() int32`

GetBlockSize returns the BlockSize field if non-nil, zero value otherwise.

### GetBlockSizeOk

`func (o *File) GetBlockSizeOk() (*int32, bool)`

GetBlockSizeOk returns a tuple with the BlockSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlockSize

`func (o *File) SetBlockSize(v int32)`

SetBlockSize sets BlockSize field to given value.

### HasBlockSize

`func (o *File) HasBlockSize() bool`

HasBlockSize returns a boolean if a field has been set.

### GetVersionNumber

`func (o *File) GetVersionNumber() int32`

GetVersionNumber returns the VersionNumber field if non-nil, zero value otherwise.

### GetVersionNumberOk

`func (o *File) GetVersionNumberOk() (*int32, bool)`

GetVersionNumberOk returns a tuple with the VersionNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersionNumber

`func (o *File) SetVersionNumber(v int32)`

SetVersionNumber sets VersionNumber field to given value.

### HasVersionNumber

`func (o *File) HasVersionNumber() bool`

HasVersionNumber returns a boolean if a field has been set.

### GetFileControlTotal

`func (o *File) GetFileControlTotal() string`

GetFileControlTotal returns the FileControlTotal field if non-nil, zero value otherwise.

### GetFileControlTotalOk

`func (o *File) GetFileControlTotalOk() (*string, bool)`

GetFileControlTotalOk returns a tuple with the FileControlTotal field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFileControlTotal

`func (o *File) SetFileControlTotal(v string)`

SetFileControlTotal sets FileControlTotal field to given value.

### HasFileControlTotal

`func (o *File) HasFileControlTotal() bool`

HasFileControlTotal returns a boolean if a field has been set.

### GetNumberOfGroups

`func (o *File) GetNumberOfGroups() int32`

GetNumberOfGroups returns the NumberOfGroups field if non-nil, zero value otherwise.

### GetNumberOfGroupsOk

`func (o *File) GetNumberOfGroupsOk() (*int32, bool)`

GetNumberOfGroupsOk returns a tuple with the NumberOfGroups field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumberOfGroups

`func (o *File) SetNumberOfGroups(v int32)`

SetNumberOfGroups sets NumberOfGroups field to given value.

### HasNumberOfGroups

`func (o *File) HasNumberOfGroups() bool`

HasNumberOfGroups returns a boolean if a field has been set.

### GetNumberOfRecords

`func (o *File) GetNumberOfRecords() int32`

GetNumberOfRecords returns the NumberOfRecords field if non-nil, zero value otherwise.

### GetNumberOfRecordsOk

`func (o *File) GetNumberOfRecordsOk() (*int32, bool)`

GetNumberOfRecordsOk returns a tuple with the NumberOfRecords field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumberOfRecords

`func (o *File) SetNumberOfRecords(v int32)`

SetNumberOfRecords sets NumberOfRecords field to given value.

### HasNumberOfRecords

`func (o *File) HasNumberOfRecords() bool`

HasNumberOfRecords returns a boolean if a field has been set.

### GetGroups

`func (o *File) GetGroups() []Group`

GetGroups returns the Groups field if non-nil, zero value otherwise.

### GetGroupsOk

`func (o *File) GetGroupsOk() (*[]Group, bool)`

GetGroupsOk returns a tuple with the Groups field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGroups

`func (o *File) SetGroups(v []Group)`

SetGroups sets Groups field to given value.

### HasGroups

`func (o *File) HasGroups() bool`

HasGroups returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


