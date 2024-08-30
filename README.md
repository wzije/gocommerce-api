# go-ecommerce

## Requirement

- golang v1.22
- postgresql v13
- docker

## RUNNING

1. Locally
    - clone project
    - go mod install
    - setup database in .env file
2. docker
   - docker build . 

## migration DB

Use this command to execute migration

```bash
cd database
make migrate_make
make migrate_up
make migrate_down
make migrate_drop
make migrate_seed
make migrate_reset
```

## DEPLOYMENT

1. Dev git push origin main:deploy-dev
2. Stage/Prod manually add tag
