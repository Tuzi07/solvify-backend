FROM golang:latest

WORKDIR /app

# Copy the root folder and its files
COPY . .

# Copy the subfolders and their files
COPY build /app/build
COPY cmd /app/cmd
COPY internal /app/internal
RUN go mod download

EXPOSE 8080

CMD ["make", "run"]