FROM node:18-alpine

COPY package*.json /usr/app/

COPY index.js /usr/app/

COPY tracing.js /usr/app/

WORKDIR /usr/app

RUN npm install

CMD ["node", "index.js"]
