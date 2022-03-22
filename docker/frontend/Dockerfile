FROM node:17.7-alpine

ENV APP_HOME /var/www/frontend
ENV LANG C.UTF-8
ENV EDITOR vim
ENV TZ Asia/Tokyo

RUN apk --update add vim

WORKDIR $APP_HOME

COPY ./frontend $APP_HOME

RUN yarn install
