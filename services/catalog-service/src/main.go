package main

import (
  "encoding/json"
  "errors"
  "io"
  "net/http"
  "os"
  "strings"
  "sync"
)

type Book struct {
  ID        string  `json:"ID"`
  Title     string  `json:"Title"`
  Author    string  `json:"Author"`
  Price     float64 `json:"Price"`
  Available bool    `json:"Available"`
}

var (
  mu       sync.Mutex
  dataFile = func() string {
    if v := os.Getenv("CATALOG_FILE"); v != "" {
      return v
    }
    return "/data/catalog.json"
  }()
)

func readAll() ([]Book, error) {
  b, err := os.ReadFile(dataFile)
  if errors.Is(err, os.ErrNotExist) {
    return []Book{}, nil
  }
  if err != nil {
    return nil, err
  }
  if len(b) == 0 {
    return []Book{}, nil
  }
  var bs []Book
  if err := json.Unmarshal(b, &bs); err != nil {
    return nil, err
  }
  return bs, nil
}

func writeAll(bs []Book) error {
  bb, err := json.MarshalIndent(bs, "", "  ")
  if err != nil {
    return err
  }
  return os.WriteFile(dataFile, bb, 0644)
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
  // only exact /books here
  if r.URL.Path != "/books" {
    http.NotFound(w, r)
    return
  }
  switch r.Method {
  case http.MethodGet:
    mu.Lock()
    bs, _ := readAll()
    mu.Unlock()
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(bs)
  case http.MethodPost:
    body, _ := io.ReadAll(r.Body)
    var b Book
    if err := json.Unmarshal(body, &b); err != nil || strings.TrimSpace(b.ID) == "" {
      w.WriteHeader(http.StatusBadRequest)
      return
    }
    mu.Lock()
    bs, _ := readAll()
    for _, x := range bs {
      if x.ID == b.ID {
        mu.Unlock()
        w.WriteHeader(http.StatusConflict)
        return
      }
    }
    bs = append(bs, b)
    _ = writeAll(bs)
    mu.Unlock()
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(b)
  default:
    w.WriteHeader(http.StatusMethodNotAllowed)
  }
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
  // normalize and extract id from /books/{id}
  path := strings.TrimSuffix(r.URL.Path, "/")
  if !strings.HasPrefix(path, "/books/") {
    http.NotFound(w, r)
    return
  }
  id := strings.TrimPrefix(path, "/books/")
  if id == "" {
    http.NotFound(w, r)
    return
  }

  switch r.Method {
  case http.MethodGet:
    mu.Lock()
    bs, _ := readAll()
    mu.Unlock()
    for _, b := range bs {
      if b.ID == id {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(b)
        return
      }
    }
    w.WriteHeader(http.StatusNotFound)

  case http.MethodPut:
    body, _ := io.ReadAll(r.Body)
    var in Book
    if err := json.Unmarshal(body, &in); err != nil {
      w.WriteHeader(http.StatusBadRequest)
      return
    }
    mu.Lock()
    bs, _ := readAll()
    updated := false
    for i := range bs {
      if bs[i].ID == id {
        if strings.TrimSpace(in.Title) != "" {
          bs[i].Title = in.Title
        }
        if strings.TrimSpace(in.Author) != "" {
          bs[i].Author = in.Author
        }
        if in.Price != 0 {
          bs[i].Price = in.Price
        }
        // Available has a default false; we still set it explicitly
        bs[i].Available = in.Available
        updated = true
        break
      }
    }
    if updated {
      _ = writeAll(bs)
      mu.Unlock()
      w.Header().Set("Content-Type", "application/json")
      json.NewEncoder(w).Encode(map[string]string{"message": "updated"})
      return
    }
    mu.Unlock()
    w.WriteHeader(http.StatusNotFound)

  case http.MethodDelete:
    mu.Lock()
    bs, _ := readAll()
    out := make([]Book, 0, len(bs))
    found := false
    for _, b := range bs {
      if b.ID == id {
        found = true
        continue
      }
      out = append(out, b)
    }
    if found {
      _ = writeAll(out)
      mu.Unlock()
      w.WriteHeader(http.StatusNoContent)
      return
    }
    mu.Unlock()
    w.WriteHeader(http.StatusNotFound)

  default:
    w.WriteHeader(http.StatusMethodNotAllowed)
  }
}

func main() {
  os.MkdirAll("/data", 0755)

  mux := http.NewServeMux()
  mux.HandleFunc("/books", booksHandler)
  mux.HandleFunc("/books/", bookHandler)
  mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"ok":true}`))
  })

  http.ListenAndServe(":4000", mux)
}
