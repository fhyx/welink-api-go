package welink

import (
	// "time"

	// "github.com/fhyx/welink-api-go/client"
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
	CorpUID      string        `json:"corpUserId,omitempty"`
	CorpDeptID   string        `json:"corpDeptCode,omitempty"`
	UserID       string        `json:"userId"`
	NameCN       string        `json:"userNameCn"`
	NameEN       string        `json:"userNameEn"`
	DepartmentID int           `json:"deptCode,string,omitempty"`
	Title        string        `json:"position,omitempty"`
	Mobile       string        `json:"mobileNumber,omitempty"`
	Email        string        `json:"userEmail,omitempty"`
	Tel          string        `json:"phoneNumber,omitempty"`
	Gender       gender.Gender `json:"sex,string,omitempty"`
	Status       Status        `json:"userStatus,omitempty"`
	Remark       string        `json:"remark,emitempty"`
	Address      string        `json:"address,emitempty"`
	IsActivated  uint8         `json:"isActivated,omitempty"`
	Created      string        `json:"creationTime,omitempty"`
	Updated      string        `json:"lastUpdatedTime,omitempty"`
}
