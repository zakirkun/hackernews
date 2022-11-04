package links

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/zakirkun/hackernews/database"
	"github.com/zakirkun/hackernews/graph/model"
	"github.com/zakirkun/hackernews/internal/pkg/redis"
	"github.com/zakirkun/hackernews/internal/users"
)

type Link struct {
	ID      string
	Title   string
	Address string
	User    *users.User
}

//#2
func (link Link) Save() int64 {
	//#3
	stmt, err := database.Db.Prepare("INSERT INTO Links(Title,Address, UserID) VALUES(?,?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	//#4
	res, err := stmt.Exec(link.Title, link.Address, link.User.ID)
	if err != nil {
		log.Fatal(err)
	}
	//#5
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}

	toJson, _ := json.Marshal(link)
	cache := redis.New()
	_ = cache.Set(strconv.FormatInt(id, 10), string(toJson))

	log.Print("Row inserted!")
	return id
}

func GetByID(Id string) (*model.Link, error) {
	cache := redis.New()

	val, err := cache.Get(Id)
	if err != nil {
		log.Fatal(err)
	}

	var result model.Link
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		log.Fatal(err)
		return &model.Link{}, err
	}

	return &result, nil
}

func GetAll() []Link {
	stmt, err := database.Db.Prepare("select L.id, L.title, L.address, L.UserID, U.Username from Links L inner join Users U on L.UserID = U.ID") // changed
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var links []Link
	var username string
	var id string
	for rows.Next() {
		var link Link
		err := rows.Scan(&link.ID, &link.Title, &link.Address, &id, &username) // changed
		if err != nil {
			log.Fatal(err)
		}
		link.User = &users.User{
			ID:       id,
			Username: username,
		} // changed
		links = append(links, link)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return links
}
