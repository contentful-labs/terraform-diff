# Terraform-config-deps

Terraform-config-deps is a thin wrapper around [terraform-config-inspect](https://github.com/hashicorp/terraform-config-inspect)
to help compute, for each folder passed as argument, a list of folders where changes might trigger a plan for
said folder.

Given a git diff in a Terraform repository, this will help compute a list of folders where a Terraform plan & apply
need to be run.

## Example

Given the following Terraform setup: 

``` 
.
├── modules
│   └── module1
│   │   └── main.tf
│   └── module2
│   │   └── main.tf
│   └── module3
│       └── main.tf
└── project1
│   └─── main.tf
└── project2
    └─── main.tf
```

And the following dependencies:
 * project1 depends on module1
 * module1 depends on module2
 * project2 depends on module3

The logic would be: 
 * if there is a change in modules/module1/, modules/module2/, or project1/, you would want to run `make plan` in project1/.
 * if there is a change in modules/module3/, or project2/, you would want to run `make plan` in project2/

``terraform-config-deps`` can help identify what folders/modules a terraform project depends on:

```
$ terraform-config-deps project1 project2

{
  "project1": [
    "project1",
    "modules/module1",
    "modules/module2"
 ], 
  "project2": [
    "project2",
    "modules/module3",
 ], 
{


```

