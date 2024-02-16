package src

import "github.com/spf13/viper"

var pagesIndex = RexFile{viper.GetString("pages.index"), `---
layout: default
---

# Rex

## ADRs

{% assign pages = site.pages -%}

| Title | Link |
| ----- | ---- |{% for page in pages -%}
{% if page.path contains 'adr' %}
|{{ page.title }} |[Click Here]({{ page.url | relative_url }}) |
{%- endif %}
{%- endfor -%}`}

var pagesDefaultADR = RexFile{viper.GetString("templates.default.adr"), `---
permalink: /:path/:basename:output_ext
title: {{ .WebTitle }}
layout: adr
---


# {{ .ADR.Title }}

| Status | Author         | Date       |
| ------ | -------------- | ---------- |
| {{ .ADR.Status }} | {{ .ADR.Author }} | {{ .ADR.Date }} |

## Context and Problem Statement

## Decision Drivers


## Considererd Options


## Decision Outcome
`}

var pagesConfig = RexFile{viper.GetString("pages.web.config"), `baseurl: "/rex"`}

var pagesWebLayoutAdr = RexFile{viper.GetString("pages.web.layout.adr"), `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>{{ page.title }}</title>
    <!-- <link rel="stylesheet" href="{{ "/assets/css/style.css?v=" | append: site.github.build_revision | relative_url }}"> -->
    <!-- <link rel="stylesheet" href="{{ "assets/css/style.css" | relative_url }}" /> -->

    <link rel="stylesheet" href="{{ "/assets/css/style.css?v=" | append: site.github.build_revision | relative_url }}">
  </head>
  <body>
    <div class="container-lg px-3 my-5 markdown-body">
      <nav>
        <a href="/rex/">Home</a>
      </nav>
      <h1>{{ page.title }}</h1>
      <section>{{ content }}</section>
      <!-- <footer>&copy; to me</footer> -->
    </div>
  </body>
</html>`}

var pagesWebLayoutDefault = RexFile{viper.GetString("pages.web.layout.default"), `<!DOCTYPE html>
<html lang="{{ site.lang | default: "en-US" }}">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <link rel="stylesheet" href="{{ "/assets/css/style.css?v=" | append: site.github.build_revision | relative_url }}">
</head>
<body>
<div class="container-lg px-3 my-5 markdown-body">
    {{ content }}

</div>
</body>
</html>`}
