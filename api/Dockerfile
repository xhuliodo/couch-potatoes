FROM node:12

RUN mkdir -p /app
WORKDIR /app

COPY package.json yarn.lock ./
RUN yarn install --silent
COPY . .

EXPOSE 4001

CMD ["yarn", "run", "start"]