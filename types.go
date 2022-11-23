package main

import "time"

type Config struct {
	Instance string `json:"instance"`
	Token    string `json:"token"`
}

type User struct {
	// Info meta
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Profile
	Name        string `json:"name"`
	Username    string `json:"username"`
	Description string `json:"description"`
	AvatarUrl   string `json:"avatarUrl"`
	BannerUrl   string `json:"bannerUrl"`
	Location    string `json:"location"`
	IsBot       bool   `json:"isBot"`
	IsCat       bool   `json:"isCat"`
	Fields      []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"fields"`

	// Relation
	FollowersCount uint `json:"followersCount"`
	FollowingCount uint `json:"followingCount"`

	// Activity
	PinnedNoteIDs []string `json:"pinnedNoteIds"`
	// Ignore note details
	PinnedPageID string `json:"pinnedPageId"`
	// Ignore page details

	// Security
	TwoFactorEnabled bool `json:"twoFactorEnabled"`

	// Status
	IsSilenced  bool `json:"isSilenced"`
	IsSuspended bool `json:"isSuspended"`
}

type ShowUserRequest struct {
	I      string `json:"i"`
	Origin string `json:"origin"`
	Offset int    `json:"offset"`
	Limit  uint   `json:"limit"`
}

type IAndUserID struct {
	I      string `json:"i"`
	UserID string `json:"userId"`
}

type UserIDAndLimit struct {
	UserID string `json:"userId"`
	Limit  uint   `json:"limit"`
}

type NothingButID struct {
	ID string `json:"id"`
	// Ignore other fields
}
