package vuforia

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
	"encoding/json"
	"io/ioutil"
)

type VuforiaClient struct {
	AccessKey string
	SecretKey string
	Host      string
}

func Init(accessKey string, secretKey string) (*VuforiaClient, error) {

	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("accessKey or secretKey is empty")
	}

	client := &VuforiaClient{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Host:      "https://vws.vuforia.com",
	}
	return client, nil
}

func timeRFC1123() string {
	now := time.Now().UTC()
	return now.Format(time.RFC1123Z)
}

func (client *VuforiaClient) autheticatedResponse(request *http.Request) ([]byte, error) {
	dateRFC := timeRFC1123()
	dateRFC = strings.Replace(dateRFC, "+0000", "GMT", -1)
	//dateRFC = "Sat, 6 Jul 2017 11:38:19 GMT"
	md5Content, err := calculateContentMd5(request)
	contentType, err := getContentType(request)
	if err != nil {
		return nil, err
	}

	stringToSign := request.Method + "\n" + md5Content + "\n" + contentType + "\n" + dateRFC + "\n" + request.URL.Path
	signatureString := hmac_sha_base64(stringToSign, client.SecretKey)

	request.Header.Set("Date", dateRFC)

	authHeader := fmt.Sprintf("VWS %s:%s", client.AccessKey, signatureString)
	request.Header.Set("Host", "vws.vuforia.com")
	request.Header.Set("Authorization", authHeader)
	request.Proto = "HTTP/1.1"
	response, err := http.DefaultClient.Do(request)

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(responseData))

	str := string(responseData)
	str = strings.Trim(str,"")

	defer response.Body.Close()
	return responseData, nil
}

func hmac_sha_base64(content string, secret string) string {

	key := []byte(secret)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(content))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func getContentType(r *http.Request) (string, error) {
	if r.Method == "PUT" || r.Method == "POST" {
		return "application/json", nil
	}
	return "", nil
}

func (client *VuforiaClient) TargetIds() (resultItems []string, err error) {
	url, err := url.Parse(client.Host)
	if err != nil {
		log.Println(err.Error())
	}
	buffer := new(bytes.Buffer)
	url.Path = path.Join(url.Path, "targets")

	request, err := http.NewRequest(http.MethodGet, url.String(), buffer)
	responseBody, err := client.autheticatedResponse(request)



	fmt.Println(string(responseBody))
	resultModel := ResultItemsModel{}
	err = json.Unmarshal(responseBody,&resultModel)

	if 	err != nil{
		return
	}
	resultItems = resultModel.Results
	return
}

func (client *VuforiaClient) GetSummaryItem(id string)(item SummaryCloudItem,err error){
	url, err := url.Parse(client.Host)
	if err != nil {
		log.Println(err.Error())
	}
	buffer := new(bytes.Buffer)
	url.Path = path.Join(url.Path, "summary")
	url.Path = path.Join(url.Path, id)

	request, err := http.NewRequest(http.MethodGet, url.String(), buffer)
	a,err := client.autheticatedResponse(request)

	err = json.Unmarshal(a,&item)
	return
}


func (client *VuforiaClient) GetItem(id string)(item CloudItem,err error){
	url, err := url.Parse(client.Host)
	if err != nil {
		log.Println(err.Error())
	}
	buffer := new(bytes.Buffer)
	url.Path = path.Join(url.Path, "targets")
	url.Path = path.Join(url.Path, id)

	request, err := http.NewRequest(http.MethodGet, url.String(), buffer)
	a,err := client.autheticatedResponse(request)

	err = json.Unmarshal(a,&item)
	return
}

func (client *VuforiaClient) UpdateItem(targetId string,name string,width float32,imageBase64 string,activeFlag bool,metaBase64 string)(isOk bool,err error){
	url, err := url.Parse(client.Host)
	if err != nil {
		isOk = false
		log.Println(err.Error())
		return
	}
	url.Path = path.Join(url.Path, "targets")
	url.Path = path.Join(url.Path, targetId)

	newItem,err := createCloudItem(name,width,imageBase64,activeFlag,metaBase64)
	if err != nil {
		log.Println(err.Error())
	}

	bytesData,err := json.Marshal(newItem)
	if err != nil {
		isOk = false
		log.Println(err.Error())
		return
	}

	request, err := http.NewRequest(http.MethodPut, url.String(), bytes.NewBuffer(bytesData))
	request.Header.Add("Content-Type","application/json")
	_,err = client.autheticatedResponse(request)
	if err != nil {
		isOk = false
		log.Println(err.Error())
		return
	}
	isOk = true
	return
}


func (client *VuforiaClient) AddItem(name string,width float32,imageBase64 string,activeFlag bool,metaBase64 string)(targetId string,isSuccess bool,err error){
	url, err := url.Parse(client.Host)
	if err != nil {
		log.Println(err.Error())
	}

	newItem,err := createCloudItem(name,width,imageBase64,activeFlag,metaBase64)

	if err != nil{
		return "",false,err
	}
	bytesData,err := json.Marshal(newItem)
	if err != nil{
		return "",false,err
	}

	url.Path = path.Join(url.Path, "targets")

	request, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewBuffer(bytesData))
	if err != nil{
		return "",false,err
	}

	request.Header.Add("Content-Type","application/json")
	bodyData,err := client.autheticatedResponse(request)
	if err != nil{
		return "",false,err
	}

	responseModel := &ResultAddModel{}
	err = json.Unmarshal(bodyData,&responseModel)
	if err != nil{
		return "",false,err
	}

	return responseModel.TargetId,true,nil
}

func (client *VuforiaClient) DeleteId(id string) {
	url, err := url.Parse(client.Host)
	if err != nil {
		log.Println(err.Error())
	}
	buffer := new(bytes.Buffer)
	url.Path = path.Join(url.Path, "targets")
	url.Path = path.Join(url.Path, id)

	request, err := http.NewRequest(http.MethodDelete, url.String(), buffer)
	a,err := client.autheticatedResponse(request)
	_ = a
}

func calculateContentMd5(r *http.Request) (string, error) {

	dataReader, err := r.GetBody()
	if err != nil {
		return "", err
	}

	bb,err := ioutil.ReadAll(dataReader)
	hasher := md5.New()
	hasher.Write(bb)
	resultHash := hex.EncodeToString(hasher.Sum(nil))
	return resultHash, nil
}

func createCloudItem(name string,width float32,imageBase64 string,activeFlag bool,metaBase64 string)(*CloudItemNew,error){
	newItem := CloudItemNew{
		Name:name,
		Width: width,
		Image:imageBase64,
		Active_flag:activeFlag,
		Application_metadata:metaBase64,
	}
	return &newItem,nil
}