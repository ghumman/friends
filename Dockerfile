FROM node:14.21.1 as builder
WORKDIR /app
COPY . .
RUN yarn install
RUN yarn run build
# Small websites can be served using serve, for production use nginx
# CMD [ "npx", "serve", "build" ]

# Bundle static assets with nginx
FROM nginx:1.21.0-alpine as production
ENV NODE_ENV production
# Copy built assets from `builder` image
COPY --from=builder /app/build /usr/share/nginx/html
# Add your nginx.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf
# Expose port
EXPOSE 80
# Start nginx
CMD ["nginx", "-g", "daemon off;"]