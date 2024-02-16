---
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
{%- endfor -%}
