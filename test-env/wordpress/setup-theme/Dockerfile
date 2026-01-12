FROM python:3-trixie

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
# this will install some shared libraries
RUN playwright install --with-deps chromium

COPY main.py .
COPY run.sh .

# Default arguments
ARG HTTP_HOST="http://wordpress"
ARG WORDPRESS_ADMIN_USER="admin"
ARG WORDPRESS_ADMIN_PASSWORD="password"

# Environment variables
ENV HTTP_HOST=$HTTP_HOST \
    WORDPRESS_ADMIN_USER=$WORDPRESS_ADMIN_USER \
    WORDPRESS_ADMIN_PASSWORD=$WORDPRESS_ADMIN_PASSWORD \
    PYTHONUNBUFFERED=1

CMD ["./run.sh"]
