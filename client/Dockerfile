FROM node:12.18-alpine3.9 as build
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
ARG REACT_APP_API_ENDPOINT
ENV REACT_APP_API_ENDPOINT=$REACT_APP_API_ENDPOINT
COPY package.json ./
COPY yarn.lock ./
RUN yarn --silent
RUN yarn add global react-scripts@4.0.1 --silent
COPY . ./
RUN yarn run build

# production environment
FROM nginx:stable-alpine
COPY --from=build /app/build /usr/share/nginx/html
#react-router config
COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]