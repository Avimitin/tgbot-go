package mark

import (
	"errors"
	"sync"
)

type mark struct {
	id          int32
	link        string
	description string
}

var (
	innerDB = struct {
		db     map[string][]mark
		lock   sync.RWMutex
		amount int32
	}{
		db:     map[string][]mark{},
		amount: 0,
	}
)

// AddMark map the given link to the user
func AddMark(user, link, description string) error {
	innerDB.lock.Lock()
	defer innerDB.lock.Unlock()

	basicMark := mark{
		id:          0,
		link:        link,
		description: description,
	}

	if linkList, ok := innerDB.db[user]; ok {
		if length := len(linkList); length > 0 {
			basicMark.id = linkList[length-1].id + 1
			linkList = append(linkList, basicMark)
		} else {
			linkList = append(linkList, basicMark)
		}
	}

	innerDB.db[user] = []mark{basicMark}

	return nil
}

// GetMark search links from given user name.
// Return links if user exist, return error if user not found.
func GetMark(user string) ([]mark, error) {
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
		if mark.id == target {
			innerDB.db[user] = deleteElement(marks, i)
			return nil
		}
	}

	return errors.New("target not found")
}

func deleteElement(markList []mark, i int) []mark {
	if len(markList)-1 == i {
		markList = markList[:i]
	} else {
		part := markList[i+1:]
		markList = append(markList[:i], part...)
	}

	return markList
}
