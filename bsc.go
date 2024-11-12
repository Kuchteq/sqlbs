package main
import ("math/rand"; "fmt")
var bsc = map[string][]string{
    "first_name":{"Mariusz","Paweł"},
    "last_name":{"Pytel","Zajac"},
    "username":{"kuchteq","pzayac"},
}
var bsc_composite = map[string]func() string{
"fullname": func() string {
	r_0 := rand.Intn(len(bsc["first_name"]))
	r_1 := rand.Intn(len(bsc["last_name"]))
	r := fmt.Sprintf("%s %s",bsc["first_name"][r_0],bsc["last_name"][r_1])
return r}}
