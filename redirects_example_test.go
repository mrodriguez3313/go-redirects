package redirects_test

import (
	"encoding/json"
	"os"

	"github.com/fission-suite/go-redirects"
)

func Example() {
	// from [a=:save1 b=value] to [code][!] [Country=x,y,z] [Language=x,y,z]
	h := redirects.Must(redirects.ParseString(`
		# Implicit 301 redirects
		#/home              /
		#/blog/my-post.php  /blog/my-post
		# /news              /blog
		# /google            https://www.google.com

		## Redirect with a 301
		#/home         /              301
#
		## Redirect with a 302
		#/my-redirect  /              302
#
		## Rewrite a path
		#/pass-through /index.html    200
#
		## Show a custom 404 for this path
		#/ecommerce    /store-closed  404
#
		## Single page app rewrite
		#/*    /index.html   200
#
		## Proxying
		#/api/*  https://api.example.com/:splat  200
#
		## Forcing
		#/app/*  /app/index.html  200!
#
		## Params
		#/	/something	302	foo=bar
		#/	/something	302	foo=bar bar=baz
		#/	id=:id /blog/:id 302
		#/articles id=:id tag=:tag /posts/:tag/:id 301!
		#/ 	/auzy 302 Country=au,nz
		/israel/*  /israel/he/:splat  302  Country=au,nz Language=he
  `))

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(h)
	// Output:
	// 	[
	//   {
	//     "From": "/home",
	//     "To": "/",
	//     "Status": 301,
	//     "Force": false,
	//     "Params": null
	//   },
	//   {
	//     "From": "/blog/my-post.php",
	//     "To": "/blog/my-post",
	//     "Status": 301,
	//     "Force": false,
	//     "Params": null
	//   },
	//   {
	//     "From": "/news",
	//     "To": "/blog",
	//     "Status": 301,
	//     "Force": false,
	//     "Params": null
	//   },
	//   {
	//     "From": "/google",
	//     "To": "https://www.google.com",
	//     "Status": 301,
	//     "Force": false,
	//     "Params": null
	//   },
	//   {
	//     "From": "/home",
	//     "To": "/",
	//     "Status": 301,
	//     "Force": false,
	//     "Params": null
	//   },
	//   {
	//     "From": "/my-redirect",
	//     "To": "/",
	//     "Status": 302,
	//     "Force": false,
	//     "Params": null
	//   },
	//   {
	//     "From": "/pass-through",
	//     "To": "/index.html",
	//     "Status": 200,
	//     "Force": false,
	//     "Params": null
	//   },
	//   {
	//     "From": "/ecommerce",
	//     "To": "/store-closed",
	//     "Status": 404,
	//     "Force": false,
	//     "Params": null
	//   },
	//   {
	//     "From": "/*",
	//     "To": "/index.html",
	//     "Status": 200,
	//     "Force": false,
	//     "Params": null
	//   },
	//   {
	//     "From": "/api/*",
	//     "To": "https://api.example.com/:splat",
	//     "Status": 200,
	//     "Force": false,
	//     "Params": null
	//   },
	//   {
	//     "From": "/app/*",
	//     "To": "/app/index.html",
	//     "Status": 200,
	//     "Force": true,
	//     "Params": {
	//       "id":":id"
	//     }
	//   },
	//   {
	//     "From": "/",
	//     "To": "/something",
	//     "Status": 302,
	//     "Force": false,
	//     "Params": {
	//       "id": ":id"
	//       "tag": ":tag"
	//     }
	//   },
	//   {
	//     "From": "/",
	//     "To": "/something",
	//     "Status": 302,
	//     "Force": false,
	//     "Params": null
	//   }
	// ]
}
