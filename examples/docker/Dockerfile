FROM node:16-alpine AS build

WORKDIR /app
COPY ./web-sveltekit/package*.json ./
RUN npm install
COPY ./web-sveltekit ./
RUN npm run build

FROM nginx:1.23.3-alpine-slim AS deploy-static

WORKDIR /usr/share/nginx/html
RUN rm -rf ./*
COPY --from=build /app/build .
ENTRYPOINT ["nginx", "-g", "daemon off;"]

FROM node:16-alpine AS deploy-node

WORKDIR /app
RUN rm -rf ./*
COPY --from=build /app/package.json .
COPY --from=build /app/build .
ENV PORT=8080
EXPOSE 8080
CMD ["node", "index.js"]
