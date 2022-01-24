module "test_local" {
  source = "./module"
}

module "test_remote" {
  source = "github.com/something/module"
}