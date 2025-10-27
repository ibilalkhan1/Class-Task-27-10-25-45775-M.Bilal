package main
import (
  "encoding/json"; "errors"; "io"; "net/http"; "os"; "sync"
)
type Book struct{ ID,Title,Author string; Price float64; Available bool }
var mu sync.Mutex
var dataFile="/data/catalog.json"

func readAll()([]Book,error){
  b,err:=os.ReadFile(dataFile)
  if errors.Is(err,os.ErrNotExist){ return []Book{},nil }
  if err!=nil { return nil,err }
  if len(b)==0 { return []Book{},nil }
  var bs []Book; if err:=json.Unmarshal(b,&bs); err!=nil { return nil,err }
  return bs,nil
}
func writeAll(bs []Book) error {
  bb,err:=json.MarshalIndent(bs,"","  "); if err!=nil { return err }
  return os.WriteFile(dataFile,bb,0644)
}
func getBooks(w http.ResponseWriter,_ *http.Request){
  mu.Lock(); defer mu.Unlock()
  bs,_:=readAll()
  w.Header().Set("Content-Type","application/json")
  json.NewEncoder(w).Encode(bs)
}
func createBook(w http.ResponseWriter,r *http.Request){
  mu.Lock(); defer mu.Unlock()
  body,_:=io.ReadAll(r.Body); var b Book
  if err:=json.Unmarshal(body,&b); err!=nil { w.WriteHeader(400); return }
  bs,_:=readAll(); bs=append(bs,b); writeAll(bs)
  w.WriteHeader(201); json.NewEncoder(w).Encode(b)
}
func main(){
  os.MkdirAll("/data",0755)
  http.HandleFunc("/books", func(w http.ResponseWriter,r *http.Request){
    if r.Method==http.MethodGet { getBooks(w,r); return }
    if r.Method==http.MethodPost{ createBook(w,r); return }
    w.WriteHeader(405)
  })
  http.HandleFunc("/health", func(w http.ResponseWriter,_ *http.Request){
    w.Header().Set("Content-Type","application/json"); w.Write([]byte(`{"ok":true}`))
  })
  http.ListenAndServe(":4000",nil)
}
