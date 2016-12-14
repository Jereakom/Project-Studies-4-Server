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
  "github.com/jereakom/Project-Studies-4-Server/env"
  "regexp"
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
  router.DELETE("/users/:id/friends/:username", RemoveFriend)
  router.GET("/users/:id/groups", GetUserGroups)
  router.POST("/users/:id/groups", AddNewGroup)
  router.POST("/parsertest", Parsertest)
  router.POST("/users/:id/groups/:name", JoinGroup)
  router.DELETE("/users/:id/groups/:id/", LeaveGroup)

  router.GET("/posts", GetAllPosts)
  router.GET("/posts/tags/:tag", GetTaggedPosts)
  router.GET("/posts/:username", GetUsernameMentionedPosts)
  router.POST("/posts", CreateNewPost)

  router.GET("/groups", GetAllGroups)
  router.DELETE("/groups/:name", DeleteGroup)
  router.GET("/groups/:id/members", GetGroupMembers)

  router.GET("/", Index)

  // ENV VARIABLES

  envars.Setenvars()

  // END OF ENV VARIABLES

  var connstring string = "user="+os.Getenv("DBUSER")+" dbname="+os.Getenv("DB")+" host="+os.Getenv("DBHOST")+" password="+os.Getenv("DBPASS")+" sslmode=disable"

  db, err = sql.Open("postgres", connstring )
  if err != nil {
    log.Fatal(err)
  }

  log.Fatal(http.ListenAndServe(":80", router))
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

  rows, err := db.Query("SELECT id, username FROM users")

  fmt.Fprintf(w, "[")
  defer rows.Close()
    for rows.Next() {

            var id int
            var username string

            if err := rows.Scan(&id, &username); err != nil {
                    log.Println(err)
            }

            response := userResponse{
              Id: id,
              Username: username}

            responseJSON, _ := json.Marshal(response)
            if err != nil{
              log.Println(err)
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
            log.Println(err)
    }

}

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

  type registerResponse struct {
    Id int `json:"id"`
    Username string `json:"username"`
    Email string `json:"email"`
  }

  r.ParseForm()

  var id int
  var username string
  var email string
  var password string

  if len(r.Form["username"]) > 0 {
    username = r.Form["username"][0]
  }
  if len(r.Form["email"]) > 0 {
    email = r.Form["email"][0]
  }
  if len(r.Form["password"]) > 0 {
    password = r.Form["password"][0]
  }

  var insert string = "INSERT INTO users (username, password, email) VALUES ('"+username+"', '"+password+"', '"+email+"') RETURNING id, username, email"
  fmt.Println(insert)
  err := db.QueryRow(insert).Scan(&id, &username, &email)

  if err != nil{
    log.Println(err)
  }

  response := registerResponse{
    Id: id,
    Username: username,
    Email: email}

  responseJSON, _ := json.Marshal(response)
  if err != nil{
    log.Println(err)
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
  var username string
  var email string
  var password string

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
      fmt.Fprintf(w, "%s\n", err)
      return
  }

  response := loginResponse{
    Id: id,
    Username: username,
    Email: email}

  responseJSON, _ := json.Marshal(response)
  if err != nil{
    log.Println(err)
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
    log.Println(err)
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
    log.Println(err)
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

  if len(r.Form["newpassword"]) > 0 {
    updateInfo += "password='"+r.Form["newpassword"][0]+"' "
  }

  if len(r.Form["newemail"]) > 0 {
    updateInfo += "email='"+r.Form["newemail"][0]+"' "
  }

  updateInfo += "WHERE username ='"+username+"' AND password='"+password+"' RETURNING id, username, email"

  fmt.Println(updateInfo)
  err = db.QueryRow(updateInfo).Scan(&id, &username, &email)
  if err != nil{
    log.Print(err)
    fmt.Fprintf(w, "%s\n", err)
    return
  }

  var userID, interr = strconv.Atoi(id)
  if interr !=nil{
    log.Println(err)
  }

  response := editUser{
    Id: userID,
    Username: username,
    Email: email}

  responseJSON, _ := json.Marshal(response)
  if err != nil{
    log.Println(err)
  }
  fmt.Fprintf(w, "%s\n", responseJSON)// */

}

func RemoveUser(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {

  // log.Println(r.Header)
  log.Println(r.Header["Authorization"][0])

  var id  = params.ByName("id")

  var username string
  var password string
  var email string

  type removeUser struct {
    Id string `json:"id"`
    Username string `json:"username"`
    Email string `json:"email"`
  }

  if len(r.Header["Authorization"]) > 0 {
    password = r.Header["Authorization"][0]
  } else {
    log.Println("No password detected")
    return
  }

  var deleteUser string = "DELETE FROM users WHERE id ='"+id+"' AND password='"+password+"' RETURNING id, username, email "
  fmt.Println(deleteUser)
  err := db.QueryRow(deleteUser).Scan(&id, &username, &email)

  if err != nil{
    log.Print(err)
    fmt.Fprintf(w, "Invalid password or username\n")
    return
  }

  response := removeUser{
    Id: id,
    Username: username,
    Email: email}

  responseJSON, _ := json.Marshal(response)
  if err != nil{
    log.Println(err)
  }
  fmt.Fprintf(w, "%s\n", responseJSON)
}

func GetUserPosts(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {

  var counter int = 0

  var id string = params.ByName("id")

  type postResponse struct {
    Id int `json:"id"`
    Username string `json:"username"`
    Caption string `json:"caption"`
    Picture string `json:"picture"`
    Latitude string `json:"latitude"`
    Longitude string `json:"longitude"`
    Postedat string `json:"postedat"`
  }

  rows, err := db.Query("SELECT posts.* FROM posts INNER JOIN users ON posts.username=users.username WHERE users.id = "+id+"")

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
                    log.Println(err)
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
              log.Println(err)
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
            log.Println(err)
    }
}

func GetUserFriends(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {
  var counter int = 0

  var id string = params.ByName("id")

  type friends struct {
    Username string `json:"username"`
  }

  rows, err := db.Query("SELECT friends.second_username FROM users INNER JOIN friends ON users.username=friends.first_username WHERE users.id = "+id+"")

  fmt.Fprintf(w, "[")
  defer rows.Close()
    for rows.Next() {

            var username string

            if err := rows.Scan(&username); err != nil {
                    log.Println(err)
            }

            response := friends{
              Username: username}

            responseJSON, _ := json.Marshal(response)
            if err != nil{
              log.Println(err)
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
            log.Println(err)
    }
}

func AddFriend(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {

  type friends struct {
    Username string `json:"username"`
  }

  r.ParseForm()

  var first_username string
  var second_username string

  if len(r.Form["funame"]) > 0{
    first_username = r.Form["funame"][0]
  }
  if len(r.Form["suname"]) > 0{
    second_username = r.Form["suname"][0]
  }


  err := db.QueryRow("INSERT INTO friends VALUES('"+first_username+"', '"+second_username+"') RETURNING second_username").Scan(&second_username)

  if err != nil {
    log.Println(err)
    fmt.Fprintf(w, "%s\n", "Failed to add a friend")
    return
  }

  response := friends{
    Username: second_username}

  responseJSON, _ := json.Marshal(response)
  if err != nil{
    log.Println(err)
  }

  fmt.Fprintf(w, "%s\n", responseJSON)

}

func RemoveFriend(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {

  type friend struct {
    Username string `json:"username"`
  }

  r.ParseForm()

  var friendToDelete string = params.ByName("username")
  var thisAccount string

  var id = params.ByName("id")

  var queryString = "SELECT username from users where id='"+id+"'"

  err := db.QueryRow(queryString).Scan(&thisAccount)

  log.Println(thisAccount)

  if err != nil {
    log.Println(err)
    fmt.Fprintf(w, "%s\n", "Could not find user to remove friend")
    return
  }

  err = db.QueryRow("DELETE FROM friends WHERE first_username='"+thisAccount+"' AND second_username='"+friendToDelete+"' RETURNING second_username").Scan(&friendToDelete)

  if err != nil {
    log.Println(err)
    fmt.Fprintf(w, "%s\n", "Could not remove friend")
    return
  }

  response := friend{
    Username: friendToDelete}

  responseJSON, _ := json.Marshal(response)
  if err != nil{
    log.Println(err)
  }

  fmt.Fprintf(w, "%s\n", responseJSON)
}

func GetUserGroups(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {
  var counter int = 0

  var id string = params.ByName("id")

  type groups struct {
    Group string `json:"group"`
  }

  rows, err := db.Query("SELECT name FROM groups_users WHERE id = "+id+"")

  fmt.Fprintf(w, "[")
  defer rows.Close()
    for rows.Next() {

            var group string

            if err := rows.Scan(&group); err != nil {
                    log.Println(err)
            }

            response := groups{
              Group: group}

            responseJSON, _ := json.Marshal(response)
            if err != nil{
              log.Println(err)
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
            log.Println(err)
    }
}

func AddNewGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {

  r.ParseForm()

  var id = params.ByName("id")

  var name string
  var group string

  type groups struct {
    Group string `json:"group"`
  }

  if len(r.Form["name"]) > 0{
    name = r.Form["name"][0]
  }


  err := db.QueryRow("INSERT INTO groups (name) VALUES ('"+name+"') RETURNING name").Scan(&group)
    if err != nil {
      log.Println(err)
    }

  err = db.QueryRow("INSERT INTO groups_users (id, name) VALUES ('"+id+"', '"+name+"') RETURNING name").Scan(&group)
  if err != nil {
    log.Println(err)
  }

    response := groups{
      Group: group}

    responseJSON, _ := json.Marshal(response)
    if err != nil{
      log.Println(err)
    }

    fmt.Fprintf(w, "%s\n", responseJSON)
}

func JoinGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {

  var id string = params.ByName("id")
  var name string = params.ByName("name")

  type joinResponse struct {
    Group string `json:"group"`
  }

  err := db.QueryRow("INSERT INTO groups_users (id, name) VALUES ('"+id+"','"+name+"') RETURNING name").Scan(&name)
  if err != nil {
    log.Println(err)
  }

    response := joinResponse{
      Group: name}

    responseJSON, _ := json.Marshal(response)
    if err != nil{
      log.Println(err)
    }

    fmt.Fprintf(w, "%s\n", responseJSON)
}

func LeaveGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {

  var id string = params.ByName("id")
  var name string

  type joinResponse struct {
    Group string `json:"group"`
  }

  if len(r.Form["name"]) > 0{
    name = r.Form["name"][0]
  } else {
    log.Println("cannot leave group")
    return
  }

  err := db.QueryRow("DELETE FROM groups_users WHERE id='"+id+"' AND name='"+name+"' RETURNING *").Scan(&id, &name)
  if err != nil {
    log.Println(err)
  }

    response := joinResponse{
      Group: name}

    responseJSON, _ := json.Marshal(response)
    if err != nil{
      log.Println(err)
    }

    fmt.Fprintf(w, "%s\n", responseJSON)
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
                    log.Println(err)
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
              log.Println(err)
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
            log.Println(err)
    }

}

func GetTaggedPosts(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {
  var counter int = 0
  fmt.Fprintf(w, "[")

  var tag = params.ByName("tag")


  type postResponse struct {
    Id int `json:"id"`
    Username string `json:"username"`
    Caption string `json:"caption"`
    Picture string `json:"picture"`
    Latitude string `json:"latitude"`
    Longitude string `json:"longitude"`
    Postedat string `json:"postedat"`
  }

  rows, err := db.Query("SELECT id, username, caption, picture, latitude, longitude, postedat FROM posts INNER JOIN tags ON posts.id=tags.post_id WHERE tags.tagword='"+tag+"'")

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
                    log.Println(err)
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
              log.Println(err)
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
            log.Println(err)
    }
}

func GetUsernameMentionedPosts(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {
  var counter int = 0
  fmt.Fprintf(w, "[")

  var username_mention = params.ByName("username")


  type postResponse struct {
    Id int `json:"id"`
    Username string `json:"username"`
    Caption string `json:"caption"`
    Picture string `json:"picture"`
    Latitude string `json:"latitude"`
    Longitude string `json:"longitude"`
    Postedat string `json:"postedat"`
  }

  rows, err := db.Query("SELECT id, username, caption, picture, latitude, longitude, postedat FROM posts INNER JOIN username_mentions ON posts.id=username_mentions.post_id WHERE username_mentions.username_mention='"+username_mention+"'")

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
                    log.Println(err)
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
              log.Println(err)
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
            log.Println(err)
    }
}

func Parsertest (w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

  r.ParseForm()

  log.Println(r.Form["caption"])
  if len(r.Form["caption"]) > 0 {
    parseForTags(r.Form["caption"][0], 3)

  }
}

func parseForTags(textToParse string, postID int) {

  var id = strconv.Itoa(postID)

  re := regexp.MustCompile("@(\\w+)")
  var usernameMentions = re.FindAllStringSubmatch(textToParse, -1)

  re = regexp.MustCompile("#(\\w+)")
  var groupMentions = re.FindAllStringSubmatch(textToParse, -1)
  log.Println(usernameMentions)
  log.Println(groupMentions)
  var unamelen = len(usernameMentions)
  var glen = len(groupMentions)
//  var usernames = make([]string, unamelen)
//  var groups = make([]string, glen)

  for i := 0; i < unamelen; i++ {

    var postQuery = "INSERT INTO username_mentions (post_id, username_mention) VALUES('"+id+"', '"+usernameMentions[i][1]+"')"

    err := db.QueryRow(postQuery)

    if err != nil {
      log.Println("Could not add username mention")
    }
  }

  for i := 0; i < glen; i++ {

    var postQuery = "INSERT INTO tags (post_id, tagword) VALUES('"+id+"', '"+groupMentions[i][1]+"')"

    err := db.QueryRow(postQuery)

    if err != nil {
      log.Println("Could not add group tag")
    }
  }

}

func CreateNewPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {

  r.ParseForm()

  type postResponse struct {
    Id int `json:"id"`
    Username string `json:"username"`
    Caption string `json:"caption"`
    Picture string `json:"picture"`
    Latitude string `json:"latitude"`
    Longitude string `json:"longitude"`
    Postedat string `json:"postedat"`
  }

  var id int
  var username string
  var caption string
  var picture string
  var latitude string
  var longitude string
  var postedat string

  if len(r.Form["username"]) > 0{
    username = r.Form["username"][0]
  }
  if len(r.Form["caption"]) > 0{
    caption = r.Form["caption"][0]
  }
  if len(r.Form["picture"]) > 0{
    picture = r.Form["picture"][0]
  }
  if len(r.Form["latitude"]) > 0{
    latitude = r.Form["latitude"][0]
  }
  if len(r.Form["longitude"]) > 0{
    longitude = r.Form["longitude"][0]
  }

  var postQuery = "INSERT INTO posts (username, caption, picture, latitude, longitude) VALUES('"+username+"', '"+caption+"', '"+picture+"', '"+latitude+"', '"+longitude+"') RETURNING *"

  err := db.QueryRow(postQuery).Scan(&id, &username, &caption, &picture, &latitude, &longitude, &postedat );

  if err != nil {
    fmt.Fprintf(w, "%s\n", "Could not post post")
    return
  }
  parseForTags(caption, id)

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
    log.Println(err)
  }

  fmt.Fprintf(w, "%s\n", responseJSON)

}

func GetAllGroups(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {

  var counter int = 0
  type groupslist struct {
    Id int `json:"id"`
    Name string `json:"name"`
  }

  rows, err := db.Query("SELECT * FROM groups")

  fmt.Fprintf(w, "[")
  defer rows.Close()
    for rows.Next() {

            var id int
            var name string

            if err := rows.Scan(&id, &name); err != nil {
                    log.Println(err)
            }

            response := groupslist{
              Id: id,
              Name: name}

            responseJSON, _ := json.Marshal(response)
            if err != nil{
              log.Println(err)
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
            log.Println(err)
    }

}


func DeleteGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {

  var name string = params.ByName("name")

  type deletedGroup struct {
    Group string `json:"group"`
  }

  if len(name) == 0{
    log.Println("cannot leave group")
    return
  }

  err := db.QueryRow("DELETE FROM groups WHERE name='"+name+"' RETURNING name").Scan(&name)
  if err != nil {
    log.Println(err)
  }

    response := deletedGroup{
      Group: name}

    responseJSON, _ := json.Marshal(response)
    if err != nil{
      log.Println(err)
    }

    fmt.Fprintf(w, "%s\n", responseJSON)
}

func GetGroupMembers(w http.ResponseWriter, r *http.Request, params httprouter.Params)  {

  var counter int = 0
  type groupmembers struct {
    Groupname string `json:"groupname"`
    Name string `json:"name"`
    Id int `json:"id"`
  }

  var id = params.ByName("id")

  var name string

  err := db.QueryRow("SELECT name FROM groups where id='"+id+"'").Scan(&name)

  if err != nil {
    log.Println(err)
    fmt.Fprintf(w, "failed to list group members")
    return
  }


  rows, err := db.Query("SELECT groups_users.name, users.username, users.id FROM groups_users INNER JOIN users ON groups_users.id=users.id where groups_users.name='"+name+"'")

  fmt.Fprintf(w, "[")
  defer rows.Close()
    for rows.Next() {

            var userid int
            var groupname string

            if err := rows.Scan(&groupname, &name, &id); err != nil {
                    log.Println(err)
            }

            userid, err = strconv.Atoi(id)
            if err != nil {
              log.Println(err)
              return
            }

            response := groupmembers{
              Groupname: groupname,
              Name: name,
              Id: userid}

            responseJSON, _ := json.Marshal(response)
            if err != nil{
              log.Println(err)
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
            log.Println(err)
    }

}
