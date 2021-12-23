FROM node:14.4.0-alpine3.12

WORKDIR /app
COPY . .

RUN apk add git && yarn install && yarn build


#FROM node:14.4.0-alpine3.12
#WORKDIR /app
#COPY . .
#COPY --from=builder /app/dist /app/dist
CMD npm run start:prod