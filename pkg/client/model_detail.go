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

// Detail struct for Detail
type Detail struct {
	TypeCode                *string    `json:"TypeCode,omitempty"`
	Amount                  *string    `json:"Amount,omitempty"`
	FundsType               *FundsType `json:"FundsType,omitempty"`
	BankReferenceNumber     *string    `json:"BankReferenceNumber,omitempty"`
	CustomerReferenceNumber *string    `json:"CustomerReferenceNumber,omitempty"`
	Text                    *string    `json:"Text,omitempty"`
}

// NewDetail instantiates a new Detail object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDetail() *Detail {
	this := Detail{}
	return &this
}

// NewDetailWithDefaults instantiates a new Detail object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDetailWithDefaults() *Detail {
	this := Detail{}
	return &this
}

// GetTypeCode returns the TypeCode field value if set, zero value otherwise.
func (o *Detail) GetTypeCode() string {
	if o == nil || o.TypeCode == nil {
		var ret string
		return ret
	}
	return *o.TypeCode
}

// GetTypeCodeOk returns a tuple with the TypeCode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Detail) GetTypeCodeOk() (*string, bool) {
	if o == nil || o.TypeCode == nil {
		return nil, false
	}
	return o.TypeCode, true
}

// HasTypeCode returns a boolean if a field has been set.
func (o *Detail) HasTypeCode() bool {
	if o != nil && o.TypeCode != nil {
		return true
	}

	return false
}

// SetTypeCode gets a reference to the given string and assigns it to the TypeCode field.
func (o *Detail) SetTypeCode(v string) {
	o.TypeCode = &v
}

// GetAmount returns the Amount field value if set, zero value otherwise.
func (o *Detail) GetAmount() string {
	if o == nil || o.Amount == nil {
		var ret string
		return ret
	}
	return *o.Amount
}

// GetAmountOk returns a tuple with the Amount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Detail) GetAmountOk() (*string, bool) {
	if o == nil || o.Amount == nil {
		return nil, false
	}
	return o.Amount, true
}

// HasAmount returns a boolean if a field has been set.
func (o *Detail) HasAmount() bool {
	if o != nil && o.Amount != nil {
		return true
	}

	return false
}

// SetAmount gets a reference to the given string and assigns it to the Amount field.
func (o *Detail) SetAmount(v string) {
	o.Amount = &v
}

// GetFundsType returns the FundsType field value if set, zero value otherwise.
func (o *Detail) GetFundsType() FundsType {
	if o == nil || o.FundsType == nil {
		var ret FundsType
		return ret
	}
	return *o.FundsType
}

// GetFundsTypeOk returns a tuple with the FundsType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Detail) GetFundsTypeOk() (*FundsType, bool) {
	if o == nil || o.FundsType == nil {
		return nil, false
	}
	return o.FundsType, true
}

// HasFundsType returns a boolean if a field has been set.
func (o *Detail) HasFundsType() bool {
	if o != nil && o.FundsType != nil {
		return true
	}

	return false
}

// SetFundsType gets a reference to the given FundsType and assigns it to the FundsType field.
func (o *Detail) SetFundsType(v FundsType) {
	o.FundsType = &v
}

// GetBankReferenceNumber returns the BankReferenceNumber field value if set, zero value otherwise.
func (o *Detail) GetBankReferenceNumber() string {
	if o == nil || o.BankReferenceNumber == nil {
		var ret string
		return ret
	}
	return *o.BankReferenceNumber
}

// GetBankReferenceNumberOk returns a tuple with the BankReferenceNumber field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Detail) GetBankReferenceNumberOk() (*string, bool) {
	if o == nil || o.BankReferenceNumber == nil {
		return nil, false
	}
	return o.BankReferenceNumber, true
}

// HasBankReferenceNumber returns a boolean if a field has been set.
func (o *Detail) HasBankReferenceNumber() bool {
	if o != nil && o.BankReferenceNumber != nil {
		return true
	}

	return false
}

// SetBankReferenceNumber gets a reference to the given string and assigns it to the BankReferenceNumber field.
func (o *Detail) SetBankReferenceNumber(v string) {
	o.BankReferenceNumber = &v
}

// GetCustomerReferenceNumber returns the CustomerReferenceNumber field value if set, zero value otherwise.
func (o *Detail) GetCustomerReferenceNumber() string {
	if o == nil || o.CustomerReferenceNumber == nil {
		var ret string
		return ret
	}
	return *o.CustomerReferenceNumber
}

// GetCustomerReferenceNumberOk returns a tuple with the CustomerReferenceNumber field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Detail) GetCustomerReferenceNumberOk() (*string, bool) {
	if o == nil || o.CustomerReferenceNumber == nil {
		return nil, false
	}
	return o.CustomerReferenceNumber, true
}

// HasCustomerReferenceNumber returns a boolean if a field has been set.
func (o *Detail) HasCustomerReferenceNumber() bool {
	if o != nil && o.CustomerReferenceNumber != nil {
		return true
	}

	return false
}

// SetCustomerReferenceNumber gets a reference to the given string and assigns it to the CustomerReferenceNumber field.
func (o *Detail) SetCustomerReferenceNumber(v string) {
	o.CustomerReferenceNumber = &v
}

// GetText returns the Text field value if set, zero value otherwise.
func (o *Detail) GetText() string {
	if o == nil || o.Text == nil {
		var ret string
		return ret
	}
	return *o.Text
}

// GetTextOk returns a tuple with the Text field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Detail) GetTextOk() (*string, bool) {
	if o == nil || o.Text == nil {
		return nil, false
	}
	return o.Text, true
}

// HasText returns a boolean if a field has been set.
func (o *Detail) HasText() bool {
	if o != nil && o.Text != nil {
		return true
	}

	return false
}

// SetText gets a reference to the given string and assigns it to the Text field.
func (o *Detail) SetText(v string) {
	o.Text = &v
}

func (o Detail) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.TypeCode != nil {
		toSerialize["TypeCode"] = o.TypeCode
	}
	if o.Amount != nil {
		toSerialize["Amount"] = o.Amount
	}
	if o.FundsType != nil {
		toSerialize["FundsType"] = o.FundsType
	}
	if o.BankReferenceNumber != nil {
		toSerialize["BankReferenceNumber"] = o.BankReferenceNumber
	}
	if o.CustomerReferenceNumber != nil {
		toSerialize["CustomerReferenceNumber"] = o.CustomerReferenceNumber
	}
	if o.Text != nil {
		toSerialize["Text"] = o.Text
	}
	return json.Marshal(toSerialize)
}

type NullableDetail struct {
	value *Detail
	isSet bool
}

func (v NullableDetail) Get() *Detail {
	return v.value
}

func (v *NullableDetail) Set(val *Detail) {
	v.value = val
	v.isSet = true
}

func (v NullableDetail) IsSet() bool {
	return v.isSet
}

func (v *NullableDetail) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDetail(val *Detail) *NullableDetail {
	return &NullableDetail{value: val, isSet: true}
}

func (v NullableDetail) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDetail) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}