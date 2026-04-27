#!/bin/sh
# Seed content/posts from the bundled image copy on first run.
# After first run, the volume takes over and persists new admin-created posts.
if [ -z "$(ls -A content/posts 2>/dev/null)" ]; then
  echo "Seeding content/posts from image..."
  mkdir -p content/posts
  cp -r content-seed/posts/. content/posts/
fi
exec ./server
