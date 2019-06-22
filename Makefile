build:
	rm -f build/*
	touch build/.keep

	env GOOS=linux go build -ldflags '-s -w' -o build/serverless_task_create  cmd/task/create/create.go
	env GOOS=linux go build -ldflags '-s -w' -o build/serverless_task_delete  cmd/task/delete/delete.go
	env GOOS=linux go build -ldflags '-s -w' -o build/serverless_task_index   cmd/task/index/index.go
	env GOOS=linux go build -ldflags '-s -w' -o build/serverless_task_migrate cmd/task/migrate/migrate.go
	env GOOS=linux go build -ldflags '-s -w' -o build/serverless_task_read    cmd/task/read/read.go
	env GOOS=linux go build -ldflags '-s -w' -o build/serverless_task_update  cmd/task/update/update.go
	
	chmod 777 build/*