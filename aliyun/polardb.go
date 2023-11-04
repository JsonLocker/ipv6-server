package aliyun

import (
	"aliyun/helper"
	"os"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	polardb20170801 "github.com/alibabacloud-go/polardb-20170801/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/joho/godotenv"
)

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func createClient() (_result *polardb20170801.Client, _err error) {

	err := godotenv.Load()
	helper.ErrorMsg(err)
	AccessKeyId := os.Getenv("ALIBABA_ACCESS_ID")
	AccessKeySecret := os.Getenv("ALIBABA_ACCESS_SECRET")
	config := &openapi.Config{
		AccessKeyId:     &AccessKeyId,
		AccessKeySecret: &AccessKeySecret,
	}

	// Endpoint 请参考 https://api.aliyun.com/product/polardb
	config.Endpoint = tea.String("polardb.aliyuncs.com")
	_result = &polardb20170801.Client{}
	_result, _err = polardb20170801.NewClient(config)
	return _result, _err
}

// 更新PolarDB 的ip白名单
// aliyun.IpWhite("117.26.247.152", "auto_flash", "pc-bp12n5p82q7z0of44")
// 参考文档 go get github.com/alibabacloud-go/polardb-20170801/v4
func IpWhite(ip string, group_name string, db_cluster_id string) (_err error) {
	client, _err := createClient()
	if _err != nil {
		return _err
	}
	modifyDBClusterAccessWhitelistRequest := &polardb20170801.ModifyDBClusterAccessWhitelistRequest{
		DBClusterId:          tea.String(db_cluster_id),
		DBClusterIPArrayName: tea.String(group_name),
		SecurityIps:          tea.String(ip),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		_, _err = client.ModifyDBClusterAccessWhitelistWithOptions(modifyDBClusterAccessWhitelistRequest, runtime)
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
