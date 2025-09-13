env "local" {
  src = "file://migrations"
  url = "mysql://root:password@localhost:3306/budget_book"
  dev = "docker://mysql/8/budget_book_dev"
}

env "docker" {
  src = "file://migrations" 
  url = "mysql://root:password@mysql:3306/budget_book"
  dev = "docker://mysql/8/budget_book_dev"
}