#!/usr/bin/bash

DOCKER_IMAGE_NAME="syspulse_db"
DB_HOST="127.0.0.1"
DB_USER="root"
DB_PASSED="123qweASD"
DB_NAME="syspulse"

cat > 01-create_db.sql << 'EOF'
CREATE DATABASE `syspulse` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
EOF

# mysqldump -h${DB_HOST} -u${DB_USER} -p${DB_PASSED} --opt -d ${DB_NAME} > 02-syspulse.sql

# sed -i "1i USE \`syspulse\`" 02-syspulse.sql

cat > 03-create_user.sql << 'EOF'
create user syspulse identified by '123qweASD';
grant all on syspulse.* to syspulse@'%';
flush privileges;
EOF

cat > 04-init_data.sql << 'EOF'
USE `syspulse`;

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (1,'admin','admin1!',1,0,0);
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;
EOF

docker build -t ${DOCKER_IMAGE_NAME} .

docker save -o ${DOCKER_IMAGE_NAME}.tar ${DOCKER_IMAGE_NAME}