# Docs


## Task Management System

There are 4 types of items:

  - RFC
  - ADR
  - Plans
  - Tasks

### RFC

RFC or Request For Comments is the highest level document in this system . It
serves as a immutable artifact that is a snapshot in time of when a non
technical decision was made. For more information refer to the RFC [README](./rfc/README.md).

### ADR

ADR or "Architectural Design Record" is the highest level technical document in
this system. Typically an RFC will be the parent document of one or many ADR's.
These serve as immutable objects that snapshot a technical decision made at a
point in time. More information is about ADR's is available in [the ADR
README](./adr/README.md).

### Plans

Plans serve a collection of Tasks, similar to a sprint but not bucketed by time
or points. They are simply a grouping of related tasks that can fit a useful
context window for Claude. More information is about ADR's is available in [the plans
README](./plans/README.md).


### Tasks

Tasks are what Claude can use as work item. This is similar to how you can
create a todo list in plan mode with Claude but these serve as a local version
of its memory. Doing it this way we get free task management system and
documentation. More information is about tasks is available in [the tasks
README](./tasks/README.md).

