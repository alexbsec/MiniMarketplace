data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "./src/scripts/atlas-gorm-loader.go",
  ]
}

env "gorm" {
    src = data.external_schema.gorm.url
    migration {
        dir = "file://src/db/migrations"
    }
    dev = "postgres://user:password@db_dev:5432/marketplace_db_dev?sslmode=disable"
    url = "postgres://user:password@db:5432/marketplace_db?sslmode=disable"

}
