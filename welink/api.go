package welink

import (
	// "encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fhyx/welink-api-go/client"
)

const (
	urlToken = "https://open.welink.huaweicloud.com/api/auth/v2/tickets"

	urlUserGet = "https://open.welink.huaweicloud.com/api/contact/v1/users"

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

// func (a *API) ListUser(deptId int, incChild bool) (data []User, err error) {
// 	var token string
// 	token, err = a.c.GetAuthToken()
// 	if err != nil {
// 		return
// 	}

// 	fc := "0"
// 	if incChild {
// 		fc = "1"
// 	}
// 	uri := fmt.Sprintf("%s?access_token=%s&department_id=%d&fetch_child=%s", urlListUser, token, deptId, fc)

// 	var ret usersResponse
// 	err = a.c.GetJSON(uri, &ret)

// 	if err == nil {
// 		data = ret.Users
// 	}

// 	return
// }

// func (a *API) GetOAuth2User(agentID int, code string) (ou *OAuth2UserInfo, err error) {
// 	var token string
// 	token, err = a.c.GetAuthToken()
// 	if err != nil {
// 		return
// 	}

// 	uri := fmt.Sprintf("%s?access_token=%s&agentid=%d&code=%s", urlOAuth2GetUser, token, agentID, code)

// 	ou = new(OAuth2UserInfo)
// 	err = a.c.GetJSON(uri, ou)

// 	return
// }
