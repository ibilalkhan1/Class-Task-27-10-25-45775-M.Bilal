package main
import ("os"; "testing")

func TestReadWriteAll(t *testing.T){
  f, err := os.CreateTemp("", "catalog-*.json"); if err != nil { t.Fatal(err) }
  defer os.Remove(f.Name())
  os.Setenv("CATALOG_FILE", f.Name()); dataFile = f.Name()

  in := []Book{{ID:"t1", Title:"T", Author:"A", Price:1.23, Available:true}}
  if err := writeAll(in); err != nil { t.Fatal(err) }

  out, err := readAll(); if err != nil { t.Fatal(err) }
  if len(out)!=1 || out[0].ID!="t1" { t.Fatalf("unexpected: %#v", out) }
}
