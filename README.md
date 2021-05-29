# FileSys404 - Handler for NotFound or 404 of Files

An asset not found or 404 Handler for static files Golang net/http default router.

Features:
- Standard plugging using `FileSystem` type of `net/http` package
- Fully compatible with `DefaultServeMux` of `net/http` package
- Protect `.dot` files or hidden files from being served
- Redirect the not found request to a pre-define custom `Handler`
- Can be used with custom routers like [httprouter](https://github.com/julienschmidt/httprouter) and [chi](https://github.com/go-chi/chi).

Docs available at : https://pkg.go.dev/github.com/boseji/filesys404

## Installing

```
go get -u github.com/boseji/filesys404
```

## Usage

```go
	dir := http.Dir("./static") // Static Assets directory
	fs := filesys404.New(dir, func(w http.ResponseWriter, r *http.Request){
		// Custom 404 Message
		http.Error(w, "We are sorry - this resource was not found", http.StatusNotImplemented)
	})
	http.HandleFunc("/static",http.StripPrefix("/static", fs))
```

## License

```
Copyright 2021 Abhijit Bose. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```