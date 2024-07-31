import {
	"recyco/config"
}

func main () {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    })

	http.HandleFunc("/index", index)

    fmt.Println("starting web server at http://localhost:8080/")
    http.ListenAndServe(":8080", nil)
}