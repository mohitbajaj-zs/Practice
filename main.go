package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Employee struct {
	EmpID   string     `json:"EmpId"`
	Name    string     `json:"Name"`
	Dept    Department `json:"Dept"`
	PhoneNo string     `json:"PhoneNo"`
}
type Department struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}

var Db *sql.DB

func GetEmployeeData(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var employees []Employee
	result, err := Db.Query("select e.emp_id, e.name,e.id,d.name,e.phone_number from employee e inner join department d on e.id=d.id;")
	if err != nil {
		log.Println(err)
	}
	defer result.Close()
	for result.Next() {
		var employee Employee
		err := result.Scan(&employee.EmpID, &employee.Name, &employee.Dept.Id, &employee.Dept.Name, &employee.PhoneNo)
		if err != nil {
			log.Println(err.Error())
		}
		employees = append(employees, employee)
	}

	json.NewEncoder(w).Encode(employees)
	//respBody, _ := json.Marshal(employees)
	//w.Write(respBody)
}
func PostEmployeeData(w http.ResponseWriter, r *http.Request) {
	var emp Employee
	w.Header().Set("Content-Type", "application/json")
	req, _ := ioutil.ReadAll(r.Body)
	_ = json.Unmarshal(req, &emp)
	_, err := Db.Exec("insert into employee (emp_id,name,id,phone_number) values (?,?,?,?)", emp.EmpID, emp.Name, emp.Dept.Id, emp.PhoneNo)
	if err != nil {
		log.Println(err)
	}
	_, _ = io.WriteString(w, "Data added successfully")
}

func main() {
	var err error
	Db, err = sql.Open("mysql",
		"root:1234@tcp(127.0.0.1:3306)/test_mohit")
	if err != nil {
		log.Println(err)
		return
	}

	defer Db.Close()

	err = Db.Ping()
	if err != nil {
		log.Println(err)
		return
	}

	http.HandleFunc("/employees", GetEmployeeData)
	http.HandleFunc("/post", PostEmployeeData)

	fmt.Println(("server at port 8080"))
	log.Fatal(http.ListenAndServe(":8080", nil))

}
