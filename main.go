package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Perfil struct {
	Cedula       string
	Nombre       string
	Apellido     string
	Email        string
	Especialidad string
}

func dbConnection() (connection *sql.DB) {
	driver := "mysql"
	user := "root"
	password := ""
	db_name := "teodiofermin"

	connection, err := sql.Open(driver, user+":"+password+"@tcp(127.0.0.1)/"+db_name)
	if err != nil {
		panic(err.Error())
	}
	return connection
}

var templates = template.Must(template.ParseGlob("src/*"))

func crearUsuario(cedula string, nombre string, apellido string, email string, especialidad string) error {

	return nil
}

// func actualizarUsuario(id int, nombre, titulo, email, rol string) error {
// 	db := dbConnection()
// 	cursor := db.Exec("UPDATE usuarios SET nombre = ?, titulo = ?, email = ?, rol = ? WHERE id = ?", nombre, titulo, email, rol, id)

// 	return nil
// }

type Registro struct {
	ID           int
	Nombre       string
	Apellido     string
	Especialidad sql.NullString
}

func mostrarFormulario(w http.ResponseWriter, r *http.Request) {
	db := dbConnection()
	query, err := db.Query("SELECT * FROM profiles")
	if err != nil {
		panic(err.Error())
	}
	profile := Perfil{}
	arrayProfile := []Perfil{}
	for query.Next() {
		var cedula, nombre, apellido, email, esp string
		var especialidad sql.NullString
		err = query.Scan(&cedula, &nombre, &apellido, &email, &especialidad)
		if err != nil {
			panic(err.Error())
		}
		if !especialidad.Valid {
			esp = " "
		} else {
			esp = especialidad.String
		}
		profile.Nombre = nombre
		profile.Apellido = apellido
		profile.Cedula = cedula
		profile.Email = email
		profile.Especialidad = esp
		arrayProfile = append(arrayProfile, profile)
	}
	templates.ExecuteTemplate(w, "index", arrayProfile)
	templates.Execute(w, nil)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "crear", nil)
	templates.Execute(w, nil)
}

func actualizarUsuario(w http.ResponseWriter, r *http.Request) {
	cedula := r.URL.Query().Get("cedula")
	if cedula == "" {
		http.Error(w, "Error en el ID proporcionado", http.StatusBadRequest)
		return
	}
	db := dbConnection()
	query, err := db.Query("SELECT * FROM profiles WHERE Cedula=?", cedula)
	if err != nil {
		panic(err.Error())
	}
	profile := Perfil{}
	for query.Next() {
		var nombre, apellido, email, esp string
		var especialidad sql.NullString
		err = query.Scan(&cedula, &nombre, &apellido, &email, &especialidad)
		if err != nil {
			panic(err.Error())
		}
		if !especialidad.Valid {
			esp = " "
		} else {
			esp = especialidad.String
		}

		profile.Nombre = nombre
		profile.Apellido = apellido
		profile.Cedula = cedula
		profile.Email = email
		profile.Especialidad = esp
	}

	templates.ExecuteTemplate(w, "actualizar", profile)
}
func update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		cedula := r.FormValue("cedula")
		apellido := r.FormValue("apellido")
		email := r.FormValue("email")
		especialidad := r.FormValue("especialidad")
		fmt.Print(nombre, cedula, apellido, email, especialidad)
		db := dbConnection()
		cursor, err := db.Prepare("UPDATE profiles SET Nombre=?,Apellido=?,Email=?,especialidad=? WHERE Cedula=?")
		if err != nil {
			panic(err.Error())
		}
		cursor.Exec(nombre, apellido, email, especialidad, cedula)

		http.Redirect(w, r, "/", 301)
	}
	fmt.Print("a")
}

func crearUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		apellido := r.FormValue("apellido")
		email := r.FormValue("email")
		cedula := r.FormValue("cedula")
		especialidad := r.FormValue("especialidad")
		db := dbConnection()
		println(db)
		cursor, err := db.Prepare("INSERT INTO profiles (Cedula, Nombre, Apellido, Email, especialidad) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
		cursor.Exec(cedula, nombre, apellido, email, especialidad)
		if err != nil {
			http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
			panic(err.Error())
		}
		http.Redirect(w, r, "/", 301)
	}

}

func eliminarUsuarioHandler(w http.ResponseWriter, r *http.Request) {
	cedula := r.URL.Query().Get("cedula")
	if cedula == "" {
		http.Error(w, "Error en el ID proporcionado", http.StatusBadRequest)
		return
	}
	db := dbConnection()

	cursor, err := db.Prepare("DELETE FROM profiles WHERE Cedula=?")
	if err != nil {
		panic(err.Error())
	}
	cursor.Exec(cedula)
	if err != nil {
		http.Error(w, "Error al eliminar el usuario", http.StatusInternalServerError)
		panic(err.Error())
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", mostrarFormulario)
	http.HandleFunc("/crear", addUser)
	http.HandleFunc("/add", crearUsuarioHandler)
	http.HandleFunc("/delete", eliminarUsuarioHandler)
	http.HandleFunc("/actualizar", actualizarUsuario)
	http.HandleFunc("/update", update)

	log.Println("Servidor escuchando en http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
