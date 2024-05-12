# Sorcerers Forge

Добро пожаловать в магическую кузницу, где искусство и волшебство переплетаются в уникальные работы и приключения! Sorcerers Forge — это мир, где каждый участник становится кузнецом своей судьбы, воплощая мечты и фантазии в форме стальных клинков, доспехов и артефактов.

## Запуск

Для запуска нужно пройти следующие шаги.

Нужно установить [docker](https://docs.docker.com/engine/install/) , [docker-compose](https://docs.docker.com/compose/install/), [go-sdk version 1.19](https://go.dev/dl/) и [migrate](https://github.com/golang-migrate/migrate)

## Старт в windows

Выполни команды
```
docker-compose up --build
migrate -path db/userMigrations -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' -verbose down
migrate -path db/userMigrations -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' -verbose up
migrate -path db/catalogMigrations -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' -verbose down
migrate -path db/catalogMigrations -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' -verbose up
migrate -path db/profileMigrations -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' -verbose down
migrate -path db/profileMigrations -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' -verbose up
migrate -path db/galleryMigrations -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' -verbose down
migrate -path db/galleryMigrations -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' -verbose up
migrate -path db/reviewsMigrations -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' -verbose down
migrate -path db/reviewsMigrations -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' -verbose up
migrate -path db/adminMigrations -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' -verbose down
migrate -path db/adminMigrations -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' -verbose up
migrate -path db/addressMigrations -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' -verbose down
migrate -path db/addressMigrations -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' -verbose up
migrate -path db/orederMigrations/ -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' -verbose down
migrate -path db/orederMigrations/ -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' -verbose up
migrate -path db/likesMigrations/ -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' -verbose down
migrate -path db/likesMigrations/ -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' -verbose up


go build ./cmd/main.go
go run ./cmd/main.go
```

## Описание APIs

```
POST /sign-up - регистрация тут происходит только отправка кода на почту для подтверждения:
поля:
     "email":"",
     "password":""
ответы:
    если успешно: "status": "success"
    если ввёл некорректный email: "error": "invalid email format"
    если пользователь с таким email уже зарегистрирован:  "error": "A user with this email is already registered"
     
POST /confirmCode - регистрация тут происходит проверка кода и регистарция пользователя:
поля:
    "email":"",
    "emailCode":"",
    "password":""
ответы:
     если успешно создал пользователя: "id": id созданного пользователя 
     если есть дубликаты:"error": "pq: duplicate key value violates unique constraint \"users_email_key\""
     если не заполнил поле email подпишет: "error": "email: cannot be blank."
     если не заполнил поле password подпишет:"error": "password: cannot be blank."
     если ничего не заполнил: "error": "email: cannot be blank; password: cannot be blank."
     если не соответсвует формат email:  "error": "email: must be a valid email address."
     если пароль меньше 6 символов: "error": "password: the length must be between 6 and 100."
     
POST /sign-in авторизация:
поля: 
    "email":"",
    "password":""
ответы:
успешный вход   "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDMwMTM2NjgsImlhdCI6MTcxMTQ3NzY2OCwidXNlcl9pZCI6NX0.EfRui1JtI89f7SK3BpuN1PslUO6dJyX4YPEy2yvQ2S0"
"error": "incorrect email or password"

GET /private/whoami - получение данных авторизованного пользователя(нужно заполнение Authorization token)
ответы:
     если авторизован: "id": 
                       "email":
     если не авторизован:
        1 не заполнил заголовок:  "error": "empty auth header"
        2 кривой заголовок:   "error": "signature is invalid" (например не до конца скопировал)
        3 написал что то не то: qwe к примеру:  "error": "token contains an invalid number of segments"     

POST /resetCode - отправка кода для восстановления пароля на почту:
поля: 
     "email":""
ответы:
    если успешно: "status": "success"
    если ввёл некорректный email: "error": "invalid email format"
    
POST /resetPassword - восстановление пароля после получения кода:
поля: 
   "email":"",
   "emailcode":"",
   "password":""   
ответы:
    если успешно сменил пароль: "status": "password successfully changed"
    если неправильно ввёл мейл или код: "error": "incorrect email or emailcode"
    
DELETE /deleteUsers - удаление пользователя: Query param id=2 (id юзера) 
ответы:
    если успешно: "status": "success"
    если нет такого email:  "error": "sql: no rows in result set"

POST /createCategory - создание категории для товара
поля: 
    "name":""
ответы:
    если успешно:  "id": , "status": "success"
    если не написать имя: "error": "invalid name" 
    
DELETE /deleteCategory - удаление категории товара: Query param id=4 (id категории)
ответы:
     "status": "delete success"
     если есть товары в категории, его надо удалить и потом категорию удалять, это норма, в бд связь "error": "pq: update or delete on table \"product_category\" violates foreign key constraint \"products_category_id_fkey\" on table \"products\""
     если нет такого id : "error": "sql: no rows in result set"

GET /categories - получение списка всех созданных категорий
ответы:
      если есть список: 
      "categories": [
        {
            "id": 1,
            "name": "Мечи"
        },
        {
            "id": 3,
            "name": ""
        },
        {
            "id": 4,
            "name": ""
        }
    ]
      если нет списка:
       "categories": []  

POST /createProduct - создание товара:
поля:
    "name": "",
    "description": "",
    "price": ,
    "quantity": ,
    "photo": "",
    "work_time":,
    "category_id": ,
    "is_active":
ответы:
     успешное создание: "id": 3, "status": "success"
     если нет категории у товара: "error": "pq: insert or update on table \"products\" violates foreign key constraint \"products_category_id_fkey\""
     если не указать обязательные поля name description price: "error": "invalid data" 

DELETE /deleteProduct - удаление товара: Query param id=10 (id товара)
ответы:
     "status": "delete success"
     если нет такого id : "error": "sql: no rows in result set"
     
GET /products - получение списка всех товаров. Два варианта 1 если не авторизоваться(не указать Authorization) то 
ответы:
    "countPages": 3,
    "countProducts": 24,
    "page": 1,
    "pageSize": 10,
    "products": [
        {
            "id": 20,
            "name": "1",
            "description": "Этот топор - верный спутник в любых трудностях. Изготовлен из высококачественной стали, он прослужит вам долгие годы. Прочный, надежный и функциональный, он станет вашим незаменимым помощником в повседневных делах и приключениях. Каждая деталь топора пропитана мастерством ковки, придавая ему особый характер.\nЭтот топор - верный спутник в любых трудностях. Изготовлен из высококачественной стали, он прослужит вам долгие годы. Прочный, надежный и функциональный, он станет вашим незаменимым помощником в повседневных делах и приключениях. Каждая деталь топора пропитана мастерством ковки, придавая ему особый характер.\nЭтот топор - верный спутник в любых трудностях. Изготовлен из высококачественной стали, он прослужит вам долгие годы. Прочный, надежный и функциональный, он станет вашим незаменимым помощником в повседневных делах и приключениях. Каждая деталь топора пропитана мастерством ковки, придавая ему особый характер.\nЭтот топор - верный спутник в любых трудностях. Изготовлен из высококачественной стали, он прослужит вам долгие годы. Прочный, надежный и функциональный, он станет вашим незаменимым помощником в повседневных делах и приключениях. Каждая деталь топора пропитана мастерством ковки, придавая ему особый характер.\nЭтот топор - верный спутник в любых трудностях. Изготовлен из высококачественной стали, он прослужит вам долгие годы. Прочный, надежный и функциональный, он станет вашим незаменимым помощником в повседневных делах и приключениях. Каждая деталь топора пропитана мастерством ковки, придавая ему особый характер.\naaaa\n1",
            "price": 1000,
            "reviews_mid": 0,
            "reviews_count": 0,
            "quantity": 5,
            "work_time": 9,
            "photo": "na servake",
            "category_id": 2,
            "is_active": false
        },
        {
            "id": 21,
            "name": "1",
            "description": "Этот топор - верный спутник в любых трудностях. Изготовлен из высококачественной стали, он прослужит вам долгие годы. Прочный, надежный и функциональный, он станет вашим незаменимым помощником в повседневных делах и приключениях. Каждая деталь топора пропитана мастерством ковки, придавая ему особый характер.\nЭтот топор - верный спутник в любых трудностях. Изготовлен из высококачественной стали, он прослужит вам долгие годы. Прочный, надежный и функциональный, он станет вашим незаменимым помощником в повседневных делах и приключениях. Каждая деталь топора пропитана мастерством ковки, придавая ему особый характер.\nЭтот топор - верный спутник в любых трудностях. Изготовлен из высококачественной стали, он прослужит вам долгие годы. Прочный, надежный и функциональный, он станет вашим незаменимым помощником в повседневных делах и приключениях. Каждая деталь топора пропитана мастерством ковки, придавая ему особый характер.\nЭтот топор - верный спутник в любых трудностях. Изготовлен из высококачественной стали, он прослужит вам долгие годы. Прочный, надежный и функциональный, он станет вашим незаменимым помощником в повседневных делах и приключениях. Каждая деталь топора пропитана мастерством ковки, придавая ему особый характер.\nЭтот топор - верный спутник в любых трудностях. Изготовлен из высококачественной стали, он прослужит вам долгие годы. Прочный, надежный и функциональный, он станет вашим незаменимым помощником в повседневных делах и приключениях. Каждая деталь топора пропитана мастерством ковки, придавая ему особый характер.\naaaa\n1",
            "price": 2333,
            "reviews_mid": 0,
            "reviews_count": 0,
            "quantity": 4,
            "work_time": 9,
            "photo": "na servake",
            "category_id": 2,
            "is_active": false
        }
    ]
    Если авторизоваться(указать токен) добаляются поля показывающие находиться ли товар в избранном или корзине(true да false нет)
    "countPages": 3,
    "countProducts": 24,
    "page": 1,
    "pageSize": 10,
    "products": [
        {
            "id": 1,
            "name": "Мастер",
            "description": "Этот уникальный нож с деревянной рукоятью - идеальный выбор для тех, кто ценит эстетику и функциональность. Рукоять изготовлена из высококачественного дерева, что придает ножу особый шарм и теплоту. Идеально подходит для повседневного использования и станет прекрасным дополнением к вашему кухонному арсеналу.",
            "price": 2990,
            "reviews_mid": 0,
            "reviews_count": 0,
            "quantity": 3,
            "work_time": 6,
            "photo": "/static/catalog/Knife2.png",
            "category_id": 3,
            "is_active": true,
            "IsFavorite": false,
            "IsCart": false
        },
        {
            "id": 3,
            "name": "Листопадный резец",
            "description": "Этот нож - истинное произведение мастерства и красоты. Его узоры на основной части напоминают о таинственности и красоте листьев деревьев в лучах солнца. Изготовленный из высококачественной стали с применением передовых технологий ковки, этот нож является прекрасным сочетанием элегантности и функциональности.",
            "price": 6500,
            "reviews_mid": 0,
            "reviews_count": 0,
            "quantity": 2,
            "work_time": 16,
            "photo": "/static/catalog/Knife4.png",
            "category_id": 3,
            "is_active": true,
            "IsFavorite": false,
            "IsCart": false
        }
        ]
    если товаров нет: "products": []
    Cписок Query Params
    page(номер страницы)
    pageSize(количество элементов на странице)
    sort(сортировки: price_asc(цена возрастания),price_desc(цена убывания),popularity_desc(по популярности))
    price(фильтр по цене, задается диапазон такого вида 3000-5000)
    active(фильтр по наличию, true или false)
    type(передается название категории, например "Нож")
    /products?&page=1&pageSize=10&type=Мачете&active=false&sort=popularity_desc - пример строки с фильром и сортировкой
    

GET /productsCategory - получение товаров по категориям Query param id (id категории)
ответы:
    "products": [
        {
            "id": 2,
            "name": "Katana",
            "description": "AAAAA",
            "price": 1143400,
            "photo": "На серваке",
            "category_id": 1
        },
        {
            "id": 3,
            "name": "Ложка",
            "description": "AAAAA",
            "price": 1143400,
            "photo": "На серваке",
            "category_id": 1
        },
        {
            "id": 4,
            "name": "Меч",
            "description": "Огромный",
            "price": 100,
            "photo": "На серваке",
            "category_id": 1
        },
        {
            "id": 5,
            "name": "",
            "description": "Огромный",
            "price": 100,
            "photo": "На серваке",
            "category_id": 1
        },
        {
            "id": 6,
            "name": "",
            "description": "",
            "price": 100,
            "photo": "На серваке",
            "category_id": 1
        }
    ]
    если нет продуктов: "products": []

POST /cart/createCart - создание корзины юзера: userID береться из jwt token
поля: 
    "product_id":3,
    "count":4
ответы: 
      успешно: "id":4,"status":"success"
      подобные ошибки на то что нет пользователя или товара:  "error": "pq: insert or update on table \"cart_items\" violates foreign key constraint \"cart_items_product_id_fkey\""   

GET /cart/getUserCart: - получение корзины юзера: userID береться из jwt token
ответы:
        успешно:  "cart": [
        {
            "id": 11,
            "user_id": 1,
            "product_id": 11,
            "user_email": "zolafarre.test@mail.ru",
            "product_name": "Тень сакуры",
            "count": 12,
            "price": 83880,
            "photo": "/static/catalog/Katana1.png"
        },
        {
            "id": 12,
            "user_id": 1,
            "product_id": 13,
            "user_email": "zolafarre.test@mail.ru",
            "product_name": "Стальной рубеж",
            "count": 1,
            "price": 5000,
            "photo": "/static/catalog/Axe1.png"
        },
        {
            "id": 13,
            "user_id": 1,
            "product_id": 4,
            "user_email": "zolafarre.test@mail.ru",
            "product_name": "Янтарный ритуал",
            "count": 1,
            "price": 40000,
            "photo": "/static/catalog/Knife5.png"
        }
    ]  
    неуспешно:
        "cart": []

GET /cart/createOrder - оформление заказа, срабатывания отправляется на почту отбивка с товарами которые лежали в корзине и они автоматически удалются, данные о пользователи беруться из Authorization
 
GET /profile/orderHistory - получение истории заказов пользователя(обязательно наличие Authorization)
Ответы:
     "orders": [
        {
            "id": 1,
            "user_id": 1,
            "Products": [
                {
                    "product_id": 24,
                    "count": 2
                },
                {
                    "product_id": 4,
                    "count": 23
                },
                {
                    "product_id": 10,
                    "count": 1
                }
            ],
            "summ": 1113380,
            "product_id": [
                24,
                4,
                10
            ],
            "product_count": [
                2,
                23,
                1
            ]
        },
        {
            "id": 2,
            "user_id": 1,
            "Products": [
                {
                    "product_id": 10,
                    "count": 1
                },
                {
                    "product_id": 4,
                    "count": 32
                },
                {
                    "product_id": 11,
                    "count": 12
                }
            ],
            "summ": 1421480,
            "product_id": [
                10,
                4,
                11
            ],
            "product_count": [
                1,
                32,
                12
            ]
        }
    ]
Если нет заказов: 
 "orders": null

GET /findProductByID - поиск товара по его id: Query param id (id продукта)
поля: 
    "id":5
ответы:
    успешно:
      "name": {
        "id": 5,
        "name": "",
        "description": "Огромный",
        "price": 100,
        "photo": "На серваке",
        "category_id": 1
    }
    неуспешно:
     "error": "record not found"

POST /favorite/addToFavorite - добавление в избранное: userID береться из jwt token
поля:
    "product_id": 3
ответы:
    успешно: "id": , "status": "success"
    нет товара:  "error": "pq: insert or update on table \"favorite_items\" violates foreign key constraint \"favorite_items_product_id_fkey\""
    нет пользователя:  "error": "pq: insert or update on table \"favorite_items\" violates foreign key constraint \"favorite_items_user_id_fkey\""
    если пользователь уже добавил в избранное "error": "user already liked this product"

GET /favorite/favorite - все избранные товары пользователя: userID береться из jwt token
ответы:
       успешно: "favorite": [
        {
            "id": 11,
            "user_id": 1,
            "product_id": 1,
            "user_email": "zolafarre.test@mail.ru",
            "product_name": "Мастер",
            "price": 2990,
            "photo": "/static/catalog/Knife2.png"
        },
        {
           "id": 12,
            "user_id": 1,
            "product_id": 1,
            "user_email": "zolafarre.test@mail.ru",
            "product_name": "Мастер",
            "price": 2990,
            "photo": "/static/catalog/Knife2.png"
        },
        {
            "id": 4,
            "user_id": 1,
            "product_id": 3,
            "user_email": "zolafarre.test@mail.ru",
            "product_name": "Ложка",
            "price": 1143400
            "photo": "/static/catalog/Knife2.png"
        }
    ]
    нет такого пользователя или товаров в корзине: "favorite": []

DELETE /cart/deleteCart, DELETE /favorite/deleteFavorite - удаление из корзины и избранного по id: userID береться из jwt token, Query param productId(id товара для удаления)   
ответы:
    если успешно: "status": "delete success"
    если нет такого id : "error": "sql: no rows in result set"

POST /profile/createProfile - заполнение информации о профиле: создается по умолчанию, все поля пустые кроме name, там лежит email написанный при регистрации
поля: userID береться из jwt token
     "name" : "zxcvzxczxczxczx",
    "surname" : "zxczxc2345234234",
    "patronymic" : "423423432423",
    "contact" : "8916758xfcvxcv32423354",
    "photo" : "на серваке"
ответы:
    успешно: "id": 11,"status": "success"
    несуществующий пользователь "error": "pq: insert or update on table \"profile\" violates foreign key constraint \"profile_user_id_fkey\""
    если не указать обязательные поля name surname contact:  "error": "invalid data"

POST /profile/updateProfile - обновление информации о пользователе: userID береться из jwt token
поля:
   "name" : "111",
    "surname" : "111",
    "patronymic" : "Vazxc222lerich",
    "contact" : "z33xc",
    "photo" : "на сервzxcаке"
ответы:
    успешно: "id": 7777,"status": "success"
   

DELETE /profile/deleteProfile - удаление всей информации профиля Query param id(id пользователя)
ответы:
    если успешно:  "status": "delete success"
    если нет такого id : "error": "sql: no rows in result set"

POST /profile/profile - отображение профиля 
поля:
     "id":17
ответы:
    если профиль есть:
   "profile": {
        "id": 1,
        "user_id": 1,
        "name": "zolafarre.test@mail.ru",
        "surname": "",
        "patronymic": "",
        "contact": "",
        "photo": ""
    }
    если профиля нет: "error": "record not found"

POST /updateProduct - обновление информации о товаре (обновить можно как одно поле так и несколько): 
поля:
    "name": "" 
    "description": "" 
    "price": "" 
    "quantity": "" 
    "work_time": "" 
    "photo": "" 
    "category_id": "" 
    "is_active": "" 
ответы:
    успешно:"id": 3, "status": "success"
    есть такие ошибки, например категория неверная "error": "pq: insert or update on table \"products\" violates foreign key constraint \"products_category_id_fkey\""

POST /sendMsg - отправка сообщения в поддержку по почте, приходит автоматическая отбивка, на почту юзера а также на почту поддержки приходит письмо с заполненными полями:
поля:
    "user_email":"eososks@bk.ru",
    "name":"Александр",
    "message":"Здравствуйте при оформлении заказа вылезает ошибка 404 помогите пожалуйста"
ответы: 
    успешно: "status": "success"
    если почту не указал: "error": "555 5.5.2 Syntax error, cannot decode response. For more information, go to\n5.5.2  https://support.google.com/a/answer/3221692 and review RFC 5321\n5.5.2 specifications. a11-20020a056512020b00b0050eae170e04sm48777lfo.81 - gsmtp"
}
  
POST /createReview - создание отзыва (обязательно наличие Authorization а также нужно заказать товар для того чтобы написать отзыв):
    "product_id":11,
    "stars" : 1,
    "message" : "Мне совершенно не понравился сервис, персонал хамил, фу не приведу сюда детей"
ответы:
    успешно: "id": 11, "status": "success"
    если товар не был заказан "error": "you can't send review"
    если нет токена  "error": "empty auth header"
    если неправильное количество звезд(верное 1-5): "error": "pq: new row for relation \"reviews\" violates check constraint \"reviews_stars_check\""
    
GET /getReviews - получение всех отзывов:
ответы:
    reviews": [
        {
            "id": 2,
            "profile_id": 7,
            "stars": 2,
            "message": "Мне совершенно не понравился сервис, персонал хамил, фу не приведу сюда детей",
            "postdate": "2024-01-10T02:21:23.939551Z",
            "photo": "На серваке",
            "name": "aa",
            "surname": "Plotnaaikov",
            "orders": "Катdfана"
        },
        {
            "id": 9,
            "profile_id": 7,
            "stars": 5,
            "message": "МНЕ ОЧЕНЬ ПОНРАВИЛОСЬ ТРАХНИТЕ МЕНЯ",
            "postdate": "2024-01-10T03:28:25.478617Z",
            "photo": "На серваке",
            "name": "aa",
            "surname": "Plotnaaikov",
            "orders": "Катdfана"
        },
        {
            "id": 10,
            "profile_id": 7,
            "stars": 5,
            "message": "МНЕ ОЧЕНЬ ПОНРАВИЛОСЬ ТРАХНИТЕ МЕНЯ",
            "postdate": "2024-01-10T03:30:13.18217Z",
            "photo": "На серваке",
            "name": "aa",
            "surname": "Plotnaaikov",
            "orders": "Катdfана"
        }
    ]
    если отзывов нет
    reviews": []
DELETE /deleteReview - удаление отзыва Query param(id отзыва)
ответы:
    если успешно:  "status": "delete success"
    если нет такого id : "error": "sql: no rows in result set"
POST /updateReviews обновление отзыва(Обязательно наличие Authorization)
ответы:
    "product_id":11,
    "stars" : 3,
    "message" : "Ну норм"
    
GET /getReview - получение отзывов конретного товара Query param (id)
POST /createGallery - создание галереи 
поля:
      "photo":"на серваке",
    "description":"Санёк с катаной не гей Алексей с мечом"
ответы:
  если все успешно:  "status": "success"
  если что то не заполнить:    "error": "invalid data"
DELETE /deleteGallery - удаление из галереи: 
поля:
    "id":4
ответы:
    если успешно:  "status": "delete success"
    если нет такого id : "error": "sql: no rows in result set"
GET /getGallery - получение всей галереи:
ответы:
    успешно:
    "gallery": [
        {
            "id": 3,
            "catalog": "photo.png",
            "description": "Санёк с катаной не гей"
        },
        {
            "id": 2,
            "catalog": "photo.png",
            "description": "Санёк с катаной не гей Алексей с мечом"
        },
        {
            "id": 8,
            "catalog": "на серваке",
            "description": "Санёк с катаной не гей Алексей с мечом"
        },
        {
            "id": 9,
            "catalog": "на серваке",
            "description": "Санёк с катаной не гей Алексей с мечом"
        }
    ]
неуспешно:
    "gallery": []    

POST /updateGallery - обновление галареи:
поля:
      "id":2,
    "photo":"photo.png",
    "description":"Санёк с катаной не гей Алексей с мечом"
ответы:
    "status": "success"
    
POST /createAddress - создание адреса для гугл мапы:
поля:
    "name":"штаб квартира",
    "latlng":"55.813580, 37.603822"
ответы:
    если успешно: "id": 14
    если неуспешно:  "error": "invalid name or latlng"
        
GET /getAllAddress - получить все адреса:
ответы: 
    "addresses": [
        {
            "id": 9,
            "name": "штаб хата",
            "latlng": "55.813580, 37.603822"
        },
        {
            "id": 10,
            "name": "главная кузница",
            "latlng": "55.813580, 37.603822"
        },
        {
            "id": 11,
            "name": "",
            "latlng": "55.813580, 37.603822"
        },
        {
            "id": 12,
            "name": "",
            "latlng": "55.813580, 37.603822"
        },
        {
            "id": 13,
            "name": "",
            "latlng": "55.813580, 37.603822"
        }
    ]
если неуспешно:
    "addresses": []
    
DELETE /deleteAddress - удаление адреса
поля:
    "id":4
ответы:
    если успешно:  "status": "delete success"
    если нет такого id : "error": "sql: no rows in result set"
POST /updateAddress - обновление адресса:
поля:
    "id":9,
    "name":"штаб хвата",
    "latlng":"55.813580, 37.603822"
ответы:
      "status": "success"
PUT /uploadCatalogPhoto - ключ photo, ошибка если не указать ключ  "error": "request Content-Type isn't multipart/form-data" если не прикрепить файл  "error": "http: no such file"
PUT /uploadProfilePhoto   ключ photo, ошибка если не указать ключ  "error": "request Content-Type isn't multipart/form-data" если не прикрепить файл  "error": "http: no such file"
PUT /uploadReviewsPhoto  ключ photo, ошибка если не указать ключ  "error": "request Content-Type isn't multipart/form-data" если не прикрепить файл  "error": "http: no such file"
PUT /uploadGalleryPhoto  ключ photo, ошибка если не указать ключ  "error": "request Content-Type isn't multipart/form-data" если не прикрепить файл  "error": "http: no such file"
PUT /uploadApks  ключ apk, ошибка если не указать ключ  "error": "request Content-Type isn't multipart/form-data" если не прикрепить файл  "error": "http: no such file"
GET /getApk - получить все апк
ответы:
    или пустой если ничего нет или "files": " 2024-01-10 22:23:35.mp4"
POST /downloadApk - скачать апк:
поля:
     "name":"1473685252276134903.jpg"
POST /deleteCatalogPhoto удаление фото каталога поля:"name":"1473685252276134903.jpg" если нет такого файла "error": "remove ./static/catalog/1473685252276134903.jpg: no such file or directory" если нет имени файла "error": "remove ./static/catalog/: directory not empty"  
POST /deleteProfilePhoto удаление фото профиля  поля:"name":"1473685252276134903.jpg" если нет такого файла "error": "remove ./static/profile/1473685252276134903.jpg: no such file or directory" если нет имени файла "error": "remove ./static/profile/: directory not empty"
POST /deleteReviewsPhoto удаление фото отзыва   поля:"name":"1473685252276134903.jpg" если нет такого файла "error": "remove ./static/reviews/1473685252276134903.jpg: no such file or directory" если нет имени файла "error": "remove ./static/reviews/: directory not empty"
POST /deleteGalleryPhoto удаление фото галереи  поля:"name":"1473685252276134903.jpg" если нет такого файла  "error": "remove ./static/gallery/1473685252276134903.jpg: no such file or directory" если нет имени файла "error": "remove./static/gallery/: directory not empty"
POST /deleteApk поля:"name":"1473685252276134903.jpg" если нет такого файла  "error": "remove ./static/apks/1473685252276134903.jpg: no such file or directory" если нет имени файла "error": "remove ./static/apks/: directory not empty"

GET /get-Photo Query param filename(имя файла которое приходит в ответе) 

POST /learn - заявка на обучение, заполняются поля и на почту приходит отбивка о успешной регистрации заявки(два варианта):
если пользователь авторизован (нужен Authorization):
Поля:
 "message":"Так ну я очень сильно хочу у вас учиться"
успешная отправка:
 "status": "successfully send request"

Если пользователь не авторизован: 
Поля:
     "user_email":"zolafarre.test@mail.ru",
    "name":"Иван",
    "message":"Так ну я очень сильно хочу у вас учиться"
успешная отправка:
 "status": "successfully send request"

Всё что для админов = для юзеров просто на других апишках, чтобы один не имел доступа к другому 

POST /admin-sign-up
POST /admin-sign-in
POST /resetAdminCode
POST /resetAdminPassword
GET /privateAdmin/whoamiAdmin
DELETE /deleteAdminUsers
```
