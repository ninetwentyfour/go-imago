// Gross hack to not have to deal with .html files after install
package main

import (
	"fmt"
)

func homeHtml() []byte {
	return []byte(`<!DOCTYPE html>
<html>
<head>
  <title>imago.rb</title>
  <meta http-equiv="content-type" content="text/html; charset=UTF-8">
  <link rel="stylesheet" media="all" href="https://imago.in/docs/public/stylesheets/normalize.css" />
  <link rel="stylesheet" media="all" href="https://imago.in/docs/docco.css" />
  <style>
    .container {
      width: 900px;
    }
    .page {
      width:780px;
    }
    h1:after {
      content: none;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="page">
      <div class="header">
          <h1>Imago</h1>
      </div>
      <p>
        Create images of websites
      </p>
      <ul>
        <li>
          <a href="/1000/800/example.com/json">/1000/800/example.com/json</a>
        </li>
        <li>
          <a href="/1000/800/example.com/html">/1000/800/example.com/html</a>
        </li>
        <li>
          <a href="/1000/800/example.com/image">/1000/800/example.com/image</a>
        </li>
      </ul>
    </div>
  </div>
</body>
</html>`)
}

func pageHtml(link, url string) []byte {
	return []byte(fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <title>imago.rb</title>
  <meta http-equiv="content-type" content="text/html; charset=UTF-8">
  <link rel="stylesheet" media="all" href="https://imago.in/docs/public/stylesheets/normalize.css" />
  <link rel="stylesheet" media="all" href="https://imago.in/docs/docco.css" />
  <style>
    .container {
      width: 900px;
    }
    .page {
      width:780px;
    }
    h1:after {
      content: none;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="page">
      <div class="header">
          <h1>Imago</h1>
      </div>
        <p style="text-align:center;">
          %s
        </p>
        <p style="text-align:center;">
          <a href="%s">%s</a>
        </p>
    </div>
  </div>
</body>
</html>`, url, link, link))
}
