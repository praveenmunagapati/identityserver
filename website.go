package website

//package the assets
//go:generate go-bindata -pkg assets -prefix assets -o ./packaged/assets/assets.go assets/...

//package 3rd party assets
//go:generate go-bindata -pkg thirdpartyassets -prefix thirdpartyassets -o ./packaged/thirdpartyassets/thirdpartyassets.go thirdpartyassets/...

//go:generate go-bindata -pkg components -prefix components -o ./packaged/components/components.go components/...

//package the html files
//go:generate go-bindata -pkg html -o ./packaged/html/html.go index.html registration.html login.html home.html error.html apidocumentation.html
