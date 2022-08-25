package wechatwork

import (
	"github.com/pkg/errors"
	"water-reminder/pkg/httpclient"
	"water-reminder/pkg/wechatwork/constant"
)

func IsError(resp httpclient.Response) error {
	code, err := resp.GetInt(constant.KeyErrcode)
	if err != nil {
		return err
	}
	if code == 0 {
		return nil
	}
	msg, err := resp.GetString(constant.KeyErrmsg)
	if err != nil {
		return err
	}
	return errors.Errorf("%s: %d, %s: %s", constant.KeyErrcode, code, constant.KeyErrmsg, msg)
}
