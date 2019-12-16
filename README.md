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
