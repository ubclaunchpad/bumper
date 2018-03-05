FROM node:carbon

WORKDIR /app
ENV NODE_ENV=production

# Copy over source code and install dependencies
COPY . .
RUN npm install --production

# Expose ports and run application
EXPOSE 8080
ENTRYPOINT npm start
