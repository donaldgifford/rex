adr:
  path: "docs/adr/"
  index_page: "README.md"
  add_to_index: true # on rex create, a new record will be added to the index page
templates:
  enabled: false # uses embedded templates by default. If true reference the paths
  path: "templates/"
  adr:
    default: "adr.tmpl"
    index: "index.tmpl"
enable_github_pages: true
pages:
  index: "index.md"
  web:
    config: "_config.yml"
    layout:
      adr: "adr.html"
      default: "default.html"
extras: true
extra_pages:
  install: install.md
  usage: usage.md