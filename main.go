package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	// Read config
	var cfg Config
	configBytes, _ := os.ReadFile("config.json")
	_ = json.Unmarshal(configBytes, &cfg)
	log.Printf("%v", cfg)

	// Get all users
	log.Printf("Get all users...")
	offset := 0
	allUserCount := 0
	emptyUserAlreadySuspendedCount := 0
	emptyUserNewlySuspendedCount := 0

	idSet := make(map[string]struct{})

	isNotEnd := true
	const limit = 100
	for isNotEnd {
		log.Printf("Grabbing users from offset %d", offset)

		reqBody := ShowUserRequest{
			I:      cfg.Token,
			Origin: "local",
			Offset: offset,
			Limit:  limit,
		}
		reqBodyBytes, _ := json.Marshal(&reqBody)

		req, _ := http.NewRequest("POST", cfg.Instance+"/api/admin/show-users", bytes.NewReader(reqBodyBytes))

		res, _ := (&http.Client{}).Do(req)

		var resU []User
		_ = json.NewDecoder(res.Body).Decode(&resU)

		offset += len(resU)

		if len(resU) < limit {
			log.Printf("Reach the end")
			isNotEnd = false
		}

		for _, u := range resU {
			// Check ID for duplicate
			if _, ok := idSet[u.ID]; ok {
				// Already analyzed
				continue
			} else {
				// New
				idSet[u.ID] = struct{}{}
			}

			allUserCount++

			if u.Name == "" && // No name
				u.Description == "" && // No description
				u.AvatarUrl == cfg.Instance+"/identicon/"+u.ID && // Default generated avatar
				u.BannerUrl == "" && // No banner
				u.Location == "" && // No location
				u.IsBot == false && // Default
				u.IsCat == false && // Default
				len(u.Fields) == 0 && // No fields
				u.FollowersCount == 0 && // No follower
				u.FollowingCount == 0 && // No following
				len(u.PinnedNoteIDs) == 0 && // No pin notes
				u.PinnedPageID == "" && // No pin page
				u.TwoFactorEnabled == false { // No 2FA
				//log.Printf("New empty user found: %v", u)
				if u.IsSuspended {
					emptyUserAlreadySuspendedCount++
				} else {
					// Check note & page
					if CheckNothingButIDCount(cfg.Instance+"/api/users/notes", u.ID, 1) == 0 && // No note
						CheckNothingButIDCount(cfg.Instance+"/api/users/pages", u.ID, 1) == 0 { // No page

						// Absolutely nothing
						log.Printf("New empty user found:\t%s ( %s )", u.ID, u.Username)

						// Suspend
						SuspendUser(cfg.Instance, cfg.Token, u.ID)

						// Add counter
						emptyUserNewlySuspendedCount++
					}
				}
			}
		}
	}

	// Log statics
	log.Printf("%d total users, %d newly suspended empty users, %d already suspended empty users", allUserCount, emptyUserNewlySuspendedCount, emptyUserAlreadySuspendedCount)

}

func CheckNothingButIDCount(endpoint string, userId string, limit uint) int {

	reqBody := UserIDAndLimit{
		UserID: userId,
		Limit:  limit,
	}

	reqBodyBytes, _ := json.Marshal(&reqBody)

	req, _ := http.NewRequest("POST", endpoint, bytes.NewReader(reqBodyBytes))

	res, _ := (&http.Client{}).Do(req)

	var resBody []NothingButID

	_ = json.NewDecoder(res.Body).Decode(&resBody)

	return len(resBody)

}

func SuspendUser(instance string, token string, userId string) {

	reqBody := IAndUserID{
		I:      token,
		UserID: userId,
	}

	reqBodyBytes, _ := json.Marshal(&reqBody)

	req, _ := http.NewRequest("POST", instance+"/api/admin/suspend-user", bytes.NewReader(reqBodyBytes))

	_, _ = (&http.Client{}).Do(req)

}
