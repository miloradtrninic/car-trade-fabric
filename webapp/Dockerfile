FROM node:carbon
# Create app directory
WORKDIR /usr/src/app
# Install app dependencies
COPY package*.json ./
RUN npm install --loglevel verbose
# Copy app source code
COPY . .
#Expose port and start application
EXPOSE 3000
CMD [ "node", "index.js" ]
