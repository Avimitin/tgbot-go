package mark

import (
	"errors"
	"math"
	"sync"
)

// Mark store information about a mark
type Mark struct {
	owner       string
	ID          int32
	Link        string
	Description string
}

func newMark(id int32, owner, link, description string) Mark {
	return Mark{
		ID:          id,
		owner:       owner,
		Link:        link,
		Description: description,
	}
}

var (
	innerDB = struct {
		db     map[string][]Mark
		lock   sync.RWMutex
		amount int32
		limit  int32
	}{
		db:     map[string][]Mark{},
		amount: 0,
		limit:  math.MaxInt32,
	}
)

// AddMark map the given link to the user
func AddMark(user, link, description string) error {
	innerDB.lock.Lock()
	defer innerDB.lock.Unlock()

	if innerDB.amount >= innerDB.limit {
		return errors.New("storage over limit")
	}

	basicMark := newMark(0, user, link, description)

	if linkList, ok := innerDB.db[user]; ok {
		if length := len(linkList); length > 0 {
			basicMark.ID = linkList[length-1].ID + 1
			linkList = append(linkList, basicMark)
		} else {
			linkList = append(linkList, basicMark)
		}
		innerDB.db[user] = linkList
	} else {
		innerDB.db[user] = []Mark{basicMark}
	}

	innerDB.amount++

	return nil
}

// GetMark search links from given user name.
// Return links if user exist, return error if user not found.
func GetMark(user string) ([]Mark, error) {
	innerDB.lock.RLock()
	defer innerDB.lock.RUnlock()

	linkList, ok := innerDB.db[user]

	if !ok {
		return nil, errors.New("user not found")
	}

	return linkList, nil
}

// DelMark delete the mark with the same id as the given target.
// Return error if user not found or the target not found
func DelMark(user string, target int32) error {
	innerDB.lock.Lock()
	defer innerDB.lock.Unlock()

	if _, ok := innerDB.db[user]; !ok {
		return errors.New("user not found")
	}

	marks := innerDB.db[user]
	for i, mark := range marks {
		if mark.ID == target {
			innerDB.db[user] = deleteElement(marks, i)
			innerDB.amount--
			return nil
		}
	}

	return errors.New("target not found")
}

func deleteElement(markList []Mark, i int) []Mark {
	if len(markList)-1 == i {
		markList = markList[:i]
	} else {
		part := markList[i+1:]
		markList = append(markList[:i], part...)
	}

	return markList
}
