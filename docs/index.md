---
layout: default
permalink: /rex
---

# Rex

## ADRs

| Title | Link |
| ----- | ---- |{% assign pages = site.pages -%}
{% for page in pages %}
{% if page.path contains 'adr' -%}
|{{ page.title }} |[Click Here]({{ page.url }}) |
{% endif %}
{% endfor -%}
