POSTGRES_USER=postgres
POSTGRES_PASSWORD=d9c3452013d1106e3be8ac25d2068cdc

docker rm -f typi-database
docker rmi -f typi-database

cd $(dirname "$0")/../../../database
docker build -t typi-database .
docker run -d --name typi-database \
    -p 4600:5432 \
    -e POSTGRES_USER=$POSTGRES_USER \
    -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
    typi-database