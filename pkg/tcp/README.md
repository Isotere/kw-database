# Описание

Простая реализация TCP Server-Client

## Описание протокола

Протокол на данный момент очень простой

### Request

MSG_LEN(4 bytes, int32)|MSG_BODY(as bytes)

### Query From Client

Write: MSG_LEN(4 bytes, int32)|MSG_BODY(as bytes)
Read: CODE(4 bytes, int32)|MSG_LEN(4 bytes, int32)|MSG_BODY(as bytes)

### Query From Server

Read: MSG_LEN(4 bytes, int32)|MSG_BODY(as bytes)
Write: CODE(4 bytes, int32)|MSG_LEN(4 bytes, int32)|MSG_BODY(as bytes)

> По сути зеркально, что логично

# НЕ РЕАЛИЗОВАНО

- пулл коннектов
- тротлинг
- таймауты
- да много чего, просто самый простой tcp клиент-сервер
