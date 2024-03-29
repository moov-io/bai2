/*
BAI2 API

Moov Bai2 ([Automated Clearing House](https://en.wikipedia.org/wiki/Automated_Clearing_House)) implements an HTTP API for creating, parsing and validating Bais files. BAI2- a widely accepted and used Bank Statement Format for Bank Reconciliation.

API version: v1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
)

// checks if the File type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &File{}

// File struct for File
type File struct {
	Sender               *string `json:"sender,omitempty"`
	Receiver             *string `json:"receiver,omitempty"`
	FileCreatedDate      *string `json:"fileCreatedDate,omitempty"`
	FileCreatedTime      *string `json:"fileCreatedTime,omitempty"`
	FileIdNumber         *string `json:"fileIdNumber,omitempty"`
	PhysicalRecordLength *int32  `json:"physicalRecordLength,omitempty"`
	BlockSize            *int32  `json:"blockSize,omitempty"`
	VersionNumber        *int32  `json:"versionNumber,omitempty"`
	FileControlTotal     *string `json:"fileControlTotal,omitempty"`
	NumberOfGroups       *int32  `json:"numberOfGroups,omitempty"`
	NumberOfRecords      *int32  `json:"numberOfRecords,omitempty"`
	Groups               []Group `json:"Groups,omitempty"`
}

// NewFile instantiates a new File object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewFile() *File {
	this := File{}
	return &this
}

// NewFileWithDefaults instantiates a new File object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewFileWithDefaults() *File {
	this := File{}
	return &this
}

// GetSender returns the Sender field value if set, zero value otherwise.
func (o *File) GetSender() string {
	if o == nil || IsNil(o.Sender) {
		var ret string
		return ret
	}
	return *o.Sender
}

// GetSenderOk returns a tuple with the Sender field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *File) GetSenderOk() (*string, bool) {
	if o == nil || IsNil(o.Sender) {
		return nil, false
	}
	return o.Sender, true
}

// HasSender returns a boolean if a field has been set.
func (o *File) HasSender() bool {
	if o != nil && !IsNil(o.Sender) {
		return true
	}

	return false
}

// SetSender gets a reference to the given string and assigns it to the Sender field.
func (o *File) SetSender(v string) {
	o.Sender = &v
}

// GetReceiver returns the Receiver field value if set, zero value otherwise.
func (o *File) GetReceiver() string {
	if o == nil || IsNil(o.Receiver) {
		var ret string
		return ret
	}
	return *o.Receiver
}

// GetReceiverOk returns a tuple with the Receiver field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *File) GetReceiverOk() (*string, bool) {
	if o == nil || IsNil(o.Receiver) {
		return nil, false
	}
	return o.Receiver, true
}

// HasReceiver returns a boolean if a field has been set.
func (o *File) HasReceiver() bool {
	if o != nil && !IsNil(o.Receiver) {
		return true
	}

	return false
}

// SetReceiver gets a reference to the given string and assigns it to the Receiver field.
func (o *File) SetReceiver(v string) {
	o.Receiver = &v
}

// GetFileCreatedDate returns the FileCreatedDate field value if set, zero value otherwise.
func (o *File) GetFileCreatedDate() string {
	if o == nil || IsNil(o.FileCreatedDate) {
		var ret string
		return ret
	}
	return *o.FileCreatedDate
}

// GetFileCreatedDateOk returns a tuple with the FileCreatedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *File) GetFileCreatedDateOk() (*string, bool) {
	if o == nil || IsNil(o.FileCreatedDate) {
		return nil, false
	}
	return o.FileCreatedDate, true
}

// HasFileCreatedDate returns a boolean if a field has been set.
func (o *File) HasFileCreatedDate() bool {
	if o != nil && !IsNil(o.FileCreatedDate) {
		return true
	}

	return false
}

// SetFileCreatedDate gets a reference to the given string and assigns it to the FileCreatedDate field.
func (o *File) SetFileCreatedDate(v string) {
	o.FileCreatedDate = &v
}

// GetFileCreatedTime returns the FileCreatedTime field value if set, zero value otherwise.
func (o *File) GetFileCreatedTime() string {
	if o == nil || IsNil(o.FileCreatedTime) {
		var ret string
		return ret
	}
	return *o.FileCreatedTime
}

// GetFileCreatedTimeOk returns a tuple with the FileCreatedTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *File) GetFileCreatedTimeOk() (*string, bool) {
	if o == nil || IsNil(o.FileCreatedTime) {
		return nil, false
	}
	return o.FileCreatedTime, true
}

// HasFileCreatedTime returns a boolean if a field has been set.
func (o *File) HasFileCreatedTime() bool {
	if o != nil && !IsNil(o.FileCreatedTime) {
		return true
	}

	return false
}

// SetFileCreatedTime gets a reference to the given string and assigns it to the FileCreatedTime field.
func (o *File) SetFileCreatedTime(v string) {
	o.FileCreatedTime = &v
}

// GetFileIdNumber returns the FileIdNumber field value if set, zero value otherwise.
func (o *File) GetFileIdNumber() string {
	if o == nil || IsNil(o.FileIdNumber) {
		var ret string
		return ret
	}
	return *o.FileIdNumber
}

// GetFileIdNumberOk returns a tuple with the FileIdNumber field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *File) GetFileIdNumberOk() (*string, bool) {
	if o == nil || IsNil(o.FileIdNumber) {
		return nil, false
	}
	return o.FileIdNumber, true
}

// HasFileIdNumber returns a boolean if a field has been set.
func (o *File) HasFileIdNumber() bool {
	if o != nil && !IsNil(o.FileIdNumber) {
		return true
	}

	return false
}

// SetFileIdNumber gets a reference to the given string and assigns it to the FileIdNumber field.
func (o *File) SetFileIdNumber(v string) {
	o.FileIdNumber = &v
}

// GetPhysicalRecordLength returns the PhysicalRecordLength field value if set, zero value otherwise.
func (o *File) GetPhysicalRecordLength() int32 {
	if o == nil || IsNil(o.PhysicalRecordLength) {
		var ret int32
		return ret
	}
	return *o.PhysicalRecordLength
}

// GetPhysicalRecordLengthOk returns a tuple with the PhysicalRecordLength field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *File) GetPhysicalRecordLengthOk() (*int32, bool) {
	if o == nil || IsNil(o.PhysicalRecordLength) {
		return nil, false
	}
	return o.PhysicalRecordLength, true
}

// HasPhysicalRecordLength returns a boolean if a field has been set.
func (o *File) HasPhysicalRecordLength() bool {
	if o != nil && !IsNil(o.PhysicalRecordLength) {
		return true
	}

	return false
}

// SetPhysicalRecordLength gets a reference to the given int32 and assigns it to the PhysicalRecordLength field.
func (o *File) SetPhysicalRecordLength(v int32) {
	o.PhysicalRecordLength = &v
}

// GetBlockSize returns the BlockSize field value if set, zero value otherwise.
func (o *File) GetBlockSize() int32 {
	if o == nil || IsNil(o.BlockSize) {
		var ret int32
		return ret
	}
	return *o.BlockSize
}

// GetBlockSizeOk returns a tuple with the BlockSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *File) GetBlockSizeOk() (*int32, bool) {
	if o == nil || IsNil(o.BlockSize) {
		return nil, false
	}
	return o.BlockSize, true
}

// HasBlockSize returns a boolean if a field has been set.
func (o *File) HasBlockSize() bool {
	if o != nil && !IsNil(o.BlockSize) {
		return true
	}

	return false
}

// SetBlockSize gets a reference to the given int32 and assigns it to the BlockSize field.
func (o *File) SetBlockSize(v int32) {
	o.BlockSize = &v
}

// GetVersionNumber returns the VersionNumber field value if set, zero value otherwise.
func (o *File) GetVersionNumber() int32 {
	if o == nil || IsNil(o.VersionNumber) {
		var ret int32
		return ret
	}
	return *o.VersionNumber
}

// GetVersionNumberOk returns a tuple with the VersionNumber field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *File) GetVersionNumberOk() (*int32, bool) {
	if o == nil || IsNil(o.VersionNumber) {
		return nil, false
	}
	return o.VersionNumber, true
}

// HasVersionNumber returns a boolean if a field has been set.
func (o *File) HasVersionNumber() bool {
	if o != nil && !IsNil(o.VersionNumber) {
		return true
	}

	return false
}

// SetVersionNumber gets a reference to the given int32 and assigns it to the VersionNumber field.
func (o *File) SetVersionNumber(v int32) {
	o.VersionNumber = &v
}

// GetFileControlTotal returns the FileControlTotal field value if set, zero value otherwise.
func (o *File) GetFileControlTotal() string {
	if o == nil || IsNil(o.FileControlTotal) {
		var ret string
		return ret
	}
	return *o.FileControlTotal
}

// GetFileControlTotalOk returns a tuple with the FileControlTotal field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *File) GetFileControlTotalOk() (*string, bool) {
	if o == nil || IsNil(o.FileControlTotal) {
		return nil, false
	}
	return o.FileControlTotal, true
}

// HasFileControlTotal returns a boolean if a field has been set.
func (o *File) HasFileControlTotal() bool {
	if o != nil && !IsNil(o.FileControlTotal) {
		return true
	}

	return false
}

// SetFileControlTotal gets a reference to the given string and assigns it to the FileControlTotal field.
func (o *File) SetFileControlTotal(v string) {
	o.FileControlTotal = &v
}

// GetNumberOfGroups returns the NumberOfGroups field value if set, zero value otherwise.
func (o *File) GetNumberOfGroups() int32 {
	if o == nil || IsNil(o.NumberOfGroups) {
		var ret int32
		return ret
	}
	return *o.NumberOfGroups
}

// GetNumberOfGroupsOk returns a tuple with the NumberOfGroups field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *File) GetNumberOfGroupsOk() (*int32, bool) {
	if o == nil || IsNil(o.NumberOfGroups) {
		return nil, false
	}
	return o.NumberOfGroups, true
}

// HasNumberOfGroups returns a boolean if a field has been set.
func (o *File) HasNumberOfGroups() bool {
	if o != nil && !IsNil(o.NumberOfGroups) {
		return true
	}

	return false
}

// SetNumberOfGroups gets a reference to the given int32 and assigns it to the NumberOfGroups field.
func (o *File) SetNumberOfGroups(v int32) {
	o.NumberOfGroups = &v
}

// GetNumberOfRecords returns the NumberOfRecords field value if set, zero value otherwise.
func (o *File) GetNumberOfRecords() int32 {
	if o == nil || IsNil(o.NumberOfRecords) {
		var ret int32
		return ret
	}
	return *o.NumberOfRecords
}

// GetNumberOfRecordsOk returns a tuple with the NumberOfRecords field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *File) GetNumberOfRecordsOk() (*int32, bool) {
	if o == nil || IsNil(o.NumberOfRecords) {
		return nil, false
	}
	return o.NumberOfRecords, true
}

// HasNumberOfRecords returns a boolean if a field has been set.
func (o *File) HasNumberOfRecords() bool {
	if o != nil && !IsNil(o.NumberOfRecords) {
		return true
	}

	return false
}

// SetNumberOfRecords gets a reference to the given int32 and assigns it to the NumberOfRecords field.
func (o *File) SetNumberOfRecords(v int32) {
	o.NumberOfRecords = &v
}

// GetGroups returns the Groups field value if set, zero value otherwise.
func (o *File) GetGroups() []Group {
	if o == nil || IsNil(o.Groups) {
		var ret []Group
		return ret
	}
	return o.Groups
}

// GetGroupsOk returns a tuple with the Groups field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *File) GetGroupsOk() ([]Group, bool) {
	if o == nil || IsNil(o.Groups) {
		return nil, false
	}
	return o.Groups, true
}

// HasGroups returns a boolean if a field has been set.
func (o *File) HasGroups() bool {
	if o != nil && !IsNil(o.Groups) {
		return true
	}

	return false
}

// SetGroups gets a reference to the given []Group and assigns it to the Groups field.
func (o *File) SetGroups(v []Group) {
	o.Groups = v
}

func (o File) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o File) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Sender) {
		toSerialize["sender"] = o.Sender
	}
	if !IsNil(o.Receiver) {
		toSerialize["receiver"] = o.Receiver
	}
	if !IsNil(o.FileCreatedDate) {
		toSerialize["fileCreatedDate"] = o.FileCreatedDate
	}
	if !IsNil(o.FileCreatedTime) {
		toSerialize["fileCreatedTime"] = o.FileCreatedTime
	}
	if !IsNil(o.FileIdNumber) {
		toSerialize["fileIdNumber"] = o.FileIdNumber
	}
	if !IsNil(o.PhysicalRecordLength) {
		toSerialize["physicalRecordLength"] = o.PhysicalRecordLength
	}
	if !IsNil(o.BlockSize) {
		toSerialize["blockSize"] = o.BlockSize
	}
	if !IsNil(o.VersionNumber) {
		toSerialize["versionNumber"] = o.VersionNumber
	}
	if !IsNil(o.FileControlTotal) {
		toSerialize["fileControlTotal"] = o.FileControlTotal
	}
	if !IsNil(o.NumberOfGroups) {
		toSerialize["numberOfGroups"] = o.NumberOfGroups
	}
	if !IsNil(o.NumberOfRecords) {
		toSerialize["numberOfRecords"] = o.NumberOfRecords
	}
	if !IsNil(o.Groups) {
		toSerialize["Groups"] = o.Groups
	}
	return toSerialize, nil
}

type NullableFile struct {
	value *File
	isSet bool
}

func (v NullableFile) Get() *File {
	return v.value
}

func (v *NullableFile) Set(val *File) {
	v.value = val
	v.isSet = true
}

func (v NullableFile) IsSet() bool {
	return v.isSet
}

func (v *NullableFile) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableFile(val *File) *NullableFile {
	return &NullableFile{value: val, isSet: true}
}

func (v NullableFile) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableFile) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
