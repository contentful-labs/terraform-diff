# Terraform-diff

Terraform-diff helps you detect what Terraform projects have changed when changes are made to Terraform modules.

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

This is where Terraform-diff is useful:

```
$ terraform-diff -h
Usage of terraform-diff:
  -output string
      output format (text or json) (default "text")
  -range string
      git commit range
$ terraform-diff project1 project2
project1
$ terraform-diff --range fbf666c786...ca37f7145f -o json project1 project2
{
  "project1",
  "project2"
}
```

## Trade-offs

``terraform-diff`` relies on git & static analysis of the Terraform files. It will **not** detect, among others:
 * changes in external datasources
 * remote states updates
