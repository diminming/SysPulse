FROM nginx

ARG ENVIRONMENT=staging

RUN mkdir -p /app/syspulse/ui
COPY env_conf/${ENVIRONMENT}/nginx.conf /etc/nginx/nginx.conf
COPY dist /app/syspulse/ui

EXPOSE 80