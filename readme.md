# Movie Festival

Simple app with feature:
- Create artist, genre, movie
- Vote movie
- Upload Movie

To run the app 
```bash
make up
```

or if you cant use makefile

```bash
docker-compose -f $(DOCKER_COMPOSE_FILE) --env-file ./.env up -d --build
```

And for the env itself, it's not good to expose it in git, but for easier setup and there is no really a secret, i think it's ok :D

For other documentation, you check it directly in postman collection which is `Movie Festival.postman_collection.json`