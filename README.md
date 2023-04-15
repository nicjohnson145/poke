# poke
Test APIs from yaml files

## Defining Tests

A test is simply a yaml file, defining a sequence of calls to execute. Below is an example sequence
file.

```yaml
vars:
  service_host: localhost:50051
  service: functionaltest.EchoService
  msg_body: "env MSG_BODY"

calls:
- name: repush
  type: grpc
  service-host: '{{ .service_host }}'
  url: '{{ .service }}/Echo'
  body:
    message: '{{ .msg_body }}'
  asserts:
  - jq: '.message'
    expected: 'foo'
- name: expectedErr
  type: grpc
  service-host: '{{ .service_host }}'
  url: '{{ .service }}/Err'
  body:
    code: 5
  want-status: 5
```

### Sequence Available Fields

| Key | Description | Required |
| --- | ----------- | -------- |
| vars | a map of variables that can be expanded using go's `text/template` syntax in calls | No |
| imports | a map of names to paths of other sequence files to import (note: imported files cannot themselves contain imports) | No |
| calls | the list of `Call` objects defining this sequence | Yes |

### Call Available Fields

| Key | Description | Required |
| --- | ----------- | -------- |
| name | The name of this call, makes logging pretty | No |
| type | The protocol type of this call (http or grpc) | No, defaults to http |
| body | the body of the request | No |
| headers | a map of headers to attach to the request | No |
| service-host | the url of a grpc service, only read if type is `grpc` | Conditionally |
| url | the http url, or the `service/Method` of a grpc request | Yes |
| method | the method to use for a http request, defaults to GET, or POST if body is given | No |
| want-status | if the expected status of call is not 200 (or 0 in GRPC), this prevents the call from being interpreted as an error | No |
| exports | A list of export directives to extract information from the returned response, this data will be made available to go `text/template` substitutions in later calls | No |
| asserts | A list of assert directives to assert information about the returned response | No |
| skip-verify | Indicates is TLS verification should be skipped for this request | No |
| from-import | Execute a call from an imported file | No |

### Export Available Fields

| Key | Description | Required |
| --- | ----------- | -------- |
| jq | the jq selector to use to get the data | Yes |
| as | what variable the data should be exported to later calls under | Yes |


### Assert Available Fields

| Key | Description | Required |
| --- | ----------- | -------- |
| jq | the jq selector to use to get the data | Yes |
| expected | the expected value for the data | Yes |

### ImportedCall Available Fields

| Key | Description | Required |
| --- | ----------- | -------- |
| name | the name of this call, makes logging nice | No |
| call | the name of the call to execute in the imported sequence | Yes |
