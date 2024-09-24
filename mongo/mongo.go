package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id          string
	Name        string
	Habits      string
	CreatedTime string
}

var tpl = <html>
<head>
<title></title>
</head>
<body>
<form action="/info" method="post">
	username:<input type="text" name="username">
	habits:<input type="password" name="habits">
	<input type="submit" value="login">
</form>
</body>
</html>

func connect(cName string) (*mgo.Session, *mgo.Collection) {
	session, err := mgo.Dial("mongodb://x.x.x.x:27017/") //Mongodb's connection
	checkErr(err)
	session.SetMode(mgo.Monotonic, true)
	//return a instantiated collect
	return session, session.DB("test").C(cName)
}

func queryByName(name string) []User {
	var user []User
	s, c := connect("user")
	defer s.Close()
	err := c.Find(bson.M{"name": name}).All(&user)
	checkErr(err)
	return user
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func submitForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	var t *template.Template
	t = template.New("Products")
	t, _ = t.Parse(tpl)
	log.Println(t.Execute(w, nil))
}

func store(user User) error {
	s, c := connect("user")
	defer s.Close()
	user.Id = bson.NewObjectId().Hex()
	return c.Insert(&user)
}

func userInfo(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	if r.Method == "POST" {
		user1 := User{Name: r.Form.Get("username"), Habits: r.Form.Get("habits")}
		store(user1)
		fmt.Fprintf(w, " %v", queryByName("aoho"))
	}
}

func main() {
	http.HandleFunc("/form", submitForm)
	http.HandleFunc("/info", userInfo)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
