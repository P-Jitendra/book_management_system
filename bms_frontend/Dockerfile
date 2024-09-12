FROM node:latest

WORKDIR /usr/src/app

COPY package*.json book_data.json ./

RUN npm install

RUN apt-get -qq -y install curl

COPY . .

EXPOSE 3000

CMD ["npm start"]
