package core

import "github.com/DisgoOrg/snowflake"

type (
	UserFindFunc func(user *User) bool

	UserCache interface {
		Get(userID snowflake.Snowflake) *User
		GetCopy(userID snowflake.Snowflake) *User
		Set(user *User) *User
		Remove(userID snowflake.Snowflake)

		Cache() map[snowflake.Snowflake]*User
		All() []*User

		FindFirst(userFindFunc UserFindFunc) *User
		FindAll(userFindFunc UserFindFunc) []*User
	}

	userCacheImpl struct {
		users      map[snowflake.Snowflake]*User
		cacheFlags CacheFlags
	}
)

func NewUserCache(cacheFlags CacheFlags) UserCache {
	return &userCacheImpl{users: map[snowflake.Snowflake]*User{}, cacheFlags: cacheFlags}
}

func (c *userCacheImpl) Get(userID snowflake.Snowflake) *User {
	return c.users[userID]
}

func (c *userCacheImpl) GetCopy(userID snowflake.Snowflake) *User {
	if user := c.Get(userID); user != nil {
		us := *user
		return &us
	}
	return nil
}

func (c *userCacheImpl) Set(user *User) *User {
	// check if we want to cache user
	usr, ok := c.users[user.ID]
	if ok {
		*usr = *user
		return usr
	}
	c.users[user.ID] = user
	return user
}

func (c *userCacheImpl) Remove(id snowflake.Snowflake) {
	delete(c.users, id)
}

func (c *userCacheImpl) Cache() map[snowflake.Snowflake]*User {
	return c.users
}

func (c *userCacheImpl) All() []*User {
	users := make([]*User, len(c.users))
	i := 0
	for _, user := range c.users {
		users[i] = user
		i++
	}
	return users
}

func (c *userCacheImpl) FindFirst(userFindFunc UserFindFunc) *User {
	for _, usr := range c.users {
		if userFindFunc(usr) {
			return usr
		}
	}
	return nil
}

func (c *userCacheImpl) FindAll(userFindFunc UserFindFunc) []*User {
	var users []*User
	for _, usr := range c.users {
		if userFindFunc(usr) {
			users = append(users, usr)
		}
	}
	return users
}
