package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

func main() {

	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	fmt.Println("servidor corriendo...")
	http.ListenAndServe(":8080", nil)

}

// Conexion con la BD
func conexionDB() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contraseña := ""
	Nombre := "golang"

	conexion, err := sql.Open(Driver, Usuario+":"+Contraseña+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

// ************************************************************************************
// ************************************* CRUD *****************************************
// ************************************************************************************

// Vista - Listar empleados
func Inicio(w http.ResponseWriter, r *http.Request) {

	conexionEstablecida := conexionDB()
	registros, err := conexionEstablecida.Query("SELECT * FROM empleados")

	if err != nil {
		panic(err.Error())
	}

	empleado := Empleado{}
	arregloEmpleado := []Empleado{}

	for registros.Next() {
		var id int
		var nombre, correo string
		err = registros.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

		arregloEmpleado = append(arregloEmpleado, empleado)
	}

	plantillas.ExecuteTemplate(w, "inicio", arregloEmpleado)
}

// Vista - Crear empleados
func Crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)
}

// Insertar registros a la BD
func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionDB()
		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO empleados(nombre, correo) VALUES (?,?)")
		if err != nil {
			panic(err.Error())
		}
		insertarRegistros.Exec(nombre, correo)
		http.Redirect(w, r, "/", 301)
	}
}

// Borrar empleados
func Borrar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")

	conexionEstablecida := conexionDB()
	borrarRegistro, err := conexionEstablecida.Prepare("DELETE FROM empleados WHERE id=?")

	if err != nil {
		panic(err.Error())
	}

	borrarRegistro.Exec(idEmpleado)
	http.Redirect(w, r, "/", 301)
}

// Vista - Editar empleados
func Editar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")

	conexionEstablecida := conexionDB()

	registro, err := conexionEstablecida.Query("SELECT * FROM empleados WHERE id=?", idEmpleado)
	if err != nil {
		panic(err.Error())
	}

	empleado := Empleado{}

	for registro.Next() {
		var id int
		var nombre, correo string
		err = registro.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
	}
	// fmt.Println(empleado)
	plantillas.ExecuteTemplate(w, "editar", empleado)
}

func Actualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionDB()

		actualizarRegistro, err := conexionEstablecida.Prepare("UPDATE empleados SET nombre=?, correo=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}

		actualizarRegistro.Exec(nombre, correo, id)
		http.Redirect(w, r, "/", 301)

	}

}
