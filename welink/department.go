package welink

import (
	"github.com/fhyx/welink-api-go/client"
)

// DepartmentUp 部门更新请求对象
// "corpDeptCode": "dddddd",
// "corpParentCode": "ddfd",
// "deptNameCn": "生产部门",
// "deptNameEn": "produce dept",
// "deptLevel": "1",
// "managerId": "",
// "valid": "1",
// "orderNo": "1000"
type DepartmentUp struct {
	CorpDeptID   int    `json:"corpDeptCode,string"`
	CorpParentID int    `json:"corpParentCode,string"`
	NameCN       string `json:"deptNameCn"`
	NameEN       string `json:"deptNameEn"`
	Level        int    `json:"deptLevel,string"`
	Leader       string `json:"managerId"`
	OrderNo      int    `json:"orderNo,string"`
	Valid        int    `json:"valid,string"`
}

// Department 部门
// "deptCode": "1",
// "deptNameCn": "产品销售部",
// "deptNameEn": "Sales Dept",
// "fatherCode": "0",
// "deptLevel": "2",
// "orderNo": 1
type Department struct {
	CorpDeptID   int    `json:"corpDeptCode,string,omitempty"`
	CorpParentID int    `json:"corpParentCode,string,omitempty"`
	ID           int    `json:"deptCode,string"`
	ParentID     int    `json:"fatherCode,string"`
	NameCN       string `json:"deptNameCn"`
	NameEN       string `json:"deptNameEn"`
	Level        int    `json:"deptLevel,string"`
	Leader       string `json:"managerId,omitempty"`
	OrderNo      int    `json:"orderNo,string,omitempty"`
	HasChild     int    `json:"hasChildDept,omitempty"`
}

type Departments []Department

// default sort
func (z Departments) Len() int      { return len(z) }
func (z Departments) Swap(i, j int) { z[i], z[j] = z[j], z[i] }
func (z Departments) Less(i, j int) bool {
	return z[i].ParentID == 0 || z[i].ParentID < z[j].ParentID ||
		z[i].Level < z[j].Level || z[i].ParentID == z[j].ParentID && z[i].OrderNo > z[j].OrderNo
}

func (z Departments) WithID(id int) *Department {
	for _, dept := range z {
		if dept.ID == id {
			return &dept
		}
	}
	return nil
}

// departmentResponse
// "offset": 100,
// "limit": 25,
// "totalCount": 327,
// "departmentInfo": []
type departmentResponse struct {
	client.Error

	TotalCount  int `json:"totalCount"`
	Departments `json:"departmentInfo"`
}

// FilterDepartment Deprecated with Departments.WithID()
func FilterDepartment(data []Department, id int) (*Department, error) {
	for _, dept := range data {
		if dept.ID == id {
			return &dept, nil
		}
	}
	return nil, ErrNotFound
}
