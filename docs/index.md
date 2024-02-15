---
layout: default
---

# Rex

## ADRs

| Title | Link |
| ----- | ---- |{% assign pages = site.pages -%}
{% for page in pages %}
{% if page.title -%}
|{{ page.title }} |[Click Here]({{ page.url }}) |
{% endif %}
{% endfor -%}
