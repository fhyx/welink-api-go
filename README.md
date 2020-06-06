# welink-api-go


## Environment for config

```
WELINK_CORP_ID=AppId
WELINK_CORP_SECRET=AppSecret

WELINK_TEST_UID='for unit test only'
```

## Usage

```go

import "fhyx.online/welink-api-go/welink"


api := NewAPI() // or New(appId, appSecret)

deptId := 0
recursive := false
data, err := api.ListDepartment(deptId, recursive)

uid := "yourUID"
at := "uid" // uid,mobile,cuid
user, err := api.GetUser(uid, at)

```

## Links

* https://open-doc.welink.huaweicloud.com/docs/serverapi/authorization/permission_internal.html?type=internal


## TODO

* Sync users
* Sync department
