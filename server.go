package main

  import (
  "fmt"
  "encoding/json"
  "net/http"
  "log"
  "database/sql"
  _ "github.com/lib/pq"
  "github.com/julienschmidt/httprouter"
)

var db *sql.DB

func main() {
  router := httprouter.New()

  var err error

  router.GET("/users", GetAllUsers)
  router.POST("/users", Register)
  router.GET("/users/:id", GetUser)
  router.PUT("/users:id", EditUser)
  router.DELETE("/users/:id", RemoveUser)
  router.GET("/users/:id/posts", GetUserPosts)
  router.GET("/users/:id/friends", GetUserFriends)
  router.POST("/users/:id/friends", AddFriend)
  router.DELETE("/users/:id/friends/:id", RemoveFriend)
  router.GET("/users/:id/groups", GetUserGroups)
  router.POST("/users/:id/groups", AddNewGroup)
  router.DELETE("/users/:id/groups/:id/", LeaveGroup)
  router.GET("/users/:id/chat", GetUserChatMessages)

  router.GET("/posts", GetAllPosts)
  router.POST("/posts", CreateNewPost)
  router.DELETE("/posts/:id", RemovePost)
  router.GET("/posts/:id", GetPost)
  router.PUT("/posts/:id", EditPost)

  router.GET("/groups", ListAllGroups)
  router.POST("/groups", CreateGroup)
  router.GET("/groups/:id", GetGroup)
  router.DELETE("/groups/:id", DeleteGroup)
  router.GET("/groups/:id/members", GetGroupMembers)
  router.POST("/groups/:id/members", AddMember)
  router.DELETE("/groups/:id/members/:id", RemoveGroupMember)

  router.GET("/chat/", GetChatHistory)
  router.POST("/chat", AddNewMessage)
  router.DELETE("/chat/:id", RemoveMessage)

  router.GET("/", Index)

  db, err = sql.Open("postgres", "user=postgres dbname=thegrid host=52.169.87.203 password=7eeGrpbyLLPgjpyu7rpZ sslmode=disable")
  if err != nil {
    log.Fatal(err)
  }

  log.Fatal(http.ListenAndServe(":8080", router))
}


func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprint(w, "INDEX ROUTE!")
  db.Ping()
}

func GetAllUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {
  type userResponse struct {
    Id int `json:"id"`
    Username string `json:"username"`
    Email string `json:"email"`
  }

  rows, err := db.Query("SELECT * FROM users")

  defer rows.Close()
    for rows.Next() {

            var id int
            var username string
            var email string
            var password string

            if err := rows.Scan(&id, &username, &email, &password); err != nil {
                    log.Fatal(err)
            }

            response := userResponse{
              Id: id,
              Username: username,
              Email: email}

            responseJSON, _ := json.Marshal(response)
            if err != nil{
              log.Fatal(err)
            }
            fmt.Fprintf(w, "%s\n", responseJSON)
    }
    if err := rows.Err(); err != nil {
            log.Fatal(err)
    }

}

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func GetUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func EditUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func RemoveUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func GetUserPosts(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func GetUserFriends(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func AddFriend(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func RemoveFriend(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func GetUserGroups(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func AddNewGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func LeaveGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func GetUserChatMessages(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func GetAllPosts(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

  type Response1 struct {
    Id int `json:"id"`
    Username string `json:"username"`
    Caption string `json:"caption"`
    Picture string `json:"picture"`
    Latitude string `json:"latitude"`
    Longitude string `json:"longitude"`
    Postedat string `json:"postedat"`
  }

  rows, err := db.Query("SELECT * FROM posts")

  defer rows.Close()
    for rows.Next() {

            var id int
            var username string
            var caption string
            var picture string
            var latitude string
            var longitude string
            var postedat string
            if err := rows.Scan(&id, &username, &caption, &picture, &latitude, &longitude, &postedat ); err != nil {
                    log.Fatal(err)
            }

            response := Response1{
              Id: id,
              Username: username,
              Caption: caption,
              Picture: picture,
              Latitude: latitude,
              Longitude: longitude,
              Postedat: postedat}

            responseJSON, _ := json.Marshal(response)
            if err != nil{
              log.Fatal(err)
            }
            fmt.Fprintf(w, "%s\n", responseJSON)
    }
    if err := rows.Err(); err != nil {
            log.Fatal(err)
    }

}

func CreateNewPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func RemovePost(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func GetPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func EditPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func ListAllGroups(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func CreateGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func GetGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func DeleteGroup(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func GetGroupMembers(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func AddMember(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func RemoveGroupMember(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func GetChatHistory(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}


func AddNewMessage(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}

func RemoveMessage(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

}
