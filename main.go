package main

import (
	"encoding/json"
	"net/http"
	"sync"
)
var(
	votes=map[string]int{"Alice":0,"Bob":0}
	voteMux sync.Mutex
)
func main(){
	http.HandleFunc("/",homeHandler)
	http.HandleFunc("/vote",voteHandler)
	http.HandleFunc("/result",resultHandler)
	http.ListenAndServe(":8080",nil )
}

func homeHandler(w http.ResponseWriter,r *http.Request){
	w.Write([]byte(`
	<html>
    <body>
        <h1>Vote for your favorite Candidate</h1>
        <form action="/vote" method="POST">
            <input type="radio" name="candidate" value="BJP"> BJP<br>    
            <input type="radio" name="candidate" value="Congress"> Congress<br>    
            <input type="submit" value="vote">
        </form>
    </body>
</html>
	
	
	
	`))
}

func voteHandler(w http.ResponseWriter, r *http.Request){
if r.Method!=http.MethodPost{
	http.Error(w,"Invalid request method",http.StatusMethodNotAllowed)
	return
}
candidate:=r.FormValue("candidate")
if candidate!="BJP" && candidate!="Congress"{
http.Error(w, "Invalid candidate",http.StatusBadRequest)
return
}
voteMux.Lock()
votes[candidate]++
voteMux.Unlock()
http.Redirect(w,r,"/result",http.StatusSeeOther)
}
func resultHandler(w http.ResponseWriter,r *http.Request){
voteMux.Lock()
results:=map[string]int{"BJP":votes["BJP"],"Congress":votes["Congress"]}
voteMux.Unlock()
w.Header().Set("Content-Type","application/json")
json.NewEncoder(w).Encode(results)
}