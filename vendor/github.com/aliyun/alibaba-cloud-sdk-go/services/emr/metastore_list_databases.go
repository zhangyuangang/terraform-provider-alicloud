package emr

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// MetastoreListDatabases invokes the emr.MetastoreListDatabases API synchronously
// api document: https://help.aliyun.com/api/emr/metastorelistdatabases.html
func (client *Client) MetastoreListDatabases(request *MetastoreListDatabasesRequest) (response *MetastoreListDatabasesResponse, err error) {
	response = CreateMetastoreListDatabasesResponse()
	err = client.DoAction(request, response)
	return
}

// MetastoreListDatabasesWithChan invokes the emr.MetastoreListDatabases API asynchronously
// api document: https://help.aliyun.com/api/emr/metastorelistdatabases.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) MetastoreListDatabasesWithChan(request *MetastoreListDatabasesRequest) (<-chan *MetastoreListDatabasesResponse, <-chan error) {
	responseChan := make(chan *MetastoreListDatabasesResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.MetastoreListDatabases(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// MetastoreListDatabasesWithCallback invokes the emr.MetastoreListDatabases API asynchronously
// api document: https://help.aliyun.com/api/emr/metastorelistdatabases.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) MetastoreListDatabasesWithCallback(request *MetastoreListDatabasesRequest, callback func(response *MetastoreListDatabasesResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *MetastoreListDatabasesResponse
		var err error
		defer close(result)
		response, err = client.MetastoreListDatabases(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// MetastoreListDatabasesRequest is the request struct for api MetastoreListDatabases
type MetastoreListDatabasesRequest struct {
	*requests.RpcRequest
	ResourceOwnerId   requests.Integer `position:"Query" name:"ResourceOwnerId"`
	DbName            string           `position:"Query" name:"DbName"`
	PageSize          requests.Integer `position:"Query" name:"PageSize"`
	FuzzyDatabaseName string           `position:"Query" name:"FuzzyDatabaseName"`
	PageNumber        requests.Integer `position:"Query" name:"PageNumber"`
}

// MetastoreListDatabasesResponse is the response struct for api MetastoreListDatabases
type MetastoreListDatabasesResponse struct {
	*responses.BaseResponse
	RequestId    string       `json:"RequestId" xml:"RequestId"`
	Description  string       `json:"Description" xml:"Description"`
	TotalCount   int          `json:"TotalCount" xml:"TotalCount"`
	PageNumber   int          `json:"PageNumber" xml:"PageNumber"`
	PageSize     int          `json:"PageSize" xml:"PageSize"`
	DbNames      DbNames      `json:"DbNames" xml:"DbNames"`
	DatabaseList DatabaseList `json:"DatabaseList" xml:"DatabaseList"`
}

// CreateMetastoreListDatabasesRequest creates a request to invoke MetastoreListDatabases API
func CreateMetastoreListDatabasesRequest() (request *MetastoreListDatabasesRequest) {
	request = &MetastoreListDatabasesRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Emr", "2016-04-08", "MetastoreListDatabases", "emr", "openAPI")
	return
}

// CreateMetastoreListDatabasesResponse creates a response to parse from MetastoreListDatabases response
func CreateMetastoreListDatabasesResponse() (response *MetastoreListDatabasesResponse) {
	response = &MetastoreListDatabasesResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
