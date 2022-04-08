# \Bai2FilesApi

All URIs are relative to *http://localhost:8208*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Health**](Bai2FilesApi.md#Health) | **Get** /health | health metro2 service
[**Parse**](Bai2FilesApi.md#Parse) | **Post** /parse | Parse bai2 file after parse bin file
[**Print**](Bai2FilesApi.md#Print) | **Post** /print | Print bai2 file after parse bin file



## Health

> string Health(ctx).Execute()

health metro2 service



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.Bai2FilesApi.Health(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `Bai2FilesApi.Health``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `Health`: string
    fmt.Fprintf(os.Stdout, "Response from `Bai2FilesApi.Health`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiHealthRequest struct via the builder pattern


### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Parse

> string Parse(ctx).Input(input).Execute()

Parse bai2 file after parse bin file



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    input := os.NewFile(1234, "some_file") // *os.File | bai2 bin file (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.Bai2FilesApi.Parse(context.Background()).Input(input).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `Bai2FilesApi.Parse``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `Parse`: string
    fmt.Fprintf(os.Stdout, "Response from `Bai2FilesApi.Parse`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiParseRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **input** | ***os.File** | bai2 bin file | 

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Print

> string Print(ctx).Input(input).Execute()

Print bai2 file after parse bin file



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    input := os.NewFile(1234, "some_file") // *os.File | bai2 bin file (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.Bai2FilesApi.Print(context.Background()).Input(input).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `Bai2FilesApi.Print``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `Print`: string
    fmt.Fprintf(os.Stdout, "Response from `Bai2FilesApi.Print`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPrintRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **input** | ***os.File** | bai2 bin file | 

### Return type

**string**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

