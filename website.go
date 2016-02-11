package website

//package the assets
//go:generate go-bindata -pkg assets -prefix assets -o ./packaged/assets/assets.go assets/...

//package the html files
//go:generate go-bindata -pkg html -o ./packaged/html/html.go index.html registration.html
