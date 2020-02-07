package welink

import (
	// "time"

	"github.com/fhyx/welink-api-go/client"
	"github.com/fhyx/welink-api-go/gender"
)

// Status 状态
type Status uint8

// 状态, 1：未开户 2：开户中 3：已开户 4:已销户
const (
	SNone     Status = 0
	SInactive Status = 1
	SActiving Status = 2
	SActived  Status = 3
	SClosed   Status = 4
)

func (s Status) String() string {
	switch s {
	case SActived:
		return "actived"
	case SInactive:
		return "inactive"
	case SActiving:
		return "activing"
	case SClosed:
		return "closed"
	default:
		return "none"
	}
}

// User 用户
// "code": "0",
// "message": "ok",
// "userStatus": "1",
// "userId": "zhangshan2@welink",
// "deptCode": "10001",
// "mobileNumber": "+86-15811847236",
// "userNameCn": "张三",
// "userNameEn": "zhangshan",
// "sex": "M",
// "corpUserId": "36188",
// "userEmail": "zhangshan4@126.com",
// "secretary": "zhangshan@welink",
// "phoneNumber": "0755-88888888",
// "address": "广东省深圳",
// "remark": "欢迎加入WeLink",
// "isActivated": 1,
// "creationTime": "2018-05-03 13:58:02",
// "lastUpdatedTime": "2018-05-03 13:58:02"
type User struct {
	CorpUID      string        `json:"corpUserId"`                // required
	CorpDeptID   int           `json:"corpDeptCode,string"`       // required
	UserID       string        `json:"userId"`                    // required
	NameCN       string        `json:"userNameCn"`                // required
	NameEN       string        `json:"userNameEn"`                // required
	DepartmentID int           `json:"deptCode,string,omitempty"` // deptCode at welink
	Mobile       string        `json:"mobileNumber"`              // required
	Phone        string        `json:"phoneNumber,omitempty"`     // required
	Email        string        `json:"userEmail"`                 // required
	Gender       gender.Gender `json:"sex,string,omitempty"`
	Status       Status        `json:"userStatus,omitempty"`
	Remark       string        `json:"remark,emitempty"`
	Address      string        `json:"address,emitempty"`
	Activated    uint8         `json:"isActivated,omitempty"`
	Createds     string        `json:"creationTime,omitempty"`
	Updateds     string        `json:"lastUpdatedTime,omitempty"`

	IsOpenAccount      int `json:"isOpenAccount,string,omitempty"`      // required
	Valid              int `json:"valid,string,string"`                 // required
	IsHideMobileNumber int `json:"isHideMobileNumber,string,omitempty"` // 1 public default, 2 hide
	OrderInDepts       int `json:"orderInDepts,string"`
}

func (u User) IsActived() bool {
	return u.Status == SActived
}

func (u User) IsEnabled() bool {
	return u.Activated == 1
}

type usersResponse struct {
	client.Error

	Total int    `json:"total"`
	Users []User `json:"data"`
}

type UserUp = User

type userBatchReq struct {
	Data []UserUp `json:"personInfo"`
}

// UserRespItem ...
type UserRespItem struct {
	client.Error

	CorpUID string `json:"corpUserId"`
}

type userBatchResp struct {
	Data []UserRespItem `json:"data"`
}

type userStatusUp struct {
	CorpUID string `json:"corpUserId,omitempty"`
	Mobile  string `json:"mobileNumber,omitempty"`
	Email   string `json:"userEmail,omitempty"`
}

type userStatusReq struct {
	Data []userStatusUp `json:"personInfo"`
}
