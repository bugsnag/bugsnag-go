FROM ruby:2.5

RUN mkdir -p /usr/src
WORKDIR /usr/src

COPY app/Gemfile /usr/src
RUN bundle install

COPY app /usr/src

CMD bundle exec ruby server.rb
