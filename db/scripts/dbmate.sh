source ./docker/.env.development

dbmate --env DATABASE_URL --schema-file ./db/schema.sql --migrations-dir ./db/migrations "$@"
