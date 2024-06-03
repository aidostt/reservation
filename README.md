# reservation-service


### To run the container you should write
docker build -t reservation-service .
docker run -p 9090:9090 --env-file .env -ti reservation-service