# Revisiting tmpl file

We define the named template using the `{{define "template-name"}} ... {{end}}` tags

We render dynamic data via the `.` character like `{{.ID}}`

the `//go:embed "<path>"` can only be used on global variables at package level, and the path should be **relative**

The embedded file system should be rooted in the directory which contains the `//go:embed` directive

You can specify multiple directories like `//go:embed "images" "styles/css" "favicon.ico" .`
