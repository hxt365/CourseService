package store

import (
	"CourseService/db"
	"CourseService/model"
	"fmt"
)

type CacheStore struct {
	client *db.Redis
}

func NewCacheStore(rd *db.Redis) *CacheStore {
	return &CacheStore{client: rd}
}

func (cs *CacheStore) DeleteCoursesListCache() {
	cs.client.DelCache("courses")
}

func (cs *CacheStore) GetCoursesListCache(offset, limit int) (string, error) {
	cached, err := cs.client.GetCacheField("courses", fmt.Sprintf("%v-%v", offset, limit))
	return cached, err
}

func (cs *CacheStore) SetCoursesListCache(offset, limit int, val string) {
	cs.client.SetCacheField("courses", fmt.Sprintf("%v-%v", offset, limit), val)
}

func (cs *CacheStore) GetMentorsCoursesListCache(offset, limit int, id uint) (string, error) {
	cached, err := cs.client.GetCacheField(fmt.Sprintf("mentors:%v:courses", id), fmt.Sprintf("%v-%v", offset, limit))
	return cached, err
}

func (cs *CacheStore) SetMentorsCoursesListCache(offset, limit int, id uint, val string) {
	cs.client.SetCacheField(fmt.Sprintf("mentors:%v:courses", id), fmt.Sprintf("%v-%v", offset, limit), val)
}

func (cs *CacheStore) DeleteMentorsCoursesListCache(id uint) {
	cs.client.DelCache(fmt.Sprintf("mentors:%v:courses", id))
}

func (cs *CacheStore) GetStudentsCoursesListCache(offset, limit int, id uint) (string, error) {
	cached, err := cs.client.GetCacheField(fmt.Sprintf("students:%v:courses", id), fmt.Sprintf("%v-%v", offset, limit))
	return cached, err
}

func (cs *CacheStore) SetStudentsCoursesListCache(offset, limit int, id uint, val string) {
	cs.client.SetCacheField(fmt.Sprintf("students:%v:courses", id), fmt.Sprintf("%v-%v", offset, limit), val)
}

func (cs *CacheStore) DeleteStudentsCoursesListCache(id uint) {
	cs.client.DelCache(fmt.Sprintf("students:%v:courses", id))
}

func (cs *CacheStore) DeleteAllStudentsCoursesListCacheByCourse(course *model.Course) {
	for _, student := range course.Students {
		cs.DeleteStudentsCoursesListCache(student.ID)
	}
}

func (cs *CacheStore) GetCourseCache(id uint) (string, error) {
	cached, err := cs.client.GetCache(fmt.Sprintf("course:%v", id))
	return cached, err
}

func (cs *CacheStore) SetCourseCache(id uint, val string) {
	cs.client.SetCache(fmt.Sprintf("course:%v", id), val)
}

func (cs *CacheStore) DeleteCourseCache(id uint) {
	cs.client.DelCache(fmt.Sprintf("course:%v", id))
}

func (cs *CacheStore) GetTakenCourseCache(courseID, studentID uint) (bool, error) {
	cached, err := cs.client.GetCacheField(fmt.Sprintf("taken_course:%v", courseID), fmt.Sprintf("%v", studentID))
	return cached == "true", err
}

func (cs *CacheStore) SetTakenCourseCache(courseID, studentID uint, val bool) {
	valStr := "true"
	if !val {
		valStr = "false"
	}
	cs.client.SetCacheField(fmt.Sprintf("taken_course:%v", courseID), fmt.Sprintf("%v", studentID), valStr)
}

func (cs *CacheStore) DeleteTakenCourseCache(courseID, studentID uint) {
	cs.client.DelCacheField(fmt.Sprintf("taken_course:%v", courseID), fmt.Sprintf("%v", studentID))
}

func (cs *CacheStore) GetCoursesReviewsListCache(offset, limit int, id uint) (string, error) {
	cached, err := cs.client.GetCacheField(fmt.Sprintf("reviews:course:%v", id), fmt.Sprintf("%v-%v", offset, limit))
	return cached, err
}

func (cs *CacheStore) SetCoursesReviewsListCache(offset, limit int, id uint, val string)  {
	cs.client.SetCacheField(fmt.Sprintf("reviews:course:%v", id), fmt.Sprintf("%v-%v", offset, limit), val)
}

func (cs *CacheStore) DeleteCoursesReviewsListCache(id uint)  {
	cs.client.DelCache(fmt.Sprintf("reviews:course:%v", id))
}

func (cs *CacheStore) GetUsersReviewsListCache(offset, limit int, id uint) (string, error) {
	cached, err := cs.client.GetCacheField(fmt.Sprintf("reviews:user:%v", id), fmt.Sprintf("%v-%v", offset, limit))
	return cached, err
}
func (cs *CacheStore) SetUsersReviewsListCache(offset, limit int, id uint, val string)  {
	cs.client.SetCacheField(fmt.Sprintf("reviews:user:%v", id), fmt.Sprintf("%v-%v", offset, limit), val)
}

func (cs *CacheStore) DeleteUsersReviewsListCache(id uint)  {
	cs.client.DelCache(fmt.Sprintf("reviews:user:%v", id))
}

func (cs *CacheStore) GetCoursesNotificationsListCache(offset, limit int, id uint) (string ,error) {
	cached, err := cs.client.GetCacheField(fmt.Sprintf("notifications:course:%v", id), fmt.Sprintf("%v-%v", offset, limit))
	return cached, err
}

func (cs *CacheStore) SetCoursesNotificationsListCache(offset, limit int, id uint, val string)  {
	cs.client.SetCacheField(fmt.Sprintf("notifications:course:%v", id), fmt.Sprintf("%v-%v", offset, limit), val)
}

func (cs *CacheStore) DeleteCoursesNotificationsListCache(id uint)  {
	cs.client.DelCache(fmt.Sprintf("notifications:course:%v", id))
}
