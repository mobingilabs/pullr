FROM node:9
WORKDIR /home/node/app
COPY ui/package.json .
RUN npm install
ENV NODE_ENV production
COPY ui .
RUN npm run build

FROM nginx:1.12.2-alpine
COPY ./conf/nginx-ui.conf /etc/nginx/conf.d/default.conf
COPY --from=0  /home/node/app/dist /www
