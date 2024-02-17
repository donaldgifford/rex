---
layout: default
---

# Rex

{% assign pages = site.pages -%}
{% for page in pages -%}
{% if page.path contains 'install' %}
|{{ page.title }} |[install]({{ page.url | relative_url }}) |
{%- endif %}
{% if page.path contains 'usage' %}
|{{ page.title }} |[usage]({{ page.url | relative_url }}) |
{%- endif %}
{%- endfor %}

## ADRs

| Title | Link |
| ----- | ---- |{% for page in pages -%}
{% if page.path contains 'adr' %}
|{{ page.title }} |[Click Here]({{ page.url | relative_url }}) |
{%- endif %}
{%- endfor -%}
