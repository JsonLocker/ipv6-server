package aliyun

import (
	"aliyun/helper"
	"encoding/json"
	"os"

	alidns20150109 "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/joho/godotenv"
)

type Record struct {
	Status   string `json: "Status"`
	RR       string `json: "RR"`
	Type     string `json: "Type"`
	Value    string `json: "Value"`
	RecordId string `json:"RecordId"`
}

// 实例化dns
func CreateClient() (_result *alidns20150109.Client, _err error) {
	err := godotenv.Load()
	helper.ErrorMsg(err)
	AccessKeyId := os.Getenv("ALIBABA_ACCESS_ID")
	AccessKeySecret := os.Getenv("ALIBABA_ACCESS_SECRET")
	config := &openapi.Config{
		AccessKeyId:     &AccessKeyId,
		AccessKeySecret: &AccessKeySecret,
	}
	config.Endpoint = tea.String("alidns.cn-hangzhou.aliyuncs.com")
	_result = &alidns20150109.Client{}
	_result, _err = alidns20150109.NewClient(config)
	return _result, _err
}

// 创建dns
func AddDns(RR string, domain string, ip string, dns_type string) (_err error) {
	client, _err := CreateClient()
	helper.ErrorMsg(_err)

	addDomainRecordRequest := &alidns20150109.AddDomainRecordRequest{
		DomainName: tea.String(domain),
		RR:         tea.String(RR),
		Type:       tea.String(dns_type),
		Value:      tea.String(ip),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		_, _err = client.AddDomainRecordWithOptions(addDomainRecordRequest, runtime)
		if _err != nil {
			return _err
		}
		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			return _err
		}
	}
	return _err
}

// 创建dns
func UpdateDns(RR string, recordId string, ip string, dns_type string) (_err error) {
	client, _err := CreateClient()
	helper.ErrorMsg(_err)

	updateDomainRecordRequest := &alidns20150109.UpdateDomainRecordRequest{
		RecordId: tea.String(recordId),
		RR:       tea.String(RR),
		Type:     tea.String(dns_type),
		Value:    tea.String(ip),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		_, _err = client.UpdateDomainRecordWithOptions(updateDomainRecordRequest, runtime)
		if _err != nil {
			return _err
		}

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			return _err
		}
	}
	return _err
}

// 查询DNS列表
func RecordList(domain string) (records []Record, _err error) {
	client, _err := CreateClient()
	helper.ErrorMsg(_err)
	records = []Record{}

	describeDomainRecordsRequest := &alidns20150109.DescribeDomainRecordsRequest{
		DomainName: tea.String(domain),
	}
	runtime := &util.RuntimeOptions{}

	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		domain_list, _err := client.DescribeDomainRecordsWithOptions(describeDomainRecordsRequest, runtime)
		helper.ErrorMsg(_err)
		list := domain_list.Body.DomainRecords.Record

		for _, v := range list {
			recode := Record{}
			json.Unmarshal([]byte(v.String()), &recode)
			records = append(records, recode)

		}
		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		_, _err = util.AssertAsString(error.Message)
		if _err != nil {
			return records, _err
		}
	}
	return records, _err
}

// 更新DNS
