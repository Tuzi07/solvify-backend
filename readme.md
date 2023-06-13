## Docker Container Creation

- docker build --tag solvify-api .
- docker run --name Solvify-Api --env MONGODB_URI="mongodb+srv://tuzi:GJKX2T1ybZU8W6wT@solvify.r5ehjjx.mongodb.net" --publish 8080:8080 solvify-api
