package welink

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fhyx/welink-api-go/client"
)

const (
	urlToken = "https://open.welink.huaweicloud.com/api/auth/v2/tickets"

	urlUserGet      = "https://open.welink.huaweicloud.com/api/contact/v1/users"
	urlUserList     = "https://open.welink.huaweicloud.com/api/contact/v1/user/users"
	urlUserListSimp = "https://open.welink.huaweicloud.com/api/contact/v2/user/userid"
	urlUserBulk     = "https://open.welink.huaweicloud.com/api/contact/v1/users/bulk"

	urlDeptList = "https://open.welink.huaweicloud.com/api/contact/v2/departments/list"
	urlDeptSync = "https://open.welink.huaweicloud.com/api/contact/v2/departments/bulk"
)

type API struct {
	corpID     string
	corpSecret string
	c          *client.Client
}

func NewAPI() *API {
	return New(os.Getenv("WELINK_CORP_ID"), os.Getenv("WELINK_CORP_SECRET"))
}

// New ...
func New(corpID, corpSecret string) *API {
	if corpID == "" || corpSecret == "" {
		log.Printf("corpID or corpSecret are empty or not found")
	}
	c := client.NewClient(urlToken)
	c.SetContentType("application/json")
	c.SetCorp(corpID, corpSecret)
	return &API{
		corpID:     corpID,
		corpSecret: corpSecret,
		c:          c,
	}
}

func (a *API) CorpID() string {
	return a.corpID
}

func uriForUserGet(uid, at string) string {

	switch at {
	case "uid":
		return fmt.Sprintf("%s?userId=%s", urlUserGet, uid)
	case "mobile":
		return fmt.Sprintf("%s?mobileNumber=%s", urlUserGet, uid)
	default:
		return fmt.Sprintf("%s?corpUserId=%s", urlUserGet, uid)
	}
}

// GetUser get user with uid,mobile,cuid
func (a *API) GetUser(uid, at string) (*User, error) {
	user := new(User)
	err := a.c.GetJSON(uriForUserGet(uid, at), user)
	if err != nil {
		logger().Infow("get user fail", "at:"+at, uid, "err", err)
		return nil, err
	}
	return user, nil
}

func (a *API) ListDepartment(id int, recursive bool) (data Departments, err error) {

	var recursiveflag int
	if id > 0 && recursive {
		recursiveflag = 1
	}
	uri := fmt.Sprintf("%s?deptCode=%d&recursiveflag=%d&limit=100", urlDeptList, id, recursiveflag)

	var ret departmentResponse
	err = a.c.GetJSON(uri, &ret)

	if err == nil {
		data = ret.Departments
	}

	if recursive && id == 0 {
		for _, dept := range data {
			var child Departments
			child, err = a.ListDepartment(dept.ID, true)
			if err != nil {
				return
			}
			data = append(data, child...)
		}
	}

	return
}

// SyncDepartment ...
func (a *API) SyncDepartment(data []DepartmentUp) error {
	var req deptBatchReq
	req.Data = data

	buf, err := json.Marshal(&req)
	if err != nil {
		return err
	}
	var resp deptBatchResp
	err = a.c.PostJSON(urlDeptSync, buf, &resp)
	if err != nil {
		logger().Infow("sync department fail", "err", err)
		return err
	}
	logger().Infow("sync department ok", "resp", resp)
	return nil
}

// ListUser ...
func (a *API) ListUser(deptId int) (data []User, err error) {
	limit := 50
	uri := fmt.Sprintf("%s?&deptCode=%d&pageSize=%d", urlUserList, deptId, limit)

	var ret usersResponse
	err = a.c.GetJSON(uri, &ret)

	if err == nil {
		data = ret.Users
	}

	return
}
