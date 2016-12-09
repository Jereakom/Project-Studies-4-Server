package main

  import (
  "fmt"
  "encoding/json"
  "net/http"
  "log"
  "database/sql"
  _ "github.com/lib/pq"
  "github.com/julienschmidt/httprouter"
  "os"
  "strconv"
)

var db *sql.DB

func main() {
  router := httprouter.New()

  var err error

  //ROUTES
  router.POST("/login", Login)

  router.GET("/users", GetAllUsers)
  router.POST("/users", Register)
  router.GET("/users/:id", GetUser)
  router.PUT("/users/:id", EditUser)
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

  // ENV VARIABLES



  // END OF ENV VARIABLES
  var connstring string = "user="+os.Getenv("DBUSER")+" dbname="+os.Getenv("DB")+" host="+os.Getenv("DBHOST")+" password="+os.Getenv("DBPASS")+" sslmode=disable"

  db, err = sql.Open("postgres", connstring )
  if err != nil {
    log.Fatal(err)
  }

  log.Fatal(http.ListenAndServe(":8080", router))
}


func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprint(w, "INDEX ROUTE!")
}

func GetAllUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {
  var counter int = 0
  type userResponse struct {
    Id int `json:"id"`
    Username string `json:"username"`
    Email string `json:"email"`
  }

  rows, err := db.Query("SELECT * FROM users")

  fmt.Fprintf(w, "[")
  defer rows.Close()
    for rows.Next() {

            var id int
            var username string
            var password string
            var email string

            if err := rows.Scan(&id, &username, &password, &email); err != nil {
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

            if counter == 0{
              fmt.Fprintf(w, "%s\n", responseJSON)
              counter ++
            } else{
              fmt.Fprintf(w, ",%s\n", responseJSON)
            }
    }
    fmt.Fprintf(w, "]")
    if err := rows.Err(); err != nil {
            log.Fatal(err)
    }

}

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

  type registerResponse struct {
    Id int `json:"id"`
    Username string `json:"username"`
    Email string `json:"email"`
  }

  r.ParseForm()
  fmt.Println("username:", r.Form["username"])
  fmt.Println("password:", r.Form["password"])
  fmt.Println("email:", r.Form["email"])

  var id int
  var username string = r.Form["username"][0]
  var email string = r.Form["email"][0]
  var password string = r.Form["password"][0]

  var insert string = "INSERT INTO users (username, password, email) VALUES ('"+username+"', '"+password+"', '"+email+"') RETURNING id, username, email"
  fmt.Println(insert)
  err := db.QueryRow(insert).Scan(&id, &username, &email)

  if err != nil{
    log.Fatal(err)
  }

  response := registerResponse{
    Id: id,
    Username: username,
    Email: email}

  responseJSON, _ := json.Marshal(response)
  if err != nil{
    log.Fatal(err)
  }
  fmt.Fprintf(w, "%s\n", responseJSON)

}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

  type loginResponse struct {
    Id int `json:"id"`
    Username string `json:"username"`
    Email string `json:"email"`
  }

  r.ParseForm()
  fmt.Println("username:", r.Form["username"])
  fmt.Println("password:", r.Form["password"])

  var id int
  var username string = r.Form["username"][0]
  var email string
  var password string = r.Form["password"][0]

  var get string = "SELECT id, username, email FROM users WHERE username ='"+username+"' AND password='"+password+"'"
  fmt.Println(get)
  err := db.QueryRow(get).Scan(&id, &username, &email)

  if err != nil{
    log.Fatal(err)
  }

  response := loginResponse{
    Id: id,
    Username: username,
    Email: email}

  responseJSON, _ := json.Marshal(response)
  if err != nil{
    log.Fatal(err)
  }
  fmt.Fprintf(w, "%s\n", responseJSON)

}

func GetUser(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {

  type User struct {
    Id int `json:"id"`
    Username string `json:"username"`
  }

  r.ParseForm()
  fmt.Println("username:", r.Form["username"])
  fmt.Println("password:", r.Form["password"])

  var id  = params.ByName("id")
  var intID, err = strconv.Atoi(id)
  if err !=nil{
    log.Fatal(err)
  }
  fmt.Println(id)

  var username string

  var get string = "SELECT id, username FROM users WHERE id="+id+""
  fmt.Println(get)
  err = db.QueryRow(get).Scan(&id, &username)
  if err != nil{
    log.Print(err)
    fmt.Fprintf(w, "%s,\n", err)
    return
  }

  response := User{
    Id: intID,
    Username: username}

  responseJSON, _ := json.Marshal(response)
  if err != nil{
    log.Fatal(err)
  }
  fmt.Fprintf(w, "%s\n", responseJSON)// */
}

func EditUser(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {

  var id  = params.ByName("id")

  var username string
  var password string
  var email string

  r.ParseForm()

  type editUser struct {
    Id int `json:"id"`
    Username string `json:"username"`
    Email string `json:"email"`
  }

  fmt.Println("username:", len(r.Form["username"]))
  fmt.Println("password:", len(r.Form["password"]))


  if len(r.Form["password"]) > 0 {
    password = r.Form["password"][0]
  }

  if len(r.Form["username"]) > 0 {
    username = r.Form["username"][0]
  }

  var get string = "SELECT id, username, email FROM users WHERE username ='"+username+"' AND password='"+password+"'"
  fmt.Println(get)
  err := db.QueryRow(get).Scan(&id, &username, &email)

  if err != nil{
    log.Print(err)
    fmt.Fprintf(w, "Invalid password\n")
    return
  }

  var updateInfo string = "UPDATE users SET "

  r.ParseForm()
  fmt.Println("email:", r.Form["email"])
  fmt.Println("password:", r.Form["password"])
  fmt.Println("newpassword:", len(r.Form["newpassword"]))

  if len(r.Form["newpassword"]) > 0 {
    updateInfo += "password='"+r.Form["newpassword"][0]+"' "
  }

  if len(r.Form["newemail"]) > 0 {
    updateInfo += "email='"+r.Form["newemail"][0]+"' "
  }

  updateInfo += "WHERE username ='"+username+"' AND password='"+password+"' RETURNING id, username, email"
//  newpassword := r.Form["newpassword"][0]
//  newemail := r.Form["newemail"][0]

  fmt.Println(updateInfo)
  err = db.QueryRow(updateInfo).Scan(&id, &username, &email)
  if err != nil{
    log.Print(err)
    fmt.Fprintf(w, "%s\n", err)
    return
  }

  var userID, interr = strconv.Atoi(id)
  if interr !=nil{
    log.Fatal(err)
  }

  response := editUser{
    Id: userID,
    Username: username,
    Email: email}

  responseJSON, _ := json.Marshal(response)
  if err != nil{
    log.Fatal(err)
  }
  fmt.Fprintf(w, "%s,\n", responseJSON)// */

}

func RemoveUser(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {

  /* TODO: STUDY SENSIBLE IMPLEMENTATION

  var id  = params.ByName("id")

  var username string
  var password string
  var email string

  type removeUser struct {
    Id int `json:"id"`
    Username string `json:"username"`
    Email string `json:"email"`
  }

  r.ParseForm()

  fmt.Println("username:", len(r.Form["username"]))
  fmt.Println("password:", len(r.Form["password"]))


  if len(r.Form["password"]) > 0 {
    password = r.Form["password"][0]
  }

  if len(r.Form["username"]) > 0 {
    username = r.Form["username"][0]
  }

  var deleteUser string = "DELETE FROM users WHERE username ='"+username+"' AND password='"+password+"' RETURNING id, username, email "
  fmt.Println(deleteUser)
  err := db.QueryRow(deleteUser).Scan(&id, &username, &email)

  if err != nil{
    log.Print(err)
    fmt.Fprintf(w, "Invalid password\n")
    return
  }

  var removedID, remerr = strconv.Atoi(id)
  if remerr !=nil{
    log.Fatal(err)
  }

  response := removeUser{
      Id: removedID,
    Username: username,
    Email: email}

  responseJSON, _ := json.Marshal(response)
  if err != nil{
    log.Fatal(err)
  }
  fmt.Fprintf(w, "%s,\n", responseJSON) */
}

func GetUserPosts(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {

  var counter int = 0


  type postResponse struct {
    Id int `json:"id"`
    Username string `json:"username"`
    Caption string `json:"caption"`
    Picture string `json:"picture"`
    Latitude string `json:"latitude"`
    Longitude string `json:"longitude"`
    Postedat string `json:"postedat"`
  }

  rows, err := db.Query("SELECT * FROM posts WHERE username =")

  fmt.Fprintf(w, "[")
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

            response := postResponse{
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
            if counter == 0{
                fmt.Fprintf(w, "%s\n", responseJSON)
                counter ++
              } else{
                fmt.Fprintf(w, ",%s\n", responseJSON)
              }
      }
    fmt.Fprintf(w, "]")
    if err := rows.Err(); err != nil {
            log.Fatal(err)
    }
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

  var counter int = 0
  fmt.Fprintf(w, "[")


  type postResponse struct {
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

            response := postResponse{
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
            if counter == 0{
              fmt.Fprintf(w, "%s\n", responseJSON)
              counter ++
              } else{
                fmt.Fprintf(w, ",%s\n", responseJSON)
              }
            }
    fmt.Fprintf(w, "]")
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
