# anviz-rpc
Http wrapper for anviz devices

## Запуск: 

    anviz-rpc -host [[host]:port]

*Пример*: 

    anviz-rpc -host localhost:8081

    anviz-rpc -host :8081

## Методы:
  Во все методы необходимо добавить заголовок авторизации: **X-API-Key** 
  
  #### POST /connect
Тело запроса: 
        
    {"host": "192.168.0.1:5010"}
    
Ответ:

    {"data":{"id":5956},"error":null} 
    
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
    
  - new_only - если равно 1, то возвращаются только новые записи. Иначе: все
 
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
  
  - count - если больше 0, то снимает флаг **new** у заданного количества записей. Иначе: у всех
  
Ответ:

- count > 0:
  
      {"data":1,"error":null}
    
- count == 0:
 
       {"data":0,"error":null}
