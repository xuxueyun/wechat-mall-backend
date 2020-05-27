package dbops

import (
	"strconv"
	"time"
	"wechat-mall-backend/model"
)

type categoryDao struct{}

var CategoryDao = new(categoryDao)

const categoryColumnList = `
id, parent_id, name, sort, online, picture, description, is_del, create_time, update_time
`

func (*categoryDao) List(pid, page, size int) (*[]model.WechatMallCategoryDO, error) {
	sql := "SELECT " + categoryColumnList + " FROM wechat_mall_category WHERE is_del = 0 AND parent_id = " + strconv.Itoa(pid) + " ORDER BY sort"
	if page > 0 && size > 0 {
		sql += " LIMIT " + strconv.Itoa((page-1)*size) + ", " + strconv.Itoa(size)
	}
	rows, err := dbConn.Query(sql)
	if err != nil {
		return nil, err
	}
	var cateList []model.WechatMallCategoryDO
	for rows.Next() {
		category := model.WechatMallCategoryDO{}
		err := rows.Scan(&category.Id, &category.ParentId, &category.Name, &category.Sort, &category.Online, &category.Picture, &category.Description, &category.Del, &category.CreateTime, &category.UpdateTime)
		if err != nil {
			return nil, err
		}
		cateList = append(cateList, category)
	}
	return &cateList, nil
}

func (*categoryDao) CountByPid(pid int) (int, error) {
	sql := "SELECT COUNT(*) FROM wechat_mall_category WHERE is_del = 0 AND parent_id = " + strconv.Itoa(pid)
	rows, err := dbConn.Query(sql)
	if err != nil {
		return 0, err
	}
	total := 0
	for rows.Next() {
		err := rows.Scan(&total)
		if err != nil {
			return 0, err
		}
	}
	return total, nil
}

func (*categoryDao) QueryById(id int) (*model.WechatMallCategoryDO, error) {
	sql := "SELECT " + categoryColumnList + " FROM wechat_mall_category WHERE id = " + strconv.Itoa(id)
	rows, err := dbConn.Query(sql)
	if err != nil {
		return nil, err
	}
	category := model.WechatMallCategoryDO{}
	for rows.Next() {
		err := rows.Scan(&category.Id, &category.ParentId, &category.Name, &category.Sort, &category.Online, &category.Picture, &category.Description, &category.Del, &category.CreateTime, &category.UpdateTime)
		if err != nil {
			return nil, err
		}
	}
	return &category, nil
}

func (*categoryDao) QueryByName(name string) (*model.WechatMallCategoryDO, error) {
	sql := "SELECT " + categoryColumnList + " FROM wechat_mall_category WHERE is_del = 0 AND name = '" + name + "'"
	rows, err := dbConn.Query(sql)
	if err != nil {
		return nil, err
	}
	category := model.WechatMallCategoryDO{}
	for rows.Next() {
		err := rows.Scan(&category.Id, &category.ParentId, &category.Name, &category.Sort, &category.Online, &category.Picture, &category.Description, &category.Del, &category.CreateTime, &category.UpdateTime)
		if err != nil {
			return nil, err
		}
	}
	return &category, nil
}

func (*categoryDao) Insert(category *model.WechatMallCategoryDO) error {
	sql := "INSERT INTO wechat_mall_category ( " + categoryColumnList[4:] + " ) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmt, err := dbConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(category.ParentId, category.Name, category.Sort, category.Online, category.Picture, category.Description, 0, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (*categoryDao) Update(category *model.WechatMallCategoryDO) error {
	sql := `
UPDATE wechat_mall_category
SET parent_id = ?, name = ?, sort = ?, online = ?, picture = ?, description = ?, is_del = ?, update_time = ? 
WHERE id = ?
`
	stmt, err := dbConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(category.ParentId, category.Name, category.Sort, category.Online, category.Picture, category.Description, category.Del, time.Now(), category.Id)
	if err != nil {
		return err
	}
	return nil
}

func (*categoryDao) QuerySubCategoryByParentId(categoryId int) (*[]int, error) {
	sql := "SELECT id FROM wechat_mall_category WHERE is_del = 0 AND parent_id = " + strconv.Itoa(categoryId)
	rows, err := dbConn.Query(sql)
	if err != nil {
		return nil, err
	}
	ids := []int{}
	for rows.Next() {
		id := 0
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return &ids, nil
}

func (*categoryDao) UpdateSubCategoryOnline(categoryId, online int) error {
	sql := "UPDATE wechat_mall_category SET update_time = now(), online = " + strconv.Itoa(online) + " WHERE is_del = 0 AND parent_id = " + strconv.Itoa(categoryId)
	_, err := dbConn.Exec(sql)
	return err
}

// 查询-所有二级分类
func (*categoryDao) QueryAllSubCategory() (*[]model.WechatMallCategoryDO, error) {
	sql := "SELECT " + categoryColumnList + " FROM wechat_mall_category WHERE is_del = 0 AND online = 1 AND parent_id != 0"
	rows, err := dbConn.Query(sql)
	if err != nil {
		return nil, err
	}
	categoryList := []model.WechatMallCategoryDO{}
	for rows.Next() {
		category := model.WechatMallCategoryDO{}
		err := rows.Scan(&category.Id, &category.ParentId, &category.Name, &category.Sort, &category.Online,
			&category.Picture, &category.Description, &category.Del, &category.CreateTime, &category.UpdateTime)
		if err != nil {
			panic(err)
		}
		categoryList = append(categoryList, category)
	}
	return &categoryList, nil
}
