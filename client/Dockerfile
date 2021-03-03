FROM node:12.18 as build
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
ARG REACT_APP_API_ENDPOINT
ENV REACT_APP_API_ENDPOINT=$REACT_APP_API_ENDPOINT
ARG REACT_APP_DOMAIN
ENV REACT_APP_DOMAIN=$REACT_APP_DOMAIN
COPY package.json ./
COPY yarn.lock ./
RUN yarn --silent
COPY . ./
RUN yarn run build 
# auto-update failed, or rather a bit too finicky
# RUN echo "self.addEventListener(\"install\", function() { self.skipWaiting(); });" >> ./build/service-worker.js


# production environment
FROM nginx:stable-alpine
COPY --from=build /app/build /usr/share/nginx/html
#react-router config
COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
