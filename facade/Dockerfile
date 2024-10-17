FROM nginx

RUN mkdir -p /app/syspulse/ui
COPY nginx.conf /etc/nginx/nginx.conf
COPY dist /app/syspulse/ui

EXPOSE 80