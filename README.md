# anviz-rpc
Http wrapper for anviz devices

## Конфигурация:

   config.json
   
    {
       "host": "localhost:8081",                       // интерфейс и порт, на котором будет запущен демон 
       "api-key": "DeLr8bFkhmnOJMxRz9xoekzYGUmvef4C"   // токен доступа, который будет использоваться в заголовке "X-API-Key"
    } 

## Запуск: 

    anviz-rpc -с [config file path]

*Пример*: 

    anviz-rpc -c config.json

## Windows сервис
#### Установка (выполняется с правами администратора):
    
    anviz-rpc-x64.exe install -c  path\to\config.json
    anviz-rpc-x64.exe start
    
#### Удалениe (выполняется с правами администратора):
    
    anviz-rpc-x64.exe stop
    anviz-rpc-x64.exe remove
    
#### Запуск приложения в режиме отладки:

    anviz-rpc-x64.exe debug -c  path\to\config.json

## Методы:
  Во все методы необходимо добавить заголовок авторизации: **X-API-Key** 
  
  #### POST /connect
Тело запроса: 
        
    {"host": "192.168.0.1:5010"}
    
Ответ:

    {"data":{"id":5956, "code":"FACE7EI", "type": "face"},"error":null}
типы устройств: face, eye, finger

  #### GET /[device_id:]/disconnect
Ответ:
  
    {"data":null,"error":null}

  #### GET /[device_id:]/state
Ответ:
  
    {"data":{"code":1,"text":"connected"},"error":null}
коды состояний: 0-disconnected, 1-connected, 2-busy

  #### GET /[device_id:]/status
Ответ:

    {"data":{"records":{"users":34,"fingerprints":0,"passwords":2,"cards":33,"all":1475,"new":0},"capacity":{"users":3000,"fingerprints":3000,"records":100000},"state":{"code":1,"text":"connected"}},"error":null}
    
  #### GET /[device_id:]/datetime
Ответ:

    {"data":{"datetime":"16-12-2019 04:46:50"},"error":null}
   
  #### POST /[device_id:]/datetime
Тело запроса: 

    {"datetime":"16-12-2019 04:46:50"}
Ответ:
     
    {"data":null,"error":null}
  #### GET /[device_id:]/records
Параметры:
    
  - **new_only** - если равно ***1***, то возвращаются только новые записи. Иначе: ***все***
 
Ответ:

    {
        "data": [
        {
                "user_id": 2,
                "datetime": "29-08-2019 18:19:56",
                "backup_code": 4,
                "type": "out",
                "attendance_mode": 0,
                "work_types": 0
        },
        {
                "user_id": 2,
                "datetime": "29-08-2019 18:20:10",
                "backup_code": 4,
                "type": "in",
                "attendance_mode": 0,
                "work_types": 0
        }
        ],
        "error":null
    }
  #### GET /[device_id:]/clear_new_records
Параметры:
  
  - **count** - если больше 0, то снимает флаг ***new*** у заданного количества записей. Иначе: ***у всех***
  
Ответ:

- **count > 0**:
  
      {"data":1,"error":null}
    
- **count == 0**:
 
       {"data":0,"error":null}
  #### GET /[device_id:]/users
  
 Ответ:
 
     {
        "data": [
            {            
                "id": 11,
                "password": 0,
                "card_code": 12912225,
                "name": "Тест",
                "department": 0,
                "group": 0,
                "attendance_mode": 248,
                "registered_fp": 0,
                "keep": 45,
                "special_info": 64,
                "is_admin": false
            },
            {
                "id": 12,
                "password": 130397,
                "card_code": 13039720,
                "name": "Виктория",
                "department": 0,
                "group": 0,
                "attendance_mode": 248,
                "registered_fp": 0,
                "keep": 8,
                "special_info": 64,
                "is_admin": false
            }
        ],
        "error":null
    }
  #### POST /[device_id:]/users/add   
  
  Тело запроса:  
  
    {            
        "id": 11,
        "password": 0,
        "card_code": 12912225,
        "name": "Тест",
        "department": 0,
        "group": 0,
        "attendance_mode": 248,
        "registered_fp": 0,
        "keep": 45,
        "special_info": 64,
        "is_admin": false
    }
   
  Ответ:
  
    {"data":null,"error":null}
    
  #### POST /[device_id:]/users/[user_id:]/modify
  
  Тело запроса:
  
    {            
        "password": 0,
        "card_code": 12912225,
        "name": "Тест",
        "department": 0,
        "group": 0,
        "attendance_mode": 248,
        "registered_fp": 0,
        "keep": 45,
        "special_info": 64,
        "is_admin": false
    }
   
  Ответ:
  
    {"data":null,"error":null}
    
   #### GET /[device_id:]/users/[user_id:]/delete
   
Параметры:
  
  - **backup** - если равен ***0xff***, то пользователь удаляется полностью. Иначе: сохраняюся данные соглаcно значению параметра  (подробнее в документации протокола)
   
  Ответ:
  
    {"data":null,"error":null}
  
  
