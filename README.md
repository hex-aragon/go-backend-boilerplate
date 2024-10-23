# go-backend-boilerplate

# 1. db

## postgres

- docker image = postgres 15-alpine
- using sqlc

## Ref

sqlc :
https://github.com/sqlc-dev/sqlc/tree/v1.4.0
https://github.com/sqlc-dev/sqlc/blob/v1.4.0/docs/query_one.md
https://github.com/sqlc-dev/sqlc/blob/v1.4.0/docs/insert.md
https://github.com/sqlc-dev/sqlc/blob/v1.4.0/docs/update.md
https://github.com/sqlc-dev/sqlc/blob/v1.4.0/docs/delete.md

simple bank:
https://github.com/techschool/simplebank/blob/master/db/query/entry.sql/

## docker command

```
각 컨테이너 네트워크 설정 확인
docker container inspect "컨테이너"
docker container inspect "go-boiler-postgres"

도커 네트워크 생성
docker network create "네트워크명"
docker network create bank-network

도커 네트워크 조회
docker network ls

도커 네트워크 검사
docker network inspect "네트워크명"
docker network inspect bridge
docker network inspect bank-network

도커 네트워크에서 컨테이너를 연결
docker network connect "네트워크명" "컨테이너 네트워크 이름"
docker network connect bank-network go-boiler-postgres


```
