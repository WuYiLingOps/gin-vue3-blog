pipeline {
    agent any

    environment {
        appName = 'gin-vue3-blog'
        harborServer = 'harbor.huang.org'
        repo = 'devops'
        imageUrl = "${harborServer}/${repo}/${appName}"
        imageTag = "${BUILD_NUMBER}"
        credential = 'harbor-user-credential'
        hostList = '10.0.0.112 10.0.0.113'
        port = '8080'

        // 数据库配置（默认值，可通过 Jenkins Parameters 覆盖）
        dbHost = "${params.DB_HOST ?: 'pg.huang.org'}"
        dbPort = '5432'
        dbUser = 'postgres'
        dbPassword = '123456ok!'
        dbName = 'blogdb'
        redisHost = "${params.REDIS_HOST ?: 'redis.huang.org'}"
        redisPort = '6379'
        redisPassword = '123456'
    }

    parameters {
        string(name: 'DB_HOST', defaultValue: '', description: '数据库地址（留空使用默认值）')
        string(name: 'REDIS_HOST', defaultValue: '', description: 'Redis地址（留空使用默认值）')
    }

    triggers {
        gitlab(triggerOnPush: true,
            acceptMergeRequestOnSuccess: false,
            branchFilterType: 'All',
            secretToken: '62dad2cd1d9ae62686ada8dc4cd0ae66')
    }

    stages {
        stage('Source') {
            steps {
                echo "📥 代码拉取完成"
            }
        }

        stage('Generate Config') {
            steps {
                echo "🔧 生成环境配置文件"
                sh """
                    umask 077
                    cat > blog-backend/.env.config.prod <<EOF
BLOG_URL=https://huangjingblog.cn
DB_HOST=${dbHost}
DB_PORT=${dbPort}
DB_USER=${dbUser}
DB_PASSWORD=${dbPassword}
DB_NAME=${dbName}
REDIS_HOST=${redisHost}
REDIS_PORT=${redisPort}
REDIS_PASSWORD=${redisPassword}
JWT_SECRET=your-jwt-secret-key
JWT_EXPIRE_HOURS=72
GITEE_CALENDAR_API_URL=https://huangjingblog.cn/gitee-calendar-api
EMAIL_HOST=smtp.qq.com
EMAIL_PORT=587
EMAIL_USERNAME=
EMAIL_PASSWORD=
EOF
                """
            }
        }

        stage('Build Docker Image') {
            steps {
                echo "🐳 构建 Docker 镜像"
                sh 'docker build -f deploy/Dockerfile -t "${imageUrl}:${imageTag}" .'
            }
        }

        stage('Push Docker Image') {
            steps {
                echo "📤 推送镜像到 Harbor"
                withCredentials([usernamePassword(credentialsId: "${credential}", \
                    passwordVariable: 'harborPassword', usernameVariable: 'harborUserName')]) {
                    sh "echo ${harborPassword} | docker login -u ${env.harborUserName} --password-stdin ${harborServer}"
                    sh "docker push ${imageUrl}:${imageTag}"
                }
            }
        }

        stage('Deploy') {
            steps {
                echo "🚀 部署到远程服务器"
                script {
                    // 修改 config.yml 为生产环境
                    sh "sed -i 's/^env:.*/env: prod/' blog-backend/config/config.yml"

                    def hosts = hostList.trim().split(/\s+/)
                    hosts.each { host ->
                        echo "📦 发布到 ${host}"
                        // 创建目录并复制配置文件
                        sh "ssh root@${host} \"mkdir -p /web/${appName}-data/config\""
                        sh "scp blog-backend/.env.config.prod root@${host}:/web/${appName}-data/config/.env.config.prod"
                        sh "scp -r blog-backend/config/. root@${host}:/web/${appName}-data/config/"
                        // 部署容器（添加 --add-host 解析内网域名）
                        sh """ssh root@${host} "docker rm -f ${appName} 2>/dev/null || true && docker run -d --restart unless-stopped --name ${appName} -p ${port}:80 \
                            --add-host pg.huang.org:10.0.0.1 \
                            --add-host redis.huang.org:10.0.0.1 \
                            -v /web/${appName}-data/config:/app/config \
                            -v /web/${appName}-data/config/.env.config.prod:/app/.env.config.prod:ro \
                            -v /web/${appName}-data/uploads:/app/uploads \
                            -v /var/log/${appName}:/var/log/${appName} \
                            ${imageUrl}:${imageTag}\""""
                    }
                }
            }
        }
    }

    post {
        always {
            sh 'rm -f blog-backend/.env.config.prod'
        }
        success {
            mail to: '2794998160@qq.com',
                 subject: "构建成功：${env.JOB_NAME} #${env.BUILD_NUMBER}",
                 body: "✅ 构建成功\n\n任务：${env.JOB_NAME}\n编号：${env.BUILD_NUMBER}\n地址：${env.BUILD_URL}"
            dingtalk(
                robot: "dingtalk",
                type: "MARKDOWN",
                title: "构建成功：${env.JOB_NAME} #${env.BUILD_NUMBER}",
                text: [
                    "### ✅ 构建成功",
                    "- **任务名称**：${env.JOB_NAME}",
                    "- **构建编号**：${env.BUILD_NUMBER}",
                    "- **构建地址**：${env.BUILD_URL}"
                ]
            )
        }
        failure {
            mail to: '2794998160@qq.com',
                 subject: "构建失败：${env.JOB_NAME} #${env.BUILD_NUMBER}",
                 body: "❌ 构建失败\n\n任务：${env.JOB_NAME}\n编号：${env.BUILD_NUMBER}\n地址：${env.BUILD_URL}"
            dingtalk(
                robot: "dingtalk",
                type: "MARKDOWN",
                title: "构建失败：${env.JOB_NAME} #${env.BUILD_NUMBER}",
                text: [
                    "### ❌ 构建失败",
                    "- **任务名称**：${env.JOB_NAME}",
                    "- **构建编号**：${env.BUILD_NUMBER}",
                    "- **构建地址**：${env.BUILD_URL}"
                ]
            )
        }
    }
}