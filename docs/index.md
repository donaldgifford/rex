---
layout: default
---

# Rex

## ADRs

| Title | Link |
| ----- | ---- |{% assign pages = site.pages -%}
{% for page in pages %}
{% if page.path contains 'adr' -%}
|{{ page.title }} |[Click Here]({{ page.url | relative_url }}) |
{% endif %}
{% endfor -%}
