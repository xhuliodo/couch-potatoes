FROM node:20 AS build
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
# webpack 4 (CRA 4) needs the legacy OpenSSL provider on Node 17+
ENV NODE_OPTIONS=--openssl-legacy-provider

ARG REACT_APP_AUTH0_DOMAIN
ENV REACT_APP_AUTH0_DOMAIN=$REACT_APP_AUTH0_DOMAIN
ARG REACT_APP_AUTH0_CLIENT_ID
ENV REACT_APP_AUTH0_CLIENT_ID=$REACT_APP_AUTH0_CLIENT_ID
ARG REACT_APP_REDIRECT_BASE
ENV REACT_APP_REDIRECT_BASE=$REACT_APP_REDIRECT_BASE
ARG REACT_APP_API_ENDPOINT
ENV REACT_APP_API_ENDPOINT=$REACT_APP_API_ENDPOINT
ARG REACT_APP_THUMBNAILS_ENDPOINT
ENV REACT_APP_THUMBNAILS_ENDPOINT=$REACT_APP_THUMBNAILS_ENDPOINT

COPY package.json ./
COPY yarn.lock ./
RUN yarn install --frozen-lockfile --network-timeout 600000
COPY . ./
RUN yarn run build

# production environment
FROM nginx:stable-alpine
COPY --from=build /app/build /usr/share/nginx/html
# react-router config
COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 3000
CMD ["nginx", "-g", "daemon off;"]
