module.exports = {
  apps: [
    {
      name: "blog-server",
      script: "./blog-server",
      instances: 1,
      autorestart: true,
      watch: false,
      max_memory_restart: "1G",
      env: {
        DATABASE_URL: "postgres://blog_user:abdulloh_009@localhost/blog_system?sslmode=disable",
        PORT: 8080,
      },
      error_file: "./logs/err.log",
      out_file: "./logs/out.log",
      log_file: "./logs/combined.log",
      time: true,
    },
  ],
}
