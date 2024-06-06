Сборка и запуск с помощью docker-compose
-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
Запуск с in-memory хранилищем: >docker-compose build app_im
>docker-compose run -p 127.0.0.1:8080:8080 -d app_im

Запуск с postgresql хранилищем: >docker-compose build app_db
>docker-compose run -p 127.0.0.1:8080:8080 -d app_db

Работа с сервисом
-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
Все query запросы, mutation register и login, subscription commentAdded можно отправлять без авторизации.
Остальные mutation запросы, нужно отправлять с токеном авторизации в headers. 
