package handler

import "CourseService/store"

type Handler struct {
	courseStore       *store.CourseStore
	notificationStore *store.NotificationStore
	reviewStore       *store.ReviewStore
	cacheStore        *store.CacheStore
}

func NewHandler(cs *store.CourseStore, ns *store.NotificationStore, rs *store.ReviewStore, ca *store.CacheStore) *Handler {
	return &Handler{
		courseStore:       cs,
		notificationStore: ns,
		reviewStore:       rs,
		cacheStore:        ca,
	}
}
